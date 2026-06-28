create table rules_boards (
	rule_id 	integer references rules(id),
	board_id 	integer references boards(id),
	created_at	timestamp default now(),
	deleted_at	timestamp default null
);