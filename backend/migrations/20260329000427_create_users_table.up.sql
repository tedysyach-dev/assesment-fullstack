CREATE TABLE
    IF NOT EXISTS users (
        id VARCHAR(36) PRIMARY KEY,
        email VARCHAR(100) NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        CONSTRAINT uq_users_email UNIQUE (email)
    )