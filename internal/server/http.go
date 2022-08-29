package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type httpServer struct {
	Log *Log
}

func newHTTPServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

// Handler function to handle when Logs are produced by clients
func (s *httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	var request ProduceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offset, err := s.Log.Append(request.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := ProduceResponse{Offset: offset}
	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Handler function to handle when Logs are consumed by clients
func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	var req ConsumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	record, err := s.Log.Read(req.Offset)
	if errors.Is(err, ErrOffsetNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return 
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

	resp := ConsumeResponse{Record: *record}
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HTTP request struct for the 'Append' to 'Log'
type ProduceRequest struct {
	Record Record `json:"record"`
}

//HTTP response struct for the 'Append' to 'Log'
type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

// HTTP request struct for the 'Read' from 'Log'
type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

//HTTP response struct for the 'Read' from 'Log'
type ConsumeResponse struct {
	Record Record `json:"record"`
}

func NewHTTPServer(addr string) *http.Server {
	httpsrv := newHTTPServer()
	r := mux.NewRouter()

	r.HandleFunc("/", httpsrv.handleProduce).Methods(http.MethodPost)
	r.HandleFunc("/", httpsrv.handleConsume).Methods(http.MethodGet)

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}
