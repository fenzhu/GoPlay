create database wiki;
use wiki;

drop table if exists article;
create table article (
	title varchar(128) not null,
	body text not null,
	primary key(title)
) engine = InnoDB;

insert into article (title, body) values ("test", "test body");