
CREATE TABLE
IF NOT EXISTS urls
(
  id           SERIAL PRIMARY KEY,
  short_code   VARCHAR
(255) NOT NULL UNIQUE,
  original_url TEXT NOT NULL,
  click_count  INTEGER DEFAULT 0,
  created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX
IF NOT EXISTS idx_urls_short_code ON urls
(short_code);
CREATE INDEX
IF NOT EXISTS idx_urls_created_at ON urls
(created_at);
