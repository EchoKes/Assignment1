-- create database db_assignment1;

use db_assignment1; 

-- -- (for creating Passengers table)

-- create table Passengers
-- (
-- 	PassengerID binary(16) primary key, 
-- 	FirstName varchar(30),
-- 	LastName varchar(30),
-- 	MobileNo varchar(8),
-- 	EmailAddr varchar(45),
-- 	Availability tinyint(0) default 0
-- );


-- -- (for creating Drivers table)

-- create table Drivers
-- (
-- 	DriverID binary(16) primary key,
-- 	FirstName varchar(30),
-- 	LastName varchar(30),
-- 	MobileNo varchar(8),
-- 	EmailAddr varchar(45),
-- 	IdNum varchar(9),
-- 	CarLicenseNum varchar(8),
-- 	Availability tinyint(0) default 0
-- );


-- -- (for creating Trips table)

-- create table Trips
-- (
-- 	TripID binary(16),
-- 	PassengerID binary(16),
-- 	DriverID binary(16),
-- 	PickUpCode varchar(6),
-- 	DropOffCode varchar(6),
--     TripStartDT datetime,
--     TripEndDT datetime
-- );


-- -- (for searching available drivers)

-- select FirstName, LastName, CarLicenseNum from db_assignment1.drivers
-- where Availability is true
-- order by rand() limit 1;


-- -- (for updating passenger and driver availability to false [when trip start])

-- update db_assignment1.drivers
-- set Availability=false
-- where DriverID='<To be set>';

-- update db_assignment1.passengers
-- set Availability=false
-- where PassengerID='<To be set>';


-- -- (for updating passenger and driver availability to true [when trip end])

-- update db_assignment1.drivers
-- set Availability=true
-- where DriverID='<To be set>';

-- update db_assignment1.passengers
-- set Availability=true
-- where PassengerID='<To be set>';

-- insert into db_assignment1.trips()
-- values(UUID_TO_BIN(UUID()), UUID_TO_BIN(passengerid), UUID_TO_BIN(driverid), '<pucode>', 'docode', datetime(), datetime());


-- -- (for viewing trip history)

-- select BIN_TO_UUID(t.TripID) as TripID, d.FirstName, d.LastName, t.PickUpCode, t.DropOffCode, t.TripStartDT, t.TripEndDT
-- from db_assignment1.drivers as d inner join db_assignment1.trips as t
-- on t.DriverID = d.DriverID
-- order by t.TripEndDT DESC;



-- create table test (
-- 	id binary(16) primary key,
--     name varchar(255)
-- );

-- insert into test(id, name)
-- values(UUID_TO_BIN(UUID()), 'john');


-- -- for creating passenger account

-- insert into passengers()
-- values(UUID_TO_BIN(UUID()), 'kester','yeo','86882678','kesteryeo@hotmail.com', false)


-- -- for creating driver account

-- insert into drivers()
-- values(UUID_TO_BIN(UUID()), 'kester','yeo','86882678','kesteryeo@hotmail.com', 'T0114959D', 'SLC2973C', true);


-- -- for updating passenger account

-- update passengers 
-- set FirstName = "John", LastName = "Cena", MobileNo = "86882678", EmailAddr = "johncena@hotmail.com" 
-- where PassengerID = UUID_TO_BIN("8365e0c8-49ca-11ec-94a8-049226daf8e1");


-- -- for updating driver account

-- update drivers 
-- set FirstName = "John", LastName = "Cena", MobileNo = "86882678", EmailAddr = "johncena@hotmail.com", CarLicenseNum = "SLC2973C" 
-- where DriverID = UUID_TO_BIN("");



-- -- for checking if passenger account exist

-- select count(*)  from passengers
-- where MobileNo = '86882678' or EmailAddr = 'kesteryeo@hotmail.com'


-- -- for checking if driver account exist

-- select count(*)  from drivers
-- where IdNum = ''

-- -- for checking driver's availability

-- select Availability from drivers
-- where DriverID = UUID_TO_BIN("d1609e2e-4adf-11ec-9339-049226daf8e1")

-- select BIN_TO_UUID(DriverID) from drivers