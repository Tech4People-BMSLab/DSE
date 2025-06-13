package ws

import (
	"dse/src/core/models"
	"dse/src/core/services/db"
	"dse/src/utils"
	"dse/src/utils/hashmap"
	"encoding/json"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

// ------------------------------------------------------------
// : Aliases
// ------------------------------------------------------------
type ResponseWriter = http.ResponseWriter
type Request        = http.Request
type User    		= models.User
type Packet         = models.Packet
type Connection     = websocket.Conn
type Address        = net.Addr
// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var (
	logger  = utils.NewLogger()
	clients = hashmap.NewHashMap[*Connection, *User]()

	upgrader = websocket.Upgrader{
		ReadBufferSize   : 1024,
		WriteBufferSize  : 1024,
		EnableCompression: true,
		CheckOrigin      : func(r *http.Request) bool {
			return true
		},
	}

	version string = "3.0.5"
)

// ------------------------------------------------------------
// : Methods
// ------------------------------------------------------------
func SendToUser(token string) {

}

// ------------------------------------------------------------
// : Handlers
// ------------------------------------------------------------
func OnUserUpdate(user *User, packet *Packet) {
	b, err := json.Marshal(packet.Data)
	if err != nil {
		logger.Error().Err(err).Msg("Error marshalling packet data")
		return
	}

	parsed := gjson.ParseBytes(b)

	user.SetVersion(packet.Version)
	
	user.State.Client.Background.State = parsed.Get("background.state").String()
	user.State.Client.Browser = map[string]interface{}{
		"name"       : parsed.Get("browser.name").String(),
		"os"         : parsed.Get("browser.os").String(),
		"os_version" : parsed.Get("browser.os_version").String(),
		"version"    : parsed.Get("browser.version").String(),
	}

	user.State.Client.Crawler.State       = parsed.Get("crawler.state").String()
	user.State.Client.Crawler.Window      = parsed.Get("crawler.window").String()
	user.State.Client.Crawler.StartedAt   = parsed.Get("crawler.started_at").String()
	user.State.Client.Crawler.CompletedAt = parsed.Get("crawler.completed_at").String()

	user.State.Client.Extension.State    = parsed.Get("extension.state").String()
	user.State.Client.Extension.Version  = parsed.Get("extension.version").String()
	user.State.Client.Extension.Language = parsed.Get("extension.language").String()

	user.State.Client.User.Type  = parsed.Get("user.type").String()
	user.State.Client.User.Token = parsed.Get("user.token").String()
	user.State.Client.User.Popup = parsed.Get("user.popup").Bool()

	user.State.Client.User.Form.Age = parsed.Get("user.form.age").String()

	user.State.Client.User.Form.Browser.Brave         = parsed.Get("user.form.browser.brave").Bool()
	user.State.Client.User.Form.Browser.Chrome        = parsed.Get("user.form.browser.chrome").Bool()
	user.State.Client.User.Form.Browser.Firefox       = parsed.Get("user.form.browser.edge").Bool()
	user.State.Client.User.Form.Browser.MicrosoftEdge = parsed.Get("user.form.browser.firefox").Bool()
	user.State.Client.User.Form.Browser.Opera         = parsed.Get("user.form.browser.opera").Bool()
	user.State.Client.User.Form.Browser.Safari        = parsed.Get("user.form.browser.safari").Bool()
	user.State.Client.User.Form.Browser.Unselected	  = parsed.Get("user.form.browser.unselected").Bool()

	user.State.Client.User.Form.Education  = parsed.Get("user.form.education").String()
	user.State.Client.User.Form.Employment = parsed.Get("user.form.employment").String()
	user.State.Client.User.Form.Income     = parsed.Get("user.form.income").String()

	user.State.Client.User.Form.Language.Duits      = parsed.Get("user.form.language.duits").Bool()
	user.State.Client.User.Form.Language.Engels     = parsed.Get("user.form.language.engels").Bool()
	user.State.Client.User.Form.Language.Frans      = parsed.Get("user.form.language.frans").Bool()
	user.State.Client.User.Form.Language.Italiaans  = parsed.Get("user.form.language.italiaans").Bool()
	user.State.Client.User.Form.Language.Nederlands = parsed.Get("user.form.language.nederlands").Bool()
	user.State.Client.User.Form.Language.Spaans     = parsed.Get("user.form.language.spaans").Bool()
	user.State.Client.User.Form.Language.Unselected = parsed.Get("user.form.language.unselected").Bool()

	user.State.Client.User.Form.Political      = parsed.Get("user.form.political").String()
	user.State.Client.User.Form.Postcode.Value = parsed.Get("user.form.postcode.value").String()
	user.State.Client.User.Form.Resident       = parsed.Get("user.form.resident").String()

	user.State.Client.User.Form.SearchEngine.Anders     = parsed.Get("user.form.search_engine.anders").Bool()
	user.State.Client.User.Form.SearchEngine.Bing       = parsed.Get("user.form.search_engine.bing").Bool()
	user.State.Client.User.Form.SearchEngine.Duckduckgo = parsed.Get("user.form.search_engine.duckduckgo").Bool()
	user.State.Client.User.Form.SearchEngine.Ecosia     = parsed.Get("user.form.search_engine.ecosia").Bool()
	user.State.Client.User.Form.SearchEngine.Google     = parsed.Get("user.form.search_engine.google").Bool()
	user.State.Client.User.Form.SearchEngine.Startpage  = parsed.Get("user.form.search_engine.startpage").Bool()
	user.State.Client.User.Form.SearchEngine.Yahoo      = parsed.Get("user.form.search_engine.yahoo").Bool()
	user.State.Client.User.Form.SearchEngine.Unselected = parsed.Get("user.form.search_engine.unselected").Bool()
	user.Save()
}

func OnReceive(conn *Connection, packet *Packet) {
	var err   error
	var user *User
	
	exists, _ := db.HasUser(packet.From)
	if !exists {
		user, err = db.CreateUser(packet.From)
		if err != nil {
			logger.Error().Err(err).Msg("Error creating user")
			return
		}
	}

	user, err = db.GetUser(packet.From)
	if err != nil {
		logger.Error().Err(err).Msg("Error getting user")
		return
	}

	clients.Set(conn, user)
	user.SetWS(conn)

	switch packet.Action {
		case "update": {
			go OnUserUpdate(user, packet)
		}
	}
}

func HandleWS(w ResponseWriter, r *Request) {
    defer func() {
        if rec := recover(); rec != nil {
            logger.Error().Interface("recover", rec).Msg("Recovered from panic in WebSocket handler")
        }
    }()


    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        logger.Error().Err(err).Msg("WebSocket Upgrade Error")
        return
    }
    defer conn.Close()


    for {
		msgtype, msg, err := conn.ReadMessage()
		if err != nil {
			// Handle WebSocket closure or unexpected error
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error().Err(err).Msg("Unexpected WebSocket error")
			} else {
				logger.Info().Msg("WebSocket connection closed")
			}
	
			// Check and clean up client on disconnect
			exists := clients.Has(conn)
			if exists {
				user := clients.MustGet(conn)
				user.State.Server.Online = false
				user.Save()
				clients.Delete(conn)
			}
			break
		}
	
		// Log the raw message for debugging
		logger.Debug().Discard().Str("raw_msg", string(msg)).Msg("Received raw WebSocket message")
	
		switch msgtype {
		case websocket.TextMessage:
			// Use json.Valid for accurate validation
			if !json.Valid(msg) {
				logger.Error().Discard().Str("msg", string(msg)).Msg("Invalid JSON format in WebSocket message")
				continue
			}
	
			var packet Packet
			err = json.Unmarshal(msg, &packet)
			if err != nil {
				logger.Error().Err(err).Str("msg", string(msg)).Msg("JSON Unmarshal Error")
				continue
			}
			OnReceive(conn, &packet)
	
		default:
			logger.Warn().Int("type", msgtype).Msg("Unhandled message type")
		}
	}
}
