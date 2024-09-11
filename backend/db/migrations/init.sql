CREATE TABLE employee
(
    id         SERIAL PRIMARY KEY,
    username   VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name  VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
    );

CREATE TABLE organization
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    type        organization_type,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organization_responsible
(
    id              SERIAL PRIMARY KEY,
    organization_id INT REFERENCES organization (id) ON DELETE CASCADE,
    user_id         INT REFERENCES employee (id) ON DELETE CASCADE
);

CREATE TYPE tender_status AS ENUM (
    'CREATED',
    'PUBLISHED',
    'CLOSED',
    'OPEN'
    );

CREATE TABLE tender
(
    id              SERIAL PRIMARY KEY,
    description     VARCHAR(255),
    organization_id INT references organization (id) ON DELETE cascade,
    employee_id     INT references employee (id) ON DELETE cascade,
    version_uid     SERIAL NOT null,
    prev_version_id INT references tender (id) ON DELETE cascade,
    status          tender_status                          NOT NULL
);

CREATE TYPE bid_status AS ENUM (
    'SUBMITTED'
    );

CREATE TABLE bid
(
    id              SERIAL PRIMARY KEY,
    description     VARCHAR(255),
    status          tender_status,
    tender_id       INT references tender (id) ON DELETE cascade,
    organization_id INT references organization (id) ON DELETE cascade,
    employ_id       INT references employee (id) ON DELETE cascade
);


-- Вставка моковых данных в таблицу employee
INSERT INTO employee (username, first_name, last_name)
VALUES ('jdoe', 'John', 'Doe'),
       ('asmith', 'Alice', 'Smith'),
       ('bwhite', 'Bob', 'White'),
       ('ljones', 'Lisa', 'Jones'),
       ('mjones', 'Mike', 'Jones');

-- Вставка моковых данных в таблицу organization
INSERT INTO organization (name, description, type)
VALUES ('Tech Corp', 'A leading tech company.', 'LLC'),
       ('Green Solutions', 'Sustainable energy solutions.', 'JSC'),
       ('Health Plus', 'Healthcare services and products.', 'IE'),
       ('Edu World', 'Educational resources and platforms.', 'LLC'),
       ('Finance Hub', 'Financial consulting and services.', 'JSC');

-- Вставка моковых данных в таблицу organization_responsible
INSERT INTO organization_responsible (organization_id, user_id)
VALUES (1, 1), -- John Doe is responsible for Tech Corp
       (2, 2), -- Alice Smith is responsible for Green Solutions
       (3, 3), -- Bob White is responsible for Health Plus
       (1, 4), -- Lisa Jones is also responsible for Tech Corp
       (4, 5); -- Mike Jones is responsible for Edu World