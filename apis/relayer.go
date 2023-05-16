package apis

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type StoreData struct {
	Data    string `json:"data"`
	Message string `json:"message"`
	Site    string `json:"site"`
	ID      string `json:"id"`
}

func (r *Router) PostData(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	var data StoreData
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Bad Request",
		})
		return
	}
	if data.ID == "" || data.Data == "" {
		c.JSON(400, gin.H{
			"message": "Bad Request",
		})
		return
	}
	if origin == "" {
		origin = c.Request.Header.Get("origin")
		if origin == "" {
			c.JSON(400, gin.H{
				"message": "Bad Request",
			})
			return
		}
	}
	data.Site = origin

	err = r.Redis.SetDataWithExpireTime(data.ID+"_data", data, 120)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	c.JSON(200, gin.H{"message": "OK"})
}

func (r *Router) GetData(c *gin.Context) {
	id := c.Query("id")

	if id == "" {
		c.JSON(400, gin.H{
			"message": "Bad Request",
		})
		return
	}
	data, err := r.Redis.GetData(id + "_data")
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	if data == nil {
		c.JSON(404, gin.H{
			"message": "Not Found",
		})
		return
	}
	var storeData StoreData
	err = json.Unmarshal([]byte(*data), &storeData)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	c.JSON(200, gin.H{"data": storeData})
}

func (r *Router) PostResult(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	var data StoreData
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Bad Request",
		})
		return
	}
	if data.ID == "" || data.Data == "" {
		c.JSON(400, gin.H{
			"message": "Bad Request",
		})
		return
	}
	if origin == "" {
		origin = c.Request.Header.Get("origin")
		if origin == "" {
			c.JSON(400, gin.H{
				"message": "Bad Request",
			})
			return
		}
	}
	data.Site = origin

	err = r.Redis.SetDataWithExpireTime(data.ID+"_result", data, 120)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	c.JSON(200, gin.H{"message": "OK"})
}

func (r *Router) GetResult(c *gin.Context) {
	id := c.Query("id")

	if id == "" {
		c.JSON(400, gin.H{
			"message": "Bad Request",
		})
		return
	}
	data, err := r.Redis.GetData(id + "_result")
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	if data == nil {
		c.JSON(404, gin.H{
			"message": "Not Found",
		})
		return
	}
	var storeData StoreData
	err = json.Unmarshal([]byte(*data), &storeData)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	c.JSON(200, gin.H{"data": storeData})
}
