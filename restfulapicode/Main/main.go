// Created by Kartik and Sai
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

//Person struct (model)
type Person struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Country string `json:"country"`
}

// init persons variable as slice
var persons []Person

func deletePerson(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	for index, value := range persons {
		// It checks for id and if the id is found delete operation is done here
		if value.ID == params["id"] {
			persons = append(persons[:index], persons[index+1:]...)
			reqBodyBytes := new(bytes.Buffer)
			json.NewEncoder(reqBodyBytes).Encode(persons)

			reqBodyBytes.Bytes()
			//After operation the data is written to the file
			writeToFile([]byte(reqBodyBytes.Bytes()))

			break
		}
	}
	json.NewEncoder(writer).Encode(persons)
}

func updatePerson(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	for index, value := range persons {
		// Update person by matching id which is received by request parameter
		if value.ID == params["id"] {
			persons = append(persons[:index], persons[index+1:]...)
			var person Person
			_ = json.NewDecoder(request.Body).Decode(&person)
			person.ID = params["id"]
			persons = append(persons, person)
			json.NewEncoder(writer).Encode(person)
			reqBodyBytes := new(bytes.Buffer)
			json.NewEncoder(reqBodyBytes).Encode(persons)

			reqBodyBytes.Bytes()
			//Update data in file
			writeToFile([]byte(reqBodyBytes.Bytes()))

		}
	}
	json.NewEncoder(writer).Encode(persons)

}

func createPerson(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var person Person
	_ = json.NewDecoder(request.Body).Decode(&person)
	person.ID = strconv.Itoa(rand.Intn(1000000))
	persons = append(persons, person)
	json.NewEncoder(writer).Encode(person)

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(persons)

	reqBodyBytes.Bytes()
	//Write data to file 
	writeToFile([]byte(reqBodyBytes.Bytes()))

}

func getPerson(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	//get params
	params := mux.Vars(request)

	//find book with id by looping through persons
	for _, value := range persons {
		if value.ID == params["id"] {
			json.NewEncoder(writer).Encode(value)
			return
		}
	}
}

func getPersons(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(persons)
}

func main() {
	//Initialize the router
	router := mux.NewRouter()

	//mock data
	persons = append(persons, Person{ID: "1", Name: "Kartik", Email: "ksethi1@umbc.edu", Country: "INDIA"})
	persons = append(persons, Person{ID: "2", Name: "Sai", Email: "sai@umbc.edu", Country: "USA"})
	//route handler / endpoints
	router.HandleFunc("/api/persons", getPersons).Methods("GET")
	router.HandleFunc("/api/persons/{id}", getPerson).Methods("GET")
	router.HandleFunc("/api/persons", createPerson).Methods("POST")
	router.HandleFunc("/api/persons/{id}", updatePerson).Methods("PUT")
	router.HandleFunc("/api/persons/{id}", deletePerson).Methods("DELETE")

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(persons)

	reqBodyBytes.Bytes()

	writeToFile([]byte(reqBodyBytes.Bytes()))
	log.Fatal(http.ListenAndServe(":8000", router))

}

func writeToFile(jsonBlob []byte) {
	err := json.Unmarshal(jsonBlob, &persons)
	personsJson, _ := json.Marshal(persons)
	err = ioutil.WriteFile("output.json", personsJson, 0644)
	fmt.Printf("The records are:")
	fmt.Printf("%+v", persons)
}
