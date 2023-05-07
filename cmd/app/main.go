package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/anupx73/go-bpcalc-backend-k8s/pkg/models/mongodb"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type application struct {
	errorLog   *log.Logger
	infoLog    *log.Logger
	bpReadings *mongodb.BPReadingModel
}

// TODO - return a new object for 'client' in heap
// func getMongoDBConnection(url string) (*mongo.Client, string, error) {
// 	// Create mongo client configuration
// 	opts := options.Client().ApplyURI(url)

// 	// Create a new client and connect to the server
// 	client, err := mongo.Connect(context.TODO(), opts)
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	defer func() {
// 		if err = client.Disconnect(context.TODO()); err != nil {
// 			panic(err)
// 		}
// 	}()
// 	// Send a ping to confirm a successful connection
// 	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
// 		return nil, "", err
// 	}

// 	return client, "Database deployment is reachable!!", nil
// }

func funnyParsingForVaultIssue853(rawFileData string) (string, string, string) {
	split1 := strings.Split(rawFileData, "[")
	split2 := strings.Split(split1[1], "]")
	split3 := strings.Split(split2[0], " ")

	pwd := strings.Split(split3[0], ":")
	url := strings.Split(split3[1], ":")
	user := strings.Split(split3[2], ":")

	return user[1], pwd[1], url[1]
}

func main() {
	// Read Config file
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.ReadInConfig()

	// Get Config data
	serverAddr := viper.GetString("dev.serverAddr")
	serverPort := viper.GetInt("dev.serverPort")
	mongoDatabase := viper.GetString("dev.mongoDatabase")
	mongoCollection := viper.GetString("dev.mongoCollection")

	// Create logger for writing information and error messages.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Get mongo connection string from Vault config
	rawBytes, err := os.ReadFile("/vault/secrets/database-config.txt")
	if err != nil {
		errLog.Panic(err)
	}
	// mongoUri := string(rawBytes[:])
	// workaround to bypass vault-helm issue 853
	user, pass, url := funnyParsingForVaultIssue853(string(rawBytes[:]))
	mongoUri := "mongodb+srv://" + user + ":" + pass + "@" + url + "/?retryWrites=true&w=majority"

	// Database connection
	opts := options.Client().ApplyURI(mongoUri)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		errLog.Panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			errLog.Panic(err)
		}
	}()
	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		errLog.Panic(err)
	}
	infoLog.Println("Database deployment is reachable!!")

	// Initialize a new instance of application containing the dependencies.
	app := &application{
		infoLog:  infoLog,
		errorLog: errLog,
		bpReadings: &mongodb.BPReadingModel{
			C: client.Database(mongoDatabase).Collection(mongoCollection),
		},
	}

	// Initialize a new http.Server struct.
	serverURI := fmt.Sprintf("%s:%d", serverAddr, serverPort)
	srv := &http.Server{
		Addr:         serverURI,
		ErrorLog:     errLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting backend server on %s", serverURI)
	err = srv.ListenAndServe()
	errLog.Fatal(err)
}
