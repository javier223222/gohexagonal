// src/api/main.go
package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "os"
    "danierpclone/src/admin"
	"danierpclone/shared/database"
	"danierpclone/shared/middleware"
)

func main() {
    // Cargar variables de entorno
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Configurar la conexión a la base de datos
    dbHost := os.Getenv("DB_HOST")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbPort := os.Getenv("DB_PORT")

    db, err := database.ConnectMySQL(dbHost, dbUser, dbPassword, dbName, dbPort)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Inyección de dependencias
    adminRepo := admin.NewMySQLAdminRepository(db)
    adminService := admin.NewAdminService(adminRepo)
    adminHandler := admin.NewAdminHandler(adminService)

    // Configuración del servidor Gin
    r := gin.Default()
    r.POST("/admin",middleware.AuthMiddleware(),middleware.AdminMiddleware(), adminHandler.Save)
	r.GET("/admin",middleware.AuthMiddleware(),middleware.AdminMiddleware(), adminHandler.Get)
	r.DELETE("/admin/:id",middleware.AuthMiddleware(),middleware.AdminMiddleware(), adminHandler.Delete)
	r.DELETE("/admin",middleware.AuthMiddleware(), adminHandler.DeleteMYOnCount)
	r.PATCH("/admin:id",middleware.AuthMiddleware(),middleware.AdminMiddleware(), adminHandler.UpdatePasswordByID)
	r.PATCH("/admin",middleware.AuthMiddleware(), adminHandler.UpdatePassword)
	
	
	r.POST("/admin/login", adminHandler.Login)


    r.Run(":8080")
}
