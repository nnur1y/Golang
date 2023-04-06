create table comments (
	id INT,
	userid INT,
	recipeid INT,
	comment VARCHAR(8)
);
insert into comments (id, userid, recipeid, comment) values (1, 4, 6, 'in');
insert into comments (id, userid, recipeid, comment) values (2, 2, 3, 'in');
insert into comments (id, userid, recipeid, comment) values (3, 2, 2, 'in');
insert into comments (id, userid, recipeid, comment) values (4, 2, 4, 'in');
insert into comments (id, userid, recipeid, comment) values (5, 2, 3, 'comments');
insert into comments (id, userid, recipeid, comment) values (6, 4, 6, 'in');
insert into comments (id, userid, recipeid, comment) values (7, 4, 1, 'comments');
insert into comments (id, userid, recipeid, comment) values (8, 1, 7, 'list');
insert into comments (id, userid, recipeid, comment) values (9, 3, 2, 'comments');
insert into comments (id, userid, recipeid, comment) values (10, 4, 8, 'list');