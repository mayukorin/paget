CREATE TABLE USER (
    id BIGINT PRIMARY KEY,
    slack_user_id VARCHAR(20) NOT NULL,
    slack_channel_id VARCHAR(20) NOT NULL
);

CREATE TABLE KEYWORD (
    id BIGINT PRIMARY KEY,
    content VARCHAR(20) NOT NULL
);

CREATE TABLE USER_KEYWORD (
    id BIGINT PRIMARY KEY,
    paget_user_id BIGINT,
    keyword_id BIGINT,
    foreign key (paget_user_id) references PAGET_USER(id),
    foreign key (keyword_id) references KEYWORD(id)
);
