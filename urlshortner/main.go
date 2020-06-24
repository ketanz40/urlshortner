package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	hashids "github.com/speps/go-hashids"
	gocb "gopkg.in/couchbase/gocb.v1"
)

//User defined data type for the user's desired short url with JSON mappings to the long url and the hashed id
type userUrl struct {
	ID       string `json:"id,omitempty"`
	LongUrl  string `json:"longUrl,omitempty"`
	ShortUrl string `json:"shortUrl,omitempty"`
}

func initialEndpoint(w http.ResponseWriter, req *http.Request) {
	var url userUrl
	_ = json.NewDecoder(req.Body).Decode(&url)
	var n1qlParams []interface{}
	n1qlParams = append(n1qlParams, url.LongUrl)
	query := gocb.NewN1qlQuery("SELECT `" + bucketName + "`.* FROM `" + bucketName + "` WHERE longUrl = $1")
	rows, err := bucket.ExecuteN1qlQuery(query, n1qlParams)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	var row userUrl
	rows.One(&row)
	if row == (userUrl{}) {
		hashcode := hashids.NewData()
		hToEncode, _ := hashids.NewWithData(hashcode)
		now := time.Now()
		url.ID, _ = hToEncode.Encode([]int{int(now.Unix())})
		url.ShortUrl = "http://localhost:8091/" + url.ID //Change the handler later
		bucket.Insert(url.ID, url, 60 /* <- Expiration Timer*/)  
		//Experiment with URL expiration with 0 (for first test run, have it be infinite/ 0 = infinite)
	} else {
		url = row
	}
	json.NewEncoder(w).Encode(url)
}

func extendEndpoint(w http.ResponseWriter, req *http.Request) {
	var n1qlParams []interface{}
	query := gocb.NewN1qlQuery("SELECT `" + bucketName + "`.* FROM `" + bucketName + "` WHERE shortUrl = $1")
	params := req.URL.Query()
	n1qlParams = append(n1qlParams, params.Get("shortUrl"))
	rows, _ := bucket.ExecuteN1qlQuery(query, n1qlParams)
	var row userUrl
	rows.One(&row)
	json.NewEncoder(w).Encode(row)
}

func terminalEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var url userUrl
	bucket.Get(params["id"],&url)
	http.Redirect(w, req, url.LongUrl, 301)
}

var bucket *gocb.Bucket
var bucketName string

func main() {
	router := mux.NewRouter()
	cluster, err := gocb.Connect("couchbase://127.0.0.1") //Connects to the database (My server)
	if err != nil {
		panic(err)
	}
	//To resolve any authentication errors while logging into the database (My database)
	err = cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "Skrj2468.",
	})
	if err != nil {
		panic(err)
	}
	bucketName = "shortUrls"
	bucket, err = cluster.OpenBucket(bucketName, "")
	if err != nil {
		panic(err)
	}
	router.HandleFunc("/create", initialEndpoint).Methods("PUT") //To take in a long url
	router.HandleFunc("/expand/", extendEndpoint).Methods("GET")
	router.HandleFunc("/{id}", terminalEndpoint).Methods("GET")
	fmt.Println("Running on port 8091") //Will print if running correctly
	log.Fatal(http.ListenAndServe(":8091", router)) //Reports an error if exists with server
}
