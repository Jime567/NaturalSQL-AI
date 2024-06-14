package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jime567/NaturalSQL-AI/structs"
	openai "github.com/sashabaranov/go-openai"
)

var (
	db         *sql.DB
	client1    *openai.Client
	client2    *openai.Client
	background string
)

func main() {
	// Initialize database and OpenAI client
	initDB()
	initOpenAIClient()
	content, err := ioutil.ReadFile("background.txt")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	background = string(content)

	// Initialize Gin router
	router := gin.Default()

	// Enable CORS
	router.Use(cors.Default())

	// Routes
	router.GET("/getCyclists", getCyclistsHandler)
	router.GET("/getBikes", getBikesHandler)
	router.GET("/getAddresses", getAddressesHandler)
	router.POST("/ask", askHandler)

	// Start server
	port := ":8080"
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(router.Run(port))
}

func initDB() {
	// Credentials
	dbUser := os.Getenv("AI_SQL_USER")
	if dbUser == "" {
		log.Fatal("SQL_USER environment variable not set")
	}
	dbPass := os.Getenv("AI_SQL_PASS")
	if dbPass == "" {
		log.Fatal("SQL_PASS environment variable not set")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/cyclists_db", dbUser, dbPass)

	// Open a database connection
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Unable to open connection: ", err)
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		log.Fatal("Connection Failed: ", err)
	}

	fmt.Println("Successfully connected to the database!")
}

func initOpenAIClient() {
	// Credentials
	AIKEY := os.Getenv("OPENAI_KEY")
	if AIKEY == "" {
		log.Fatal("OPENAI_KEY environment variable not set")
	}

	// Initialize OpenAI client
	client1 = openai.NewClient(AIKEY)
	client2 = openai.NewClient(AIKEY)
}

func getCyclistsHandler(c *gin.Context) {
	cyclists := getCyclists()
	c.JSON(http.StatusOK, cyclists)
}

func getBikesHandler(c *gin.Context) {
	bikes := getBikes()
	c.JSON(http.StatusOK, bikes)
}

func getAddressesHandler(c *gin.Context) {
	addresses := getAddresses()
	c.JSON(http.StatusOK, addresses)
}

func askHandler(c *gin.Context) {
	var requestBody struct {
		Text string `json:"text"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	question := requestBody.Text
	query := askChatbotSQL(client1, question, background)
	sqlReturn, _ := executeSQL(query)
	humanReturn := askChatbot(client2, question, sqlReturn)
	response := humanReturn + " \n\n SQL:\n" + query + "\n\n " + sqlReturn

	c.JSON(http.StatusOK, gin.H{"response": response})
}

func executeSQL(query string) (string, error) {
	var result string

	rows, err := db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Iterate over the rows
	for rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			return "", err
		}

		// print the values from 'values' slice
		var row string
		for i, col := range values {
			if col == nil {
				row += fmt.Sprintf("%s: NULL, ", columns[i])
			} else {
				row += fmt.Sprintf("%s: %s, ", columns[i], col)
			}
		}
		// Remove trailing comma and space
		row = row[:len(row)-2] + "\n"
		result += row
	}

	if err := rows.Err(); err != nil {
		return "", err
	}

	return result, nil
}

func getCyclists() []structs.Cyclist {
	var cyclists []structs.Cyclist

	rows, err := db.Query("SELECT * FROM cyclists")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		var cyclist structs.Cyclist
		if err := rows.Scan(&cyclist.ID, &cyclist.Name, &cyclist.PhoneNumber, &cyclist.SkillLevel, &cyclist.AddressID, &cyclist.BikeID); err != nil {
			log.Fatal(err)
		}
		cyclists = append(cyclists, cyclist)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return cyclists
}

func getBikes() []structs.Bike {
	var bikes []structs.Bike

	rows, err := db.Query("SELECT * FROM bikes")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		var bike structs.Bike
		if err := rows.Scan(&bike.ID, &bike.Nickname, &bike.SerialNumber, &bike.Year, &bike.Model, &bike.Make, &bike.Mileage); err != nil {
			log.Fatal(err)
		}
		bikes = append(bikes, bike)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return bikes
}

func getAddresses() []structs.Address {
	var addresses []structs.Address

	rows, err := db.Query("SELECT * FROM addresses")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		var address structs.Address
		if err := rows.Scan(&address.ID, &address.Street, &address.Zip, &address.State); err != nil {
			log.Fatal(err)
		}
		addresses = append(addresses, address)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return addresses
}

func askChatbotSQL(client *openai.Client, question string, background string) string {
	question = background + question
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: question,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return ""
	}

	return resp.Choices[0].Message.Content
}

func askChatbot(client *openai.Client, question string, response string) string {
	question = "The following question or request was made in regards to an sql database: " + question + " The response to this question or request is: " + response + " Give me the human readable version of this response without telling me that you are."
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: question,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return ""
	}

	return resp.Choices[0].Message.Content
}
