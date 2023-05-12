package person

var EntityTableDDL = `
CREATE TABLE person (
id varchar(30) primary key,
lastname  varchar(30) not null,
nickname  varchar(30),
age int,
consensus bool,
creation_tm timestamp
);
`
var EntityTableDropDDL = `DROP TABLE person`
