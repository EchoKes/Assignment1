use db_assignment1; 

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

-- -- for inserting trip before start trip end trip

-- insert into trips(TripID, PassengerID, DriverID, PickUpCode, DropOffCode)
-- values(UUID_TO_BIN(UUID()), UUID_TO_BIN(passengerid), UUID_TO_BIN(driverid), '460160', '510294');


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


-- select BIN_TO_UUID(PassengerID), FirstName, Availability from passengers;

-- insert into trips(TripID, PassengerID, DriverID, PickUpCode, DropOffCode) 
-- values(UUID_TO_BIN(UUID()), UUID_TO_BIN('0d2a0132-4e07-11ec-87b6-049226daf8e1'), UUID_TO_BIN('d1609e2e-4adf-11ec-9339-049226daf8e1'),'460100','521111');


-- delete from trips where TripID = UUID_TO_BIN('0089a83a-4e0f-11ec-87b6-049226daf8e1');
-- select BIN_TO_UUID(TripID) from trips

-- select BIN_TO_UUID(DriverID), FirstName, LastName, CarLicenseNum 
-- from drivers where Availability is false
-- order by rand() limit 1;

-- select BIN_TO_UUID(TripID) as TripID from trips;
-- delete from trips where TripID = uuid_to_bin("c6bc6211-4fe6-11ec-8928-0433c219c4d5");

-- select BIN_TO_UUID(TripID) from trips 
-- where TripStartDT is null or TripEndDT is null
-- and (PassengerID = uuid_to_bin("9dc1cddb-4ea0-11ec-8928-0433c219c4d5") or
-- DriverID = uuid_to_bin("9dc1cddb-4ea0-11ec-8928-0433c219c4d5")); 

select * from trips

-- select BIN_TO_UUID(TripID) from trips 
-- where PassengerID = "9dc1cddb-4ea0-11ec-8928-0433c219c4d5" or DriverID = "9dc1cddb-4ea0-11ec-8928-0433c219c4d5" 
-- and if(TripStartDT is null,1,0)=1;

-- select bin_to_uuid(PassengerID) from passengers