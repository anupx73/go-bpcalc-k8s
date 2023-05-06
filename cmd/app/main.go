package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/anupx73/go-bpcalc-backend-k8s/pkg/models/mongodb"
	vault "github.com/hashicorp/vault/api"
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

func getVaultSecrets() (string, string, string, error) {
	config := vault.DefaultConfig()
	config.Address = os.Getenv("VAULT_ADDR")

	client, err := vault.NewClient(config)
	if err != nil {
		return "", "", "", errors.New("Vault init failed; err: " + fmt.Sprintf("%v", err))
	}

	client.SetToken(os.Getenv("VAULT_TOKEN"))
	secret, err := client.KVv2("tudublin").Get(context.Background(), "mongo-atlas")
	if err != nil {
		return "", "", "", errors.New("Unable to read Vault secret; err: " + fmt.Sprintf("%v", err))
	}

	dbUri, ok := secret.Data["url"].(string)
	if !ok {
		return "", "", "", errors.New("db url type assertion failed")
	}
	dbUsername, ok := secret.Data["username"].(string)
	if !ok {
		return "", "", "", errors.New("username type assertion failed")
	}
	dbPassword, ok := secret.Data["password"].(string)
	if !ok {
		return "", "", "", errors.New("password type assertion failed")
	}

	return dbUri, dbUsername, dbPassword, nil
}

func getMongoDBConnection(url string) (*mongo.Client, string, error) {
	// Create mongo client configuration
	opts := options.Client().ApplyURI(url)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, "", err
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		return nil, "", err
	}

	return client, "Database deployment is reachable!!", nil
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

	// Get Vault secrets for mongo
	uri, user, pass, err := getVaultSecrets()
	if err != nil {
		errLog.Panic(err)
	}
	fullMongoURI := "mongodb+srv://" + user + ":" + pass + "@" + uri

	// Database connection
	client, status, err := getMongoDBConnection(fullMongoURI)
	if err != nil {
		errLog.Panic(err)
	}
	infoLog.Printf(status)

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
