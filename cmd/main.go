package main

import (
	"flag"
	"fmt"
	"role-helper/cfg"
	"role-helper/internal/delivery/http"
	"role-helper/service"
)

func main() {
	configPath := flag.String("config-path", "./cfg/cfg.yaml", "path to config file")
	flag.Parse()

	config, err := cfg.GetConfig(*configPath)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("config parsed")
	}

	dbPostgres, err := service.InitPostgres(config)
	if err != nil {
		panic(err)
	}
	fmt.Println("postgres connected")

	redisClient, err := service.InitRedis(config, 0)
	if err != nil {
		panic(err)
	}
	fmt.Println("redis session connected")

	server := httpserver.NewHTTPServer()
	err = server.Start(config, dbPostgres, redisClient)
	if err != nil {
		panic(err)
	}
}
