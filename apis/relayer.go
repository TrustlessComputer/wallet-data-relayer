package apis

import "github.com/gin-gonic/gin"

type StoreData struct {
	Data string `json:"data"`
	ID   string `json:"id"`
}

func (r *Router) PostData(c *gin.Context) {
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

	err = r.Redis.SetStringDataWithExpTime(data.ID, data.Data, 600)
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
	data, err := r.Redis.GetData(id)
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
	c.JSON(200, gin.H{"data": data})
}
