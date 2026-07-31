create table configs (
	id serial primary key,
	posting_enabled boolean default true
);

insert into configs default values;