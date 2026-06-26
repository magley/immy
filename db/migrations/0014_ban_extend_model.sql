alter table bans
	add column board_code varchar,
	add column creator_username varchar;

update bans
set board_code = boards.code
from boards
where bans.board_id = boards.id;

update bans
set creator_username = users.username
from users
where bans.creator_id = users.id;

