-- Пользователи (Users)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255), 
    login VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE
);

-- Услуги (Languages)
CREATE TABLE langs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    short_description VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    img_link VARCHAR(255),
    author VARCHAR(255),
    year CHAR(4),
    version VARCHAR(50),
    list TEXT,
    status BOOLEAN NOT NULL DEFAULT TRUE
);

-- Заявки (Projects)
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    creation_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deletion_time TIMESTAMP,
    completion_time TIMESTAMP,
    status INT NOT NULL, 
    moderator_id INT REFERENCES users(id) ON DELETE SET NULL,
    count INT DEFAULT 0
);

-- Файлы (Files)
CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    lang_id INT REFERENCES langs(id) ON DELETE CASCADE NOT NULL,
    project_id INT REFERENCES projects(id) ON DELETE CASCADE NOT NULL,
    code TEXT,
    autocheck INT DEFAULT 0
);
