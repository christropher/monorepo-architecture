package main

import (
	"context"
	"log"
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	product "Product/proto/products"

	faker "github.com/go-faker/faker/v4"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"os"
)

type server struct {
	product.UnimplementedProductServiceServer
	dynamoDBClient *dynamodb.DynamoDB
}

func (s *server) Check(ctx context.Context, in *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	log.Printf("Received Check request: %v", in)
	return &health.HealthCheckResponse{Status: health.HealthCheckResponse_SERVING}, nil
}

func (s *server) Watch(in *health.HealthCheckRequest, _ health.Health_WatchServer) error {
	log.Printf("Received Watch request: %v", in)
	return status.Error(codes.Unimplemented, "unimplemented")
}

func (s *server) GetRecords(ctx context.Context, req *product.ProductRequest) (*product.ProductResponse, error) {
	log.Println("Received request for product ID: " + req.Id)
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Products"),
		Key: map[string]*dynamodb.AttributeValue{
			"ProductId": {
				N: &req.Id,
			},
		},
	}
	result, err := s.dynamoDBClient.GetItem(input)
	if err != nil {
		log.Println("Error getting item from DynamoDB")
		return &product.ProductResponse{
			Id:          req.Id,
			Description: faker.Paragraph(),
			ProductType: faker.Word(),
		}, nil
	}

	return &product.ProductResponse{
		Id:          *result.Item["product_id"].S,
		Description: *result.Item["description"].S,
		ProductType: *result.Item["product_type"].S,
	}, nil
}

func getEnv(name string, defaultValue string) string {
	value, exists := os.LookupEnv(name)
	if !exists {
		return defaultValue
	}
	return value
}

func main() {
	port := getEnv("PORT", "50051")

	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	cs := server{}
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-west-2")})
	if err != nil {
		log.Println("Error with AWS session creation")
	}
	_, err = sess.Config.Credentials.Get()
	if err != nil {
		log.Println("Error with AWS credentials")
	}
	dynamoDBClient := dynamodb.New(sess)

	product.RegisterProductServiceServer(s, &server{
		dynamoDBClient: dynamoDBClient,
	})
	log.Println("Dynamic server started")
	health.RegisterHealthServer(s, &cs)
	reflection.Register(s)

	e := s.Serve(lis)
	if e != nil {
		log.Printf("Failed to start server: %v", e)
	}
}
