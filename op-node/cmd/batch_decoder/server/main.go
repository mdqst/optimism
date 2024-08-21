package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/ethereum-optimism/optimism/op-node/cmd/batch_decoder/utils"

	"github.com/gorilla/mux"
)

// Span batch request is a request to find all span batches in a given block range.
type SpanBatchRequest struct {
	StartBlock  uint64 `json:"startBlock"`
	EndBlock    uint64 `json:"endBlock"`
	L2ChainID   uint64 `json:"l2ChainID"`
	L2Node      string `json:"l2Node"`
	L1RPC       string `json:"l1RPC"`
	L1Beacon    string `json:"l1Beacon"`
	BatchSender string `json:"batchSender"`
}

// Response to a span batch request.
type SpanBatchResponse struct {
	Ranges []utils.SpanBatchRange `json:"ranges"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/span-batch-ranges", handleSpanBatchRanges).Methods("POST")

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Return all of the span batches in a given L2 block range.
func handleSpanBatchRanges(w http.ResponseWriter, r *http.Request) {
	var req SpanBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config := utils.BatchDecoderConfig{
		L2ChainID:   new(big.Int).SetUint64(req.L2ChainID),
		L2Node:      req.L2Node,
		L1RPC:       req.L1RPC,
		L1Beacon:    req.L1Beacon,
		BatchSender: req.BatchSender,
		StartBlock:  req.StartBlock,
		EndBlock:    req.EndBlock,
		DataDir:     "/tmp/batch_decoder/transactions_cache_new", // You might want to make this configurable
	}

	ranges, err := utils.GetAllSpanBatchesInBlockRange(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := SpanBatchResponse{
		Ranges: ranges,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
