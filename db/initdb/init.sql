create database db;
	
create table boards (
	id 			serial primary key,
	name 		varchar(255) unique not null,
	code 		varchar(16) unique not null,
	description varchar(255),
	created_at 	timestamp default now(),
	locked 		bool default false,
	hidden 		bool default false
);

create table threads (
	id 			serial primary key,
	board_id 	integer references boards(id),
	locked 		bool default false,
	sticky 		bool default false
);

create table posts (
	id 			serial primary key,
	thread_id 	integer references threads(id),
	num 		integer,
	ipv4 		varchar(16),
	user_name 	varchar(128) default 'Anonymous',
	tripcode 	varchar(128),
	created_at 	timestamp default now(),
	sage 		bool default false,
	post 		varchar,
	filename 	varchar,
	html 		varchar
);

