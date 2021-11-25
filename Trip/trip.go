package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

type Trip struct {
	Passengerid string `json:"PassengerID"`
	Driverid    string `json:"DriverID"`
	Pickupcode  string `json:"PickUpCode"`
	Dropoffcode string `json:"DropOffCode"`
	Tripstartdt string `json:"TripStartDT"`
	Tripenddt   string `json:"TripEndDT"`
}

type TripDetails struct {
	Driverid      string `json:"DriverID"`
	Driverfn      string `json:"DriverFirstName"`
	Driverln      string `json:"DriverLastName"`
	Carlicensenum string `json:"CarLicenseNum"`
	Passengerid   string `json:"PassengerID"`
	Passengerfn   string `json:"PassengerFirstName"`
	Passengerln   string `json:"PassengerLastName"`
	Pickupcode    string `json:"PickUpCode"`
	Dropoffcode   string `json:"DropOffCode"`
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "~ Ride-Sharing Platform ~")
}

func requestTrip(w http.ResponseWriter, r *http.Request) {
	// mysql init
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/db_assignment1")

	// handle db error
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after main function has finished executing
	defer db.Close()

	params := mux.Vars(r)

	if r.Header.Get("Content-type") == "application/json" {
		// POST method to insert trip row into db
		if r.Method == "POST" {
			var trip Trip
			regBody, err := ioutil.ReadAll(r.Body)
			trip.Passengerid = params["passengerid"]

			if err == nil {
				json.Unmarshal(regBody, &trip)
				// check if there are available drivers
				if AvailDriver() == "" {
					fmt.Println("dope")
				}

				// if AvailDriver() != "nil" {
				// 	// insert trip into db
				// 	InsertTripDB(db, trip)
				// 	w.WriteHeader(http.StatusCreated)
				// 	//w.Write([]byte("201 - Trip Created!"))
				// 	// update driver & passenger availability to false

				// } else {
				// 	// no available driver
				// 	w.WriteHeader(http.StatusConflict)
				// 	w.Write([]byte("409 - No Available Drivers"))
				// }

			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please enter account details in JSON format"))
			}
		}
	}
}

// returns DriverID from driver microservice if there found, otherwise "nil".
func AvailDriver() string {
	type Driver struct {
		Driverid      string `json:"DriverID"`
		Firstname     string `json:"FirstName"`
		Lastname      string `json:"LastName"`
		Carlicensenum string `json:"CarLicenseNum"`
	}

	url := "http://localhost:2000/availabledrivers"

	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()

		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			var d Driver
			json.Unmarshal(body, &d)
			return d.Driverid
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
	return "nil"
}

// encodes json string of passenger details after system assigned
// post it to driver microservice
func UpdateDriverMS(details TripDetails) {
	url := fmt.Sprintf("http://localhost:2000/driver/%s", details.Driverid)
	jsonValue, _ := json.Marshal(details)

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

// encodes json string of driver details after system assigned
// post it to passenger microservice
func UpdatePassengerMS(details TripDetails) {
	url := fmt.Sprintf("http://localhost:1000/Passenger/%s", details.Passengerid)
	jsonValue, _ := json.Marshal(details)

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

// db function to insert a trip row
func InsertTripDB(db *sql.DB, trip Trip) {
	pid := trip.Passengerid
	did := trip.Driverid
	puc := trip.Pickupcode
	doc := trip.Dropoffcode

	query := fmt.Sprintf(
		`insert into trips(TripID, PassengerID, DriverID, PickUpCode, DropOffCode) 
		 values(UUID_TO_BIN(UUID()), UUID_TO_BIN('%s'), UUID_TO_BIN('%s'),'%s','%s')`,
		pid, did, puc, doc)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
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
	router.HandleFunc("/passenger/{passengerid}", requestTrip).Methods("POST")

	fmt.Println("listening at port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
