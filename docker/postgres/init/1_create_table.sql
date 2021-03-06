CREATE TABLE SLACK_USER (
    id  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    slack_id VARCHAR(20) NOT NULL,
    slack_channel_id VARCHAR(20) NOT NULL
);

CREATE TABLE KEYWORD (
    id  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    content VARCHAR(20) NOT NULL
);

CREATE TABLE USER_KEYWORD (
    id  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    slack_user_id BIGINT,
    keyword_id BIGINT,
    foreign key (slack_user_id) references SLACK_USER(id),
    foreign key (keyword_id) references KEYWORD(id)
);
