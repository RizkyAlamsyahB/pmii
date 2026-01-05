CREATE TABLE inboxes (
    id SERIAL PRIMARY KEY,
    sender_id INT NOT NULL,
    receiver_id INT NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_sender FOREIGN KEY(sender_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_receiver FOREIGN KEY(receiver_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Index untuk mempercepat query pencarian pesan antar user dan inbox list
CREATE INDEX idx_inboxes_participants ON inboxes (sender_id, receiver_id);
CREATE INDEX idx_inboxes_created_at ON inboxes (created_at DESC);