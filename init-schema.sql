--If the conditional table below is changed for any reason, delete the table once for the updated conditions to be used

-- DROP TABLE IF EXISTS events;
-- DROP TABLE IF EXISTS rsvp;

CREATE TABLE IF NOT EXISTS events (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL
        CHECK (char_length(title) < 51 AND char_length(title)>5 ),
    location TEXT NOT NULL
        CHECK (char_length(location) < 51 AND char_length(location)>5),
    image TEXT NOT NULL
        CHECK (image ~ '^https?://' AND image ~ '\.(png|jpg|jpeg|gif|gifv)$' AND char_length(image) < 51 AND char_length(image)>5),
    date TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS rsvp (
    event_id INT REFERENCES events(id),
    email_address TEXT NOT NULL
        CHECK (email_address ~ '^[a-zA-Z0-9_\-\.]+@yale\.edu$')
);

-- ALTER TABLE events DROP CONSTRAINT title ;


-- "Upsert" our default events. This will overwrite your
-- default events every time this code is run.
INSERT INTO events (id, date, image, title, location) VALUES
('1', '2022-10-17 16:30:00', 'http://i.imgur.com/pXjrQ.gif', 'SOM House Party', 'Kyle''s house'),
('3', '2022-10-19 19:00:00', 'http://i.imgur.com/7pe2k.gif', 'BBQ party for hackers and nerds', 'Sharon Oster''s house'),
('4', '2022-12-02 18:00:00', 'http://i.imgur.com/CJLrRqh.gif', 'BBQ for managers', 'Barry Nalebuff''s house'),
('5', '2022-12-21 19:00:00', 'http://i.imgur.com/02KT9.gif', 'Cooking lessons for the busy business student', 'Yale Farm')
ON CONFLICT ON CONSTRAINT events_pkey DO UPDATE
SET (id, date, image, title, location) = (EXCLUDED.id, EXCLUDED.date, EXCLUDED.image, EXCLUDED.title, EXCLUDED.location) ;

-- Notice that I did not insert an event with id 2. This is
-- just so you see that not all ids will be present.

-- Reset the id sequence because we did some manual insertions.
SELECT setval('events_id_seq', (SELECT MAX(id)+1 FROM events));
