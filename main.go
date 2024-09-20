package main

import (
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "github.com/lincolnaleixo/smith/db"
    "github.com/lincolnaleixo/smith/handlers"
    "github.com/rs/cors"
)

func main() {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found. Proceeding with environment variables.")
    }

    // Get MongoDB URI from environment variables
    mongoURI := os.Getenv("MONGODB_URI")
    if mongoURI == "" {
        log.Fatal("MONGODB_URI not set in environment variables")
    }

    // Initialize MongoDB
    db.InitMongoDB(mongoURI)

    // Get server port from environment variables, default to 8080
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // Initialize Router
    router := mux.NewRouter()

    // API Routes
    router.HandleFunc("/api/logs", handlers.CreateLog).Methods("POST")
    router.HandleFunc("/api/logs", handlers.GetLogs).Methods("GET")
    router.HandleFunc("/api/servers", handlers.GetServers).Methods("GET")

    // Serve Static Files
    staticFileDirectory := http.Dir("./static/")
    staticFileHandler := http.StripPrefix("/", http.FileServer(staticFileDirectory))
    router.PathPrefix("/").Handler(staticFileHandler).Methods("GET")

    // Handle CORS if necessary
    handler := cors.Default().Handler(router)

    // Start Server
    log.Printf("Smith Logger is running on port %s", port)
    if err := http.ListenAndServe(":"+port, handler); err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}