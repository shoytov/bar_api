package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

type Guest struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
	Items  []Item `json:"items"`
}

type Item struct {
	Id    string `json:"id"`
	Value uint8  `json:"value"`
}

var guests []Guest

func index(w http.ResponseWriter, r *http.Request) {
	// главная страница
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, nil)
}

func addGuest(w http.ResponseWriter, r *http.Request) {
	// добавление гостя
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusCreated)

	var guest Guest

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &guest)

	_id, _ := exec.Command("uuidgen").Output()
	guest.Id = strings.Replace(string(_id), "\n", "", 1)
	guest.Active = true

	guests = append(guests, guest)

	json.NewEncoder(w).Encode(guest)
}

func getGuests(w http.ResponseWriter, r *http.Request) {
	// вывод спсика гостей
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	json.NewEncoder(w).Encode(guests)
}

func deleteGuests(w http.ResponseWriter, r *http.Request) {
	// удаление всех гостей
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusNoContent)

	guests = nil
}

func makeGuestInactive(w http.ResponseWriter, r *http.Request) {
	// отметка гостя как неактивного/активного
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusNoContent)

	params := mux.Vars(r)
	for index, item := range guests {
		if item.Id == params["id"] {
			reqBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(reqBody, &guests[index])

			return
		}
	}
}

func deleteGuest(w http.ResponseWriter, r *http.Request) {
	// удаление гостя
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusNoContent)

	params := mux.Vars(r)
	for index, item := range guests {
		if item.Id == params["id"] {
			guests = append(guests[:index], guests[index+1:]...)
			return
		}
	}
}

func addGuestItem(w http.ResponseWriter, r *http.Request) {
	// добавление покупки гостю
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusCreated)

	params := mux.Vars(r)
	for index, item := range guests {
		if item.Id == params["guest_id"] {
			reqBody, _ := ioutil.ReadAll(r.Body)
			var guestItem Item
			json.Unmarshal(reqBody, &guestItem)
			guestItem.Id = strconv.Itoa(rand.Intn(1000000))
			guests[index].Items = append(guests[index].Items, guestItem)

			return
		}
	}
}

func deleteGuestItem(w http.ResponseWriter, r *http.Request) {
	// удаление покупки у гостя
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusNoContent)

	params := mux.Vars(r)
	for index, guest := range guests {
		if guest.Id == params["guest_id"] {
			for index1, guestItem := range guests[index].Items {
				if guestItem.Id == params["id"] {
					guests[index].Items = append(guests[index].Items[:index1], guests[index].Items[index1+1:]...)
					return
				}
			}
		}
	}
}

func handlerRequests() {
	// обработчик маршрутов
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/guests", addGuest).Methods("POST")
	r.HandleFunc("/guests", getGuests).Methods("GET")
	r.HandleFunc("/guests", deleteGuests).Methods("DELETE")
	r.HandleFunc("/guests/{id}", makeGuestInactive).Methods("PATCH")
	r.HandleFunc("/guests/{id}", deleteGuest).Methods("DELETE")
	r.HandleFunc("/items/{guest_id}", addGuestItem).Methods("POST")
	r.HandleFunc("/items/{guest_id}/{id}", deleteGuestItem).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func main() {
	handlerRequests()
}
