package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	mutex      sync.Mutex
	method1Ch  = make(chan bool, 1)
	method2Ch  = make(chan bool, 1)
	userDB     = make(map[int]string)
)

func main() {
	r := gin.Default()

	r.POST("/methods", methodHandler)

	// Start the server
	r.Run(":8080")
}


func methodHandler(c *gin.Context) {
	var requestData struct {
		Method   int `json:"method"`
		WaitTime int `json:"waitTime"`
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	switch requestData.Method {
	case 1:
		method1(c, requestData.WaitTime)
	case 2:
		method2(c, requestData.WaitTime)
	default:
		c.JSON(400, gin.H{"error": "Invalid method"})
	}
}


func method1(c *gin.Context, waitTime int) {
	mutex.Lock()
	defer mutex.Unlock()

	select {
	case method1Ch <- true:
		
	default:
		
		<-method1Ch
		method1Ch <- true
	}

	
	time.Sleep(time.Duration(waitTime) * time.Second)

	userData := getUsersData()

	<-method1Ch

	c.JSON(200, gin.H{"users": userData})
}

func method2(c *gin.Context, waitTime int) {
	select {
	case method2Ch <- true:
		
	default:
		
	}

	
	time.Sleep(time.Duration(waitTime) * time.Second)

	
	userData := getUsersData()

	<-method2Ch

	c.JSON(200, gin.H{"users": userData})
}

func getUsersData() map[int]string {
	
	userData := make(map[int]string)
	for id, name := range userDB {
		userData[id] = name
	}
	return userData
}
