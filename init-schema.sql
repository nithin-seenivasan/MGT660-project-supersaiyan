CREATE EXTENSION IF NOT EXISTS "pgcrypto";

defaultEvents := []Event{
                Event{
                        ID:        1,
                        Title:     "SOM House Party",
                        Date:      time.Date(2019, 10, 17, 16, 30, 0, 0, newYork),
                        Image:     "http://i.imgur.com/pXjrQ.gif",
                        Location:  "Kyle's house",
                        Attending: []string{"kyle.jensen@yale.edu", "kim.kardashian@yale.edu"},
                },
                Event{
                        ID:        2,
                        Title:     "BBQ party for hackers and nerds",
                        Date:      time.Date(2019, 10, 19, 19, 0, 0, 0, newYork),
                        Image:     "http://i.imgur.com/7pe2k.gif",
                        Location:  "Sharon Oster's house",
                        Attending: []string{"kyle.jensen@yale.edu", "kim.kardashian@yale.edu"},
                },
                Event{
                        ID:        3,
                        Title:     "BBQ for managers",
                        Date:      time.Date(2019, 12, 2, 18, 0, 0, 0, newYork),
                        Image:     "http://i.imgur.com/CJLrRqh.gif",
                        Location:  "Barry Nalebuff's house",
                        Attending: []string{"kim.kardashian@yale.edu"},
                },
                Event{
                        ID:        4,
                        Title:     "Cooking lessons for the busy business student",
                        Date:      time.Date(2019, 12, 21, 19, 0, 0, 0, newYork),
                        Image:     "http://i.imgur.com/02KT9.gif",
                        Location:  "Yale Farm",
                        Attending: []string{"homer.simpson@yale.edu"},
                },
        }

-- Create a new table if it doesn't exist
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

-- Alter if it already exists (make changes to the DB constraints here)
ALTER TABLE IF EXISTS events
    DROP CONSTRAINT IF EXISTS title,
    ADD  CONSTRAINT title CHECK (char_length(title) < 51 AND char_length(title)>5 ), 
    DROP CONSTRAINT IF EXISTS location,
    ADD  CONSTRAINT location CHECK (char_length(location) < 51 AND char_length(location)>5),
    DROP CONSTRAINT IF EXISTS image, 
    ADD  CONSTRAINT image CHECK (image ~ '^https?://' AND image ~ '\.(png|jpg|jpeg|gif|gifv)$' AND char_length(image) < 51 AND char_length(image)>5);


CREATE TABLE IF NOT EXISTS rsvp (
    event_id INT REFERENCES events(id),
    email_address TEXT NOT NULL
        CHECK (email_address ~ '^[a-zA-Z0-9_\-\.]+@yale\.edu$'),
    confirmation_code text GENERATED ALWAYS AS (encode(digest(email_address,'sha256')::bytea,'hex')) STORED,
    CONSTRAINT unique_rsvp PRIMARY KEY (event_id, email_address)
);

ALTER TABLE IF EXISTS rsvp
    DROP CONSTRAINT IF EXISTS email_address,
    ADD CONSTRAINT email_address CHECK (email_address ~ '^[a-zA-Z0-9_\-\.]+@yale\.edu$');
        
-- "Upsert" our default events. This will overwrite your
-- default events every time this code is run.




INSERT INTO events (id, date, image, title, location) VALUES
('1', '2022-10-17 16:30:00', 'http://i.imgur.com/pXjrQ.gif', 'SOM House Party', 'Kyle''s house'),
('2', '2022-10-19 19:00:00', 'http://i.imgur.com/7pe2k.gif', 'BBQ party for hackers and nerds', 'Sharon Oster''s house'),
('3', '2022-12-02 18:00:00', 'http://i.imgur.com/CJLrRqh.gif', 'BBQ for managers', 'Barry Nalebuff''s house'),
('4', '2022-12-21 19:00:00', 'http://i.imgur.com/02KT9.gif', 'Cooking lessons for the busy business student', 'Yale Farm')
ON CONFLICT ON CONSTRAINT events_pkey DO UPDATE
SET (id, date, image, title, location) = (EXCLUDED.id, EXCLUDED.date, EXCLUDED.image, EXCLUDED.title, EXCLUDED.location) ;

-- Notice that I did not insert an event with id 2. This is
-- just so you see that not all ids will be present.

-- Reset the id sequence because we did some manual insertions.
SELECT setval('events_id_seq', (SELECT MAX(id)+1 FROM events));

INSERT INTO rsvp (event_id, email_address) VALUES
('1', 'kim.kardashian@yale.edu'),
('2', 'kim.kardashian@yale.edu'),
('3', 'kim.kardashain@yale.edu'),
('4', 'kim.kardashian@yale.edu')