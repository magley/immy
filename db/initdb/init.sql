create database db;
	
create table boards (
	id 			serial primary key,
	name 		varchar(255) unique not null,
	code 		varchar(16) unique not null,
	description varchar(255),
	created_at 	timestamp default now(),
	deleted_at	timestamp default null,
	-- ==================================== --
	locked 				bool 		default false,
	hidden 				bool 		default false,
	max_file_size 		integer 	default 2097152, -- 2MB
	reply_files_allowed bool 		default true,
	mime_types_allowed	varchar array,
	bump_limit			integer		default 250,
	image_limit			integer		default 150,
	flags_enabled		bool		default false,
	ids_enabled			bool		default false,
	code_enabled		bool		default false,
	math_enabled		bool		default false,
	max_threads			integer		default 100,
	allow_spoilers		bool		default false,
	spoiler_image		varchar		default null,
	-- ==================================== --
	post_count 		integer default 0,
	bytes_uploaded 	integer default 0
);

create table threads (
	id  		serial primary key,
	deleted_at	timestamp default null,
	board_id 	integer references boards(id),
	post_num 	integer,  -- Redundant field to simplify search and avoid joins.
	subject		varchar(128),
	locked 		bool default false,
	sticky 		bool default false,
	archived 	bool default false,
	archived_at timestamp default now(),
	auto_cycle 	integer default 0
);

create table posts (
	id 			serial primary key,
	thread_id 	integer references threads(id) on delete cascade,
	thread_num  integer,						-- Redundant field to avoid joins
	board_id    integer references boards(id),  -- Redundant field to avoid joins
	num 		integer,
	ipv4 		varchar(16),
	user_id		integer references users(id) on delete cascade,
	user_role	varchar default null,			-- Redundant field to avoid joins
	public_id	varchar(7),
	name 		varchar(128) default 'Anonymous',
	tripcode 	varchar(128),
	created_at 	timestamp default now(),
	deleted_at	timestamp default null,
	sage 		bool default false,
	capcode 	bool default false,
	content		varchar,
	filename 	varchar,
	filesize	integer,
	img_width	integer,
	img_height	integer,
	md5 		varchar,						-- Stored as base64
	src_filename varchar,
	spoiler		bool default false,
	html 		varchar
);

create type user_role as enum (
	'admin',
	'moderator',
	'janitor'
);

create table users (
	id 			serial primary key,
	username 	varchar(32) unique not null,
	password 	varchar(255) not null,
	role 		user_role not null,
	created_at 	timestamp default now()
);

create table bans (
	id			serial primary key,
	ipstart		bigint not null,
	ipend		bigint default null,
	created_at	timestamp default now(),
	expires		timestamp default null,			-- If null, ban is permanent
	board_id	integer references boards(id), 	-- If null, banned from all boards
	creator_id	integer references users(id),
	reason		varchar,
	warning		bool default false,
	seen		bool default false
);