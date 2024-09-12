
create table users (
    "id" serial primary key not null,
    "username" character varying(60) not null
);

create table musics(
    "id" serial primary key not null,
    "title" character varying(60) not null
);

create table ratings (
    "music_id" int not null,
    "user_id" int not null,
    "rating" numeric(2,1),
    primary key ("music_id", "user_id"),
    CONSTRAINT ratings_user_id_fk
        FOREIGN KEY ("user_id")
            REFERENCES users("id")
                ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT ratings_music_id_fk
        FOREIGN KEY ("music_id")
            REFERENCES musics("id")
                ON UPDATE CASCADE ON DELETE CASCADE
);


create table l_musics(
    "user_id" int not null,
    "music_ids" int[] default '{}'::int[],
    "similar_user_ids" int[] default '{}'::int[],
    primary key ("user_id"),
    CONSTRAINT l_musics_user_id_fk
        FOREIGN KEY ("user_id")
            REFERENCES users("id")
                ON UPDATE CASCADE ON DELETE CASCADE
);
