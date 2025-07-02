
package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
	"seckill/internal/api"
	"seckill/internal/repository"
	"seckill/internal/service"
)

func main() {
	// Load config
	config, err := ioutil.ReadFile("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var configInterfaceMap map[interface{}]interface{}
	if err := yaml.Unmarshal(config, &configInterfaceMap); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	configMap := convertKeysToStrings(configInterfaceMap)

	// Init DB
	if err := repository.InitDB(configMap); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Init Redis
	if err := repository.InitRedis(configMap); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Sync product stock to Redis
	if err := service.SyncProductStockToRedis(); err != nil {
		log.Fatalf("Failed to sync product stock to Redis: %v", err)
	}

	router := api.NewRouter()
	router.Run(":8080")
}

func convertKeysToStrings(m map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range m {
		switch v2 := v.(type) {
		case map[interface{}]interface{}:
			res[k.(string)] = convertKeysToStrings(v2)
		default:
			res[k.(string)] = v
		}
	}
	return res
}
