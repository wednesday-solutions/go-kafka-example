echo "name of the service is: $1 $2"

copilot init -a "go-kafka-example" -t "Load Balanced Web Service" -n "producer-svc" -d ./Dockerfile

copilot env init --name "dev" --profile default --default-config

copilot storage init -n "gke-producer-svc-cluster" -t Aurora -w "producer-svc" --engine PostgreSQL --initial-db "gke_producer_db"

copilot deploy --name "producer-svc" -e "dev"

