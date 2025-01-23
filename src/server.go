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
	"database/sql"
	_ "github.com/lib/pq"
	"scmlike/src/compiler"
)

type ScmLike struct {
	Exp string `json:"exp"`
}

type Result struct {
	Exp string `json:"exp"`
}

func Compile(exp string) (string, error) {
	ast, err := compiler.Parse(exp)
	if err != nil {
		return "", err
	}

	monAst := compiler.ToAnf(ast, 0)
	ss := compiler.SelectInstructions(monAst)
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

	connStr := "user=compiler dbname=compiler password=hello123 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	compilerTable := `
        CREATE TABLE IF NOT EXISTS compiler (
           id SERIAL PRIMARY KEY,
           ir TEXT NOT NULL,
           assembly TEXT NOT NULL,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
          );`

	submissionTable := `
        CREATE TABLE IF NOT EXISTS submission (
        id SERIAL PRIMARY KEY,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        compiler TEXT NOT NULL,
        result TEXT NOT NULL,
        cpu_time FLOAT NOT NULL, -- CPU time in seconds
        memory_usage BIGINT NOT NULL -- Memory usage in bytes
       );`

	userTable := `
        CREATE TABLE IF NOT EXISTS "user" (
         id SERIAL PRIMARY KEY,
         username VARCHAR(50) UNIQUE NOT NULL,
         email VARCHAR(100) UNIQUE NOT NULL,
         password VARCHAR(255) NOT NULL,
         name VARCHAR(100),
         bio TEXT
         );`
	postTable := `
        CREATE TABLE IF NOT EXISTS "post" (
         postid SERIAL PRIMARY KEY,
         userid INT NOT NULL,
         content TEXT NOT NULL,
         mediatype VARCHAR(20),
         mediaurl VARCHAR(255),
         timestamp TIMESTAMPTZ DEFAULT NOW(),
         FOREIGN KEY (userid) REFERENCES "user"(id) ON DELETE CASCADE
         );`

	commentTable := `
        CREATE TABLE IF NOT EXISTS "comment" (
         commentid SERIAL PRIMARY KEY,
         postid INT NOT NULL,
         userid INT NOT NULL,
         content TEXT NOT NULL,
         timestamp TIMESTAMPTZ DEFAULT NOW(),
         FOREIGN KEY (postid) REFERENCES "post"(postid) ON DELETE CASCADE,
         FOREIGN KEY (userid) REFERENCES "user"(id) ON DELETE CASCADE
         );`

	likeTable := `
        CREATE TABLE IF NOT EXISTS "likes" (
          likeid SERIAL PRIMARY KEY,
          postid INT,
          commentid INT,
          userid INT NOT NULL,
          timestamp TIMESTAMPTZ DEFAULT NOW(),
          FOREIGN KEY (postid) REFERENCES "post"(postid) ON DELETE CASCADE,
          FOREIGN KEY (commentid) REFERENCES "comment"(commentid) ON DELETE CASCADE,
          FOREIGN KEY (userid) REFERENCES "user"(id) ON DELETE CASCADE
         );`

	friendTable := `
        CREATE TABLE IF NOT EXISTS "friendship" (
          friendshipid SERIAL PRIMARY KEY,
          userid1 INT NOT NULL,
          userid2 INT NOT NULL,
          timestamp TIMESTAMPTZ DEFAULT NOW(),
          FOREIGN KEY (userid1) REFERENCES "user"(id) ON DELETE CASCADE,
          FOREIGN KEY (userid2) REFERENCES "user"(id) ON DELETE CASCADE,
          UNIQUE (userid1, userid2),
          CHECK (userid1 != userid2)
         );`

	messageTable := `
        CREATE TABLE IF NOT EXISTS "message" (
          messageid SERIAL PRIMARY KEY,
          senderid INT NOT NULL,
          receiverid INT NOT NULL,
          content TEXT NOT NULL,
          timestamp TIMESTAMPTZ DEFAULT NOW(),
          FOREIGN KEY (senderid) REFERENCES "user"(id) ON DELETE CASCADE,
          FOREIGN KEY (receiverid) REFERENCES "user"(id) ON DELETE CASCADE
         );`
         
         
	
	_, err = db.Exec(compilerTable)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	fmt.Println("Table created successfully!")

	_, err = db.Exec(submissionTable)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	fmt.Println("Table created sucessfully!")

	_, err = db.Exec(userTable)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	fmt.Println("Table created successfully!")

	_, err = db.Exec(postTable)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	fmt.Println("Table created sucessfully!")
	
	_, err = db.Exec(commentTable)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	fmt.Println("Table created successfully!")

	_, err = db.Exec(likeTable)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	fmt.Println("Table created sucessfully!")

	_, err = db.Exec(friendTable)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	fmt.Println("Table created successfully!")

	_, err = db.Exec(messageTable)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	fmt.Println("Table created sucessfully!")
	
	
	mux.Handle("/api/compiler", http.HandlerFunc(CompileHandler))
	mux.HandleFunc("/", defaultHandler)
	log.Println("listening on port:", PORT)
	log.Fatal(s.ListenAndServe())
}
