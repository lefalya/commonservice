package connection

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func ConnectRedis(redisHostAddr string, password string, isClustered bool) redis.UniversalClient {
	if redisHostAddr == "" {
		log.Fatal("REDIS_HOST environment variable not set")
	}

	if isClustered {
		clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    []string{redisHostAddr},
			Password: password,
		})

		_, err := clusterClient.Ping(context.Background()).Result()
		if err != nil {
			log.Fatal(err)
		}

		return clusterClient
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisHostAddr,
		Password: password,
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func ConnectMongo(mongoHostAddr string) *mongo.Client {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	client, errConnect := mongo.Connect(ctx, options.Client().ApplyURI(mongoHostAddr))

	if errConnect != nil {
		panic(errConnect)
	}

	if errPing := client.Ping(ctx, readpref.Primary()); errPing != nil {
		panic(errPing)
	}

	return client
}
