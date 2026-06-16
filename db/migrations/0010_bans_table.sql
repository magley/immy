create table bans (
	id			serial primary key,
	ipstart		bigint not null,
	ipend		bigint default null,
	created_at	timestamp default now(),
	expires_at	timestamp default null,
	board_id	integer references boards(id),
	creator_id	integer references users(id),
	reason		varchar,
	warning		bool default false,
	seen		bool default false
);