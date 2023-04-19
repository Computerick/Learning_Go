package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Computerick/api-crud-persistencia-arquivo/domain"
	"github.com/Computerick/api-crud-persistencia-arquivo/domain/person"
)

func main() {
	personService, err := person.NewService("person.json")
	if err != nil {
		fmt.Println("Erro trying to create person service")
		return
	}

	// Handle recebe response e request, o padrão de rota é person
	http.HandleFunc("/person/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var person domain.Person
			// Decodificador do Body (corpo)
			err := json.NewDecoder(r.Body).Decode(&person)
			if err != nil {
				fmt.Printf("Error trying to decode body. Body should be a json. Error: %s", err.Error())
				http.Error(w, "Error trying to create person", http.StatusBadRequest)
				return
			}
			// Verifica id nulo
			if person.Id <= 0 {
				http.Error(w, "Error trying to create person. Id should be a positive integer", http.StatusBadRequest)
				return
			}
			//Criar pessoa
			err = personService.Create(person)
			if err != nil {
				fmt.Printf("Error trying to create person. %s", err.Error())
				http.Error(w, "Error trying to create person", http.StatusBadRequest)
				return
			}

			// resposta (w)
			w.WriteHeader(http.StatusCreated)
			return
		}
		if r.Method == "GET" {
			path := strings.TrimPrefix(r.URL.Path,"/person/")
			if path == "" {
				// /person list all
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusOK)
				people := personService.List()
				err := json.NewEncoder(w).Encode(people)
				if err != nil{
					http.Error(w, "Error trying to list people", http.StatusInternalServerError)
					return
				}
				// person/{id}, precisando converter o path(id) com "strconv" pra inteiro
			}else {
				personId, err := strconv.Atoi(path)
				if err != nil{
					http.Error(w, "Invalid id given. Person Id must be an integer", http.StatusBadRequest)
				}
				person, err := personService.GetbyId(personId)
				if err != nil{
					http.Error(w,err.Error(), http.StatusNotFound)
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-type","application/json")
				err = json.NewEncoder(w).Encode(person)
				if err != nil{
					http.Error(w, "Error trying to list people", http.StatusBadRequest)
					return
				}
			}	
		}
		if r.Method == "PUT" {
			var person domain.Person
			// Decodificador do Body (corpo)
			err := json.NewDecoder(r.Body).Decode(&person)
			if err != nil {
				fmt.Printf("Error trying to decode body. Body should be a json. Error: %s", err.Error())
				http.Error(w, "Error trying to update person", http.StatusBadRequest)
				return
			}
			// Verifica id nulo
			if person.Id <= 0 {
				http.Error(w, "Error trying to update person. Id should be a positive integer", http.StatusBadRequest)
				return
			}
			//Atualizar pessoa
			err = personService.Update(person)
			if err != nil {
				fmt.Printf("Error trying to update person. %s", err.Error())
				http.Error(w, "Error trying to update person", http.StatusBadRequest)
				return
			}

			// resposta (w)
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == "DELETE" {
			path := strings.TrimPrefix(r.URL.Path, "/person/")
			if path == "" {
				http.Error(w, "Id must be provided in the url", http.StatusBadRequest)
				return
			}else{
				personId, err := strconv.Atoi(path)
				if err != nil {
					http.Error(w, "Invalid id given, person Id must be an integer", http.StatusBadRequest)
				}
				err = personService.DeleteById(personId)
				if err != nil {
					fmt.Printf("Error trying to delete person. %s", err.Error())
					http.Error(w, "Error trying to delete person", http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusOK)
			}	
		}
	})
	http.ListenAndServe(":8080", nil)
}
