package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	Username        string `json:"username"`
	Coins           int    `json:"coins" `
	Point_of_3cards int    `json:"point_Of_3Cards"`
}
type Deck struct {
	Deck_id   int `json:"deck_id"`
	Remaining int `json:"remaining"`
}

type Card struct {
	Card_value int    `json:"card_value"`
	Card_image string `json:"card_image"`
	Status     bool   `json:"status"`
	Deck_id    int    `json:"deck_id"`
}

func (User) TableName() string { return "users" }
func (Deck) TableName() string { return "decks" }
func (Card) TableName() string { return "cards" }
func main() {
	dsn := "root:109339Lam@@tcp(127.0.0.1:3306)/GameBaiCao?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("test")
	}
	router := gin.Default()
	v1 := router.Group("/v1") // API For User
	{
		v1.POST("/users", createUser(db))
		v1.GET("/users", getListOfUser(db))
		v1.GET("/users/:username", readUserByUsername(db))
		v1.PUT("/users/:username", editUserByUsername(db))
		v1.DELETE("/users/:username", deleteUserByUsername(db))
	}

	v2 := router.Group("/v2") // API For Deck
	{
		v2.POST("/decks", createDeck(db))
		v2.GET("/decks", getListOfDeck(db))
		v2.GET("/decks/:deck_id", readDeckByDeckId(db))
		v2.PUT("/decks/:deck_id", editDeckByDeckId(db))
		v2.DELETE("/decks/:deck_id", deleteDeckByDeckId(db))
	}

	v3 := router.Group("/v3") // API For Card
	{
		v3.POST("/cards", createCard(db)) // Mỗi khi tạo 1 deck thì sẽ đồng thời tạo 52 thẻ bài cho deck đó với deck_id tương ứng
		v3.GET("/cards/:card_value", ReadCardByCardValue(db))
		v3.PUT("/cards/:card_value", editCardByCardValue(db))
	}
	router.Run()
}
func createCard(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dataCard Card

		if err := c.ShouldBind(&dataCard); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		dataCard.Status = false // set default - does not belong to anyone
		if err := db.Create(&dataCard).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"notify": "add deck successfully"})
	}
}

func createDeck(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dataDeck Deck

		if err := c.ShouldBind(&dataDeck); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		dataDeck.Remaining = 52 // set default
		if err := db.Create(&dataDeck).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"notify": "add deck successfully"})
	}
}
func createUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dataUser User

		if err := c.ShouldBind(&dataUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if dataUser.Username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username cannot be blank"})
			return
		}

		dataUser.Coins = 5000        // set default
		dataUser.Point_of_3cards = 0 // set default
		if err := db.Create(&dataUser).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": dataUser.Username})
	}
}
func getListOfUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var Users []User
		if err := db.Table(User{}.TableName()).Find(&Users).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": Users})
	}
}
func getListOfDeck(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var Decks []Deck
		if err := db.Table(Deck{}.TableName()).Find(&Decks).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": Decks})
	}
}

func readUserByUsername(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dataUser User

		userName := c.Param("username")
		if err := db.Where("username = ?", userName).First(&dataUser).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": dataUser})
	}
}
func readDeckByDeckId(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dataDeck Deck

		deckId, err := strconv.Atoi(c.Param("deck_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Where("deck_id = ?", deckId).First(&dataDeck).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": dataDeck})
	}
}
func ReadCardByCardValue(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dataCard Card

		cardValue, err := strconv.Atoi(c.Param("card_value"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Where("card_value = ?", cardValue).First(&dataCard).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": dataCard})
	}
}
func editUserByUsername(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userName := c.Param("username")

		var dataUser User
		if err := c.ShouldBind(&dataUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Where("username = ?", userName).Updates(&dataUser).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": true})
	}
}
func editDeckByDeckId(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		deckId, err := strconv.Atoi(c.Param("deck_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var dataDeck Deck
		if err := c.ShouldBind(&dataDeck); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Where("deck_id = ?", deckId).Updates(&dataDeck).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": true})
	}
}
func editCardByCardValue(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		cardValue, err := strconv.Atoi(c.Param("card_value"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var dataCard Card
		if err := c.ShouldBind(&dataCard); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Where("card_value = ?", cardValue).Updates(&dataCard).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": true})
	}
}
func deleteUserByUsername(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userName := c.Param("username")
		if err := db.Table(User{}.TableName()).
			Where("username = ?", userName).
			Delete(nil).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": true})
	}
}
func deleteDeckByDeckId(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		deckId, err := strconv.Atoi(c.Param("deck_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Table(Deck{}.TableName()).
			Where("deck_id = ?", deckId).
			Delete(nil).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": true})
	}
}
