package events

import "time"

type CmdEvent struct {
	Type       string    `json:"type"`
	Ts         time.Time `json:"ts"`
	Shell      string    `json:"shell"`
	Host       string    `json:"host"`
	User       string    `json:"user"`
	Cwd        string    `json:"cwd"`
	Cmd        string    `json:"cmd"`
	Exit       int       `json:"exit"`
	DurationMs int64     `json:"duration_ms,omitempty"`
}
