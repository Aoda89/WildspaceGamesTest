package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	urlCash = "http://localhost:4040"
	urlApp  = "localhost:8080"
)

func main() {
	router := gin.Default()
	go func() {
		router.GET("/", getData)
	}()
	router.Run(urlApp)
}

func getData(c *gin.Context) {

	limit := c.Query("limit")
	offset := c.Query("offset")
	getCash(urlCash, limit, offset)
	log.Println(limit, offset)
}

func getCash(url, limit, offset string) {

	jquery := fmt.Sprintf("%s/?limit=%s&offset=%s", url, limit, offset)
	response, err := http.Get(jquery)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
}
