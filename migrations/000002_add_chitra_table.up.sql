CREATE TABLE IF NOT EXISTS chitras (
    id SERIAL PRIMARY KEY,
    chitra_url varchar(200) NOT NULL,
    user_id INT NOT NULL,
    FOREIGN KEY(user_id) references users(id) ON DELETE CASCADE
)