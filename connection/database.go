package connection

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

func ConnectRedis() redis.UniversalClient {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		log.Fatal("REDIS_HOST environment variable not set")
	}

	if os.Getenv("APP_ENV") == "production" {
		clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    []string{redisHost},
			Password: os.Getenv("REDIS_PASS"),
		})

		_, err := clusterClient.Ping(context.Background()).Result()
		if err != nil {
			log.Fatal(err)
		}

		return clusterClient
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func ConnectMongo() *mongo.Client {
	URI := os.Getenv("MONGO_HOST")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	client, errConnect := mongo.Connect(ctx, options.Client().ApplyURI(URI))

	if errConnect != nil {
		panic(errConnect)
	}

	// Auto-reconnect
	go func() {
		for {
			select {
			case <-time.After(1 * time.Minute):
				// Ping
				err := client.Ping(context.Background(), nil)
				if err != nil {
					log.Println("Reconnecting to MongoDB...")
					err = client.Connect(context.Background())
					if err != nil {
						log.Println("Failed to reconnect:", err)
					} else {
						log.Println("Reconnected to MongoDB successfully.")
					}
				}
			}
		}
	}()

	return client
}
