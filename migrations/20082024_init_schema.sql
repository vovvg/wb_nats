-- +goose Up
-- +goose StatementBegin

-- Создание таблицы delivery перед request
CREATE TABLE delivery
(
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(255) NOT NULL,
    phone   VARCHAR(255) NOT NULL,
    zip     VARCHAR(255) NOT NULL,
    city    VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    region  VARCHAR(255) NOT NULL,
    email   VARCHAR(255) NOT NULL
);

-- Создание таблицы payment перед request
CREATE TABLE payment
(
    id            SERIAL PRIMARY KEY,
    transaction   VARCHAR(255) NOT NULL,
    request_id    VARCHAR(255) NOT NULL,
    currency      VARCHAR(255) NOT NULL,
    provider      VARCHAR(255) NOT NULL,
    amount        NUMERIC      NOT NULL,
    payment_dt    DATE         NOT NULL,
    bank          VARCHAR(255) NOT NULL,
    delivery_cost NUMERIC      NOT NULL,
    goods_total   NUMERIC      NOT NULL,
    custom_fee    NUMERIC      NOT NULL
);

-- Создание таблицы request после зависимых таблиц
CREATE TABLE request
(
    id                 SERIAL PRIMARY KEY,
    order_uid          VARCHAR(255) NOT NULL UNIQUE,
    track_number       VARCHAR(255) NOT NULL,
    entry              VARCHAR(255) NOT NULL,
    locale             VARCHAR(50)  NOT NULL,
    internal_signature VARCHAR(255),
    customer_id        VARCHAR(255) NOT NULL,
    delivery_service   VARCHAR(255) NOT NULL,
    shard_key          VARCHAR(50)  NOT NULL,
    sm_id              INT          NOT NULL,
    date_created       TIMESTAMP    NOT NULL,
    oof_shard          VARCHAR(50)  NOT NULL,
    delivery_id        INT          NOT NULL,
    payment_id         INT          NOT NULL,
    CONSTRAINT fk_delivery FOREIGN KEY (delivery_id) REFERENCES delivery (id) ON DELETE CASCADE,
    CONSTRAINT fk_payment FOREIGN KEY (payment_id) REFERENCES payment (id) ON DELETE CASCADE
);

-- Создание таблицы items
CREATE TABLE items
(
    id           SERIAL PRIMARY KEY,
    chrt_id      NUMERIC      NOT NULL,
    track_number VARCHAR(255) NOT NULL,
    price        NUMERIC      NOT NULL,
    rid          VARCHAR(255) NOT NULL,
    name         VARCHAR(255) NOT NULL,
    sale         NUMERIC      NOT NULL,
    size         VARCHAR(50)  NOT NULL,
    total_price  NUMERIC      NOT NULL,
    nm_id        NUMERIC      NOT NULL,
    brand        VARCHAR(255) NOT NULL,
    status       NUMERIC      NOT NULL,
    request_id   INT          NOT NULL,
    CONSTRAINT fk_request FOREIGN KEY (request_id) REFERENCES request (id) ON DELETE CASCADE
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE request;
DROP TABLE delivery;
DROP TABLE payment;
DROP TABLE items;
-- +goose StatementEnd