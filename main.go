package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	method := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origin := handlers.AllowedOrigins([]string{"*"})
	headers := handlers.AllowedHeaders([]string{"X-Request-With", "Content-Type", "Authorization"})
	r.HandleFunc("/submit", submitHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(headers, method, origin)(r)))
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	p, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	var node []Node
	err = json.Unmarshal(p, &node)
	if err != nil {
		fmt.Println(err)
	}
	s := generateHTML(node[1])
	w.Write([]byte("yes"))
	fmt.Println(node)
	fmt.Println(s)
}

type Node struct {
	NodeElement string   `json:"nodeName"`
	NodeContent string   `json:"nodeContent"`
	NodeAtr     []string `json:"nodeAtr"`
	NodeChild   []*Node  `json:"nodeChild"`
}

// var Temple = template.New("my templatre").Parse(`
// 	<section>
// 	<{{.NodeElement}}>

// 	</>section
// `)

func generateHTML(node Node) string {
	var a string
	if node.NodeElement == "p" {
		a = fmt.Sprintf("<%s %v>%s</%s>", node.NodeElement, node.NodeAtr, node.NodeContent, node.NodeElement)
	}
	return a
}
