CREATE USER 'user'@'%' IDENTIFIED BY 'password';
GRANT ALL ON *.* TO 'user'@'%';

create database db_assignment1;
use db_assignment1; 

-- -- (for creating Passengers table)

create table passengers
(
	PassengerID binary(16) primary key, 
	FirstName varchar(30),
	LastName varchar(30),
	MobileNo varchar(8),
	EmailAddr varchar(45),
	Availability tinyint(0) default 0
);


-- -- (for creating Drivers table)

create table drivers
(
	DriverID binary(16) primary key,
	FirstName varchar(30),
	LastName varchar(30),
	MobileNo varchar(8),
	EmailAddr varchar(45),
	IdNum varchar(9),
	CarLicenseNum varchar(8),
	Availability tinyint(0) default 0
);


-- -- (for creating Trips table)

create table trips
(
	TripID binary(16) primary key,
	PassengerID binary(16),
	DriverID binary(16),
	PickUpCode varchar(6),
	DropOffCode varchar(6),
    TripStartDT datetime,
    TripEndDT datetime
);

-- -- (for creating trigger that runs on every insert into trips table) : currently not in use as unnecessary

-- create trigger uuid_getter
-- after insert on trips
-- for each row
-- set @uuid = new.TripID;
