package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)


//Document struct
type Document struct {
	ID   string
	Name string
	Size int64
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/documents", getDocuments).Methods("GET")
	log.Fatal(http.ListenAndServe(":9000", router))
}

func getDocuments(w http.ResponseWriter, r *http.Request) {
	var docs []Document
	fileDir := "./FileToTest/"
	fileInfos, err := ioutil.ReadDir(fileDir)
	if err != nil {
		fmt.Println("Error in accessing directory:", err)
	}

	for _, file := range fileInfos {
		fileHash, err := hashFileToMD5CheckSum(fileDir + "/" + file.Name())
		if err == nil {
			docs = append(docs, Document{ID: fileHash, Name: file.Name(), Size: file.Size()})
			fmt.Printf("MD5: %s  -  Name: %s  -  Size: %d Kbps \n", fileHash, file.Name(), file.Size())
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}


func hashFileToMD5CheckSum(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil

}
