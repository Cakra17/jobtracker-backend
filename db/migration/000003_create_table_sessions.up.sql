CREATE TABLE IF NOT EXISTS `sessions` (
  id VARCHAR(36) NOT NULL,
  user_id VARCHAR(36) NOT NULL,
  refresh_token VARCHAR(255) NOT NULL,
  expire_at TIMESTAMP NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  INDEX idx_user_refresh_token (user_id, refresh_token)
) ENGINE = INNODB;