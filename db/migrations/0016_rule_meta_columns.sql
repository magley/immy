alter table rules
	add column created_at timestamp default now(),
	add column deleted_at timestamp default null
;