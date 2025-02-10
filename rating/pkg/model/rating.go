package model

// RecordID defined a record id. Together with RecordType
// Identifies unique records across all types.
type RecordID string

type RecordType string

const (
	RecordTypeMovie = RecordType("movie")
)

type UserID string

type RatingValue int

type Rating struct {
	RecordID    `json:"record_id"`
	RecordType  `json:"record_type"`
	UserID      `json:"user_id"`
	RatingValue `json:"rating_value"`
}

type RatingEvent struct {
	UserID     UserID
	RecordID   RecordID
	RecordType RecordType
	Value      RatingValue
	EventType  RatingEventType
}

type RatingEventType string

const (
	RatingEventTypePut    = "put"
	RatingEventTypeDelete = "delete"
)
