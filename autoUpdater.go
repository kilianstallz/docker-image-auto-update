package main

import (
	"fmt"
	"regexp"
	"bufio"
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

	log.Println("Stack Name:", r.URL.Query().Get("stack"))
	// get the stack name from query
	var stackName = r.URL.Query().Get("stack")

	// define command
	var dockerCommand = "docker stack ps " + stackName

	// run command
	cmd := exec.Command("/bin/sh", "-c", dockerCommand)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Start()

	buf := bufio.NewReader(stdout) // Notice that this is not in a loop
	num := 1
	serviceSet := make(map[string]bool)

	// filter all services and parse them in a map
	for {
		line, _, _ := buf.ReadLine()
		if num > 3 {
			break
		}
		num++
		var regex = fmt.Sprintf(`(?m)%s_([^.]+)`, stackName)
		var name = regexp.MustCompile(regex)
		result := name.FindStringSubmatch(string(line))

		// get distinct services
		for i, serviceName := range result {
			if (i == 1) {
				_, exists := serviceSet[serviceName]
				if (!exists) {
					serviceSet[serviceName] = true
				}
			}
		}
	}

	for service := range serviceSet {
		log.Println(service)
	}

	io.WriteString(w, `{"status": "ok"}`)
	return

	// cmd := exec.Command("/bin/sh", "./scripts/script.sh")
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// if err := cmd.Run(); err != nil {
	// 	log.Fatalln(err)
	// 	io.WriteString(w, `{"status": "error"}`)	
	// } else {
	// 	log.Println("Update Executed")
	// 	io.WriteString(w, `{"status": "ok"}`)
	// }
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