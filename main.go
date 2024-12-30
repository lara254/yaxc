package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type ScmLike struct {
	Exp string `json:"exp"`
}

type Result struct {
	Exp string `json:"exp"`
}

func Compile(exp string) (string, error) {
	ast, err := Parse(exp)
	if err != nil {
		return "", err
	}

	monAst := ToAnf(ast)
	ss := SelectInstructions(monAst)
	return SelectInsToString(ss.Instructs), nil
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusNotFound)
	body := "Thanks for visiting!\n"
	fmt.Fprintf(w, "%s", body)
}

func SelectInsToString(arr [][]string) string {
	rows := make([]string, len(arr))
	for i, row := range arr {
		rows[i] = "[" + strings.Join(row, " ") + "]"
	}
	return "[" + strings.Join(rows, " ") + "]"
}

func CompileHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host, r.Method)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var exp ScmLike
	err = json.Unmarshal(d, &exp)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	selectIns, err := Compile(exp.Exp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Result{Exp: selectIns})
}

func main() {
	var PORT = ":1234"
	if len(os.Args) > 1 {
		PORT = ":" + os.Args[1]
	}
	mux := http.NewServeMux()
	s := &http.Server{
		Addr:         PORT,
		Handler:      mux,
		IdleTimeout:  10 * time.Second,
		WriteTimeout: time.Second,
	}

	mux.Handle("/api/compiler", http.HandlerFunc(CompileHandler))
	mux.HandleFunc("/", defaultHandler)
	log.Println("listening on port:", PORT)
	log.Fatal(s.ListenAndServe())
}
