package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kyle8615/ethereum-parser/v1/internal/parser"
)

type Server struct {
	httpServer *http.Server
	Parser     parser.Parser
}

func NewServer(p parser.Parser) *Server {
	server := &http.Server{Addr: ":8080"}
	return &Server{httpServer: server, Parser: p}
}

func (s *Server) Start() error {
	http.HandleFunc("/subscribe", s.handleSubscribe)
	http.HandleFunc("/transactions", s.handleGetTransactions)
	http.HandleFunc("/currentblock", s.handleGetCurrentBlock)

	return s.httpServer.ListenAndServe()
}

func (s *Server) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var requestData struct {
		Address string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	success := s.Parser.Subscribe(requestData.Address)
	if !success {
		sendErrorResponse(w, "Failed to subscribe", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, nil, http.StatusOK)
}

func (s *Server) handleGetTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendErrorResponse(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	address := r.URL.Query().Get("address")
	if address == "" {
		sendErrorResponse(w, "Address parameter is missing", http.StatusBadRequest)
		return
	}

	transactions := s.Parser.GetTransactions(address)
	if len(transactions) == 0 {
		sendJSONResponse(w, nil, http.StatusOK)
	} else {
		sendJSONResponse(w, transactions, http.StatusOK)
	}
}

func (s *Server) handleGetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendErrorResponse(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	currentBlock := s.Parser.GetCurrentBlock()
	sendJSONResponse(w, currentBlock, http.StatusOK)
}

func sendJSONResponse(w http.ResponseWriter, result interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		Success: statusCode >= 200 && statusCode < 300,
		Result:  result,
	}

	json.NewEncoder(w).Encode(response)
}

func sendErrorResponse(w http.ResponseWriter, errorMessage string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		Success: false,
		Error:   errorMessage,
	}

	json.NewEncoder(w).Encode(response)
}

func (s *Server) Shutdown(ctx context.Context) error {
	<-ctx.Done()
	return s.httpServer.Shutdown(ctx)
}
