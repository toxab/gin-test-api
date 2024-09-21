package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type User struct {
    ID    uint   `json:"id" gorm:"primaryKey"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func connectToDatabase(dsn string) (*gorm.DB, error) {
    var db *gorm.DB
    var err error
    retries := 5

    for i := 0; i < retries; i++ {
        db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
        if err == nil {
            fmt.Println("Database connection established")
            return db, nil
        }

        fmt.Printf("Could not connect to the database: %v. Retrying...\n", err)
        time.Sleep(5 * time.Second)
    }

    return nil, fmt.Errorf("failed to connect to the database after %d attempts", retries)
}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    mysqlHost := os.Getenv("MYSQL_HOST")
    mysqlPort := os.Getenv("MYSQL_PORT")
    mysqlUser := os.Getenv("MYSQL_USER")
    mysqlPassword := os.Getenv("MYSQL_PASSWORD")
    mysqlDB := os.Getenv("MYSQL_DB")

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDB)

    db, err := connectToDatabase(dsn)
    if err != nil {
        log.Fatalf("Could not connect to the database: %v", err)
    }

    db.AutoMigrate(&User{})

    // Set up the Gin router
    r := gin.Default()

    // Trust all proxies (this is the default behavior, can be customized)
    r.SetTrustedProxies(nil)

    r.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Welcome to the Gin API!",
        })
    })

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "pong",
        })
    })

    r.POST("/users", func(c *gin.Context) {
        var user User

        // Bind the JSON payload to the user struct
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if err := db.Create(&user).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, user)
    })

    r.PUT("/users/:id", func(c *gin.Context) {
        var user User

        if err := db.First(&user, c.Param("id")).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }

        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        db.Save(&user)
        c.JSON(http.StatusOK, user)
    })

    r.GET("/users", func(c *gin.Context) {
        var users []User
        if result := db.Find(&users); result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
            return
        }
        c.JSON(http.StatusOK, users)
    })

    r.Run()
}
