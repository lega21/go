package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Region struct {
	IDREGION int    `json:"idregion"`
	REGION   string `json:"region"`
}

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/almacen")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/region/", all).Methods("GET")
	router.HandleFunc("/region/", save).Methods("POST")
	router.HandleFunc("/region/{id}", get).Methods("GET")
	router.HandleFunc("/region/{id}", update).Methods("PUT")
	router.HandleFunc("/region/{id}", delete).Methods("DELETE")

	http.ListenAndServe(":4200", router)
}

func all(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "aplication/json")

	var regions []Region

	result, err := db.Query("SELECT * FROM region")
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	for result.Next() {
		var region Region
		err := result.Scan(&region.IDREGION, &region.REGION)
		if err != nil {
			panic(err.Error())
		}
		regions = append(regions, region)
	}
	json.NewEncoder(w).Encode(regions)
}

func save(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare("INSERT INTO region(region) VALUES (?)")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	region := keyVal["region"]

	_, err = stmt.Exec(region)
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "Region Created!")
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	result, err := db.Query("SELECT * FROM region WHERE idregion = ?", params["idregion"])
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var region Region

	for result.Next() {
		err := result.Scan(&region.IDREGION, &region.REGION)
		if err != nil {
			panic(err.Error())
		}
	}

	json.NewEncoder(w).Encode(region)
}

func update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	stmt, err := db.Prepare("UPDATE region SET region = ? WHERE idregion = ?")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newRegion := keyVal["region"]

	_, err = stmt.Exec(newRegion, params["id"])
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "Region Updated")
}

func delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	stmt, err := db.Prepare("DELETE FROM region WHERE idregion = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "Region Deleted")
}
