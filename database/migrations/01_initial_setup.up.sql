--- avatars url
CREATE TABLE IF NOT EXISTS avatars
(
    id SERIAL PRIMARY KEY,
    avatar_url TEXT NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE
);



--- parents details
CREATE TABLE IF NOT EXISTS users
(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    year_of_birth INT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    firebase_uid TEXT NOT NULL,
    avatar_id INT REFERENCES avatars (id),
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS user_firebase_uid_index on users(firebase_uid);



--- kids details
CREATE TABLE IF NOT EXISTS kids_profile
(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    age INT NOT NULL,
    avatar_id INT REFERENCES avatars (id),
    user_id INT REFERENCES users (id),
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS kid_parent_index on kids_profile(user_id);



--- story searched by user (parent)
CREATE TABLE IF NOT EXISTS story
(
    id SERIAL PRIMARY KEY,
    searched_text TEXT NOT NULL,
    user_id INT REFERENCES users (id),
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS searched_story_user_index on story(user_id);
