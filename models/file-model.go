package models

import "encoding/json"

type FileObject struct {
	Id string `json:"id" db:"id"`
	FileName string `json:"file_name" db:"file_name"`
	AWS_Key string `json:"aws_key" db:"aws_key"`
	Etag string `json:"etag" db:"etag"`
	SecureURL string `json:"secure_url" db:"secure_url"`
}

func (f FileObject) MarshalBinary() ([]byte, error) {
	return json.Marshal(f)
}

func (f *FileObject) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &f)
}