-- Drop tables in reverse order of creation to respect foreign key constraints
DROP TABLE IF EXISTS book_tags;
DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS genres;
DROP TABLE IF EXISTS users;
