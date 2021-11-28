package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

type PassengerAccount struct {
	Passengerid  string `json:"PassengerID"`
	Firstname    string `json:"FirstName"`
	Lastname     string `json:"LastName"`
	Mobile       string `json:"MobileNo"`
	Email        string `json:"EmailAddr"`
	Availability bool   `json:"Availability"`
}

type Trips struct {
}

type TripRequest struct {
	Passengerid string `json:"PassengerID"`
	Firstname   string `json:"FirstName"`
	Lastname    string `json:"LastName"`
	PickUpCode  string `json:"PickUpCode"`
	DropOffCode string `json:"DropOffCode"`
}

type Trip struct {
	Passengerid string `json:"PassengerID"`
	Driverid    string `json:"DriverID"`
	Pickupcode  string `json:"PickUpCode"`
	Dropoffcode string `json:"DropOffCode"`
	Tripstartdt string `json:"TripStartDT"`
	Tripenddt   string `json:"TripEndDT"`
}

// // stuct NOT IN USE
// type TripDetails struct {
// 	Tripid        string `json:"TripID"`
// 	Driverid      string `json:"DriverID"`
// 	Firstname     string `json:"FirstName"`
// 	Lastname      string `json:"LastName"`
// 	Carlicensenum string `json:"CarLicenseNum"`
// }

func landing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "~ Ride-Sharing Platform ~")
}

func passenger(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Create an account")

	// mysql init
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/db_assignment1")

	// handle db error
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after main function has finished executing
	defer db.Close()

	if r.Header.Get("Content-type") == "application/json" {
		// using POST to create new account
		if r.Method == "POST" {
			// reading user input
			var newAcc PassengerAccount
			regBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to obj
				json.Unmarshal(regBody, &newAcc)
				// check if account already exist
				if !PassengerExist(db, newAcc.Mobile, newAcc.Email) {
					// TDL check if user input correct
					InsertPassengerDB(db, newAcc)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Passenger Account Created"))
				} else {
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte("409 - Account Already Exists!"))
				}

			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please enter account details in JSON format"))
			}
		}

		// using PUT to update
		if r.Method == "PUT" {
			var acc PassengerAccount
			regBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to obj
				json.Unmarshal(regBody, &acc)
				// check if account already exist. if exist:update, else:create
				// TDL separate function. one for checking existence thru id, the other from mobile and email
				if PassengerExist(db, acc.Mobile, acc.Email) {
					//TDL check if user input correct
					// update account
					UpdatePassengerDB(db, acc)
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte("202 - Passenger Account Updated"))

				} else {
					// create account
					InsertPassengerDB(db, acc)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Passenger Account Created"))
				}
			}
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "~ Welcome Passenger! ~")

	params := mux.Vars(r)
	passengerid := params["passengerid"]

	// mysql init
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/db_assignment1")

	// handle db error
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after main function has finished executing
	defer db.Close()

	// using GET to retrieve trip details
	if r.Method == "GET" {
		// retrieve details if not available
		result := IsPassengerAvail(db, passengerid)
		if !result {
			// TDL retrieve details using http.get

			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Passenger on trip"))
		} else {
			// passenger available
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("409 - Passenger not on trip"))
		}
	}

	if r.Header.Get("Content-type") == "application/json" {
		// updates availability of passenger when trip is successfully created in trip microservice
		// this method is triggered automatically
		if r.Method == "POST" {
			res := strconv.FormatBool(!IsPassengerAvail(db, passengerid))
			if UpdateAvailabilityDB(db, passengerid, res) {
				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte("202 - Passenger availability updated"))
			} else {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("409 - Unable to update"))
			}
		}
	}
}

func trips(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	passid := params["passengerid"]

	// mysql init
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/db_assignment1")

	// handle db error
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after main function has finished executing
	defer db.Close()

	// using GET to retrieve trip details
	if r.Method == "GET" {
		triparray := GetAllTrips(passid)
		for _, t := range triparray {
			json.NewEncoder(w).Encode(t)
		}
	}
}

// // function NOT IN USE
// func GetTripDetails(id string) TripDetails {
// 	url := "http://localhost:1000/passenger/" + id
// 	var t TripDetails

// 	if resp, err := http.Get(url); err == nil {
// 		defer resp.Body.Close()

// 		if body, err := ioutil.ReadAll(resp.Body); err == nil {
// 			json.Unmarshal(body, &t)
// 		} else {
// 			log.Fatal(err)
// 		}
// 	} else {
// 		log.Fatal(err)
// 	}
// 	return t
// }

func GetAllTrips(id string) []Trip {
	url := "http://localhost:3000/trips/" + id
	var t []Trip
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()

		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			json.Unmarshal(body, &t)
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
	return t
}

func InsertPassengerDB(db *sql.DB, PA PassengerAccount) {
	fn := PA.Firstname
	ln := PA.Lastname
	mn := PA.Mobile
	ea := PA.Email

	query := fmt.Sprintf(
		`insert into passengers() 
		 values(UUID_TO_BIN(UUID()), '%s','%s','%s','%s', true)`,
		fn, ln, mn, ea)
	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

func UpdatePassengerDB(db *sql.DB, PA PassengerAccount) bool {
	id := PA.Passengerid
	fn := PA.Firstname
	ln := PA.Lastname
	mn := PA.Mobile
	ea := PA.Email

	query := fmt.Sprintf(
		`update passengers set FirstName = '%s', LastName = '%s',
		 MobileNo = '%s', EmailAddr = '%s' 
		 where PassengerID = UUID_TO_BIN('%s');`, fn, ln, mn, ea, id)
	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())

	} else {
		return true
	}
}

func UpdateAvailabilityDB(db *sql.DB, id string, avail string) bool {
	query := fmt.Sprintf(
		`update passengers set Availability = %s 
		 where PassengerID = UUID_TO_BIN('%s');`, avail, id)
	res, err := db.Exec(query)

	if err != nil {
		panic(err.Error())
	}

	rows, _ := res.RowsAffected()
	return rows == 1
}

// validation check which returns true if record already exists in database
// returns false by default, true if passenger record exists
func PassengerExist(db *sql.DB, mn string, ea string) bool {
	query := fmt.Sprintf(
		`select count(*) from passengers 
		 where MobileNo = '%s' or EmailAddr = '%s'`, mn, ea)

	results, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}

	if results.Next() {
		var r int
		results.Scan(&r)
		if r > 0 {
			return true
		}
	}
	return false
}

// retrives passenger availability
// returns bool based on passengerid
func IsPassengerAvail(db *sql.DB, id string) bool {
	r := false
	query := fmt.Sprintf(
		`select Availability from passengers where 
		 PassengerID = UUID_TO_BIN('%s')`, id)

	results, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}

	if results.Next() {
		results.Scan(&r)
	}
	return r
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
	router.HandleFunc("/", landing)

	router.HandleFunc("/passenger", passenger).Methods("POST", "PUT")
	router.HandleFunc("/passenger/{passengerid}", home).Methods("GET", "POST")
	router.HandleFunc("/passenger/{passengerid}/trips", trips).Methods("GET")

	fmt.Println("listening at port 1000")
	log.Fatal(http.ListenAndServe(":1000", router))
}
