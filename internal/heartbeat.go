package internal

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type HeartbeatError struct {
	Msg string
	Err error
}

func (e *HeartbeatError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Err)
	}
	return e.Msg
}

func (e *HeartbeatError) Unwrap() error {
	return e.Err
}

type Heartbeat struct {
	Packet       []byte
	Port         int
	Host         string
	RecvResponse bool
	Timeout      time.Duration
}

func Send(h Heartbeat) error {
	addr := net.JoinHostPort(h.Host, strconv.Itoa(h.Port))

	conn, err := net.DialTimeout("udp", addr, h.Timeout)

	if err != nil {
		return &HeartbeatError{
			Msg: fmt.Sprintf("socket error heartbeating with port %d", h.Port),
			Err: err,
		}
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(h.Timeout)); err != nil {
		return &HeartbeatError{
			Msg: fmt.Sprintf("failed to set deadline for port %d", h.Port),
			Err: err,
		}
	}

	if _, err := conn.Write(h.Packet); err != nil {
		return &HeartbeatError{
			Msg: fmt.Sprintf("socket error heartbeating with port %d", h.Port),
			Err: err,
		}
	}

	if h.RecvResponse {
		buffer := make([]byte, 4096)

		_, err := conn.Read(buffer)

		if err != nil {
			return &HeartbeatError{
				Msg: fmt.Sprintf("socket error heartbeating with port %d", h.Port),
				Err: err,
			}
		}
	}

	return nil
}
