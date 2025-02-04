package types

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullString struct {
	sql.NullString
}

func NewNullString(str sql.NullString) NullString {
	return NullString{
		NullString: str,
	}
}

func NewNullStringFromString(str string) NullString {
	return NullString{
		NullString: sql.NullString{
			String: str,
			Valid:  true,
		},
	}
}

func (s *NullString) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	}
	return []byte("null"), nil
}

func (s *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		s.Valid = false
		return nil
	}
	s.Valid = true
	return json.Unmarshal(data, &s.String)
}

type NullInt struct {
	sql.NullInt32
}

func NewNullInt(i sql.NullInt32) NullInt {
	return NullInt{
		NullInt32: i,
	}
}

func NewNullIntFromInt(i int32) NullInt {
	return NullInt{
		NullInt32: sql.NullInt32{
			Int32: i,
			Valid: true,
		},
	}
}

func (i *NullInt) MarshalJSON() ([]byte, error) {
	if i.Valid {
		return json.Marshal(i.Int32)
	}
	return []byte("null"), nil
}

func (i *NullInt) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		i.Valid = false
		return nil
	}
	i.Valid = true
	return json.Unmarshal(data, &i.Int32)
}

type NullTime struct {
	sql.NullTime
}

func NewNullTime(t sql.NullTime) NullTime {
	return NullTime{
		NullTime: t,
	}
}

func NewNullTimeFromTime(t time.Time) NullTime {
	return NullTime{
		NullTime: sql.NullTime{
			Time:  t,
			Valid: true,
		},
	}
}

func (t *NullTime) MarshalJSON() ([]byte, error) {
	if t.Valid {
		return json.Marshal(t.Time)
	}
	return []byte("null"), nil
}

func (t *NullTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		t.Valid = false
		return nil
	}
	t.Valid = true
	return json.Unmarshal(data, &t.Time)
}
