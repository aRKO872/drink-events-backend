CREATE TABLE pair_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sender_id UUID NOT NULL REFERENCES users(id),
    receiver_id UUID NOT NULL REFERENCES users(id),
    is_active BOOLEAN DEFAULT TRUE,
    CONSTRAINT sender_receiver_not_equal CHECK (sender_id <> receiver_id)
);

CREATE UNIQUE INDEX idx_unique_sender_receiver_pair ON pair_requests (
    LEAST(sender_id, receiver_id),
    GREATEST(sender_id, receiver_id)
);

CREATE INDEX idx_pair_requests_sender_id ON pair_requests (sender_id);
CREATE INDEX idx_pair_requests_receiver_id ON pair_requests (receiver_id);

CREATE TABLE friends (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    friend1_id UUID NOT NULL REFERENCES users(id),
    friend2_id UUID NOT NULL REFERENCES users(id),
    is_active BOOLEAN DEFAULT TRUE,
    either_acc_deactivated BOOLEAN DEFAULT FALSE,
    CONSTRAINT friend_ids_not_equal CHECK (friend1_id <> friend2_id)
);

CREATE UNIQUE INDEX idx_unique_friends_pair ON friends (
    LEAST(friend1_id, friend2_id),
    GREATEST(friend1_id, friend2_id)
);

CREATE INDEX idx_friends_friend1_id ON friends (friend1_id);
CREATE INDEX idx_friends_friend2_id ON friends (friend2_id);

CREATE TABLE friend_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    friends_chat_id UUID NOT NULL REFERENCES friends(id),
    message VARCHAR NOT NULL,
    sent_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_friend_messages_friends_chat_id ON friend_messages (friends_chat_id);
CREATE INDEX idx_friend_messages_sent_by ON friend_messages (sent_by);
CREATE INDEX idx_friend_messages_created_at ON friend_messages (created_at);
