package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "net/http"

    "github.com/gin-gonic/gin"
    "google.golang.org/grpc"

    "myapp/models"
    "myapp/grpc_server"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    pb "myapp/pkg/.." // импорт сгенерированного proto пакета
)

func main() {
    // Подключение к PostgreSQL через GORM
    dsn := "host=localhost user=youruser password=yourpass dbname=yourdb port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("failed to connect database:", err)
    }

    // Автоматическая миграция таблицы User
    db.AutoMigrate(&models.User{})

    // Запуск gRPC сервера в отдельной горутине
    go startGRPCServer(db)

    // Создаем Gin роутер
    r := gin.Default()

    // REST API: получить всех пользователей
    r.GET("/users", func(c *gin.Context) {
        var users []models.User
        db.Find(&users)
        c.JSON(http.StatusOK, users)
    })

    // REST API: добавить пользователя
    r.POST("/users", func(c *gin.Context) {
        var user models.User
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        db.Create(&user)
        c.JSON(http.StatusCreated, user)
    })

    // Запуск HTTP сервера на порту 8080
    r.Run(":8080")
}

func startGRPCServer(db *gorm.DB) {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

   grpcServer := grpc.NewServer()
   pb.RegisterUserServiceServer(grpcServer, &services.GRPCServer{DB: db})

   fmt.Println("gRPC server listening on :50051")
   if err := grpcServer.Serve(lis); err != nil {
       log.Fatalf("failed to serve gRPC: %v", err)
   }
}
