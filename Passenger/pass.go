package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
	// trip request should parse in PassengerID, PickUpCode, DropOffCode
	// if system succesfully assigned driver then insert trip, update availability
	// when driver starts trip, alter db to include datetime
	Passengerid string `json:"PassengerID"`
	Driverid    string `json:"DriverID"`
	PickUpCode  string `json:"PickUpCode"`
	DropOffCode string `json:"DropOffCode"`
}

type Driver struct {
	Driverid      string `json:"DriverID"`
	Firstname     string `json:"FirstName"`
	Lastname      string `json:"LastName"`
	Mobile        string `json:"MobileNo"`
	Email         string `json:"EmailAddr"`
	IdNum         string `json:"IdNum"`
	Carlicensenum string `json:"CarLicenseNum"`
	Availability  bool   `json:"Availability"`
}

//var users map[string]PassengerAccount

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "~ Ride-Sharing Platform ~")
}

func passenger(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Create an account")
	//params := mux.Vars(r)

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

func request(w http.ResponseWriter, r *http.Request) {
	// mysql init
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/db_assignment1")

	// handle db error
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after main function has finished executing
	defer db.Close()

	if r.Header.Get("Content-type") == "application/json" {
		// using post to request a trip
		if r.Method == "POST" {
			var tripReq TripRequest
			regBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to obj
				json.Unmarshal(regBody, &tripReq)
				// retrieve driver details if trip request successful
				driver := DriversAvail(db)
				if len(driver.Driverid) > 0 {
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Driver Found!"))
				} else {
					w.WriteHeader(http.StatusNotFound)
					//TDL might have to change status
					w.Write([]byte("404 - No Available Driver"))
				}
			}
		}
	}
}

// test function to return driver if avail
func getDriver(w http.ResponseWriter, r *http.Request) {
	// mysql init
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/db_assignment1")

	// handle db error
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after main function has finished executing
	defer db.Close()

	driver := DriversAvail(db)
	if len(driver.Driverid) > 0 {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("201 - Driver found!"))
		json.NewEncoder(w).Encode(driver)
	} else {
		w.WriteHeader(http.StatusNotFound)
		//TDL might have to change status
		w.Write([]byte("404 - No Available Driver"))
	}

}

func InsertPassengerDB(db *sql.DB, PA PassengerAccount) {
	fn := PA.Firstname
	ln := PA.Lastname
	mn := PA.Mobile
	ea := PA.Email

	query := fmt.Sprintf(
		`insert into passengers() 
		 values(UUID_TO_BIN(UUID()), '%s','%s','%s','%s', false)`,
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

// retrieves available driver
// returns driver's name and car license number if found
func DriversAvail(db *sql.DB) Driver {
	var driver Driver

	query := fmt.Sprintln(
		`select BIN_TO_UUID(DriverID), FirstName, LastName, CarLicenseNum 
		 from db_assignment1.drivers where Availability is true 
		 order by rand() limit 1;`)

	results, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}

	if results.Next() {
		results.Scan(&driver.Driverid, &driver.Firstname, &driver.Lastname, &driver.Carlicensenum)
	}

	return driver
}

// TDL Remove this function after updating passengeraccount struct to include PassengerID
// retrieval of PassengerID to aid UpdatePassengerDB() function
// func GetCurrentPID(db *sql.DB, ea string) string {
// 	id := ""
// 	query := fmt.Sprintf(
// 		`select BIN_TO_UUID(PassengerID) from passengers
// 		 where EmailAddr = '%s'`, ea)

// 	results, err := db.Query(query)

// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	for results.Next() {
// 		var r string
// 		results.Scan(&r)
// 		id = r
// 	}
// 	return id
// }

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

	router.HandleFunc("/passenger", passenger).Methods("POST", "PUT")
	router.HandleFunc("/passenger/request", request).Methods("POST")
	router.HandleFunc("/getDriver", getDriver)

	fmt.Println("listening at port 1000")
	log.Fatal(http.ListenAndServe(":1000", router))
}
