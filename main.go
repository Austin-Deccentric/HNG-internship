package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"os"
	"github.com/joho/godotenv"
	"html/template"
	"errors"
	"log"
	"strconv"
	"fmt"
	"bytes"
	"encoding/json"
	//"io/ioutil"
)

type ClientRequest struct {
	Code string
	Amount int64
	Phonenumber string
	SecretKey string
}

var (
	templates = template.Must(template.ParseGlob("views/*html"))
)

func main () {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(8000)
	}

	r := mux.NewRouter()
	
	r.Handle("/", http.FileServer(http.Dir("./views"))).Methods("GET") // root address renders homepage
	r.HandleFunc("/success", page).Methods("GET") // 
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.HandleFunc("/purchase", NotImplemented).Methods("GET")

	r.HandleFunc("/transact", handlePost).Methods("POST")

	fmt.Printf("Listening and serving on port %s.....\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Not Implemented"))
  })

func handlePost(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("network")
	phoneNumber := r.FormValue("number")
	amount := r.FormValue("amount")
	//fmt.Println(code, phoneNumber, amount)
	
	purchase(w, r, code, phoneNumber, amount)
}

const secretKey = "hfucj5jatq8h"
var url = "https://sandbox.wallets.africa/bills/airtime/purchase"


func purchase(w http.ResponseWriter, r *http.Request,code, number, amt string) {
	//fmt.Println("From form",amt)
	parseamt,_ := strconv.ParseInt(amt, 10, 64)
	//fmt.Println(parseamt)
	NewClient := &ClientRequest {
		Code: code,
		Amount: parseamt,
		Phonenumber: number,
		SecretKey: secretKey,
	}
	const publicToken = "uvjqzm5xl6bw"  // Todo: to change
 	var bearer = "Bearer " + publicToken
	
	requestBody, err := json.Marshal(NewClient); if err!= nil{
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

  client := &http.Client {
  }
  req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))

  if err != nil {
	  w.WriteHeader(http.StatusInternalServerError)
    fmt.Println(err)
  }
  req.Header.Add("Content-Type", "application/json")
  req.Header.Add("Authorization", bearer)

  res, err := client.Do(req); if err!= nil{
	w.WriteHeader(http.StatusBadRequest)
    log.Println(err)
  }

  defer res.Body.Close()

  responsebody := response{}
  json.NewDecoder(res.Body).Decode(&responsebody)

  if res.Status == "200 OK" {
	http.Redirect(w, r, "/success", http.StatusFound)  // popup message showing success
  }else {
	  w.WriteHeader(http.StatusBadRequest)
	  fmt.Fprintln(w, responsebody.Message)	// popup message showing error message
  }

  
  //fmt.Println(responsebody.Message)

  //body, err := ioutil.ReadAll(res.Body)

  //fmt.Println("Response body:",string(body))
  //fmt.Println(string(body))
  fmt.Println("response Status:", res.Status)
}

func page(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf8")
	err := templates.ExecuteTemplate(w, "success.html", nil); if err != nil{
		http.Error(w, errors.New("Something went wrong. If this continues contact an admin").Error(), http.StatusInternalServerError)
		log.Println("error loading template",err)
	}
}

type response struct {
	ResponseCode string
	Message string
}