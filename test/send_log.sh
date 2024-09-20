curl -X POST http://localhost:8080/api/logs \
     -H "Content-Type: application/json" \
     -d '{
           "level": "ERROR",
           "server_name": "server-1",
           "application": "auth-service",
           "environment": "production",
           "message": "Failed to connect to database",
           "error_code": "DB_CONN_TIMEOUT"
         }'