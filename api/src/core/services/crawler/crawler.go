package crawler

import (
	"dse/src/core/models"
	"dse/src/core/services/db"
	"dse/src/utils"
	"dse/src/utils/event"
	"os"
	"time"

	"github.com/dromara/carbon/v2"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/tidwall/gjson"
)

// ------------------------------------------------------------
// : Aliases
// ------------------------------------------------------------
type User  = models.User
type State = models.State

type Keyword = models.Keyword
type Website = models.Website
type Task    = models.Task

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var (
	logger = utils.NewLogger()
	
	hash string
	keywords = []string{}
	websites = []Website{}

	queue = cmap.New[*User]()
)
// ------------------------------------------------------------
// : Helpers
// ------------------------------------------------------------
func Populate(user *User) {
	if user.CrawlingFlag().Load() { return }

	counter := 0

	var state = user.State.Server
	state.UpdatedAt   = carbon.Now(carbon.UTC).ToIso8601String()
	state.StartedAt   = ""
	state.CompletedAt = ""
	state.TaskHash    = hash
	state.ClearTasks()

	for _, keyword := range keywords {
		for _, website := range websites {
			task := &Task{
				ID  : counter,
				Type: "search",

				Keyword: keyword,
				Website: website,

				CreatedAt: carbon.Now(carbon.UTC).ToIso8601String(),
			}
			*user.State.Server.Tasks = append(*user.State.Server.Tasks, task)

			counter += 1
		}
	}
	user.Save()
}

// ------------------------------------------------------------
// : Listeners
// ------------------------------------------------------------
func OnConnect(user *User) {
	// TODO: TBI
}

func OnUpdate(user *User) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error().Interface("recover", r).Msg("Recovered")
		}
	}()

	if !user.RequiersScraping() { return }

	if queue.Has(user.Token) {
		return
	}

	queue.Set(user.Token, user)
	defer queue.Remove(user.Token)

	// Generate tasks if user is new
	if len(*user.State.Server.Tasks) == 0 {
		Populate(user)
	}
	
	user.Start()
}

func OnReset(user *User) {
	user.Reset()
	Populate(user)
}

func OnResetAll() {
	logger.Info().Msg("Resetting all users")
	users, err := db.GetUsers()
	if err != nil {
		logger.Error().Err(err).Msg("Failed	to get users")
		return
	}

	users.Each(func(index int, user *User) bool {
		Populate(user)
		return true
	})
}

// ------------------------------------------------------------
// : Monnitor
// ------------------------------------------------------------
func Listen() {
	go func() {
		ch := event.On("user.connected")

		for e := range ch {
			user, ok := e.Args[0].(*User)
			if !ok { continue }
			go OnConnect(user)
		}
	}()

	go func() {
		ch := event.On("user.updated")

		for e := range ch {
			user, ok := e.Args[0].(*User)
			if !ok { continue }
			go OnUpdate(user)
		}
	}()

	go func() {
		ch := event.On("user.reset")

		for e := range ch {
			user, ok := e.Args[0].(*User)
			if !ok { continue }
			go OnReset(user)
		}
	}()
}

func Monitor() {
	time.Sleep(5 * time.Second) // TODO: Change to 10 seconds

	go func() {
		for {
			users, err := db.GetUsers()
			if err != nil { continue }

			counter := 0

			users.Each(func(index int, user *User) bool {
				counter += 1

				if !user.RequiersScraping() { return true }
				if !user.IsOnline()         { return true }

				return false
			})

			time.Sleep(5 * time.Second)
		}
	}()
}

// ------------------------------------------------------------
// : Init
// ------------------------------------------------------------
func Init() {
	b, err := os.ReadFile("config/searches.json")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to read searches")
		return
	}

	hash = utils.GenerateHash(string(b))
	parsed := gjson.ParseBytes(b)
	parsed.Get("keywords").ForEach(func(_, v gjson.Result) bool {
		keywords = append(keywords, v.String())
		return true
	})

	parsed.Get("websites").ForEach(func(_, v gjson.Result) bool {
		website := Website{
			Name : v.Get("name").String(),
			Query: v.Get("query").String(),
			Url  : v.Get("url").String(),
		}
		websites = append(websites, website)

		return true
	})

	go Listen()
	go Monitor()
}
