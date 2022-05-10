package main

import (
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
)

type product struct {
	ID          string `json:"id"`
	ProductName string `json:"name"`
	Provider    string `json:"provider"`
	Quantity    int    `json:"quantity"`
}

var products = []product{
	{ID: "1", ProductName: "Coca-cola", Provider: "Coca-cola", Quantity: 3},
	{ID: "2", ProductName: "Kinder", Provider: "KinderBueno", Quantity: 5},
	{ID: "3", ProductName: "Ferrero", Provider: "KinderBueno", Quantity: 1},
}

func getProducts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, products)
}

func productById(c *gin.Context) {
	id := c.Param("id")
	product, err := getProductById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, product)
}

func checkOutProduct(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id Query"})
		return
	}

	product, err := getProductById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Project not found."})
	}

	if product.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Project not available"})
		return
	}

	product.Quantity -= 1
	c.IndentedJSON(http.StatusOK, product)
}

func returnProduct(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id Query"})
		return
	}

	product, err := getProductById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Project not found."})
	}

	product.Quantity += 1
	c.IndentedJSON(http.StatusOK, product)
}

func getProductById(id string) (*product, error) {
	for i, p := range products {
		if p.ID == id {
			return &products[i], nil
		}
	}

	return nil, errors.New("project not found")
}

func createProduct(c *gin.Context) {
	var newProduct product

	if err := c.BindJSON(&newProduct); err != nil {
		return
	}

	products = append(products, newProduct)
	c.IndentedJSON(http.StatusCreated, newProduct)
}

func main() {
	router := gin.Default()
	router.GET("/products", getProducts)
	router.GET("/products/:id", productById)
	router.POST("/products", createProduct)
	router.PATCH("/checkout", checkOutProduct)
	router.PATCH("/return", returnProduct)
	router.Run("localhost:8080")
}
