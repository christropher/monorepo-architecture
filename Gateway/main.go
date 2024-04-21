package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	product "Gateway/proto/products"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func loadCertificate() (credentials.TransportCredentials, error) {
	certFile := "ssl/ca.crt"
	return credentials.NewClientTLSFromFile(certFile, "")
}

func buildClientOpts() (grpc.DialOption, error) {
	tls := false
	if tls {
		creds, sslErr := loadCertificate()
		if sslErr != nil {
			log.Fatalf("Failed to load certificate: %v", sslErr)
			return nil, sslErr
		}
		return grpc.WithTransportCredentials(creds), nil
	} else {
		return grpc.WithInsecure(), nil
	}
}

func getEnv(name string, defaultValue string) string {
	value, exists := os.LookupEnv(name)
	if !exists {
		return defaultValue
	}
	return value
}

var productClient product.ProductServiceClient

func main() {
	hostname := getEnv("PRODUCT_HOST", "localhost:50051")
	log.Println("Product Hostname: " + hostname)
	port := getEnv("PORT", "8080")

	conn, err := grpc.Dial(hostname, grpc.WithInsecure())
	if err != nil {
		log.Println("Error connecting to product service: " + err.Error())
		return
	}

	defer conn.Close()

	productClient = product.NewProductServiceClient(conn)

	log.Println("Listening to port: " + port)
	http.HandleFunc("/", handle)
	http.ListenAndServe("0.0.0.0:"+port, nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("Id")
	log.Println("Id: " + id)
	log.Println("URL" + r.URL.String())
	message := getProduct(id)
	fmt.Fprint(w, message)
}

func getProduct(id string) string {
	if id == "" {
		id = "1"
	}
	for {
		res, err := productClient.GetRecords(context.Background(), &product.ProductRequest{Id: id})
		if err == nil {
			fmt.Println("Description:", res.Description)
			return res.Description
		} else {
			log.Println("Error getting product: " + err.Error())
			return "(Error)"
		}
	}
}
