package models

import (
	"dse/src/utils/datetime"
	"dse/src/utils/gatekeeper"
	"sync"
	"time"

	"github.com/dromara/carbon/v2"
)

// ------------------------------------------------------------
// : State
// ------------------------------------------------------------
type State struct {
	Client *ClientState `json:"client"`
	Server *ServerState `json:"server"`
}

func NewState() *State {
	state := &State{}
	sc    := NewClientState()
	ss    := NewServerState()

	state.Client = sc
	state.Server = ss

	return state
}

// ------------------------------------------------------------
// : Client
// ------------------------------------------------------------
type ClientState struct {
	Browser    map[string]interface{} `json:"browser"`
	User       *ClientUser            `json:"user"`
	Background *ClientBackground      `json:"background"`
	Crawler    *ClientCrawler         `json:"crawler"`
	Extension  *ClientExtension       `json:"extension"`
}

func (c *ClientState) GetCrawler() *ClientCrawler {
	return c.Crawler
}

type ClientUser struct {
	Token  string   `json:"token"`
	Form   *Form	`json:"form"`
	Popup  bool     `json:"popup"`
	Type   string   `json:"type"`
}

type ClientBackground struct {
	State string `json:"state"`
}

type ClientCrawler struct {
	State       string `json:"state"`
	Window      string `json:"window"`
	StartedAt   string `json:"started_at"`
	CompletedAt string `json:"completed_at"`

	mutex sync.Mutex
	cond  sync.Cond
}

func (c *ClientCrawler) SetState(state string) {
	c.State = state
}

func (c *ClientCrawler) GetState() string {
	return c.State
}

type ClientExtension struct {
	State    string	`json:"state"`
	Version  string `json:"version"`
	Language string `json:"language"`
}

// Constructor
func NewClientState() *ClientState {
	c := &ClientState{
		Browser: map[string]interface{}{},
		User: &ClientUser{
			Token: "",
			Form : NewForm(),
			Popup: false,
			Type : "",
		},
		Background: &ClientBackground{
			State: "",
		},
		Crawler: &ClientCrawler{
			State      : "",
			StartedAt  : "",
			CompletedAt: "",
		},
		Extension: &ClientExtension{
			State   : "",
			Version : "",
			Language: "",
		},
	}
	return c
}
// ------------------------------------------------------------
// : Server > State
// ------------------------------------------------------------
type ServerState struct {
	State     string    `json:"state"`
	Online    bool      `json:"online"`
	LastPing  string    `json:"last_ping"`

	UpdatedAt   string `json:"updated_at"`
	StartedAt   string `json:"started_at"`
	CompletedAt string `json:"completed_at"`
	
	Tasks  *[]*Task `json:"tasks"`
	TaskHash string `json:"task_hash"`

	GK *gatekeeper.GateKeeper `json:"-"`
}

func NewServerState() *ServerState {
	ss := &ServerState{
		Online  : false,
		LastPing: "",

		UpdatedAt  : "",
		StartedAt  : "",
		CompletedAt: "",

		Tasks : &[]*Task{},
		TaskHash: "",
	}

	return ss
}

func (s *ServerState) ClearTasks() {
	s.Tasks = &[]*Task{}
}

func (s *ServerState) GetTasks() []*Task {
	return *s.Tasks
}

func (s *ServerState) AddTask(task *Task) {
	*s.Tasks = append(*s.Tasks, task)
}
// ------------------------------------------------------------
// : Server > State > Crawler > Tasks
// ------------------------------------------------------------
type Task struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`

	Keyword     string      `json:"keyword"`
	Website     interface{} `json:"website"`

	Filepath    string `json:"filepath"`

	CreatedAt   string `json:"created_at"`
	StartedAt   string `json:"started_at"`
	CompletedAt string `json:"completed_at"`
}

func (t *Task) Reset() {
	t.StartedAt   = ""
	t.CompletedAt = ""
}

func (t *Task) GetCompletedAt() time.Time {
	return carbon.Parse(t.CompletedAt).StdTime()
}


func (t *Task) IsStale() bool {
	monday  := datetime.ToTime(datetime.StartOfWeek())
	updated := t.GetCompletedAt()

	return updated.IsZero() || updated.Before(monday)
}

func (t *Task) IsCompleted() bool {
	monday  := datetime.ToTime(datetime.StartOfWeek())
	return  !t.GetCompletedAt().After(monday)
}
