alter table boards
add max_threads	integer	default 100
;

alter table threads
add archived bool default false,
add archived_at timestamp default now()
;
