CREATE TABLE SLACK_USER (
    id BIGINT PRIMARY KEY,
    slack_id VARCHAR(20) NOT NULL,
    slack_channel_id VARCHAR(20) NOT NULL
);

CREATE TABLE KEYWORD (
    id BIGINT PRIMARY KEY,
    content VARCHAR(20) NOT NULL
);

CREATE TABLE USER_KEYWORD (
    id BIGINT PRIMARY KEY,
    slack_user_id BIGINT,
    keyword_id BIGINT,
    foreign key (slack_user_id) references SLACK_USER(id),
    foreign key (keyword_id) references KEYWORD(id)
);
