package models

type PairRequestWSEvent struct {
	To string `json:"to"`
}

type PairRequestOutputWSEvent struct {
	FromName string `json:"name"`
	ProfilePic string `json:"profile_pic"`
}

type PRAcceptOutputWSEvent struct {
	FriendRecordId string `json:"friend_record_id"`
	AccepterName string `json:"accepter_name"`
	ProfilePicAccepter string `json:"profile_pic"`
}

type NotifyMsgEvent struct {
	Msg string `json:"msg"`
}