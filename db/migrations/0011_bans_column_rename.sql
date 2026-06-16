alter table bans
	rename column ipstart to ip_start
;
alter table bans
	rename column ipend to ip_end
;
alter table bans
	add column deleted_at timestamp default null
;