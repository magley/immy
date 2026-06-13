alter table posts
	add spoiler bool default false
;

alter table boards
	add	allow_spoilers		bool		default false,
	add spoiler_image		varchar		default null
;