# cec-gateway
Close Encounters Corps Gateway

## Installation

### Kubernetes

`helm install cec-gw ./helm/`

## Development

### If you change swagger.yaml

`swagger generate server -t gen -f swagger.yaml --exclude-main -A cec-gw`

## Usage

Just open in your browser `/api/docs`. 