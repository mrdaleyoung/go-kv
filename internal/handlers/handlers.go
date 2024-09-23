package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-kv/internal/services"
	"io"
	"net/http"
)

func HandleGet(kvService services.KVServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Param("key")
		value, err := kvService.Get(key)
		if err != nil || value == nil {
			//Return 404 if key is not found
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		//Return 202
		c.JSON(http.StatusOK, value)
	}
}

// HandlePut accepts raw data, validates it, converts to a string, and stores in the KV store
func HandlePut(kvService services.KVServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Param("key")

		// Read the raw body
		rawData, err := io.ReadAll(c.Request.Body)
		if err != nil || len(rawData) == 0 {
			//Return 400 if we cannot read the body or it's empty
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Check if the body is valid JSON, otherwise treat it as a string
		var bodyContent string
		if json.Valid(rawData) {
			var jsonBody interface{}
			if err := json.Unmarshal(rawData, &jsonBody); err == nil {
				// Convert JSON content to string
				bodyContent = fmt.Sprintf("%v", jsonBody)
			} else {
				//Return 400 if we cannot read the body
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
				return
			}
		} else {
			bodyContent = string(rawData)
		}

		// Store the value in the KV store
		if err := kvService.Put(key, bodyContent); err != nil {
			//Return 500 if server cannot store the key
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to store value"})
			return
		}

		// Respond with an empty JSON object
		c.JSON(http.StatusAccepted, gin.H{})
	}
}

func HandleDelete(kvService services.KVServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Param("key")

		// Check if the key exists by inspecting the value instead of relying on error
		value, _ := kvService.Get(key)
		if value == nil {
			// If the key doesn't exist, return 404
			c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
			return
		}

		// Proceed to delete the key
		err := kvService.Delete(key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete key"})
			return
		}

		// Return 200 on successful deletion
		c.Status(http.StatusOK)
	}
}

func HandleListKeys(kvService services.KVServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		keys := kvService.ListKeys()
		c.JSON(http.StatusOK, keys)
	}
}
