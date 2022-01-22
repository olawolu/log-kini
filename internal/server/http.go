package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type httpServer struct {
	Log *Log
}

type CommitRequest struct {
	Record Record `json:"record"`
}

type CommitResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record Record `json:"record"`
}

func NewHTTPServer(addr string) *http.Server {
	httpsrv := newHTTPServer()
	router := mux.NewRouter()

	router.HandleFunc("/", httpsrv.handleCommit).Methods("POST")
	router.HandleFunc("/", httpsrv.handleConsume).Methods("Get")
	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}

func newHTTPServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

func (s *httpServer) handleCommit(w http.ResponseWriter, r *http.Request) {
	var req CommitRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleHTTPErr(w, http.StatusBadRequest, err)
		return
	}

	off, err := s.Log.Append(req.Record)
	if err != nil {
		handleHTTPErr(w, http.StatusInternalServerError, err)
		return
	}

	res := CommitResponse{
		Offset: off,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		handleHTTPErr(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	var req ConsumeRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleHTTPErr(w, http.StatusBadRequest, err)
		return
	}

	record, err := s.Log.Read(req.Offset)
	if err == ErrOffsetNotFound {
		handleHTTPErr(w, http.StatusNotFound, err)
		return
	}
	if err != nil {
		handleHTTPErr(w, http.StatusInternalServerError, err)
		return
	}

	res := ConsumeResponse{
		Record: record,
	}
	
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		handleHTTPErr(w, http.StatusInternalServerError, err)
		return
	}
}

func handleHTTPErr(w http.ResponseWriter, status int, err error) {
	log.Println(err)
	http.Error(w, err.Error(), status)
}
