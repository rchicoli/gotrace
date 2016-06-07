package main

import (
	"encoding/json"
	"fmt"
)

type Commands []*Command

const (
	CmdCreate = "create goroutine"
	CmdStop   = "stop goroutine"
	CmdSend   = "send to channel"
)

// Command is a common structure for all
// types of supported events (aka 'commands').
// It's main purpose to handle JSON marshalling.
type Command struct {
	Time     int64       "json:\"t\""
	Command  string      "json:\"command\""
	Name     string      "json:\"name,omitempty\""
	Parent   string      "json:\"parent,omitempty\""
	Channels []string    "json:\"channels,omitempty\""
	From     string      "json:\"from,omitempty\""
	To       string      "json:\"to,omitempty\""
	Channel  string      "json:\"ch,omitempty\""
	Value    interface{} "json:\"value,omitempty\""
	EventID  string      "json:\"eid,omitempty\""
	Duration int64       "json:\"duration,omitempty\""
}

func (c *Commands) toJSON() []byte {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic(err)
	}

	return data
}

func (c *Commands) StartGoroutine(ts int64, name string, gid, pid uint64) {
	parent := fmt.Sprintf("#%d", pid)
	// ignore parent for 'main()' which has pid 0
	if pid == 0 {
		parent = ""
	}
	cmd := &Command{
		Time:    ts,
		Command: CmdCreate,
		Name:    fmt.Sprintf("#%d", gid),
		Parent:  parent,
	}
	*c = append(*c, cmd)
}

func (c *Commands) StopGoroutine(ts int64, name string, gid uint64) {
	cmd := &Command{
		Time:    ts,
		Command: CmdStop,
		Name:    fmt.Sprintf("#%d", gid),
	}
	*c = append(*c, cmd)
}

func (c *Commands) ChanSend(ts int64, cid, fgid, tgid, val uint64) {
	cmd := &Command{
		Time:    ts,
		Command: CmdSend,
		From:    fmt.Sprintf("#%d", fgid),
		To:      fmt.Sprintf("#%d", tgid),
		Channel: fmt.Sprintf("#%d", cid),
		Value:   fmt.Sprintf("%d", val),
	}
	*c = append(*c, cmd)
}

//ByTimestamp implements sort.Interface for sorting command by timestamp.
type ByTimestamp Commands

func (a ByTimestamp) Len() int           { return len(a) }
func (a ByTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimestamp) Less(i, j int) bool { return a[i].Time < a[j].Time }

// Count counts total number of commands
func (c Commands) Count() int {
	return len(c)
}

// CountCreateGoroutine counts total number of CreateGoroutine commands.
func (c Commands) CountCreateGoroutine() int {
	var count int
	for _, cmd := range c {
		if cmd.Command == CmdCreate {
			count++
		}
	}
	return count
}

// CountStopGoroutine counts total number of StopGoroutine commands.
func (c Commands) CountStopGoroutine() int {
	var count int
	for _, cmd := range c {
		if cmd.Command == CmdStop {
			count++
		}
	}
	return count
}

// CountSendToChannel counts total number of SendToChannel commands.
func (c Commands) CountSendToChannel() int {
	var count int
	for _, cmd := range c {
		if cmd.Command == CmdSend {
			count++
		}
	}
	return count
}
