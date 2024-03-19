CREATE TYPE otp_event_types AS ENUM (
  'verify_email',
  'verify_phone',
  'login_email',
  'login_phone'
);

CREATE TABLE otps (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  number INT NOT NULL CHECK (number >= 100000 AND number <= 999999), 
  event otp_event_types NOT NULL,
  expiry TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);