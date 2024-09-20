package models

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type LogLevel string

const (
    DEBUG LogLevel = "DEBUG"
    INFO  LogLevel = "INFO"
    WARN  LogLevel = "WARN"
    ERROR LogLevel = "ERROR"
    FATAL LogLevel = "FATAL"
)

type Log struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Level       LogLevel           `bson:"level" json:"level"`
    ServerName  string             `bson:"server_name" json:"server_name"`
    Application string             `bson:"application" json:"application"`
    Environment string             `bson:"environment" json:"environment"`
    Message     string             `bson:"message" json:"message"`
    ErrorCode   string             `bson:"error_code,omitempty" json:"error_code,omitempty"`
    Timestamp   time.Time          `bson:"timestamp" json:"timestamp"`
}