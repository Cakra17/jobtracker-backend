-- auth
CREATE TABLE IF NOT EXISTS users(
  id             VARCHAR(36)  PRIMARY KEY DEFAULT UUID(),
  email          VARCHAR(255) UNIQUE NOT NULL,
  username       VARCHAR(255) UNIQUE NOT NULL,
  display_name   VARCHAR(255),
  avatar_url     TEXT,
  email_verified BOOL DEFAULT FALSE,
  password_hash  VARCHAR(255) NOT NULL,
  created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
) ENGINE=INNODB;

CREATE TABLE IF NOT EXISTS oauth_accounts (
  id                  VARCHAR(36) PRIMARY KEY DEFAULT UUID(),
  user_id             VARCHAR(36) NOT NULL,
  `provider`          VARCHAR(50) NOT NULL,
  provider_account_id VARCHAR(255) NOT NULL,
  access_token        TEXT,
  refresh_token       TEXT,
  token_type          VARCHAR(50) DEFAULT 'Bearer',
  expires_at          TIMESTAMP,
  scope               TEXT,
  created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE KEY unique_provider_account (provider, provider_account_id)
) ENGINE=INNODB;