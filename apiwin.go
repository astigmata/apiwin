package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"

	"gopkg.in/yaml.v3"
)

// Config Struct for the yaml file
type Config struct {
	APIPort string `yaml:"apiPort"`
}

type Results struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

var p = fmt.Printf
var l = log.Fatalf

var c Config

func listDirectory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) // to get from id in param
	cmd := params["folder"]
	folderToScan := fmt.Sprint(strings.ReplaceAll(cmd, "|", "\\"))

	file, err := os.Open(folderToScan) // For read access.
	if err != nil {
		//log.Fatal(err)
		p("pas glop")
		return
	}
	defer file.Close()

	var results []Results

	err = filepath.Walk(folderToScan, func(path string, info os.FileInfo, err error) error {
		a := Results{Filename: info.Name(), Size: info.Size()}
		results = append(results, a)

		return nil
	})
	if err != nil {
		panic(err)
	}

	p("GET /directory/%s\n", folderToScan)
	for _, result := range results {
		fmt.Println(result)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)

}

func main() {
	loadconf()

	r := mux.NewRouter()

	r.HandleFunc("/v1/list/{folder}", listDirectory).Methods("GET")

	p("listening port %s...", c.APIPort)
	port := fmt.Sprintf(":%s", c.APIPort)
	http.ListenAndServe(port, r)

}

func loadconf() {
	// load conf
	filename, _ := filepath.Abs("./apiwin.yml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		panic(err)
	}
}
