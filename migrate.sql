CREATE TABLE era
(
  id INTEGER NOT NULL PRIMARY KEY,
  title TEXT NOT NULL,
  min_year SMALLINT,
  max_year SMALLINT
);

CREATE TABLE size
(
  id INTEGER NOT NULL PRIMARY KEY,
  title TEXT NOT NULL,
  min_pages SMALLINT,
  max_pages SMALLINT
);

CREATE TABLE genre
(
  id INTEGER NOT NULL PRIMARY KEY,
  title TEXT NOT NULL
); 

CREATE TABLE author
(
  id INTEGER NOT NULL PRIMARY KEY,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL
);

CREATE TABLE book
(
  id INTEGER NOT NULL PRIMARY KEY,
  title TEXT NOT NULL,
  year_published SMALLINT NOT NULL,
  rating NUMERIC(3, 2) NOT NULL,
  pages SMALLINT NOT NULL,
  genre_id INTEGER REFERENCES genre(id),
  author_id INTEGER REFERENCES author(id)
);

CREATE INDEX book_published ON book USING btree (year_published);
CREATE INDEX book_rating ON book USING btree (rating);
CREATE INDEX book_pages ON book USING btree (pages);
CREATE INDEX book_genre_id ON book USING btree (genre_id);
CREATE INDEX book_author_id ON book USING btree (author_id);

INSERT INTO era (id, title, min_year, max_year)
VALUES
  (0, 'Any', NULL, NULL),
  (1, 'Classic', NULL, 1969),
  (2, 'Modern', 1970, NULL);

INSERT INTO size (id, title, min_pages, max_pages)
VALUES
  (0, 'Any', NULL, NULL),
  (1, 'Short story – up to 35 pages', NULL, 34),
  (2, 'Novelette – 35 to 85 pages', 35, 84),
  (3, 'Novella – 85 to 200 pages', 85, 199),
  (4, 'Novel – 200 to 500 pages', 200, 499),
  (5, 'Brick – 500 to 800 pages', 500, 799),
  (6, 'Monument – 800 pages and up', 800, NULL);

INSERT INTO genre (id, title)
VALUES
  (1, 'Young Adult'),
  (2, 'SciFi/Fantasy'),
  (3, 'Romance'),
  (4, 'Nonfiction'),
  (5, 'Mystery'),
  (6, 'Memoir'),
  (7, 'Fiction'),
  (8, 'Childrens');

INSERT INTO author (id, first_name, last_name)
VALUES
  (1, 'Wendell', 'Stackhouse'),
  (2, 'Amelia', 'Wangerin, Jr.'),
  (3, 'Anastasia', 'Inez'),
  (4, 'Arthur', 'McCrumb'),
  (5, 'Arturo', 'Hijuelos'),
  (6, 'Bernard', 'Hopf'),
  (7, 'Bianca', 'Thompson'),
  (8, 'Bravig', 'Lewisohn'),
  (9, 'Burton', 'Malamud'),
  (10, 'Carolyn', 'Segal'),
  (11, 'Charles', 'Fenimore'),
  (12, 'Clifford', 'Wolitzer'),
  (13, 'Darryl', 'Fleischman'),
  (14, 'David', 'Beam'),
  (15, 'Elizabeth', 'Herbach'),
  (16, 'Elmer', 'Komroff'),
  (17, 'Gloria', 'Green'),
  (18, 'Grace', 'Harrison'),
  (19, 'Hamlin', 'Myrer'),
  (20, 'Hillary', 'Barnhardt'),
  (21, 'Jill', 'Hergesheimer'),
  (22, 'John W.', 'Spanogle'),
  (23, 'Jonathan', 'Kotzwinkle'),
  (24, 'Kathy', 'Yglesias'),
  (25, 'Kenneth', 'Douglas'),
  (26, 'Kris', 'Elegant'),
  (27, 'Langston', 'Lippman'),
  (28, 'Leonard', 'Nabokov'),
  (29, 'Lori', 'Kaan'),
  (30, 'Lynne', 'Danticat'),
  (31, 'Malin', 'Wolff'),
  (32, 'Oliver', 'Lowry'),
  (33, 'Patricia', 'Hazzard'),
  (34, 'Philip', 'Antrim'),
  (35, 'Phoebe', 'Brown'),
  (36, 'R.M.', 'Larner'),
  (37, 'Robert', 'Plimpton'),
  (38, 'Robert', 'Milofsky'),
  (39, 'Ursula', 'Karénine'),
  (40, 'Ward', 'Haigh'),
  (41, 'Abraham', 'Barton');

INSERT INTO book (id, title, year_published, rating, pages, genre_id, author_id)
VALUES
  (1, 'Alanna Saves the Day', 1972, 1.62, 169, 8, 6),
  (2, 'Adventures of Kaya', 1999, 2.13, 619, 1, 40),
  (3, 'A Horrible Human with the Habits of a Monster', 1976, 1.14, 258, 7, 25),
  (4, 'And I Said Yes', 1954, 3.3, 183, 7, 16),
  (5, 'Ballinby Boys', 1960, 1.88, 205, 2, 4),
  (6, 'Banana Slug and the Lost Cow', 1983, 2.53, 527, 8, 20),
  (7, 'Banana Slug and Xyr Friends', 1989, 3.64, 558, 8, 20),
  (8, 'Banana Slug and the Glass Half Full', 1952, 4.51, 796, 8, 17),
  (9, 'Banana Slug and the Mossy Rock', 2006, 4.43, 70, 8, 31),
  (10, 'Burnished Silver', 1932, 1.2, 202, 3, 30),
  (11, 'Cimornul', 1942, 1.08, 791, 2, 21),
  (12, 'Can I Be Honest?', 2007, 4.77, 542, 1, 11),
  (13, 'Concerning Prophecy', 1944, 3.8, 155, 2, 18),
  (14, 'Don''t Check your Ego', 1993, 3.02, 100, 4, 36),
  (15, 'The Deep Grey', 1931, 3.94, 43, 7, 37),
  (16, 'Dust on the Rim', 1946, 4.24, 38, 2, 24),
  (17, 'Did You Hear?', 1954, 2.48, 887, 7, 30),
  (18, 'Heliotrope Pajamas', 1952, 3.74, 16, 8, 31),
  (19, 'Hashtag QuokkaSelfie', 1995, 3.42, 417, 4, 27),
  (20, 'Interrobangs for All', 2011, 3.37, 677, 7, 16),
  (21, 'Inconvenient Confessions: a 6', 1972, 4.11, 766, 6, 32),
  (22, 'It''s Never Just a Glass', 1956, 3.55, 305, 1, 28),
  (23, 'Kalakalal Avenue', 2016, 4.27, 26, 7, 16),
  (24, 'Lace and Brandy', 1967, 4.13, 158, 3, 30),
  (25, 'Land Water Sky Space', 1983, 1.64, 320, 4, 15),
  (26, '(im)Mortality', 1985, 1.72, 214, 1, 12),
  (27, 'Muddy Waters', 2020, 4.76, 594, 3, 30),
  (28, 'Not to Gossip, But', 1958, 3.96, 537, 7, 17),
  (29, 'Nothing But Capers', 2004, 3.87, 347, 4, 1),
  (30, 'No More Lightning', 1978, 3.16, 99, 7, 11),
  (31, 'Natural Pamplemousse', 1957, 4.66, 886, 4, 35),
  (32, '9803 North Millworks Road', 1935, 4.76, 449, 5, 10),
  (33, 'Post Alley', 2014, 1.63, 374, 7, 9),
  (34, 'Portmeirion', 2020, 2.11, 277, 2, 7),
  (35, 'Quiddity and Quoddity', 2005, 2.42, 318, 1, 21),
  (36, 'Rystwyth', 1930, 1.6, 59, 2, 7),
  (37, 'Saint Esme', 1949, 1.84, 196, 3, 30),
  (38, 'Some Eggs or Something?', 1997, 3.24, 12, 7, 29),
  (39, 'Say it with Snap!', 1989, 3.77, 499, 4, 22),
  (40, 'Soft, Pliable Truth', 1933, 3.28, 453, 2, 38),
  (41, 'She Also Tottered', 2010, 2.09, 225, 2, 38),
  (42, 'The Spark and The Ashes', 2000, 2.71, 721, 1, 39),
  (43, 'Thatchwork Cottage', 1986, 2.43, 667, 7, 9),
  (44, 'Tales of the Compass', 1945, 4.22, 570, 2, 24),
  (45, 'The Elephant House', 1979, 3.95, 349, 4, 22),
  (46, 'The Winchcombe Railway Museum Heist', 2004, 3.04, 731, 5, 10),
  (47, 'The Startling End of Mr. Hidhoo', 1986, 1.59, 842, 7, 23),
  (48, 'The Thing Is', 1988, 2.83, 115, 7, 17),
  (49, 'The Mallemaroking', 1970, 1.95, 418, 2, 7),
  (50, 'The Scent of Oranges', 2006, 2.37, 264, 3, 30),
  (51, 'the life and times of an utterly inconsequential person', 1992, 1, 509, 7, 14),
  (52, 'The Seawitch Sings', 1977, 4.62, 90, 3, 30),
  (53, 'Turn Left Til You Get There', 1985, 4.54, 331, 7, 26),
  (54, 'The Triscanipt', 2018, 2.26, 16, 2, 39),
  (55, 'Whither Thou Goest', 1963, 4.44, 146, 3, 30),
  (56, 'Who Did You Think You Were Kidding?', 1986, 4.6, 867, 6, 34),
  (57, 'We''re Sisters and We Kinda Like Each Other', 1989, 4.71, 67, 6, 33),
  (58, 'Zero over Twelve', 1981, 1.01, 287, 5, 9);
