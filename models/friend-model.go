package models

import "encoding/json"

type Friends struct {
	Id        string `json:"id,omitempty" db:"id"`
	Friend1ID string `json:"friend1_id,omitempty" db:"friend1_id"`
	Friend2ID string `json:"friend2_id,omitempty" db:"friend2_id"`
	IsActive  bool   `json:"is_active,omitempty" db:"is_active"`
	CreatedAt string `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty" db:"updated_at"`
}

func (f Friends) MarshalBinary() ([]byte, error) {
	return json.Marshal(f)
}

func (f *Friends) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &f)
}
