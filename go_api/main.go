package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// --- 1. Definição da Entidade (Modelo) ---
// As "tags" (ex: `json:"name"`) definem como os dados aparecem no JSON e no Banco.
type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"not null" json:"name"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	User     string `gorm:"uniqueIndex;not null" json:"user"`
	Password string `gorm:"not null" json:"password"`
}

var db *gorm.DB

// --- 2. Conexão Otimizada com o Banco ---
func connectDatabase() {
	// Lê as variáveis de ambiente que definiremos no docker-compose
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	// Loop de retry: Tenta conectar 5 vezes caso o banco demore a subir
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Tentando conectar ao banco (%d/5)...", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		panic("Erro fatal: Não foi possível conectar ao PostgreSQL!")
	}

	// Cria a tabela 'users' automaticamente
	db.AutoMigrate(&User{})

	// --- PERFORMANCE TUNING ---
// --- PERFORMANCE TUNING ---
    sqlDB, _ := db.DB()

    // MELHORIA 4: Aumentar conexões em espera e máximas
    sqlDB.SetMaxIdleConns(20)   // Era 10
    sqlDB.SetMaxOpenConns(80)   // Era 100 (Reduzi um pouco por segurança pois temos 4 replicas: 4*80=320)
    sqlDB.SetConnMaxLifetime(time.Hour)
}

// --- 3. Handlers (Funções das Rotas) ---

func createUser(c *gin.Context) {
	var input User
	// Valida o JSON recebido
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Tenta salvar no banco
	if result := db.Create(&input); result.Error != nil {
		// Retorna erro se email/user já existirem
		c.JSON(http.StatusConflict, gin.H{"error": "User or Email already exists"})
		return
	}
	c.JSON(http.StatusCreated, input)
}

func getUsers(c *gin.Context) {
	var users []User
	db.Find(&users)
	c.JSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {
	var user User
	// Busca pelo ID passado na URL
	if err := db.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func updateUser(c *gin.Context) {
	var user User
	if err := db.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&user).Updates(input)
	c.JSON(http.StatusOK, user)
}

func deleteUser(c *gin.Context) {
	var user User
	if err := db.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	db.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// --- 4. Função Principal ---
func main() {
	connectDatabase()

	// Define modo de produção (remove logs de debug, melhora performance)
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()        // Cria router sem middlewares padrão
	r.Use(gin.Recovery()) // Adiciona apenas recuperação de pânico (mais leve)

	// Rotas
	r.POST("/users", createUser)
	r.GET("/users", getUsers)
	r.GET("/users/:id", getUser)
	r.PUT("/users/:id", updateUser)
	r.DELETE("/users/:id", deleteUser)

	// Roda na porta 8080
	r.Run(":8080")
}