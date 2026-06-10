alter table boards
	add deleted_at timestamp default null
;

alter table threads
	add deleted_at timestamp default null
;

alter table posts
	add deleted_at timestamp default null
;