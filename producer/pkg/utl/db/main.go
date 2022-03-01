package db

import (
	"encoding/json"
	"os"
)

type DbDetails struct {
	DbClusterIdentifier string `json:"dbClusterIdentifier"`
	Password            string `json:"password"`
	Dbname              string `json:"dbname"`
	Engine              string `json:"engine"`
	Port                int64  `json:"port"`
	Host                string `json:"host"`
	Username            string `json:"username"`
}

func GetDbDetails() (DbDetails, error) {
	dbDetails := os.Getenv("GKEPRODUCERSVCCLUSTER_SECRET")
	dbDetailsMap := DbDetails{}
	err := json.Unmarshal([]byte(dbDetails), &dbDetailsMap)
	if err != nil {
		return DbDetails{}, err
	}
	return dbDetailsMap, nil
}
