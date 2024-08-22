package literals

const (
	// Header Keys
	HEADER_USER_ID = "user_id"
	HEADER_USER_TYPE = "user_type"
	HEADER_USER_SEARCH_DISTANCE = "user_search_distance"

	// Redis Keys
	USER_LOCATION_KEY = "user_location"

	// Socket Message Type
	SOCKET_LOCATION_SET_TYPE = "socket_set_location_type"
	SOCKET_MESSAGE_TYPE = "socket_message_type"

	// REDIS TABLE KEYS 
	USER_INFO_REDIS_KEY = "%s_user_info"
	FILE_INFO_REDIS_KEY = "%s_file_info"
	PAIR_REQUEST_REDIS_KEY = "%s_%s_pair_request"
	FRIEND_RECORD_REDIS_KEY = "%s_%s_friend_record"

	// Context literals
	CTX_USER_ID = "USER_ID"

	//DATE Format
	DATE_FORMAT = "2006-01-02 15:04:05.999999"

	// Event Handlers
	SEND_PAIR_REQUEST = "send_pair_request"
	SEND_PAIR_REQUEST_SUCCESS = "receive_pair_request"
	NOTIFY_MESSAGE = "notify_msg"
	PAIR_REQUEST_ACCEPTED = "pair_request_accepted"

	// Worker Job Names
	ADD_PAIR_REQUEST_JOB = "add_pair_request"

	// Queue Namespace
	WORKER_QUEUE_NAMESPACE = "drink_events_namespace"
)