DROP TABLE IF EXISTS CENKeys;
DROP TABLE IF EXISTS CENReport;

CREATE TABLE `CENKeys` (
   `cenKey`   varchar(32) DEFAULT "", 
   `reportID` varchar(64) DEFAULT "",
   `reportTS` int,
   PRIMARY KEY (`cenKey`, `reportID`),
   KEY (`reportID`),
   KEY (`reportTS`),
   KEY (`cenKey`)
);

CREATE TABLE `CENReport` (
   `reportID` varchar(64) DEFAULT "",
   `report`     varchar(4000) DEFAULT "",
   `reportMimeType` varchar(64) DEFAULT "",
   `reportTS` int,
   PRIMARY KEY (`reportID`),
   KEY (`reportTS`)
);


