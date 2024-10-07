# create the database 

create database library;

use library;

create table book(
    id int primary key,
    name  varchar(100) not null,
    isbn varchar(100) not null
)

insert  into book (id,name, isbn) values (1,'Java', 'Everythinbg about Java');

psql -U your_username -h localhost -d your_password


psql -U your_username -h localhost -d library

