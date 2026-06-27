create table rules (
	id 			serial primary key,
	title 		varchar not null,
	description varchar not null,
	is_global 	bool default false,
	danger 		integer default 1
);