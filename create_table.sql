CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY,
    username VARCHAR(30) UNIQUE NOT NULL,
    email VARCHAR(30) UNIQUE NOT NULL,
    password text NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS items (
    item_id UUID PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    quantity INT NOT NULL,
    price INT NOT NULL,
    description text,
    owner VARCHAR(30) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    user_id UUID,
    CONSTRAINT fk_users
        FOREIGN KEY (user_id)
        REFERENCES "users" (user_id)
        ON DELETE SET NULL
);
CREATE TABLE IF NOT EXISTS posts (
    post_id UUID PRIMARY KEY,
    content TEXT NOT NULL,
    owner varchar(50) not null,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    user_id UUID,
    CONSTRAINT fk_users
        FOREIGN KEY(user_id)
        REFERENCES "users" (user_id)
        ON DELETE SET NULL
);