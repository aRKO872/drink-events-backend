CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE user_types AS ENUM ('user', 'admin', 'pub_representative');

CREATE TABLE files (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  file_name VARCHAR NOT NULL,
  etag VARCHAR NOT NULL,
  secure_url VARCHAR NOT NULL
);

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR NOT NULL,
  user_type user_types NOT NULL,
  email VARCHAR UNIQUE NOT NULL,
  phone VARCHAR UNIQUE,
  bio VARCHAR,
  profile_picture UUID REFERENCES files(id),
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
  latitude DECIMAL(10, 6),
  longitude DECIMAL(10, 6),
  email_last_changed TIMESTAMP,
  phone_last_changed TIMESTAMP,
  search_radius INT DEFAULT 5000 NOT NULL
);

CREATE TABLE blocked_users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id),
  blocked_user_id UUID NOT NULL REFERENCES users(id),
  time_of_blocking TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE TABLE drinks (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  drink_name VARCHAR NOT NULL,
  drink_pic UUID REFERENCES files(id)
);

CREATE TABLE user_drinks_mappings (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  drink_id UUID NOT NULL REFERENCES drinks(id),
  user_id UUID NOT NULL REFERENCES users(id)
);
