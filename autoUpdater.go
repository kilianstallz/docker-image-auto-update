package main

import (
	"net/http"
	"os/exec"
	"os"
	"io"
	"log"
)

const port = ":8090"

func main() {
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/rollback", rollbackHandler)
	http.HandleFunc("/prune", pruneHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	log.Println("** Server Started on Port " + port + " **")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}	
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	cmd := exec.Command("/bin/sh", "./scripts/script.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
		io.WriteString(w, `{"status": "error"}`)	
	} else {
		log.Println("Update Executed")
		io.WriteString(w, `{"status": "ok"}`)
	}
}

func rollbackHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	cmd := exec.Command("/bin/sh", "./scripts/rollback.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
		io.WriteString(w, `{"status": "error"}`)	
	} else {
		log.Println("Update Executed")
		io.WriteString(w, `{"status": "ok"}`)
	}
}

func pruneHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	cmd := exec.Command("/bin/sh", "./scripts/prune.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
		io.WriteString(w, `{"status": "error"}`)	
	} else {
		log.Println("Update Executed")
		io.WriteString(w, `{"status": "ok"}`)
	}
}