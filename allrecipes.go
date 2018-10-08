package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
	
	"github.com/go-redis/redis"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

// Recipe is a recipe found in search results
type Recipe struct {
	Name        string   `json:"name,omitempty"`
	Rating      uint8    `json:"rating,omitempty"`
	URL         string   `json:"url,omitempty"`
	Thumb       string   `json:"thumb,omitempty"`
	Author      string   `json:"author,omitempty"`
	Times       string   `json:"times,omitempty"`
	Servings    string   `json:"servings,omitempty"`
	Calories    uint16   `json:"calories,omitempty"`
	Ingredients []string `json:"ingredients,omitempty"`
	Steps       []string `json:"steps,omitempty"`
}

func main() {
	// Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	if pong == "PONG" {
		fmt.Println("Redis client connected")
	}
	// Routing
	router := routing.New()
	router.Get("/recipes/<q>", func(c *routing.Context) error {
		query := c.Param("q")
		var key bytes.Buffer
		key.WriteString("allrecipes:")
		key.WriteString(query)
		c.Response.Header.Set("Content-Type", "application/json")
		exists, err := client.Exists(key.String()).Result()
		if err != nil {
			panic(err)
		} else if exists != 1 {
			results := ParseResults(query)
			marshaledresults, err := json.Marshal(results)
			if err != nil {
				panic(err)
			}
			err = client.Set(key.String(), string(marshaledresults), 72*time.Hour).Err()
			if err != nil {
				panic(err)
			}
			json.NewEncoder(c).Encode(results)
			fmt.Printf("Request: %s\nCached: %s\n", key.String(), "false")
		} else {
			val, err := client.Get(key.String()).Result()
			if err != nil {
				panic(err)
			}
			unmar := new([]Recipe)
			_ = json.Unmarshal([]byte(val), &unmar)
			json.NewEncoder(c).Encode(unmar)
			fmt.Printf("Request: %s\nCached: %s\n", key.String(), "true")
		}
		return nil
	})
	panic(fasthttp.ListenAndServe(":5557", router.HandleRequest))
}
