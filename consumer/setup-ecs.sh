echo "name of the service is: $1 $2"

copilot init -a "go-kafka-example" -t "Load Balanced Web Service" -n "consumer-svc" -d ./Dockerfile

copilot env init --name "dev" --profile default --default-config

copilot storage init -n "gke-consumer-svc-cluster" -t Aurora -w "consumer-svc" --engine PostgreSQL --initial-db "gke_consumer_db"

copilot deploy --name "consumer-svc" -e "dev"
