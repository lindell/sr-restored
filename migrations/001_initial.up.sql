CREATE TABLE programs (
    id INT PRIMARY KEY,
    name TEXT,
    description TEXT,
    email TEXT,
    copyright TEXT,
    url TEXT,
    image_url TEXT
);

CREATE TABLE episodes (
    id INT PRIMARY KEY,
    program_id INT,
    title TEXT,
    description TEXT,
    url TEXT,
    publish_date timestamptz,
    image_url TEXT,
    file_url TEXT,
    file_duration INT,
    file_bytes INT
);

CREATE INDEX episodes_program_id ON episodes(program_id);
