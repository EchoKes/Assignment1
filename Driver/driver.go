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

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "~ Ride-Sharing Platform ~\n")

	// mysql init
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/db_assignment1")

	// handle db error
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after main function has finished executing
	defer db.Close()

	if r.Header.Get("Content-type") == "application/json" {
		// using POST to start trip
		// if availability is false, driver has been assigned a trip
		params := mux.Vars(r)

		if !GetDriverStatus(db, params["driverid"]) {
			var s string
			r.URL.Query().Get(s)
			fmt.Println(s)
			// insert tripstartdt
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Trip Start DateTime Stored"))
		}
	}
}

func driver(w http.ResponseWriter, r *http.Request) {
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
			var newAcc DriverAccount
			regBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to obj
				json.Unmarshal(regBody, &newAcc)
				// check if account already exist
				if !DriverExist(db, newAcc.IdNum) {
					InsertDriverDB(db, newAcc)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Driver Account Created"))
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
				if DriverExist(db, acc.IdNum) {
					//TDL check if user input correct
					// update account
					UpdatePassengerDB(db, acc)
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte("202 - Driver Account Updated"))

				} else {
					// create account
					InsertDriverDB(db, acc)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Driver Account Created"))
				}
			}
		}
	}
}

// validation check which returns true if record already exists in database
// returns false by default, true if driver record exists
func DriverExist(db *sql.DB, idNum string) bool {
	query := fmt.Sprintf(
		`select count(*) from drivers 
		 where IdNum = '%s'`, idNum)

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

func GetDriverStatus(db *sql.DB, driverid string) bool {

	query := fmt.Sprintf(
		`select Availability from drivers
		 where DriverID = UUID_TO_BIN('%s')`, driverid)
	results, err := db.Query(query)

	if err != nil {
		return false
		//panic(err.Error())
	}

	if results.Next() {
		var r bool
		results.Scan(&r)
		fmt.Println(r)
		if r {
			return true
		}
	}
	return false
}

func InsertDriverDB(db *sql.DB, DA DriverAccount) {
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

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
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
	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())

	} else {
		return true
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

	//start router

	router := mux.NewRouter()
	router.HandleFunc("/{driverid}", home)

	router.HandleFunc("/driver", driver).Methods("POST", "PUT")

	fmt.Println("listening at port 2000")
	log.Fatal(http.ListenAndServe(":2000", router))
}
