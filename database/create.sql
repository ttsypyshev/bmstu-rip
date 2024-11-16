-- Создание ENUM типа для статуса проекта
CREATE TYPE project_status AS ENUM ('draft', 'deleted', 'created', 'completed', 'rejected');

-- Пользователи (Users)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    email TEXT UNIQUE,
    login VARCHAR(50) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE
);

-- Услуги (Languages)
CREATE TABLE langs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    short_description VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    img_link VARCHAR(255),
    author VARCHAR(50),
    year CHAR(4),
    version VARCHAR(50),
    list JSONB,
    status BOOLEAN NOT NULL DEFAULT TRUE
);

-- Заявки (Projects)
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    creation_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    formation_time TIMESTAMP,
    completion_time TIMESTAMP,
    status project_status NOT NULL,  -- Использование ENUM типа
    moderator_id INT REFERENCES users(id) ON DELETE SET NULL,
    moderator_comment TEXT,
    count INT DEFAULT 0  -- Количество файлов в проекте
);

-- Файлы (Files)
CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    lang_id INT REFERENCES langs(id) ON DELETE CASCADE NOT NULL,
    project_id INT REFERENCES projects(id) ON DELETE CASCADE NOT NULL,
    code TEXT,
    file_name VARCHAR(255),
    file_size BIGINT DEFAULT 0,
    comment TEXT,
    auto_check INT DEFAULT 0
);
