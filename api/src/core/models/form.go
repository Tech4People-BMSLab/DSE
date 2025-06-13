package models

// ------------------------------------------------------------
// : Form
// ------------------------------------------------------------
type Form struct {
	Age          string        `json:"age"`
	Browser      *Browser      `json:"browser"`
	Education    string        `json:"education"`
	Employment   string        `json:"employment"`
	Income       string        `json:"income"`
	Language     *Language     `json:"language"`
	Political    string        `json:"political"`
	Postcode     *Postcode     `json:"postcode"`
	Resident     string        `json:"resident"`
	SearchEngine *SearchEngine `json:"search_engine"`
	Sex          string        `json:"sex"`
	Social       *Social       `json:"social"`
}

func NewForm() *Form {
	return &Form{
		Age          : "",
		Browser      : &Browser{
			Brave         : false,
			Chrome        : false,
			Firefox       : false,
			MicrosoftEdge : false,
			Opera         : false,
			Safari        : false,
			Unselected    : false,
		},
		Education    : "",
		Employment   : "",
		Income       : "",
		Language     : &Language{
			Duits      : false,
			Engels     : false,
			Frans      : false,
			Italiaans  : false,
			Nederlands : false,
			Spaans     : false,
			Unselected : false,
		},
		Political    : "",
		Postcode     : &Postcode{
			Value: "",
		},
		Resident     : "",
		SearchEngine : &SearchEngine{
			Anders     : false,
			Bing       : false,
			Duckduckgo : false,
			Ecosia     : false,
			Google     : false,
			Startpage  : false,
			Yahoo      : false,
			Unselected : false,
		},
	}
}

// ------------------------------------------------------------
// : Browser
// ------------------------------------------------------------
type Browser struct {
	Brave         bool `json:"brave"`
	Chrome        bool `json:"chrome"`
	Firefox       bool `json:"firefox"`
	MicrosoftEdge bool `json:"microsoft-edge"`
	Opera         bool `json:"opera"`
	Safari        bool `json:"safari"`
	Unselected    bool `json:"unselected"`
}

// ------------------------------------------------------------
// : Language
// ------------------------------------------------------------
type Language struct {
	Duits      bool `json:"duits"`
	Engels     bool `json:"engels"`
	Frans      bool `json:"frans"`
	Italiaans  bool `json:"italiaans"`
	Nederlands bool `json:"nederlands"`
	Spaans     bool `json:"spaans"`
	Unselected bool `json:"unselected"`
}

// ------------------------------------------------------------
// : SearchEngine
// ------------------------------------------------------------
type SearchEngine struct {
	Anders     bool `json:"anders"`
	Bing       bool `json:"bing"`
	Duckduckgo bool `json:"duckduckgo"`
	Ecosia     bool `json:"ecosia"`
	Google     bool `json:"google"`
	Startpage  bool `json:"startpage"`
	Yahoo      bool `json:"yahoo"`
	Unselected bool `json:"unselected"`
}

// ------------------------------------------------------------
// : Social
// ------------------------------------------------------------
type Social struct {
	Anders         bool `json:"anders"`
	DeKrant        bool `json:"de-krant"`
	Facebook       bool `json:"facebook"`
	WhatsApp	   bool `json:"whatsapp"`
	Instagram      bool `json:"instagram"`
	Linkedin       bool `json:"linkedin"`
	Nieuwswebsites bool `json:"nieuwswebsites"`
	YouTube	       bool `json:"youtube"`
	Radio          bool `json:"radio"`
	Reddit         bool `json:"reddit"`
	Telegram       bool `json:"telegram"`
	TV             bool `json:"tv"`
	Twitter        bool `json:"twitter"`
	Unselected     bool `json:"unselected"`
}

// ------------------------------------------------------------
// : Postcode
// ------------------------------------------------------------
type Postcode struct {
	Value string `json:"value"`
}
