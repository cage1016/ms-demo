# gokit microservice demo

| Service | Description           |
| ------- | --------------------- |
| add     | Expose Sum method     |
| tictac  | Expose Tic/Tac method |

## Features

- **[Kubernetes](https://kubernetes.io)/[GKE](https://cloud.google.com/kubernetes-engine/):**
  The app is designed to run on Kubernetes (both locally on "Docker for
  Desktop", as well as on the cloud with GKE).
- **[gRPC](https://grpc.io):** Microservices use a high volume of gRPC calls to
  communicate to each other.
- **[Istio](https://istio.io):** Application works on Istio service mesh.
- **[Skaffold](https://skaffold.dev):** Application
  is deployed to Kubernetes with a single command using Skaffold.
- **[go-kit/kit](https://github.com/go-kit/kit):** Go kit is a programming toolkit for building microservices (or elegant monoliths) in Go. We solve common problems in distributed systems and application architecture so you can focus on delivering business value.

## Install

1. Run `skaffold run` (first time will be slow)
2. Set the `ADD_HTTP_LB_URL/ADD_GRPC_LB_URL` & `TICTAC_HTTP_LB_URL/TICTAC_GRPC_LB_URL` environment variable in your shell to the public IP/port of the Kubernetes loadBalancer
    ```sh
    export ADD_HTTP_LB_PORT=$(kubectl get service add-external -o jsonpath='{.spec.ports[?(@.name=="http")].port}')
    export ADD_GRPC_LB_PORT=$(kubectl get service add-external -o jsonpath='{.spec.ports[?(@.name=="grpc")].port}')
    export ADD_LB_HOST=$(kubectl get service add-external -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
    export ADD_HTTP_LB_URL=$ADD_LB_HOST:$ADD_HTTP_LB_PORT
    export ADD_GRPC_LB_URL=$ADD_LB_HOST:$ADD_GRPC_LB_PORT
    echo $ADD_HTTP_LB_URL
    echo $ADD_GRPC_LB_URL

    export TICTAC_HTTP_LB_PORT=$(kubectl get service tictac-external -o jsonpath='{.spec.ports[?(@.name=="http")].port}')
    export TICTAC_GRPC_LB_PORT=$(kubectl get service tictac-external -o jsonpath='{.spec.ports[?(@.name=="grpc")].port}')
    export TICTAC_LB_HOST=$(kubectl get service tictac-external -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
    export TICTAC_HTTP_LB_URL=$TICTAC_LB_HOST:$TICTAC_HTTP_LB_PORT
    export TICTAC_GRPC_LB_URL=$TICTAC_LB_HOST:$TICTAC_GRPC_LB_PORT
    echo $TICTAC_HTTP_LB_URL
    echo $TICTAC_GRPC_LB_URL
    ```
3. Access by command
    - sum method
    ```sh
    curl -X POST $ADD_HTTP_LB_URL/sum -d '{"a": 1, "b":1}'
    
    or
    
    grpcurl -d '{"a": 1, "b":1}' -plaintext -proto ./pb/add/add.proto $ADD_GRPC_LB_URL pb.Add.Sum
    ```
    - tic method
    ```sh
    curl -X POST $TICTAC_HTTP_LB_URL/tic
    
    or
    
    grpcurl -plaintext -proto ./pb/tictac/tictac.proto $TICTAC_GRPC_LB_URL pb.Tictac.Tic
    ```
    - tac method
    ```sh
    curl $TICTAC_HTTP_LB_URL/tac
    
    or
    
    grpcurl -plaintext -proto ./pb/tictac/tictac.proto $TICTAC_GRPC_LB_URL pb.Tictac.Tac
    ```
4. Apply istio manifests `kubectl apply -f deployments/istio-manifests`
5. Set the `GATEWAY_HTTP_URL/GATEWAY_GRPC_URL` environment variable in your shell to the public IP/port of the Istio Ingress gateway.
    ```sh
    export INGRESS_HTTP_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].port}')
    export INGRESS_GRPC_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].port}')
    export INGRESS_HOST=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
    export GATEWAY_HTTP_URL=$INGRESS_HOST:$INGRESS_HTTP_PORT
    export GATEWAY_GRPC_URL=$INGRESS_HOST:$INGRESS_GRPC_PORT
    echo $GATEWAY_HTTP_URL
    echo $GATEWAY_GRPC_URL
    ```
7. Access by command
    - sum method
    ```sh
    curl -X POST $GATEWAY_HTTP_URL/api/v1/add/sum -d '{"a": 1, "b":1}'
    
    or
    
    grpcurl -d '{"a": 1, "b":1}' -plaintext -proto ./pb/add/add.proto $GATEWAY_GRPC_URL pb.Add.Sum
    ```
    - tic method
    ```sh
    curl -X POST $GATEWAY_HTTP_URL/api/v1/tictac/tic
    
    or
    
    grpcurl -plaintext -proto ./pb/tictac/tictac.proto $GATEWAY_GRPC_URL pb.Tictac.Tic
    ```
    - tac method
    ```sh
    curl $GATEWAY_HTTP_URL/api/v1/tictac/tac
    
    or
    
    grpcurl -plaintext -proto ./pb/tictac/tictac.proto $GATEWAY_GRPC_URL pb.Tictac.Tac
    ```


## CleanUP

`skaffold delete`

or 

`kubectl delete -f deployments/istio-manifests`