CREATE TABLE post
(
    id serial primary key,
    author varchar(100),
    header varchar(100),
    content varchar(2000),
    comments_allowed boolean,
    created_at varchar(100) default to_char(now(), 'YYYY-MM-DDTHH24:MI:SS')
);  

CREATE TABLE comment
(
    id serial primary key,
    post int not null,
    FOREIGN KEY (post) REFERENCES post(id),
    author varchar(100),
    content varchar(2000),
    created_at varchar(100) default to_char(now(), 'YYYY-MM-DDTHH24:MI:SS'),
    reply_to int,
    FOREIGN KEY (reply_to) REFERENCES comment(id)
);
