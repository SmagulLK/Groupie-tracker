CREATE TABLE Users(
        userid serial not null,
        username varchar(100) not null unique ,
        user_password varchar(100) not null,
        email varchar(150) not null
);