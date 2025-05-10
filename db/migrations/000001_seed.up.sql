-- Enable foreign key support in SQLite
PRAGMA foreign_keys = ON;

-- User table
CREATE TABLE users (
    id                TEXT PRIMARY KEY,
    email             TEXT NOT NULL UNIQUE,
    hashed_password   TEXT NOT NULL,
    created_at        DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Genre table
CREATE TABLE genres (
    code TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

-- Tag table
CREATE TABLE tags (
    code TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

-- Book table
CREATE TABLE books (
    id           TEXT PRIMARY KEY,
    title        TEXT NOT NULL,
    author       TEXT NOT NULL,
    genre_code   TEXT,
    release_date DATETIME,
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (genre_code) REFERENCES genres(code) ON DELETE SET NULL
);

-- Join table for many-to-many book-tags
CREATE TABLE book_tags (
    book_id TEXT NOT NULL,
    tag_code TEXT NOT NULL,
    PRIMARY KEY (book_id, tag_code),
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_code) REFERENCES tags(code) ON DELETE CASCADE
);

-- Seed data for genres
INSERT INTO genres (code, name) VALUES
    ('SCI_FI', 'Science Fiction'),
    ('FANTASY', 'Fantasy'),
    ('MYSTERY', 'Mystery'),
    ('ROMANCE', 'Romance'),
    ('THRILLER', 'Thriller'),
    ('NON_FICTION', 'Non-Fiction');

-- Seed data for tags
INSERT INTO tags (code, name) VALUES
    ('BESTSELLER', 'Bestseller'),
    ('AWARD_WINNER', 'Award Winner'),
    ('NEW_RELEASE', 'New Release'),
    ('CLASSIC', 'Classic'),
    ('RECOMMENDED', 'Recommended');

-- Seed data for books
INSERT INTO books (id, title, author, genre_code, release_date) VALUES
    ('1e7d0448-77d6-40a8-bfd6-0a4bea231171', 'Dune', 'Frank Herbert', 'SCI_FI', '1965-08-01'),
    ('87540847-7a97-4ffa-84f1-a1e652ee1385', 'The Hobbit', 'J.R.R. Tolkien', 'FANTASY', '1937-09-21'),
    ('8e572764-d35c-47be-8f86-d1e627093e86', 'Murder on the Orient Express', 'Agatha Christie', 'MYSTERY', '1934-01-01'),
    ('5b424281-6c31-420a-9383-9747b4629519', 'Pride and Prejudice', 'Jane Austen', 'ROMANCE', '1813-01-28'),
    ('266a1fdf-8abb-4116-95b6-f5dc207ced08', 'The Da Vinci Code', 'Dan Brown', 'THRILLER', '2003-03-18'),
    ('43242307-285a-4acf-a94e-dd1c9ea357d2', 'Sapiens: A Brief History of Humankind', 'Yuval Noah Harari', 'NON_FICTION', '2011-01-01'),
    ('635d66fa-4859-438d-8512-eb2d26f91dc8', 'Foundation', 'Isaac Asimov', 'SCI_FI', '1951-05-01'),
    ('0e426b67-aa38-488f-8fa6-6ec3152b3a69', 'The Lord of the Rings', 'J.R.R. Tolkien', 'FANTASY', '1954-07-29'),
    ('d0e2009a-8fc0-4497-acb3-ad9d7d61130b', 'The Silent Patient', 'Alex Michaelides', 'THRILLER', '2019-02-05'),
    ('1066eefb-8b13-4fe5-8a35-6c65cea3cfa5', 'Becoming', 'Michelle Obama', 'NON_FICTION', '2018-11-13');

-- Seed data for book_tags
INSERT INTO book_tags (book_id, tag_code) VALUES
    ('1e7d0448-77d6-40a8-bfd6-0a4bea231171', 'CLASSIC'),
    ('1e7d0448-77d6-40a8-bfd6-0a4bea231171', 'AWARD_WINNER'),
    ('87540847-7a97-4ffa-84f1-a1e652ee1385', 'CLASSIC'),
    ('87540847-7a97-4ffa-84f1-a1e652ee1385', 'BESTSELLER'),
    ('8e572764-d35c-47be-8f86-d1e627093e86', 'CLASSIC'),
    ('5b424281-6c31-420a-9383-9747b4629519', 'CLASSIC'),
    ('266a1fdf-8abb-4116-95b6-f5dc207ced08', 'BESTSELLER'),
    ('43242307-285a-4acf-a94e-dd1c9ea357d2', 'BESTSELLER'),
    ('43242307-285a-4acf-a94e-dd1c9ea357d2', 'RECOMMENDED'),
    ('635d66fa-4859-438d-8512-eb2d26f91dc8', 'CLASSIC'),
    ('0e426b67-aa38-488f-8fa6-6ec3152b3a69', 'CLASSIC'),
    ('0e426b67-aa38-488f-8fa6-6ec3152b3a69', 'BESTSELLER'),
    ('d0e2009a-8fc0-4497-acb3-ad9d7d61130b', 'NEW_RELEASE'),
    ('1066eefb-8b13-4fe5-8a35-6c65cea3cfa5', 'BESTSELLER'),
    ('1066eefb-8b13-4fe5-8a35-6c65cea3cfa5', 'RECOMMENDED');
