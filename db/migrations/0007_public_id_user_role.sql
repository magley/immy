alter table posts
	RENAME COLUMN user_id to public_id
;

alter table posts
	add user_id integer references users(id) on delete cascade,
	add user_role varchar default null
;

ALTER TYPE user_type
	RENAME TO user_role
;

alter table users
	RENAME COLUMN type to role
;