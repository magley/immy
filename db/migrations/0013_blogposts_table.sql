create table blogposts (
	id			serial primary key,
	title 		varchar not null,
	html		varchar not null,
	author_id	integer references users(id),
	author_name	varchar							-- Redundant
	created_at	timestamp default now(),
	deleted_at	timestamp default null
);