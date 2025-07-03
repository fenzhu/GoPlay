package main

import (
	"io/ioutil"
	"log"
	"net/http" // Import net/http for pprof
	_ "net/http/pprof" // Import pprof for HTTP endpoints

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"seckill/internal/api"
	"seckill/internal/repository"
	"seckill/internal/service"
)

func main() {
	// Set Gin to Release Mode for production
	gin.SetMode(gin.ReleaseMode)

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

	// Register pprof handlers on a separate port
	// In production, consider exposing this on an internal network or with authentication.
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil)) // pprof will listen on port 6060
	}()

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