package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/handlers"
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

type Driver struct {
	Tripid        string `json:"TripID"`
	Driverid      string `json:"DriverID"`
	Firstname     string `json:"FirstName"`
	Lastname      string `json:"LastName"`
	Carlicensenum string `json:"CarLicenseNum"`
}

type Passenger struct {
	Tripid      string `json:"TripID"`
	Passengerid string `json:"PassengerID"`
	Firstname   string `json:"FirstName"`
	Lastname    string `json:"LastName"`
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
	passid := params["passengerid"]

	if r.Header.Get("Content-type") == "application/json" {
		// POST method to insert trip row into db
		if r.Method == "POST" {
			var trip Trip
			regBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(regBody, &trip)
				trip.Passengerid = passid
				driver := GetAvailDriver()
				driverid := driver.Driverid
				// check if there are available drivers
				if driver.Driverid != "nil" {
					// insert trip into db
					trip.Driverid = driverid
					if InsertTripDB(db, trip) {
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte("201 - Trip created"))
						// update driver & passenger availability to false
						UpdatePassenger(passid)
						UpdateDriver(driverid)
					} else {
						w.WriteHeader(http.StatusBadRequest)
						w.Write([]byte("400 - Unable to create trip"))
					}
					// tripid := RetrieveTripID(db)
					// if tripid != "nil" {
					// 	// TDL posts trip information to both parties
					// 	driver.Tripid = tripid
					// 	var pass Passenger
					// 	pass.Tripid = tripid
					// 	pass.Passengerid = params["passenderid"]
					// }
				} else {
					// no available driver
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte("409 - No Available Drivers"))
				}

			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please enter account details in JSON format"))
			}
		}
	}
}

func tripStartEnd(w http.ResponseWriter, r *http.Request) {
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
	tripid := TempRetrieveTripID(db, driverid)
	passid := RetrievePassID(db, tripid)
	act := r.URL.Query().Get("action")

	if r.Header.Get("Content-type") == "application/json" {
		// POST method to insert trip row into db
		if r.Method == "POST" {
			switch act {
			case "start":
				// insert tripstartdt into db
				if InsertTripTime(db, tripid, "TripStartDT") {
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte("202 - Trip start DateTime stored"))
				}
			case "end":
				// insert tripenddt into db & update users' availability
				if InsertTripTime(db, tripid, "TripEndDT") {
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte("202 - Trip end DateTime stored"))
					UpdatePassenger(passid)
					UpdateDriver(driverid)
				}
			default:
				// error of unknown action parsed in
				w.WriteHeader(http.StatusNotAcceptable)
				w.Write([]byte("406 - Unknown action"))
			}
		}
	}
}

func getPassengerTrips(w http.ResponseWriter, r *http.Request) {
	// mysql init
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/db_assignment1")

	// handle db error
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after main function has finished executing
	defer db.Close()

	params := mux.Vars(r)
	passid := params["passengerid"]

	if r.Method == "GET" {
		triplist := GetAllTrips(db, passid)
		json.NewEncoder(w).Encode(triplist)
	}

}

// returns DriverID from driver microservice if there found, otherwise "nil".
func GetAvailDriver() Driver {
	url := "http://localhost:2000/availabledrivers"
	var d Driver
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()

		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			json.Unmarshal(body, &d)
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
	return d
}

// returns details of passenger from driver microservice using id
func PassengerDetails(id string) Passenger {
	url := "http://localhost:1000/passenger/details/" + id
	var p Passenger
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()

		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			json.Unmarshal(body, &p)
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}

	return p
}

// post it to driver microservice to update availability
func UpdateDriver(driverid string) {
	url := fmt.Sprintf("http://localhost:2000/driver/%s", driverid)
	response, err := http.Post(url, "application/json", nil)
	if err != nil {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

// post it to passenger microservice to update availability
func UpdatePassenger(passengerid string) {
	url := fmt.Sprintf("http://localhost:1000/passenger/%s", passengerid)
	response, err := http.Post(url, "application/json", nil)
	if err != nil {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

// db function to insert a trip row
// returns true if insert is successful
func InsertTripDB(db *sql.DB, trip Trip) bool {
	pid := trip.Passengerid
	did := trip.Driverid
	puc := trip.Pickupcode
	doc := trip.Dropoffcode

	query := fmt.Sprintf(
		`insert into trips(TripID, PassengerID, DriverID, PickUpCode, DropOffCode) 
		 values(UUID_TO_BIN(UUID()), UUID_TO_BIN('%s'), UUID_TO_BIN('%s'),'%s','%s');`,
		pid, did, puc, doc)

	res, err := db.Exec(query)

	if err != nil {
		panic(err.Error())
	}

	rows, _ := res.RowsAffected()
	return rows == 1
}

// db function to retrieve tripid after insertion
func RetrieveTripID(db *sql.DB) string {
	id := "nil"

	query := `select BIN_TO_UUID(@uuid) as TripID;`
	result, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	if result.Next() {
		result.Scan(&id)
	}

	return id
}

// temp db function to retrieve tripid
func TempRetrieveTripID(db *sql.DB, userid string) string {
	id := "nil"

	query := fmt.Sprintf(
		`select BIN_TO_UUID(TripID) from trips 
		 where TripStartDT is null or TripEndDT is null
		 and (PassengerID = UUID_TO_BIN('%s') or
		 DriverID = UUID_TO_BIN('%s'));`, userid, userid)

	result, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	if result.Next() {
		result.Scan(&id)
	}

	return id
}

// retrieves passengerid using trip id
func RetrievePassID(db *sql.DB, tripid string) string {
	id := "nil"

	query := fmt.Sprintf(
		`select BIN_TO_UUID(PassengerID) from trips 
		 where TripID = UUID_TO_BIN('%s');`, tripid)

	result, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	if result.Next() {
		result.Scan(&id)
	}

	return id
}

// db function to insert either trip start or end datetime
// returns true if successful; false otherwise
func InsertTripTime(db *sql.DB, id string, action string) bool {
	query := fmt.Sprintf(
		`update trips set %s = NOW() 
		 where TripID = UUID_TO_BIN('%s');`,
		action, id)

	res, err := db.Exec(query)

	if err != nil {
		panic(err.Error())
	}

	rows, _ := res.RowsAffected()

	return rows == 1
}

func GetAllTrips(db *sql.DB, passid string) []Trip {
	var triplist []Trip

	query := fmt.Sprintf(
		`select BIN_TO_UUID(TripID), BIN_TO_UUID(PassengerID), BIN_TO_UUID(DriverID)
		 , PickUpCode, DropOffCode, TripStartDT, TripEndDT  
		 from Trips where PassengerID = UUID_TO_BIN('%s')
		 order by TripStartDT desc;`, passid)

	res, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}

	for res.Next() {
		var t Trip
		res.Scan(&t.Driverid, &t.Passengerid, &t.Driverid, &t.Pickupcode, &t.Dropoffcode, &t.Tripstartdt, &t.Tripenddt)
		triplist = append(triplist, t)
	}

	return triplist
}

func main() {
	// start router
	router := mux.NewRouter()
	// specify allowed headers, methods, & origins to allow CORS
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})

	router.HandleFunc("/", home)
	router.HandleFunc("/passenger/{passengerid}", requestTrip).Methods("POST")
	router.HandleFunc("/trip/{driverid}", tripStartEnd).Methods("POST")
	router.HandleFunc("/trips/{passengerid}", getPassengerTrips).Methods("GET")

	fmt.Println("listening at port 3000")
	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(headers, origins, methods)(router)))
}

// VALIDATIONS (put a 'V' to those done)
/*
check driver availability (V)

check action is appropriate in tripStartEnd function (V)

check creation of trip is successful in requestTrip function (V)

*/
