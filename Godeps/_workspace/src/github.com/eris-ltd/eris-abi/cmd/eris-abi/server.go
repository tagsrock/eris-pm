package main

import (
	"encoding/json"
	"fmt"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/eris-abi"
	"net/http"
)

//------------------------------------------------------------------------
// http server exports same commands as the cli
// all request arguments are keyed and passed through header
// body is ignored

func ListenAndServe(host, port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/pack", packHandler)
	return http.ListenAndServe(host+":"+port, mux)
}

// dead simple response struct
type HTTPResponse struct {
	Response string
	Error    string
}

func WriteResult(w http.ResponseWriter, result string) {
	resp := HTTPResponse{result, ""}
	b, _ := json.Marshal(resp)
	w.Write(b)
}

func WriteError(w http.ResponseWriter, err error) {
	resp := HTTPResponse{"", err.Error()}
	b, _ := json.Marshal(resp)
	w.Write(b)
}

//------------------------------------------------------------------------
// handlers

func packHandler(w http.ResponseWriter, r *http.Request) {
	input := r.Header.Get("input") //json, hash, or index. "FILE" NOT VALID

	//tx argument parsing
	argStr := r.Header.Get("args")
	if argStr == "" {
		argStr = `[]`
	}

	var args []string
	err := json.Unmarshal([]byte(argStr), &args)
	if err != nil {
		WriteError(w, err)
		return
	}

	//input method switch
	if input == "json" {
		jsonabi := []byte(r.Header.Get("json"))

		tx, err := ebi.Packer(jsonabi, args...)
		if err != nil {
			WriteError(w, err)
			return
		}

		WriteResult(w, fmt.Sprintf("%s", tx))

	} else if input == "hash" {
		hash := r.Header.Get("hash")

		tx, err := ebi.HashPack(hash, args...)
		if err != nil {
			WriteError(w, err)
			return
		}

		WriteResult(w, fmt.Sprintf("%s", tx))
		return

	} else if input == "index" {
		index := r.Header.Get("index")
		if index == "" {
			index = DefaultIndex
		}

		key := r.Header.Get("key")
		if key == "" {
			WriteError(w, fmt.Errorf("A key for the index MUST be specified"))
		}

		tx, err := ebi.IndexPack(index, key, args...)
		if err != nil {
			WriteError(w, err)
			return
		}

		WriteResult(w, fmt.Sprintf("%s", tx))
		return

	} else {
		WriteError(w, fmt.Errorf("Unrecoginized abi specification method"))
	}
}
