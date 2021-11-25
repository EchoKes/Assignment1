package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

type Trip struct {
	// trip request should parse in PassengerID, PickUpCode, DropOffCode
	// if system succesfully assigned driver then insert trip, update availability
	// when driver starts trip, alter db to include datetime
	Passengerid string `json:"PassengerID"`
	Driverid    string `json:"DriverID"`
	PickUpCode  string `json:"PickUpCode"`
	DropOffCode string `json:"DropOffCode"`
	Tripstartdt string `json:"TripStartDT"`
	Tripenddt   string `json:"TripEndDT"`
}

func home(w http.ResponseWriter, r *http.Request) {

}

func getAvailDriver(w http.ResponseWriter, r *http.Request) {
	url := "http://localhost:1000/getDriver"

	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			fmt.Println(string(body))
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
}

func main() {
	// // mysql init
	// db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/db_assignment1")

	// // handle db error
	// if err != nil {
	// 	panic(err.Error())
	// }

	// //insert db test function below

	// // defer the close till after main function has finished executing
	// defer db.Close()

	// start router
	router := mux.NewRouter()
	router.HandleFunc("/", home)

	router.HandleFunc("/getAvailDriver", getAvailDriver).Methods("GET")

	fmt.Println("listening at port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
