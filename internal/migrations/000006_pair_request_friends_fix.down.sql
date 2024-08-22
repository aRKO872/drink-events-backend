DROP INDEX IF EXISTS idx_unique_sender_receiver_pair;
DROP INDEX IF EXISTS idx_unique_friends_pair;

CREATE UNIQUE INDEX idx_unique_sender_receiver_pair ON pair_requests (
    LEAST(sender_id, receiver_id),
    GREATEST(sender_id, receiver_id)
);

CREATE UNIQUE INDEX idx_unique_friends_pair ON friends (
    LEAST(friend1_id, friend2_id),
    GREATEST(friend1_id, friend2_id)
);