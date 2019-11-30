package main

import (
	"log"
	"os"

	"github.com/go-redis/redis"
)

var sessionLog *log.Logger
var sessionStore *redis.Client

func init() {
	sessionLog = log.New(os.Stdout, "SRS Redis: ", log.LstdFlags)
}

func ConnectRedis() {
	sessionStore = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PW"),
		DB:       0,
	})

	_, err := sessionStore.Ping().Result()
	if err != nil {
		sessionLog.Fatalf("%s: %s\n", "Error connecting to Redis", err.Error())
	}
	sessionLog.Println("Connected to Redis.")
}
