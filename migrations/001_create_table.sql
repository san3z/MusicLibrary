CREATE TABLE IF NOT EXISTS public."musicLibrary" (
id SERIAL PRIMARY KEY AUTO INCREMENT,
group VARCHAR(20),
song VARCHAR(20),
release_date date,
album VARCHAR(20),
lyrics TEXT,
link VARCHAR(100)
);