## Requirements
1. Docker - https://docs.docker.com/get-docker/
2. AWS Account
3. AWS CLI - https://aws.amazon.com/cli/

## Deploy
1. create an .env file with
    ```sh
        export PROJECT_NAME=aws-appmesh-demo
        export AWS_DEFAULT_REGION=
        export AWS_ACCOUNT_ID=
        export ENVOY_IMAGE=public.ecr.aws/appmesh/aws-appmesh-envoy:v1.27.3.0-prod-fips
        export KEY_PAIR=
    ```
2. run `cd ./Infra && make build-infra`

## Diagrams

### High Level Architecture

![app_mesh_general drawio](https://github.com/christropher/MONOREPO/assets/19397775/d3b17eaa-800e-4064-ad66-34e6513f9430)

### CloudFormation Architecture
![application-composer-2024-04-16T175232 101Zyd3-base-infra yaml](https://github.com/christropher/MONOREPO/assets/19397775/e74489a7-4ca1-4eb9-b89f-c570978a69da)
![application-composer-2024-04-16T175319 720Zb24-app-mesh yaml](https://github.com/christropher/MONOREPO/assets/19397775/dc7d0dd7-7f0e-4db2-ad21-3f2a72a47cc1)
![application-composer-2024-04-16T175409 115Zzln-containers yaml](https://github.com/christropher/MONOREPO/assets/19397775/d1ea1ee0-8978-4c8d-ade5-8e9779aed948)

