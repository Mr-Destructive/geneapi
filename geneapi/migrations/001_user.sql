CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY AUTOINCREMENT ,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(50),
    password VARBINARY(255) NOT NULL,
    api_key VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS llmapikeys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    openai VARCHAR(255),
    palm2 VARCHAR(255),
    anthropic VARCHAR(255),
    cohereai VARCHAR(255),
    huggingchat VARCHAR(255),
    user_id INTEGER,
    FOREIGN KEY (id) REFERENCES user(id) ON DELETE CASCADE
);


INSERT INTO user (id, username, email, password)
SELECT 1, 'test', 'test@test.com', 'test123'
WHERE NOT EXISTS (SELECT 1 FROM user WHERE id = 1);
