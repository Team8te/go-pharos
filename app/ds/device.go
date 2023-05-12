package ds

import "time"

type Device struct {
	ID        int64
	UUID      string
	IP        string
	CreatedAt time.Time
	UpdatedAt time.Time
}
