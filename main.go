package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	method := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origin := handlers.AllowedOrigins([]string{"*"})
	headers := handlers.AllowedHeaders([]string{"X-Request-With", "Content-Type", "Authorization"})
	r.HandleFunc("/submit", submitHandler).Methods("POST")
	r.HandleFunc("/upload", uploadHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(headers, method, origin)(r)))
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("recieved")
	p, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	var node []*Node
	err = json.Unmarshal(p, &node)
	if err != nil {
		fmt.Println(err)
	}

	for i := range node {
		fmt.Println(i, *node[i])
	}
	f, err := os.Create("gnerated.html")

	if err != nil {
		log.Fatal(err)
	}
	generateHTML(node, f)
	w.Write([]byte("yes"))
	fmt.Println(node)
}

//Atribute is some type
type Atribute struct {
	ID     string `json:"id"`
	Class  string `json:"class"`
	Link   string `json:"href"`
	src    string `json:"src"`
	Target string `json:"target"`
}

//Node is some other type
type Node struct {
	NodeElement string   `json:"nodeName"`
	NodeContent string   `json:"nodeContent"`
	NodeAtr     Atribute `json:"nodeAtr"`
	NodeChild   []*Node  `json:"nodeChild"`
}

// var Temple = template.New("my templatre").Parse(`
// 	<section>
// 	<{{.NodeElement}}>

// 	</>section
// `)

func generateHTML(node []*Node, file *os.File) {
	for _, element := range node {
		s := fmt.Sprintf("<%s ", element.NodeElement)
		_, err := file.WriteString(s)
		if err != nil {
			log.Fatal(err)
		}
		if (Atribute{}) != element.NodeAtr {
			s = fmt.Sprintf("%s %s %s %s %s>", element.NodeAtr.Class, element.NodeAtr.ID, element.NodeAtr.Link, element.NodeAtr.src, element.NodeAtr.Target)
			fmt.Println(s)
			_, err = file.WriteString(s)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			s = ">"
			_, err = file.WriteString(s)
			if err != nil {
				log.Fatal(err)
			}
		}
		s = fmt.Sprintf("%s", element.NodeContent)
		_, err = file.WriteString(s)
		if err != nil {
			log.Fatal(err)
		}
		for len(element.NodeChild) > 0 {
			generateHTML(element.NodeChild, file)
			element.NodeChild = element.NodeChild[1:len(element.NodeChild)]
		}
		s = fmt.Sprintf("</%s>", element.NodeElement)
		_, err = file.WriteString(s)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	imageName, err := fileUpload(r)
	var response Result
	if err != nil {
		// http.Error(w, "Invalid Data", http.StatusBadRequest)
		response.Res = ""
		response.Er = err.Error()
	}
	w.Header().Set("Content-Type", "application/json")
	response.Res = imageName
	fmt.Println(imageName, "image name is here", response)
	json.NewEncoder(w).Encode(response)
}

func fileUpload(r *http.Request) (string, error) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("image")
	fmt.Println(file, handler)
	if err != nil {
		return "", err
	}
	defer file.Close()
	f, err := os.OpenFile("static/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()
	io.Copy(f, file)
	fn := "static/" + handler.Filename
	return fn, nil
}

type Result struct {
	Res string `json:"result"`
	Er  string `json:"err"`
}
