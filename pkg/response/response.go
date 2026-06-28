package response

import (
	"encoding/json"
	"net/http"
)

// Response provides a structured way to write HTTP responses
type Response struct {
	http.ResponseWriter
}

// New creates a new Response writer
func New(w http.ResponseWriter) *Response {
	return &Response{ResponseWriter: w}
}

// JSON writes a JSON response
func (r *Response) JSON(status int, data interface{}) {
	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(status)
	json.NewEncoder(r.ResponseWriter).Encode(data)
}

// Error writes an error JSON response
func (r *Response) Error(status int, message string) {
	r.JSON(status, map[string]string{
		"error": message,
	})
}

// Success writes a success response with data and optional meta
func Success(w http.ResponseWriter, data interface{}, meta interface{}) {
	resp := map[string]interface{}{
		"success": true,
		"data":    data,
	}
	if meta != nil {
		resp["meta"] = meta
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Created writes a 201 Created response
func Created(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

// NoContent writes a 204 No Content response
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Message writes a message response
func Message(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": status < 400,
		"message": message,
	})
}
