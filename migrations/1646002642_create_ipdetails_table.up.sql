create table ip_details
(
    id            serial PRIMARY KEY,
    ip_address    VARCHAR NOT NULL,
    country_code  VARCHAR NOT NULL,
    country       VARCHAR NOT NULL,
    city          VARCHAR NOT NULL,
    latitude      FLOAT   NOT NULL,
    longitude     FLOAT   NOT NULL,
    mystery_value BIGINT  NOT NULL,
    created_at    TIMESTAMP,
    updated_at    TIMESTAMP
);
CREATE UNIQUE INDEX index_name on ip_details (ip_address);
