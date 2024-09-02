-- +goose Up
-- +goose StatementBegin

create table delivery
(
    id      serial primary key,
    name    varchar(255) not null,
    phone   varchar(255) not null,
    zip     varchar(255) not null,
    city    varchar(255) not null,
    address varchar(255) not null,
    region  varchar(255) not null,
    email   varchar(255) not null
);
create table payment
(
    id            serial primary key,
    transaction   varchar(255) not null,
    request_id    varchar(255) not null,
    currency      varchar(255) not null,
    provider      varchar(255) not null,
    amount        numeric      not null,
    payment_dt    date         not null,
    bank          varchar(255) not null,
    delivery_cost numeric      not null,
    goods_total   numeric      not null,
    custom_fee    numeric      not null
);
create table items
(
    id           serial primary key,
    chrt_id      numeric      not null,
    track_number varchar(255) not null,
    price        numeric      not null,
    rid          varchar(255) not null,
    name         varchar(255) not null,
    sale         numeric      not null,
    size         varchar(255) not null,
    total_price  numeric      not null,
    nm_id        numeric      not null,
    brand        varchar(255) not null,
    status       numeric      not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE delivery;
DROP TABLE payment;
DROP TABLE items;
-- +goose StatementEnd