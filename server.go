package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ducdang91/go-backend/graph"
	"github.com/ducdang91/go-backend/graph/model"
)

const defaultPort = "8090"
var db *gorm.DB;

const DB_USERNAME = "root"
const DB_PASSWORD = "qweqwe"
const DB_NAME = "test_db"
const DB_HOST = "localhost"
const DB_PORT = "3307"

func initDB() {
	var err error
	//dsn := DB_USERNAME + ":" + DB_PASSWORD + "@tcp" + "(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?" + "parseTime=true&loc=Local"
	dsn := DB_USERNAME + ":" + DB_PASSWORD + "@tcp" + "(" + DB_HOST + ":" + DB_PORT + ")/?" + "parseTime=true&loc=Local"
	//dataSourceName := "root:@tcp(localhost:3307)/?parseTime=True"
	db, err = gorm.Open("mysql", dsn)

    if err != nil {
        fmt.Println(err)
        panic("failed to connect database")
    }

    db.LogMode(true)

    // Create the database. This is a one-time step.
    // Comment out if running multiple times - You may see an error otherwise
    db.Exec("CREATE DATABASE test_db")
    db.Exec("USE test_db")

    // Migration to create tables for Order and Item schema
    db.AutoMigrate(&model.Order{}, &model.Item{})	
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	initDB()
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		DB: db,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
