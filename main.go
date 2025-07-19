package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

type LEDStatus struct {
    Status    string    `json:"status"`   
    Timestamp time.Time `json:"timestamp"`
}

type LEDInfo struct {
    Status    string    `json:"status"`   
    Color     string    `json:"color"`
}

var ledHistoryBlue []LEDStatus
var ledHistoryRed []LEDStatus
var ledHistoryGreen []LEDStatus

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "https://gofe-production.up.railway.app")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func main() {
    r := gin.Default()
	r.Use(CORSMiddleware())
    
    r.POST("/led", func(c *gin.Context) {
        var newStatus LEDInfo
        
        if err := c.ShouldBindJSON(&newStatus); err != nil {
            c.JSON(400, gin.H{"error": "Błędny JSON"})
            return
        }

		var status LEDStatus
        status.Timestamp = time.Now()
		status.Status = newStatus.Status
        
        switch newStatus.Color {
        case "blue":
            ledHistoryBlue = append(ledHistoryBlue, status)
        case "red":
            ledHistoryRed = append(ledHistoryRed, status)
        case "green":
            ledHistoryGreen = append(ledHistoryGreen, status)
        default:
            c.JSON(400, gin.H{"error": "Nieznany kolor"})
            return
        }

        
        c.JSON(200, gin.H{
            "message": "Status zapisany!",
            "status":  newStatus.Status,
        })
    })

	r.GET("/ledHistory/blue", func(c *gin.Context) {
        c.JSON(200, ledHistoryBlue)
    })


    r.GET("/led/blue", func(c *gin.Context) {
        c.JSON(200, ledHistoryBlue[len(ledHistoryBlue)-1].Status)
    })

	r.GET("/led/red", func(c *gin.Context) {
        c.JSON(200, ledHistoryRed[len(ledHistoryRed)-1].Status)
    })

	r.GET("/led/green", func(c *gin.Context) {
        c.JSON(200, ledHistoryGreen[len(ledHistoryGreen)-1].Status)
    })

    r.Run(":8080")
}