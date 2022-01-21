DROP TABLE IF EXISTS skill;
DROP TABLE IF EXISTS city;
DROP TABLE IF EXISTS schedule;
DROP TABLE IF EXISTS schedule_item;
DROP TABLE IF EXISTS profile;
DROP TABLE IF EXISTS profile_skill;
DROP TABLE IF EXISTS connections;
DROP TABLE IF EXISTS incident;


CREATE TABLE skill
(
    id   INT not null auto_increment,
    NAME VARCHAR(255),
    PRIMARY KEY (id)
);
CREATE TABLE city
(
    CODE     VARCHAR(255),
    NAME     VARCHAR(255),
    UF       VARCHAR(2),
    ZIP_CODE VARCHAR(255)
);
CREATE TABLE profile
(
    id              INT not null auto_increment,
    city_code       VARCHAR(255),
    description     VARCHAR(255),
    is_instructor   BIT(1),
    username        VARCHAR(255),
    password        VARCHAR(255),
    first_name      VARCHAR(255),
    last_name       VARCHAR(255),
    profile_picture VARCHAR(255),
    fcm_token       VARCHAR(255),
    PRIMARY KEY (id)
);
CREATE TABLE profile_skill
(
    skill_id INT,
    user_id  INT
);
create table schedule
(
    id            INT not null auto_increment,
    instructor_id INT,
    student_id    INT,
    skill_id      int,
    accepted      BIT(1),
    updated       TIMESTAMP,
    PRIMARY KEY (id)
);
create table schedule_item
(
    id          INT not null auto_increment,
    week_day    INT,
    duration    INT,
    hour        INT,
    minutes     INT,
    schedule_id INT,
    PRIMARY KEY (id)
);
create table connections
(
    owner_id        int not null,
    contact_id      int not null,
    date_connection TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

create table incident
(
    id               INT      not null auto_increment,
    schedule_item_id INT      not null,
    requested_by     INT      not null,
    accepted         bool default false,
    answered         bool default false,
    created          datetime not null,
    day_of_change       date     not null,
    week_day         INT,
    duration         INT,
    hour             INT,
    minutes          INT,
    type             ENUM ('cancellation', 'change'),
    motive           VARCHAR(255),
    PRIMARY KEY (id)
);


INSERT INTO skill(id, name)
VALUES (1, 'Musculação');
INSERT INTO skill(id, name)
VALUES (2, 'Funcional');
INSERT INTO skill(id, name)
VALUES (3, 'Cross-fit');
INSERT INTO skill(id, name)
VALUES (4, 'Zumba');

INSERT INTO city(code, name, uf, zip_code)
values ('4314050', 'San Francisco', 'CA', '95630000');
INSERT INTO city(code, name, uf, zip_code)
values ('4314051', 'Los Angeles', 'CA', '95600000');
INSERT INTO city(code, name, uf, zip_code)
values ('4314052', 'Nova york', 'NY', '98888951');
INSERT INTO city(code, name, uf, zip_code)
values ('4314053', 'Porto Alegre', 'RS', '95790000');

INSERT INTO profile(id, city_code, description, is_instructor, username, password, first_name, last_name, fcm_token)
VALUES (1, '4314050', '', 1, 'teste@teste.com', '12345678', 'Joao ', 'Instrutor', '');
INSERT INTO profile(id, city_code, description, is_instructor, username, password, first_name, last_name, fcm_token)
VALUES (2, '4314050', '', 0, 'teste2@teste.com', '12345678', 'José', 'Aluno', '');

INSERT INTO profile_skill(skill_id, user_id)
values (1, 1);
INSERT INTO profile_skill(skill_id, user_id)
values (2, 1);

INSERT INTO schedule(id, instructor_id, student_id, skill_id, accepted, updated)
VALUES (1, 1, 2, 1, 1, now());
INSERT INTO schedule(id, instructor_id, student_id, skill_id, accepted, updated)
VALUES (2, 1, 2, 1, 0, now());

INSERT INTO schedule_item(id, week_day, duration, hour, minutes, schedule_id)
VALUES (1, 1, 60, 17, 30, 1);
INSERT INTO schedule_item(id, week_day, duration, hour, minutes, schedule_id)
VALUES (2, 2, 60, 18, 30, 1);
INSERT INTO schedule_item(id, week_day, duration, hour, minutes, schedule_id)
VALUES (3, 3, 60, 19, 30, 1);

INSERT INTO schedule(id, instructor_id, student_id, accepted, updated)
VALUES (3, 1, 2, 1, now());

INSERT INTO schedule_item(id, week_day, duration, hour, minutes, schedule_id)
VALUES (4, 4, 60, 17, 30, 2);
INSERT INTO schedule_item(id, week_day, duration, hour, minutes, schedule_id)
VALUES (5, 5, 60, 18, 30, 2);
INSERT INTO schedule_item(id, week_day, duration, hour, minutes, schedule_id)
VALUES (6, 6, 60, 19, 30, 2);