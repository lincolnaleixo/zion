package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "time"

    "github.com/lincolnaleixo/smith/db"
    "go.mongodb.org/mongo-driver/bson"
)

// GetServers handles GET /api/servers
func GetServers(w http.ResponseWriter, r *http.Request) {
    var servers []string

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Aggregate distinct server names
    distinctServers, err := db.LogCollection.Distinct(ctx, "server_name", bson.D{})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    for _, server := range distinctServers {
        if serverStr, ok := server.(string); ok {
            servers = append(servers, serverStr)
        }
    }

    json.NewEncoder(w).Encode(servers)
}