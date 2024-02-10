package main

import (
	"time"
)

var users = map[string]string{
	"user1": "a",
	"user2": "password2",
}

var sessions = map[string]session{}

type session struct {
	username string
	expiry   time.Time
}

func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}
