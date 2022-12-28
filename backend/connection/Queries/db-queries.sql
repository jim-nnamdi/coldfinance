CREATE TABLE IF NOT EXISTS users(id int PRIMARY KEY AUTO_INCREMENT, username TEXT, password TEXT, email TEXT, location TEXT, verified INT);

CREATE TABLE IF NOT EXISTS posts(id int PRIMARY KEY AUTO_INCREMENT, title TEXT, body LONGTEXT, slug TEXT, author TEXT, image TEXT);

ALTER TABLE posts ADD approved int(0);