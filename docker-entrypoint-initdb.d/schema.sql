CREATE TABLE IF NOT EXISTS customers
(
    id       BIGSERIAL PRIMARY KEY,
    name     TEXT      NOT NULL,
    phone    TEXT      NOT NULL UNIQUE,
    password TEXT      NOT NULL,
    active   BOOLEAN   NOT NULL DEFAULT TRUE,
    created  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS managers
(
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT      NOT NULL,
    salary      INTEGER   NOT NULL DEFAULT 0,
    plan        INTEGER   NOT NULL DEFAULT 0,
    BOSS_ID     BIGINT REFERENCES managers,
    department TEXT,
    phone       text      not null unique,
    password    TEXT,
    is_admin    BOOLEAN   NOT NULL DEFAULT TRUE,
    active      BOOLEAN   NOT NULL DEFAULT TRUE,
    created     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS customers_tokens
(
    token       TEXT      NOT NULL UNIQUE,
    customer_id BIGINT    NOT NULL REFERENCES customers,
    expire      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour',
    created     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS managers_tokens
(
    token      TEXT      NOT NULL UNIQUE,
    manager_id BIGINT    NOT NULL REFERENCES managers,
    expire     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour',
    created    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS products
(
    id      BIGSERIAL PRIMARY KEY,
    name    TEXT      NOT NULL,
    price   INTEGER   NOT NULL CHECK (price > 0),
    qty     INTEGER   NOT NULL DEFAULT 0 CHECK (qty >= 0),
    active  BOOLEAN   NOT NULL DEFAULT TRUE,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sales
(
    id          BIGSERIAL PRIMARY KEY,
    manager_id  BIGINT    NOT NULL REFERENCES managers,
    customer_id BIGINT    NOT NULL,
    created     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sales_positions
(
    id         BIGSERIAL PRIMARY KEY,
    product_id BIGINT    NOT NULL REFERENCES products,
    sale_id    BIGINT    NOT NULL REFERENCES sales,
    price      INTEGER   NOT NULL CHECK (price >= 0),
    qty        INTEGER   NOT NULL DEFAULT 0 CHECK (qty >= 0),
    created    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);