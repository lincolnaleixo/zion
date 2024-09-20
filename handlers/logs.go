package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "strconv"
    "time"

    "github.com/lincolnaleixo/smith/db"
    "github.com/lincolnaleixo/smith/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type GetLogsResponse struct {
    Logs []models.Log `json:"logs"`
}

// CreateLog handles POST /api/logs
func CreateLog(w http.ResponseWriter, r *http.Request) {
    var logEntry models.Log
    if err := json.NewDecoder(r.Body).Decode(&logEntry); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Set timestamp if not provided
    if logEntry.Timestamp.IsZero() {
        logEntry.Timestamp = time.Now().UTC()
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := db.LogCollection.InsertOne(ctx, logEntry)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(logEntry)
}

// GetLogs handles GET /api/logs with optional filters
func GetLogs(w http.ResponseWriter, r *http.Request) {
    var logs []models.Log
    filter := bson.M{}

    // Filtering
    if level := r.URL.Query().Get("level"); level != "" {
        filter["level"] = level
    }
    if server := r.URL.Query().Get("server"); server != "" {
        filter["server_name"] = server
    }
    if application := r.URL.Query().Get("application"); application != "" {
        filter["application"] = application
    }
    if environment := r.URL.Query().Get("environment"); environment != "" {
        filter["environment"] = environment
    }

    // Sorting
    sortParam := bson.D{}
    sortOrder := r.URL.Query().Get("sort")
    if sortOrder == "asc" {
        sortParam = bson.D{{Key: "timestamp", Value: 1}}
    } else {
        sortParam = bson.D{{Key: "timestamp", Value: -1}}
    }

    findOptions := options.Find()
    findOptions.SetSort(sortParam)

    // Pagination (optional)
    limitStr := r.URL.Query().Get("limit")
    offsetStr := r.URL.Query().Get("offset")
    if limitStr != "" {
        limit, err := strconv.Atoi(limitStr)
        if err == nil {
            findOptions.SetLimit(int64(limit))
        }
    }
    if offsetStr != "" {
        offset, err := strconv.Atoi(offsetStr)
        if err == nil {
            findOptions.SetSkip(int64(offset))
        }
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    cursor, err := db.LogCollection.Find(ctx, filter, findOptions)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var logEntry models.Log
        if err := cursor.Decode(&logEntry); err != nil {
            continue // Skip invalid entries
        }
        logs = append(logs, logEntry)
    }

    if err := cursor.Err(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(logs)
}