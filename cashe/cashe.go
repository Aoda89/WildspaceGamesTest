package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	urlCash = "localhost:4040"
)

type Product struct {
	Id    string `json:"id"`
	Price int    `json:"price"`
}

var (
	mu      sync.Mutex
	storage map[string]int
)

func main() {
	storage = make(map[string]int)
	router := gin.Default()
	go func() {
		router.GET("/", sendData)
	}()
	go func() {
		updateCashe()
	}()

	router.Run(urlCash)
}

func updateCashe() {
	config, err := pgx.ParseConfig("postgres://admin:admin@localhost:5432/postgres")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}

	var (
		storageSize = 0
		BdSize      = 0
		product     Product
	)
	for {
		err := conn.QueryRow(context.Background(), "SELECT count(*) FROM products").Scan(&BdSize)
		if err != nil {
			log.Println(error(err))
		}
		if storageSize < BdSize {
			difference := BdSize - storageSize
			rows, err := conn.Query(context.Background(), "SELECT * FROM products LIMIT $1 OFFSET $2", difference, storageSize)
			if err != nil {
				log.Println(err.Error())
			}
			for rows.Next() {
				err := rows.Scan(&product.Id, &product.Price)
				if err != nil {
					log.Fatal(err)
				}
				mu.Lock()
				storage[product.Id] = product.Price
				mu.Unlock()
			}
		}
		time.Sleep(time.Minute)
	}
}

func sendData(c *gin.Context) {
	var count int
	var product []Product
	var countProd = 0

	limit := c.Query("limit")
	offset := c.Query("offset")

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		log.Fatal(err)
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		log.Fatal(err)
	}

	end := offsetInt + limitInt
	if end > len(storage) {
		log.Fatal("Выход за границу")
	}

	mu.Lock()
	for key, value := range storage {
		if count >= offsetInt && count <= end {
			product[countProd].Id = key
			product[countProd].Price = value
		}
		count++
		countProd++
	}
	mu.Unlock()

	c.JSON(http.StatusOK, product)
}
