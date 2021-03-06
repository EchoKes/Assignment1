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

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type DriverAccount struct {
	Driverid      string `json:"DriverID"`
	Firstname     string `json:"FirstName"`
	Lastname      string `json:"LastName"`
	Mobile        string `json:"MobileNo"`
	Email         string `json:"EmailAddr"`
	IdNum         string `json:"IdNum"`
	Carlicensenum string `json:"CarLicenseNum"`
	Availability  bool   `json:"Availability"`
}

// // struct NOT IN USE
// type TripDetails struct {
// 	Tripid      string `json:"TripID"`
// 	Passengerid string `json:"PassengerID"`
// 	Firstname   string `json:"FirstName"`
// 	Lastname    string `json:"LastName"`
// }

func landing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "~ Ride-Sharing Platform ~")
}

func driver(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Create an account")

	// mysql init
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/db_assignment1")

	// handle db error
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after main function has finished executing
	defer db.Close()

	// using DELETE to delete
	// in this case, delete function is not allowed
	if r.Method == "DELETE" {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("406 - Unable to delete"))
	}

	if r.Header.Get("Content-type") == "application/json" {
		// using POST to create new account
		if r.Method == "POST" {
			// reading user input
			var newAcc DriverAccount
			regBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to obj
				json.Unmarshal(regBody, &newAcc)
				// check if account already exist
				if !DriverExist(db, newAcc) {
					if InsertDriverDB(db, newAcc) {
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte("201 - Driver Account Created"))
					} else {
						w.WriteHeader(http.StatusBadRequest)
						w.Write([]byte("400 - Unable to create account"))
					}
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
			var acc DriverAccount
			regBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to obj
				json.Unmarshal(regBody, &acc)
				// check if account already exist. if exist:update, else:create
				if DriverExist(db, acc) {
					driverid := RetrieveDriverID(db, acc.Mobile)
					acc.Driverid = driverid
					// update account
					if UpdatePassengerDB(db, acc) {
						w.WriteHeader(http.StatusAccepted)
						w.Write([]byte("202 - Driver account updated"))
					} else {
						w.WriteHeader(http.StatusBadRequest)
						w.Write([]byte("400 - Unable to update account"))
					}
				}
				// else {
				// 	// create account
				// 	if InsertDriverDB(db, acc) {
				// 		w.WriteHeader(http.StatusCreated)
				// 		w.Write([]byte("201 - Driver account created"))
				// 	} else {
				// 		w.WriteHeader(http.StatusBadRequest)
				// 		w.Write([]byte("400 - Unable to create account"))
				// 	}
				// }
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please enter account details in JSON format"))
			}
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	// mysql init
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/db_assignment1")

	// handle db error
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after main function has finished executing
	defer db.Close()

	params := mux.Vars(r)
	driverid := params["driverid"]

	if r.Method == "GET" {
		// if availability is false, driver has been assigned a trip
		if !GetDriverAvail(db, driverid) {
			// TDL retrieve details
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Driver on trip"))

		} else {
			// driver available
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("409 - Driver not assigned to a trip"))
		}
	}
}

func details(w http.ResponseWriter, r *http.Request) {
	// mysql init
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/db_assignment1")

	// handle db error
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after main function has finished executing
	defer db.Close()

	params := mux.Vars(r)
	driverid := params["driverid"]

	if r.Method == "GET" {
		// retrieve details if driver exists
		result := DriverIdExist(db, driverid)
		if result {
			dDetails := RetrieveDriverDetails(db, driverid)
			json.NewEncoder(w).Encode(dDetails)
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Driver details retrieved"))

		} else {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("409 - Unable to retrieve driver details"))
		}
	}

	if r.Header.Get("Content-type") == "application/json" {
		// updates availability of passenger when trip is successfully created in trip microservice
		// this method is triggered automatically
		if r.Method == "POST" {
			res := strconv.FormatBool(!GetDriverAvail(db, driverid))
			if UpdateAvailabilityDB(db, driverid, res) {
				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte("202 - Driver availability updated"))
			} else {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("409 - Unable to update"))
			}
		}
	}
}

// for communcating with rest api in trips microservice
func availDrivers(w http.ResponseWriter, r *http.Request) {
	// mysql init
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/db_assignment1")

	// handle db error
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after main function has finished executing
	defer db.Close()

	// using GET to retrieve a available driver
	if r.Method == "GET" {
		// sent regardless of avail or non avail
		// validation to be done in trips microservice
		driver := DriversAvail(db)
		json.NewEncoder(w).Encode(driver)

		// // available drivers
		// if driver.Driverid != "nil" {
		// 	w.WriteHeader(http.StatusAccepted)
		// 	w.Write([]byte("202 - Driver found"))
		// } else {
		// 	// no available drivers
		// 	w.WriteHeader(http.StatusConflict)
		// 	w.Write([]byte("409 - Driver not found"))
		// }
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	mobile := params["mobile"]

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
		// check that passenger is in database
		// proceeds if true
		id := RetrieveDriverID(db, mobile)
		if id != "nil" {
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte(id))
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Mobile number incorrect"))
		}

	}
}

// retrieve driver details
func RetrieveDriverDetails(db *sql.DB, id string) DriverAccount {
	var da DriverAccount
	query := fmt.Sprintf(
		`select FirstName, LastName, MobileNo, EmailAddr, IdNum, CarLicenseNum
		 from drivers where DriverID = UUID_TO_BIN('%s')`, id)

	results, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}

	if results.Next() {
		results.Scan(&da.Firstname, &da.Lastname, &da.Mobile, &da.Email, &da.IdNum, &da.Carlicensenum)
	}
	return da
}

func RetrieveDriverID(db *sql.DB, mobile string) string {
	id := "nil"
	query := fmt.Sprintf(
		`select BIN_TO_UUID(DriverID) from drivers 
		 where MobileNo = '%s';`, mobile)

	res, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}

	if res.Next() {
		res.Scan(&id)
	}

	return id
}

// validation check which returns true if record already exists in database
// returns false by default, true if driver record exists
func DriverIdExist(db *sql.DB, id string) bool {
	query := fmt.Sprintf(
		`select count(*) from drivers 
		 where DriverID = UUID_TO_BIN('%s')`, id)

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

// validation check which returns true if record already exists in database
// returns false by default, true if driver record exists
func DriverExist(db *sql.DB, acc DriverAccount) bool {
	idn := acc.Driverid
	mn := acc.Mobile
	ea := acc.Email
	cln := acc.Carlicensenum

	query := fmt.Sprintf(
		`select count(*) from drivers 
		 where IdNum = '%s' or EmailAddr = '%s' 
		 or MobileNo = '%s' or CarLicenseNum = '%s'`, idn, ea, mn, cln)

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

func GetDriverAvail(db *sql.DB, driverid string) bool {
	r := false
	query := fmt.Sprintf(
		`select Availability from drivers
		 where DriverID = UUID_TO_BIN('%s')`, driverid)
	results, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}

	if results.Next() {
		results.Scan(&r)
	}
	return r
}

func InsertDriverDB(db *sql.DB, DA DriverAccount) bool {
	fn := DA.Firstname
	ln := DA.Lastname
	mn := DA.Mobile
	ea := DA.Email
	idn := DA.IdNum
	cln := DA.Carlicensenum

	query := fmt.Sprintf(
		`insert into drivers() 
		 values(UUID_TO_BIN(UUID()), '%s','%s','%s','%s','%s','%s', true)`,
		fn, ln, mn, ea, idn, cln)

	res, err := db.Exec(query)

	if err != nil {
		panic(err.Error())
	}

	rows, _ := res.RowsAffected()
	return rows == 1
}

func UpdatePassengerDB(db *sql.DB, DA DriverAccount) bool {
	id := DA.Driverid
	fn := DA.Firstname
	ln := DA.Lastname
	mn := DA.Mobile
	ea := DA.Email
	cln := DA.Carlicensenum

	query := fmt.Sprintf(
		`update drivers set FirstName = '%s', LastName = '%s',
		 MobileNo = '%s', EmailAddr = '%s', CarLicenseNum = '%s'  
		 where DriverID = UUID_TO_BIN('%s');`, fn, ln, mn, ea, cln, id)

	res, err := db.Exec(query)

	if err != nil {
		panic(err.Error())
	}

	rows, _ := res.RowsAffected()
	return rows == 1
}

func UpdateAvailabilityDB(db *sql.DB, id string, avail string) bool {
	query := fmt.Sprintf(
		`update drivers set Availability = %s 
		 where DriverID = UUID_TO_BIN('%s');`, avail, id)
	res, err := db.Exec(query)

	if err != nil {
		panic(err.Error())
	}

	rows, _ := res.RowsAffected()
	return rows == 1
}

// retrieves available driver
// returns driver's details
func DriversAvail(db *sql.DB) DriverAccount {
	var driver DriverAccount

	query := fmt.Sprintln(
		`select BIN_TO_UUID(DriverID), FirstName, LastName, CarLicenseNum, Availability 
		 from drivers where Availability is true 
		 order by rand() limit 1;`)

	results, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}

	if results.Next() {
		results.Scan(&driver.Driverid, &driver.Firstname, &driver.Lastname, &driver.Carlicensenum, &driver.Availability)
	} else {
		driver.Driverid = "nil"
	}

	return driver
}

func main() {
	//start router
	router := mux.NewRouter()
	// specify allowed headers, methods, & origins to allow CORS
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	router.HandleFunc("/", landing)
	router.HandleFunc("/driver", driver).Methods("POST", "PUT", "DELETE")
	router.HandleFunc("/availabledrivers", availDrivers).Methods("GET")
	router.HandleFunc("/driver/{driverid}", home).Methods("GET")
	router.HandleFunc("/driver/{driverid}/details", details).Methods("GET", "POST")
	router.HandleFunc("/driver/{mobile}/id", login).Methods("GET")

	fmt.Println("listening at port 2000")
	log.Fatal(http.ListenAndServe(":2000", handlers.CORS(headers, origins, methods)(router)))
}

// VALIDATIONS (put a 'V' to those done)
/*
check driver exist using their uuid:
- in home "GET" "POST" function (V)

check driver availability (V)

check creation and update of account is successful in driver function (V)

check if any available drivers. none should set driver id as nil before returning (V)
*/
