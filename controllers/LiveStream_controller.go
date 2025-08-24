package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"My_WEB/config"
	"My_WEB/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var liveStreamCollection *mongo.Collection
var websocketServiceURL = "http://localhost:8081" // WebSocket service URL

func init() {
	if config.DB != nil {
		liveStreamCollection = config.DB.Collection("livestreams")
	}
}

// Call WebSocket service to get viewers
func GetStreamViewers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	streamID := params["stream_id"]

	// Call WebSocket service
	resp, err := http.Get(fmt.Sprintf("%s/api/livestream/%s/viewers", websocketServiceURL, streamID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get viewers"})
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

// Call WebSocket service to send announcement
func SendStreamAnnouncement(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	streamID := params["stream_id"]

	var announcement map[string]interface{}
	json.NewDecoder(r.Body).Decode(&announcement)

	jsonData, _ := json.Marshal(announcement)
	resp, err := http.Post(
		fmt.Sprintf("%s/api/livestream/%s/announcement", websocketServiceURL, streamID),
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to send announcement"})
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

// ...existing CRUD functions...
func CreateLiveStream(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	var liveStream models.LiveStream
	if err := json.NewDecoder(r.Body).Decode(&liveStream); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	liveStream.ID = primitive.NewObjectID()
	liveStream.CreatedAt = time.Now()
	liveStream.UpdatedAt = time.Now()
	liveStream.Status = "inactive"

	result, err := liveStreamCollection.InsertOne(ctx, liveStream)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create live stream"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": result.InsertedID, "message": "Live stream created successfully"})
}

// ...other existing CRUD functions remain the same...
