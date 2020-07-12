package main

import (
	"fmt"
	"os"
	)

func main() {
	a := App{}
	pwd, exists1 := os.LookupEnv("REDIS_PASSWORD")
	addrs, exists2 := os.LookupEnv("REDIS_ENDPOINT")
	port, _ := os.LookupEnv("PORT")
	if port == "" {
		port = "8080"
	}
	if exists1 && exists2 {
		a.Initialize(pwd, addrs)
		a.Run(":" + port)
	}else{
		_ = fmt.Errorf("needs env var REDIS_PASSWORD")
	}
}
