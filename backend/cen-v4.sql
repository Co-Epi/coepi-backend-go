DROP TABLE IF EXISTS TCNReport;


CREATE TABLE `TCNReport` (
   `reportVK`	varchar(32) DEFAULT "",
   `report`     varchar(4000) DEFAULT "",
   `reportTS` bigint unsigned,
   PRIMARY KEY (`reportVK`,`reportTS`),
   KEY (`reportTS`)
);


