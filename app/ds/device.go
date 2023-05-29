package ds

import "time"

type Device struct {
	ID        int64     `db:"id,omitempty" json:"id,omitempty"`
	UUID      string    `db:"uuid" json:"uuid"`
	IP        string    `db:"ip" json:"ip"`
	CreatedAt time.Time `db:"create_at" json:"create_at"`
	UpdatedAt time.Time `db:"update_at" json:"update_at"`
}
