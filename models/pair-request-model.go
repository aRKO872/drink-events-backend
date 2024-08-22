package models

import "encoding/json"

type PairRequest struct {
	Id         string `json:"id" db:"id"`
	SenderID   string `json:"sender_id" db:"sender_id"`
	ReceiverID string `json:"receiver_id" db:"receiver_id"`
	IsActive   bool   `json:"is_active" db:"is_active"`
	CreatedAt  string `json:"created_at" db:"created_at"`
	UpdatedAt  string `json:"updated_at" db:"updated_at"`
}

func (pr PairRequest) MarshalBinary() ([]byte, error) {
	return json.Marshal(pr)
}

func (pr *PairRequest) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &pr)
}

type PairRequestsSentByIdItem struct {
	Id           string `json:"id"`
	ReceiverID   string `json:"receiver_id"`
	ProfilePic   string `json:"profile_pic"`
	ReceiverName string `json:"receiver_name"`
	CreatedAt    string `json:"created_at"`
}

type PairRequestsSentToIdItem struct {
	Id         string `json:"id"`
	SenderID   string `json:"sender_id"`
	SenderName string `json:"sender_name"`
	ProfilePic string `json:"profile_pic"`
	CreatedAt  string `json:"created_at"`
}

type PairRequestsSentInput struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type PairRequestSentToIdOutput struct {
	Status       bool                       `json:"status"`
	ErrorMsg     string                     `json:"msg"`
	PairRequests []PairRequestsSentToIdItem `json:"pair-requests"`
}

type PairRequestSentByIdOutput struct {
	Status       bool                       `json:"status"`
	ErrorMsg     string                     `json:"msg"`
	PairRequests []PairRequestsSentByIdItem `json:"pair-requests"`
}
