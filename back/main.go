package main

import (
	"chart/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func update(c *gin.Context) {
	err := utils.UpdateData()
	if err == nil {
		c.JSON(200, gin.H{
			"status": "success",
		})
	} else {
		c.JSON(200, gin.H{
			"status": "fail",
		})
	}
}

func getResult(c *gin.Context) {
	dataset := c.Query("dataset")
	ratio := c.QueryArray("ratio[]")
	result := utils.GetData(dataset, ratio)

	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

func getError(c *gin.Context) {
	dataset := c.Query("dataset")
	ratio := c.QueryArray("ratio[]")
	n_clusters := c.QueryArray("n_clusters[]")
	error := utils.GetError(dataset, ratio, n_clusters)

	c.JSON(http.StatusOK, gin.H{
		"error": error,
	})
}

func getInfo(c *gin.Context) {
	info := utils.GetSlotnameDistribution()
	c.JSON(
		http.StatusOK,
		gin.H{
			"info": info,
		},
	)
}

func getCluster(c *gin.Context) {
	dataset := c.Query("dataset")
	metric := c.Query("metric")
	feature := c.Query("feature")

	cluster := utils.GetClusterInfo(dataset, metric, feature)
	c.JSON(
		http.StatusOK,
		gin.H{
			"info": cluster,
		},
	)
}

func main() {
	r := gin.Default()
	api := r.Group("/api")
	{
		api.GET("/update", update)
		api.GET("/result", getResult)
		api.GET("/error", getError)
		api.GET("/info", getInfo)
		api.GET("/cluster", getCluster)
	}
	fmt.Println()

	r.Run(":8081") // 监听并在 0.0.0.0:8080 上启动服务
}
