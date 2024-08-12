DROP INDEX IF EXISTS idx_friend_messages_created_at;
DROP INDEX IF EXISTS idx_friend_messages_sent_by;
DROP INDEX IF EXISTS idx_friend_messages_friends_chat_id;

DROP TABLE IF EXISTS friend_messages;

DROP INDEX IF EXISTS idx_friends_friend2_id;
DROP INDEX IF EXISTS idx_friends_friend1_id;
DROP INDEX IF EXISTS idx_unique_friends_pair;

DROP TABLE IF EXISTS friends;

DROP INDEX IF EXISTS idx_pair_requests_receiver_id;
DROP INDEX IF EXISTS idx_pair_requests_sender_id;
DROP INDEX IF EXISTS idx_unique_sender_receiver_pair;

DROP TABLE IF EXISTS pair_requests;
