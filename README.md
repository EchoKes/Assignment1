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

**Account page for sign up and sign in.**

![Account_Screen_Shot][account-screenshot1]
---
![Account_Screen_Shot][account-screenshot2]
---

**Passenger page for getting current status, viewing trips and requesting for trips.**

![Passenger_Screen_Shot][passenger-screenshot]
---

**Driver page for getting current status, starting and ending ride.**

![Driver_Screen_Shot][driver-screenshot]
---

The requirements were to create 2 or more microservices. 
3 microservices were created which will be ran in the background while a monolithic frontend such as website can be used to interact with the functions in the microservices. 

The 3 microservices:
(the functions will be shown in each microservice.)
1. Passenger
  * Log in
  * View all trips (in reverse chronological order) 
  * Delete
2. Trip
  * Updates users status
  * Request trip (passenger)
  * Retrieves trip details
3. Driver
  * Log in
  * Start ride
  * End ride
  * Delete

To better understand the interactions, below are 2 diagrams; The microservice architecture and frontend diagram.
![Microservice-Diagram][msdiag-screenshot]
---
![Frontend-Diagram][fediag-screenshot]
---
<p align="left">(<a href="#top">back to top</a>)</p>


### Built With

The main objective of this assignment is to put the knowledge and skills learnt about Golang to use. 
The chosen tech stack is:
- Vanilla Javascript and HTML5 frontend. (Monolithic)
- Golang for backend. (Microservices)
- MySQL for database.


* [Golang](https://go.dev/)
* [HTML](https://html.com/)
* [JavaScript](https://www.javascript.com/)
* [JQuery](https://jquery.com/)
* [MySQL](https://www.mysql.com/)
<p align="left">(<a href="#top">back to top</a>)</p>


<!-- GETTING STARTED -->
## Getting Started

### Prerequisites

Please ensure that MySQL and Golang is installed and operational on your device.

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
- [X] Set up MySQL database
- [x] Create microservices backend using Golang
- [x] Create frontend using HTML JavaScript
- [ ] Tidy up both backend and frontend

<p align="left">(<a href="#top">back to top</a>)</p>


<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="left">(<a href="#top">back to top</a>)</p>


<!-- CONTACT -->
## Contact

Kester Yeo - [School Email](mailto:s10185261@connect.np.edu.sg) 

Project Link: [https://github.com/EchoKes/Assignment1](https://github.com/EchoKes/Assignment1)

<p align="left">(<a href="#top">back to top</a>)</p>


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[account-screenshot1]: ./images/login.PNG
[account-screenshot2]: ./images/register_driver.PNG
[passenger-screenshot]: ./images/viewtrips.PNG
[driver-screenshot]: ./images/driverpage.PNG
[msdiag-screenshot]: ./images/msdiag.PNG
[fediag-screenshot]: ./images/fediag.PNG
