create database db;
	
create table boards (
	id 			serial primary key,
	name 		varchar(255) unique not null,
	code 		varchar(16) unique not null,
	description varchar(255),
	created_at 	timestamp default now(),
	locked 		bool default false,
	hidden 		bool default false,
	post_count 	integer
);

create table threads (
	id  		serial primary key,			
	board_id 	integer references boards(id),
	post_num 	integer,  -- Redundant field to simplify search and avoid joins.
	subject		varchar(128),
	locked 		bool default false,
	sticky 		bool default false
);

create table posts (
	id 			serial primary key,
	thread_id 	integer references threads(id),
	thread_num  integer,						-- Redundant field to avoid joins
	board_id    integer references boards(id),  -- Redundant field to avoid joins
	num 		integer,
	ipv4 		varchar(16),
	name 		varchar(128) default 'Anonymous',
	tripcode 	varchar(128),
	created_at 	timestamp default now(),
	sage 		bool default false,
	content		varchar,
	filename 	varchar,
	src_filename varchar,
	html 		varchar
);

create type user_type as enum (
	'admin',
	'moderator',
	'janitor'
);

create table users (
	id 			serial primary key,
	username 	varchar(32) unique not null,
	password 	varchar(255) not null,
	type 		user_type not null,
	created_at 	timestamp default now()
);

create table users_boards (
	user_id 	integer references users(id) on update cascade on delete cascade,
	board_id 	integer references boards(id) on update cascade on delete cascade,
	constraint 	user_board_pkey primary key (user_id, board_id)
);