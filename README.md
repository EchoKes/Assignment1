<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#installation">Installation</a></li>
        <li><a href="#usage">Usage</a></li>
      </ul>
    </li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>


<!-- ABOUT THE PROJECT -->
## About The Project

Account page for sign up and sign in.
[![Account_Screen_Shot][account-screenshot]]

Passenger page for getting current status, viewing trips and requesting for trips.
[![Passenger_Screen_Shot][passenger-screenshot]]

Driver page for getting current status, starting and ending ride.
[![Driver_Screen_Shot][driver-screenshot]]

This assignment's requirements are to create 2 or more microservices which will be ran in the background while a frontend such as website can be used to interact with the functions in the microservices. 

There are a total of 3 microservices:
(the functions will be shown in each microservice.)
1. Passenger
  * Log in
  * View all trips (in reverse chronological order) 
2. Trip
  * Updates users status
  * Request trip (passenger)
  * Retrieves trip details
3. Driver
  * Log in
  * Start ride
  * End ride

To better understand, below is a diagram of the assignment's structure and how communications are made.
// insert ur diagram here

<p align="left">(<a href="#top">back to top</a>)</p>


### Built With

The main objective of this assignment is to put the knowledge and skills learnt about Golang to use. 
Vanilla Javascript and HTML CSS is used for the frontend.

* [Golang](https://go.dev/)
* [HTML](https://html.com/)
* [JavaScript](https://www.javascript.com/)
* [JQuery](https://jquery.com)
<p align="left">(<a href="#top">back to top</a>)</p>


<!-- GETTING STARTED -->
## Getting Started

### Prerequisites

Make sure that MySQL and Golang is downloaded on your device.

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/EchoKes/Assignment1.git
   ```
2. Install necessary libraries
   ```sh
   go get -u github.com/go-sql-driver/mysql
   go get -u github.com/gorilla/mux
   go get -u github.com/gorilla/handlers
   ```
3. Execute database start script in `/db/ExecStartDB.sql`

<p align="left">(<a href="#top">back to top</a>)</p>


<!-- USAGE EXAMPLES -->
## Usage

To start using the ride-hailing platform, follow the below steps:
1. Run all microservices in each directory
 ```sh
 cd assignment1/passenger
 go run pass.go
 ```
 ```sh
 cd assignment1/driver
 go run driver.go
 ```
 ```sh
 cd assignment1/trip
 go run trip.go
 ```
2. Open frontend by opening `account.html` in `assignment1/Frontend_web`

<p align="left">(<a href="#top">back to top</a>)</p>


<!-- ROADMAP -->
## Roadmap

- [x] Backend using Golang
- [x] Frontend using HTML JavaScript
- [ ] Tidy up both backend and frontend

<p align="left">(<a href="#top">back to top</a>)</p>


<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="left">(<a href="#top">back to top</a>)</p>


<!-- CONTACT -->
## Contact

Kester Yeo - [School Email](mailto:s10185261@connect.np.edu.sg) - s10185261@connect.np.edu.sg

Project Link: [https://github.com/EchoKes/Assignment1](https://github.com/EchoKes/Assignment1)

<p align="left">(<a href="#top">back to top</a>)</p>


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[account-screenshot]: images/account.png
[passenger-screenshot]: images/passenger.png
[driver-screenshot]: images/driver.png
