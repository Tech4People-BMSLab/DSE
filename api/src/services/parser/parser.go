package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"

	"bms.dse/src/utils/logutil"
	"github.com/tidwall/sjson"
	"gopkg.in/xmlpath.v2"
)

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var logger = logutil.NewLogger("Parser")

// ------------------------------------------------------------
// : Helpers
// ------------------------------------------------------------
func WriteHTML(html string) {
	f, _ := os.Create("output.html")
	defer f.Close()
	f.WriteString(html)
}

func RemoveNonPrintable(s string) string {
    return strings.Map(func(r rune) rune {
        if unicode.IsPrint(r) {
            return r
        }
        return -1
    }, s)
}

func Unescape(s string) string {
	json.Unmarshal([]byte(s), &s)
	strings.TrimSpace(s)

	s = strings.ReplaceAll(s, "\u0026", "&")
	s = strings.ReplaceAll(s, "\u003c", "<")
	s = strings.ReplaceAll(s, "\u003e", ">")
	s = strings.ReplaceAll(s, "\u0022", "\"")
	s = strings.ReplaceAll(s, "\u0027", "'")
	s = strings.ReplaceAll(s, "\u002f", "/")
	s = strings.ReplaceAll(s, "\u003d", "=")
	s = strings.ReplaceAll(s, "\u0025", "%")
	s = strings.ReplaceAll(s, "\u002b", "+")
	s = strings.ReplaceAll(s, "\u0023", "#")
	s = strings.ReplaceAll(s, "\u003f", "?")
	
	return s
}


// ------------------------------------------------------------
// : Handlers
// ------------------------------------------------------------
func ParseDuckDuckGo(html string) (string, error) {
	var output string
	
	logger.Info().Msg("Parsing DuckDuckGo")

	reader    := strings.NewReader(html)
	root, err := xmlpath.ParseHTML(reader)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to parse HTML")
		return "", err
	}

	WriteHTML(html)

	{ //// searc_result
		path_nodes := xmlpath.MustCompile(`//article`)
		counter := 0
		for it := path_nodes.Iter(root); it.Next(); {
			var title       string
			var link        string
			var description string

			path_title  := xmlpath.MustCompile(`.//div//h2/a/span/text()`)
			path_link   := xmlpath.MustCompile(`.//div/h2/a/@href`)
			path_desc   := xmlpath.MustCompile(`.//span[@class="kY2IgmnCmOGjharHErah"]/text()`)
			
			title, _ = path_title.String(it.Node())
			title    = Unescape(title)

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)

			description, _ = path_desc.String(it.Node())
			description    = Unescape(description)

			value := map[string]string{
				"title"      : title,
				"link"       : link,
				"description": description,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("search_result.%d", counter), value)
			counter += 1
		}

	}

	{ //// recent_news
		path_nodes := xmlpath.MustCompile(`//div[@class="module--carousel__item has-image"]`)
		counter := 0
		for it := path_nodes.Iter(root); it.Next(); {
			var title       string
			var link        string
			var description string

			path_title       := xmlpath.MustCompile(`.//a`)
			path_link        := xmlpath.MustCompile(`.//a/@href`)
			path_desc := xmlpath.MustCompile(`.//span[@class="module--carousel__source result__url"]`)

			title, _ = path_title.String(it.Node())
			title    = Unescape(title)

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)

			description, _ = path_desc.String(it.Node())
			description    = Unescape(description)

			value := map[string]string{
				"title"      : title,
				"link"       : link,
				"description": description,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("recent_news.%d", counter), value)
			counter += 1
		}
	}

	obj := map[string]interface{}{}

	json.Unmarshal([]byte(output), &obj)

	return output, nil
}

func ParseYouTube(html string) (string, error) {
	logger.Warn().Msg("Parsing YouTube")

	/**
	* NOTE: YouTube has been removed from the current implementation (9 Apr 2024)
	*       Due to YouTube restricting adding query parameters to the URL.
	*/
	return "", errors.New("Not implemented")
}

func ParseBing(html string) (string, error) {
	var output string

	logger.Info().Msg("Parsing Bing")

	reader    := strings.NewReader(html)
	root, err := xmlpath.ParseHTML(reader)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to parse HTML")
		return "", err
	}

	{ //// search_result
		counter    := 0
		path_nodes := xmlpath.MustCompile(`//li[contains(@class, "b_algo")]`)

		path_title := xmlpath.MustCompile(`.//h2//text()`)
		path_link  := xmlpath.MustCompile(`.//h2//a/@href`)
		path_desc  := xmlpath.MustCompile(`.//p[contains(@class, "b_paractl")]//text()`)

		//// 1
		for it := path_nodes.Iter(root); it.Next(); {
			var title string
			var link  string
			var desc  string

			title, _ = path_title.String(it.Node())
			title    = Unescape(title)

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)

			desc_iter := path_desc.Iter(it.Node())
			for desc_iter.Next() {
				desc += desc_iter.Node().String()
			}
			desc    = Unescape(desc)

			value := map[string]string{
				"title"      : title,
				"link"       : link,
				"description": desc,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("search_result.%d", counter), value)
			counter += 1
		}
	}

	return output, nil
}

func ParseYahoo(html string) (string, error) {
	var output string

	logger.Info().Msg("Parsing Yahoo")

	reader    := strings.NewReader(html)
	root, err := xmlpath.ParseHTML(reader)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to parse HTML")
		return "", err
	}

	{ //// search_result
		counter    := 0
		path_nodes := xmlpath.MustCompile(`//div[contains(@class, 'algo-sr')]`)

		path_title := xmlpath.MustCompile(`./div[contains(@class, 'compTitle')]//a/text()`)
		path_link  := xmlpath.MustCompile(`./div[contains(@class, 'compTitle')]//a/@href`)
		path_desc  := xmlpath.MustCompile(`./div[contains(@class, 'compText')]//span`)

		//// 1
		for it := path_nodes.Iter(root); it.Next(); {
			var title string
			var link  string
			var desc  string
			
			title, _ = path_title.String(it.Node())
			title    = Unescape(title)

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)

			desc, _ = path_desc.String(it.Node())
			desc    = Unescape(desc)

			value := map[string]string{
				"title"       : title,
				"link"        : link,
				"description" : desc,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("search_result.%d", counter), value)
			counter += 1
		}
	}

	return output, nil
}

func ParseGoogleNews(html string) (string, error) {
	var output string

	logger.Info().Msg("Parsing Google News")

	reader    := strings.NewReader(html)
	root, err := xmlpath.ParseHTML(reader)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to parse HTML")
		return "", err
	}

	{ //// search_result
		counter    := 0
		path_nodes := xmlpath.MustCompile(`//div[contains(@class, "SoaBEf")]`)
		
		//// 1
		path_title := xmlpath.MustCompile(`.//div[contains(@class, "n0jPhd ")]`)
		path_link  := xmlpath.MustCompile(`.//a/@href`)
		path_desc  := xmlpath.MustCompile(`.//div[contains(@class, "GI74Re")]/text()`)

		for it := path_nodes.Iter(root); it.Next(); {
			var title string
			var link  string
			var desc  string

			title, _ = path_title.String(it.Node())
			title    = Unescape(title)

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)
			
			desc, _ = path_desc.String(it.Node())
			desc    = Unescape(desc)
			
			value := map[string]string{
				"title": title,
				"link" : link,
				"description" : desc,
			}	
			
			output, _ = sjson.Set(output, fmt.Sprintf("search_result.%d", counter), value)
			counter += 1
		}
	}

	return output, nil
}

func ParseGoogleVideos(html string) (string, error) {
	var output string

	logger.Info().Msg("Parsing Google Videos")

	reader    := strings.NewReader(html)
	root, err := xmlpath.ParseHTML(reader)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to parse HTML")
		return "", err
	}


	{ //// search_result
		counter    := 0
		path_nodes := xmlpath.MustCompile(`//div[contains(@class, "MjjYud")]`)

		//// 1
		path_title := xmlpath.MustCompile(`.//div[contains(@class, "nhaZ2c")]//h3`)
		path_link  := xmlpath.MustCompile(`.//div[contains(@class, "nhaZ2c")]/div/span/a/@href`)
		path_desc  := xmlpath.MustCompile(`//div[contains(@class, "ITZIwc")]//text()`)

		for it := path_nodes.Iter(root); it.Next(); {
			var title  string
			var link   string
			var desc   string

			title, _ = path_title.String(it.Node())
			title    = Unescape(title)

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)

			desc_iter := path_desc.Iter(it.Node())
			for desc_iter.Next() {
				desc += desc_iter.Node().String()
			}
			desc    = Unescape(desc)

			value := map[string]string{
				"title": title,
				"link" : link,
				"description" : desc,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("search_result.%d", counter), value)
			counter += 1
		}
	}

	return output, nil
}

func ParseGoogle(html string) (string, error) {
	var output string

	logger.Info().Msg("Parsing Google")

	reader    := strings.NewReader(html)
	root, err := xmlpath.ParseHTML(reader)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to parse HTML")
		return "", err
	}

	{ //// search_result
		counter    := 0
		path_nodes := xmlpath.MustCompile(`//*[@class='g Ww4FFb vt6azd tF2Cxc asEBEc']`)
		
		//// 1
		path_title := xmlpath.MustCompile(`.//h3[@class='LC20lb MBeuO DKV0Md']/text()`)
		path_link  := xmlpath.MustCompile(`.//h3[@class='LC20lb MBeuO DKV0Md']/ancestor::a/@href`)
		path_publ  := xmlpath.MustCompile(`.//span[@class='VuuXrf']`)
		path_desc  := xmlpath.MustCompile(`.//div[@class='VwiC3b yXK7lf lyLwlc yDYNvb W8l4ac lEBKkf']`)


		for it := path_nodes.Iter(root); it.Next(); {
			var title       string
			var link        string
			var publisher   string
			var description string

			title, _ = path_title.String(it.Node())
			title    = Unescape(title)

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)

			publisher, _ = path_publ.String(it.Node())
			publisher    = Unescape(publisher)

			description, _ = path_desc.String(it.Node())
			description    = Unescape(description)
			
			value := map[string]string{
				"title"      : title,
				"link"       : link,
				"publisher"  : publisher,
				"description": description,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("search_result.%d", counter), value)
			counter += 1
		}

		//// 2
		path_nodes = xmlpath.MustCompile(`//*[@class='g Ww4FFb vt6azd tF2Cxc asEBEc']`)

		path_title = xmlpath.MustCompile(`.//h3[@class='LC20lb MBeuO DKV0Md']/text()`)
		path_link  = xmlpath.MustCompile(`.//h3[@class='LC20lb MBeuO DKV0Md']/ancestor::a/@href`)
		path_publ  = xmlpath.MustCompile(`.//span[@class='VuuXrf']`)
		path_desc  = xmlpath.MustCompile(`.//div[@class='VwiC3b yXK7lf lyLwlc yDYNvb W8l4ac']`)
		

		for it := path_nodes.Iter(root); it.Next(); {
			var title       string
			var link        string
			var publisher   string
			var description string

			title, _ = path_title.String(it.Node())
			title    = Unescape(title)

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)

			publisher, _ = path_publ.String(it.Node())
			publisher    = Unescape(publisher)

			description, _ = path_desc.String(it.Node())
			description    = Unescape(description)

			value := map[string]string{
				"title"      : title,
				"link"       : link,
				"publisher"  : publisher,
				"description": description,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("search_result.%d", counter), value)
			counter += 1
		}

		//// 3
		path_nodes = xmlpath.MustCompile(`//div[@class="eKjLze"]`)

		path_title = xmlpath.MustCompile(`.//h3[@class='LC20lb MBeuO DKV0Md']/text()`)
		path_link  = xmlpath.MustCompile(`.//h3[@class='LC20lb MBeuO DKV0Md']/ancestor::a/@href`)
		path_publ  = xmlpath.MustCompile(`.//span[@class='VuuXrf']`)
		path_desc  = xmlpath.MustCompile(`.//div[@class="VwiC3b yXK7lf lyLwlc yDYNvb W8l4ac lEBKkf"]`)

		for it := path_nodes.Iter(root); it.Next(); {
			var title       string
			var link        string
			var publisher   string
			var description string

			title, _ = path_title.String(it.Node())
			title    = Unescape(title)

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)

			publisher, _ = path_publ.String(it.Node())
			publisher    = Unescape(publisher)

			description, _ = path_desc.String(it.Node())
			description    = Unescape(description)

			value := map[string]string{
				"title"      : title,
				"link"       : link,
				"publisher"  : publisher,
				"description": description,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("search_result3.%d", counter), value)
			counter += 1
		}
	}

	{ //// featured_snippets
		counter    := 0
		path_nodes := xmlpath.MustCompile(`//div[@class="eKjLze"]`)

		path_title  := xmlpath.MustCompile(`.//a`)
		path_link   := xmlpath.MustCompile(`.//a/@href`)
		path_desc   := xmlpath.MustCompile(`.//div[@class="zz3gNc"]`)
		

		for it := path_nodes.Iter(root); it.Next(); {
			var title       string
			var link        string
			var description string

			title, _ = path_title.String(it.Node())
			title    = Unescape(title)

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)

			description, _ = path_desc.String(it.Node())
			description    = Unescape(description)

			value := map[string]string{
				"title"      : title,
				"link"       : link,
				"description": description,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("search_result.%d", counter), value)
			counter += 1
		}
	}

	{ //// sidebar_result
		var title       string
		var link        string
		var description string

		path_node := xmlpath.MustCompile(`//div[@class="kp-wholepage ss6qqb u7yw9 zLsiYe mnr-c UBoxCb kp-wholepage-osrp Jb0Zif EyBRub"]`)
		node_iter := path_node.Iter(root)

		if node_iter.Next() {
			path_title  := xmlpath.MustCompile(`.//span[@class="yKMVIe"]`)
			path_link   := xmlpath.MustCompile(`.//a[@class="ruhjFe NJLBac fl"]`)

			title, _ = path_title.String(node_iter.Node())
			title    = Unescape(title)

			link, _ = path_link.String(node_iter.Node())
			link    = Unescape(link)
		}

		output, _ = sjson.Set(output, "sidebar_result", map[string]string{
			"title"      : title,
			"link"       : link,
			"description": description,
		})

	}

	{ //// people_also_ask
		counter    := 0
		path_nodes := xmlpath.MustCompile(`//div[@class="Wt5Tfe"]//div[@class="dnXCYb"]`)

		path_question := xmlpath.MustCompile(`.//div[@class="L3Ezfd"]`)

		for it := path_nodes.Iter(root); it.Next(); {
			var question    string

			question, _ = path_question.String(it.Node())
			question    = Unescape(question)

			value := map[string]string{
				"question": question,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("people_also_ask.%d", counter), value)
			counter += 1
		}
	}

	{ //// related_searches
		counter    := 0
		path_nodes := xmlpath.MustCompile(`//a[@class="k8XOCe R0xfCb VCOFK s8bAkb"]`)

		path_title := xmlpath.MustCompile(`.//div[@class="s75CSd u60jwe r2fjmd AB4Wff"]`)
		path_link  := xmlpath.MustCompile(`.//@href`)


		for it := path_nodes.Iter(root); it.Next(); {
			var title string
			var link  string

			title, _ = path_title.String(it.Node())
			title    = Unescape(title)

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)

			value := map[string]string{
				"title": title,
				"link" : link,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("related_searches.%d", counter), value)
			counter += 1
		}
	}

	{ //// videos
		counter    := 0
		path_nodes := xmlpath.MustCompile(`//div[@jsname="pKB8Bc"]`)

		path_title   := xmlpath.MustCompile(`.//span[@class="cHaqb"]`)
		path_link    := xmlpath.MustCompile(`.//a[@class="X5OiLe"]`)
		path_channel := xmlpath.MustCompile(`.//span[@class='pcJO7e']/span`)

		for it := path_nodes.Iter(root); it.Next(); {
			var title	   string
			var link	   string
			var channel	   string

			title, _ = path_title.String(it.Node())
			title    = Unescape(title)

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)

			channel, _ = path_channel.String(it.Node())
			channel    = Unescape(channel)

			value := map[string]string{
				"title"  : title,
				"link"   : link,
				"channel": channel,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("videos.%d", counter), value)
			counter += 1
		}

		//// 2
		path_nodes = xmlpath.MustCompile(`//div[@jscontroller="rTuANe"]`)

		path_title   = xmlpath.MustCompile(`.//h3[@class="LC20lb MBeuO DKV0Md"]`)
		path_link    = xmlpath.MustCompile(`.//a[@jsname="UWckNb"]`)
		path_channel = xmlpath.MustCompile(`.//div[@class="gqF9jc"]/span[2]`)

		for it := path_nodes.Iter(root); it.Next(); {
			var title	   string
			var link	   string
			var channel	   string

			title, _ = path_title.String(it.Node())
			title    = Unescape(title)

			link, _ = path_link.String(it.Node())
			link    = Unescape(link)

			channel, _ = path_channel.String(it.Node())
			channel    = Unescape(channel)

			value := map[string]string{
				"title"  : title,
				"link"   : link,
				"channel": channel,
			}

			output, _ = sjson.Set(output, fmt.Sprintf("videos.%d", counter), value)
			counter += 1
		}
	}
	return output, nil
}

// ------------------------------------------------------------
// : Function
// ------------------------------------------------------------
func Parse(url string, html string) (string, error) {
	var err    error
	var output string

	switch true {
		case regexp.MustCompile("duckduckgo")     .MatchString(url): output, err = ParseDuckDuckGo(html)   //// DuckDuckGo
		case regexp.MustCompile("youtube")        .MatchString(url): output, err = ParseYouTube(html)      //// YouTube
		case regexp.MustCompile("bing")           .MatchString(url): output, err = ParseBing(html)         //// Bing
		case regexp.MustCompile("yahoo")          .MatchString(url): output, err = ParseYahoo(html)        //// Yahoo
		case regexp.MustCompile("google.+tbm=nws").MatchString(url): output, err = ParseGoogleNews(html)   //// Google News
		case regexp.MustCompile("google.+tbm=vid").MatchString(url): output, err = ParseGoogleVideos(html) //// Google Videos
		case regexp.MustCompile("google")         .MatchString(url): output, err = ParseGoogle(html)       //// Google
	}

	if err != nil {
		logger.Error().
			Err(err).
			Str("url", url).
			Msg("Failed to parse")
		return "", err
	}

	return output, nil
}
