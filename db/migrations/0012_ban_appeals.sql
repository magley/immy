create type ban_appeal_status as enum (
	'pending',
	'rejected',
	'rejected_final',
	'approved'
);

create table ban_appeals (
	id			serial primary key,
	ban_id		integer references bans(id),
	message		varchar,
	status		ban_appeal_status default 'pending',
	created_at	timestamp default now(),
	deleted_at	timestamp default null,
	reviewed_at timestamp default null,
	reviewed_by integer references users(id)
);