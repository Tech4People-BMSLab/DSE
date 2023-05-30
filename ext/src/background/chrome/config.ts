export const local_config = {
	"version": "1.0.0",
	"keywords": [

        "Energietransitie",
        "Nam",
        "Windmolen",
        "Groningen",
        "Zonnepanelen",
        "Groene energie",
        "Aardgasvrij",
        "Stikstof",
        "Landbouw",
        "Eurocommissaris",
        "Bouwsector",
        "Intensieve veehouderij",
        "Huisvestingscrisis",
        "Stilleggen bouw",
        "BBB",
        "Immigratie",
        "Asielzoekers",
        "Oekraïne",
        "Gelukzoekers",
        "Woningtekort",
        "Woningnood",
        "Verkiezingen",
        "Verkiezingsuitslag",
        "Vertrouwen",
        "Provinciale Staten",
        "Burgerberaad",
        "Eerste kamer",
        "Kabinet",
        "Klimaat",
        "Wil je meer werken, laat het merken",
        "Polarisatie",
        "Grensoverschrijdend gedrag",
        "Arbeidsmarkt",
        "ChatGPT",
        "Kunstmatige intelligentie",
        "Pauze",
        "Onderwijs",
        "Lager opgeleid",
        "Hoger opgeleid",
        "Extinction rebellion",
	],
	"alarm_config": {
		"verify": {
            "period"        : "60000",
            "datetime_start": "2020-01-01T00:00:00.000Z",
            "datetime_end"  : "2020-12-31T23:59:59.999Z"
        },
        "search": {
            "period"        : "60000",
            "datetime_start": "2020-01-01T00:00:00.000Z",
            "datetime_end"  : "2020-12-31T23:59:59.999Z"
        }
	},
	"content_config": {
		"allowed_domains": [
			"google.com",
			"google.nl",
			"google.de",
			"google.be",

			"youtube.com",
			"youtube.nl",
			"youtube.de",
			"youtube.be",

			"duckduckgo.com",
			"duckduckgo.nl",
			"duckduckgo.de",
			"duckduckgo.be",

			"twitter.com"
		]
	},
	"mousecapture_config": {
		"query_parameters": [
			{
				"name": "google",
				"value": "q"
			},
			{
				"name": "youtube",
				"value": "search_query"
			},
			{
				"name": "duckduckgo",
				"value": "q"
			},
			{
				"name": "twitter",
				"value": "q"
			}
		]
	},
	"extractor_config": {
		"google.com/search?tbm=nws&q=": {
			"name"   : "Google News",
			"selectors": [
				{
					"name"  : "search_result",
					"xpath" : "//*[@id=\"rso\"]/div/div/div/div/a/div",
					"xpaths": {
						"publisher"  : ".//*[@class=\"CEMjEf NUnG9d\"]//span/text()",
						"title"      : ".//*[@class=\"mCBkyc y355M ynAwRc MBeuO nDgy9d\"]/text()",
						"description": ".//*[@class=\"GI74Re nDgy9d\"]/text()",
						"time"       : ".//*[@class=\"OSrXXb ZE0LJd\"]//span/text()",
						"link"       : ".//*[@class=\"mCBkyc y355M ynAwRc MBeuO nDgy9d\"]/ancestor::a[@href]/@href"
					}
				}
			]
		},
		"google.nl/search?tbm=nws&q=": {
			"name"     : "Google News",
			"selectors": [
				{
					"name"  : "search_result",
					"xpath" : "//*[@id=\"rso\"]/div/div/div/div/a/div",
					"xpaths": {
						"publisher"  : ".//*[@class=\"CEMjEf NUnG9d\"]//span/text()",
						"title"      : ".//*[@class=\"mCBkyc y355M ynAwRc MBeuO nDgy9d\"]/text()",
						"description": ".//*[@class=\"GI74Re nDgy9d\"]/text()",
						"time"       : ".//*[@class=\"OSrXXb ZE0LJd\"]//span/text()",
						"link"       : ".//*[@class=\"mCBkyc y355M ynAwRc MBeuO nDgy9d\"]/ancestor::a[@href]/@href"
					}
				}
			]
		},
		"google.be/search?tbm=nws&q=": {
			"name"   : "Google News",
			"selectors": [
				{
					"name" : "search_result",
					"xpath": "//*[@id=\"rso\"]/div/div/div/div/a/div",
					"xpaths": {
						"publisher"  : ".//*[@class=\"CEMjEf NUnG9d\"]//span/text()",
						"title"      : ".//*[@class=\"mCBkyc y355M ynAwRc MBeuO nDgy9d\"]/text()",
						"description": ".//*[@class=\"GI74Re nDgy9d\"]/text()",
						"time"       : ".//*[@class=\"OSrXXb ZE0LJd\"]//span/text()",
						"link"       : ".//*[@class=\"mCBkyc y355M ynAwRc MBeuO nDgy9d\"]/ancestor::a[@href]/@href"
					}
				}
			]
		},
		"google.de/search?tbm=nws&q=": {
			"name"     : "Google News",
			"selectors": [
				{
					"name": "search_result",
					"xpath": "//*[@id=\"rso\"]/div/div/div/div/a/div",
					"xpaths": {
						"publisher": ".//*[@class=\"CEMjEf NUnG9d\"]//span/text()",
						"title": ".//*[@class=\"mCBkyc y355M ynAwRc MBeuO nDgy9d\"]/text()",
						"description": ".//*[@class=\"GI74Re nDgy9d\"]/text()",
						"time": ".//*[@class=\"OSrXXb ZE0LJd\"]//span/text()",
						"link": ".//*[@class=\"mCBkyc y355M ynAwRc MBeuO nDgy9d\"]/ancestor::a[@href]/@href"
					}
				}
			]
		},
		"google.com/search?tbm=vid&q=": {
			"name"  : "Google Videos",
			"selectors": [
				{
					"name": "search_result",
					"xpath": "//*[@id=\"rso\"]/div/div/video-voyager",
					"xpaths": {
						"publisher": ".//div[@class=\"TbwUpd NJjxre\"]/cite[@class=\"iUh30 qLRx3b tjvcx\"]/text()",
						"title": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/text()",
						"link": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/ancestor::a[@href]/@href",
						"description": ".//*[@class=\"Uroaid\"]//text()",
						"subtitle": ".//*[@class=\"P7xzyf\"]/span//text()"
					}
				}
			]
		},
		"google.nl/search?tbm=vid&q=": {
			"name"  : "Google Videos",
			"selectors": [
				{
					"name": "search_result",
					"xpath": "//*[@id=\"rso\"]/div/div/video-voyager",
					"xpaths": {
						"publisher": ".//div[@class=\"TbwUpd NJjxre\"]/cite[@class=\"iUh30 qLRx3b tjvcx\"]/text()",
						"title": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/text()",
						"link": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/ancestor::a[@href]/@href",
						"description": ".//*[@class=\"Uroaid\"]//text()",
						"subtitle": ".//*[@class=\"P7xzyf\"]/span//text()"
					}
				}
			]
		},
		"google.be/search?tbm=vid&q=": {
			"name"  : "Google Videos",
			"selectors": [
				{
					"name": "search_result",
					"xpath": "//*[@id=\"rso\"]/div/div/video-voyager",
					"xpaths": {
						"publisher": ".//div[@class=\"TbwUpd NJjxre\"]/cite[@class=\"iUh30 qLRx3b tjvcx\"]/text()",
						"title": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/text()",
						"link": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/ancestor::a[@href]/@href",
						"description": ".//*[@class=\"Uroaid\"]//text()",
						"subtitle": ".//*[@class=\"P7xzyf\"]/span//text()"
					}
				}
			]
		},
		"google.de/search?tbm=vid&q=": {
			"name"  : "Google Videos",
			"selectors": [
				{
					"name": "search_result",
					"xpath": "//*[@id=\"rso\"]/div/div/video-voyager",
					"xpaths": {
						"publisher": ".//div[@class=\"TbwUpd NJjxre\"]/cite[@class=\"iUh30 qLRx3b tjvcx\"]/text()",
						"title": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/text()",
						"link": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/ancestor::a[@href]/@href",
						"description": ".//*[@class=\"Uroaid\"]//text()",
						"subtitle": ".//*[@class=\"P7xzyf\"]/span//text()"
					}
				}
			]
		},
		"google.com/search?q=": {
			"name": "Google",
			"selectors": [
				{
					"name": "search_result",
					"xpath": "//*[@id=\"search\"]//*[@class=\"g Ww4FFb vt6azd tF2Cxc\"]",
					"xpaths": {
						"publisher": ".//div[@class=\"TbwUpd NJjxre\"]//*[@class=\"iUh30 qLRx3b tjvcx\" or @class=\"iUh30 tjvcx\"]/text()[1]",
						"title": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]//text()",
						"description": ".//*[@class=\"VwiC3b yXK7lf MUxGbd yDYNvb lyLwlc lEBKkf\"]/span[last()]//text()",
						"link": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/parent::a/@href",
						"date": ".//*[@class=\"MUxGbd wuQ4Ob WZ8Tjf\"]/span/text()"
					}
				},
				{
					"name": "featured_result",
					"xpath": "//*[@class=\"BYM4Nd\"]",
					"xpaths": {
						"publisher": ".//div[@class=\"TbwUpd NJjxre\"]//*[@class=\"iUh30 qLRx3b tjvcx\" or @class=\"iUh30 tjvcx\"]/text()[1]",
						"title": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/text()"
					}
				},
				{
					"name": "featured_result_links",
					"xpath": "//*[@class=\"BYM4Nd\"]//*[@class=\"usJj9c\"]",
					"xpaths": {
						"title": ".//h3//text()",
						"link": ".//a/@href",
						"description": ".//div//text()"
					}
				},
				{
					"name": "sidebar_result",
					"xpath": "//*[@class=\"I6TXqe\"]",
					"xpaths": {
						"title": ".//*[@class=\"qrShPb kno-ecr-pt PZPZlf HOpgu q8U8x hNKfZe\"]//text()",
						"link": "//*[@class=\"wDYxhc NFQFxe\"]/div/a/@href"
					}
				},
				{
					"name": "people_also_searched",
					"xpath": "//*[@class=\"zVvuGd MRfBrb\"]/div",
					"xpaths": {
						"title": ".//a/@title",
						"link": ".//a/@href"
					}
				},
				{
					"name": "see_results_about",
					"xpath": "//*[@class=\"EfDVh wDYxhc NFQFxe\"]",
					"xpaths": {
						"title": ".//*[@class=\"RJn8N ellip tNxQIb ynAwRc\"]/text()",
						"link": ".//a/@href"
					}
				},
				{
					"name": "snippet_result",
					"xpath": "//*[@id=\"Odp5De\"]",
					"xpaths": {
						"text": ".//*[@class=\"hgKElc\"]//text()",
						"publisher": ".//*[@class=\"TbwUpd NJjxre\"]//*[@class=\"iUh30 qLRx3b tjvcx\"]/text()[1]",
						"title": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/text()",
						"link": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/ancestor::a/@href"
					}
				},
				{
					"name": "people_also_asked",
					"xpath": "//div[@class=\"iDjcJe IX9Lgd wwB5gf\"]",
					"xpaths": {
						"question": ".//span/text()"
					}
				},
				{
					"name": "top_stories",
					"xpath": "//a[@class=\"WlydOe\"]",
					"xpaths": {
						"publisher": ".//*[@class=\"CEMjEf NUnG9d\"]//text()",
						"title": ".//*[@class=\"mCBkyc tNxQIb ynAwRc nDgy9d\"]/text()",
						"link": ".//@href",
						"time": ".//*[@class=\"OSrXXb ZE0LJd\"]//text()"
					}
				},
				{
					"name": "related_searches",
					"xpath": "//*[@class=\"k8XOCe R0xfCb VCOFK s8bAkb\"]",
					"xpaths": {
						"title": ".//*[@class=\"s75CSd OhScic AB4Wff\"]//text()",
						"link": ".//@href"
					}
				},
				{
					"name": "maps_locations",
					"xpath": "//*[@class=\"yYlJEf VByer Q7PwXb VDgVie\"]/ancestor::div[@class=\"VkpGBb\"]",
					"xpaths": {
						"title": ".//*[@class=\"OSrXXb\"]/text()",
						"link": ".//*[@class=\"yYlJEf VByer Q7PwXb VDgVie\"]/@data-url",
						"properties": ".//*[@class=\"rllt__details\"]/div[position()>1]/text()",
						"rated": ".//*[@class=\"MvDXgc\"]//*[@aria-label]/@aria-label",
						"reviews_count": ".//*[@class=\"HypWnf YrbPuc\"]/text()"
					}
				},
				{
					"name": "videos",
					"xpath": "//*[@class=\"X5OiLe\"]",
					"xpaths": {
						"title": ".//*[@class=\"cHaqb\"]//text()",
						"host": ".//*[@class=\"pcJO7e\"]/cite/text()",
						"publisher": ".//*[@class=\"pcJO7e\"]/span/text()",
						"link": ".//*[@class=\"cHaqb\"]//text()/ancestor::a/@href",
						"time": ".//*[@class=\"hMJ0yc\"]/span/text()"
					}
				}
			]
		},
		"google.nl/search?q=": {
			"name": "Google",
			"selectors": [
				{
					"name": "search_result",
					"xpath": "//*[@id=\"search\"]//*[@class=\"g Ww4FFb vt6azd tF2Cxc\"]",
					"xpaths": {
						"publisher": ".//div[@class=\"TbwUpd NJjxre\"]//*[@class=\"iUh30 qLRx3b tjvcx\" or @class=\"iUh30 tjvcx\"]/text()[1]",
						"title": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]//text()",
						"description": ".//*[@class=\"VwiC3b yXK7lf MUxGbd yDYNvb lyLwlc lEBKkf\"]/span[last()]//text()",
						"link": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/parent::a/@href",
						"date": ".//*[@class=\"MUxGbd wuQ4Ob WZ8Tjf\"]/span/text()"
					}
				},
				{
					"name": "featured_result",
					"xpath": "//*[@class=\"BYM4Nd\"]",
					"xpaths": {
						"publisher": ".//div[@class=\"TbwUpd NJjxre\"]//*[@class=\"iUh30 qLRx3b tjvcx\" or @class=\"iUh30 tjvcx\"]/text()[1]",
						"title": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/text()"
					}
				},
				{
					"name": "featured_result_links",
					"xpath": "//*[@class=\"BYM4Nd\"]//*[@class=\"usJj9c\"]",
					"xpaths": {
						"title": ".//h3//text()",
						"link": ".//a/@href",
						"description": ".//div//text()"
					}
				},
				{
					"name": "sidebar_result",
					"xpath": "//*[@class=\"I6TXqe\"]",
					"xpaths": {
						"title": ".//*[@class=\"qrShPb kno-ecr-pt PZPZlf HOpgu q8U8x hNKfZe\"]//text()",
						"link": "//*[@class=\"wDYxhc NFQFxe\"]/div/a/@href"
					}
				},
				{
					"name": "people_also_searched",
					"xpath": "//*[@class=\"zVvuGd MRfBrb\"]/div",
					"xpaths": {
						"title": ".//a/@title",
						"link": ".//a/@href"
					}
				},
				{
					"name": "see_results_about",
					"xpath": "//*[@class=\"EfDVh wDYxhc NFQFxe\"]",
					"xpaths": {
						"title": ".//*[@class=\"RJn8N ellip tNxQIb ynAwRc\"]/text()",
						"link": ".//a/@href"
					}
				},
				{
					"name": "snippet_result",
					"xpath": "//*[@id=\"Odp5De\"]",
					"xpaths": {
						"text": ".//*[@class=\"hgKElc\"]//text()",
						"publisher": ".//*[@class=\"TbwUpd NJjxre\"]//*[@class=\"iUh30 qLRx3b tjvcx\"]/text()[1]",
						"title": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/text()",
						"link": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/ancestor::a/@href"
					}
				},
				{
					"name": "people_also_asked",
					"xpath": "//div[@class=\"iDjcJe IX9Lgd wwB5gf\"]",
					"xpaths": {
						"question": ".//span/text()"
					}
				},
				{
					"name": "top_stories",
					"xpath": "//a[@class=\"WlydOe\"]",
					"xpaths": {
						"publisher": ".//*[@class=\"CEMjEf NUnG9d\"]//text()",
						"title": ".//*[@class=\"mCBkyc tNxQIb ynAwRc nDgy9d\"]/text()",
						"link": ".//@href",
						"time": ".//*[@class=\"OSrXXb ZE0LJd\"]//text()"
					}
				},
				{
					"name": "related_searches",
					"xpath": "//*[@class=\"k8XOCe R0xfCb VCOFK s8bAkb\"]",
					"xpaths": {
						"title": ".//*[@class=\"s75CSd OhScic AB4Wff\"]//text()",
						"link": ".//@href"
					}
				},
				{
					"name": "maps_locations",
					"xpath": "//*[@class=\"yYlJEf VByer Q7PwXb VDgVie\"]/ancestor::div[@class=\"VkpGBb\"]",
					"xpaths": {
						"title": ".//*[@class=\"OSrXXb\"]/text()",
						"link": ".//*[@class=\"yYlJEf VByer Q7PwXb VDgVie\"]/@data-url",
						"properties": ".//*[@class=\"rllt__details\"]/div[position()>1]/text()",
						"rated": ".//*[@class=\"MvDXgc\"]//*[@aria-label]/@aria-label",
						"reviews_count": ".//*[@class=\"HypWnf YrbPuc\"]/text()"
					}
				},
				{
					"name": "videos",
					"xpath": "//*[@class=\"X5OiLe\"]",
					"xpaths": {
						"title": ".//*[@class=\"cHaqb\"]//text()",
						"host": ".//*[@class=\"pcJO7e\"]/cite/text()",
						"publisher": ".//*[@class=\"pcJO7e\"]/span/text()",
						"link": ".//*[@class=\"cHaqb\"]//text()/ancestor::a/@href",
						"time": ".//*[@class=\"hMJ0yc\"]/span/text()"
					}
				}
			]
		},
		"google.be/search?q=": {
			"name": "Google",
			"selectors": [
				{
					"name": "search_result",
					"xpath": "//*[@id=\"search\"]//*[@class=\"g Ww4FFb vt6azd tF2Cxc\"]",
					"xpaths": {
						"publisher": ".//div[@class=\"TbwUpd NJjxre\"]//*[@class=\"iUh30 qLRx3b tjvcx\" or @class=\"iUh30 tjvcx\"]/text()[1]",
						"title": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]//text()",
						"description": ".//*[@class=\"VwiC3b yXK7lf MUxGbd yDYNvb lyLwlc lEBKkf\"]/span[last()]//text()",
						"link": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/parent::a/@href",
						"date": ".//*[@class=\"MUxGbd wuQ4Ob WZ8Tjf\"]/span/text()"
					}
				},
				{
					"name": "featured_result",
					"xpath": "//*[@class=\"BYM4Nd\"]",
					"xpaths": {
						"publisher": ".//div[@class=\"TbwUpd NJjxre\"]//*[@class=\"iUh30 qLRx3b tjvcx\" or @class=\"iUh30 tjvcx\"]/text()[1]",
						"title": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/text()"
					}
				},
				{
					"name": "featured_result_links",
					"xpath": "//*[@class=\"BYM4Nd\"]//*[@class=\"usJj9c\"]",
					"xpaths": {
						"title": ".//h3//text()",
						"link": ".//a/@href",
						"description": ".//div//text()"
					}
				},
				{
					"name": "sidebar_result",
					"xpath": "//*[@class=\"I6TXqe\"]",
					"xpaths": {
						"title": ".//*[@class=\"qrShPb kno-ecr-pt PZPZlf HOpgu q8U8x hNKfZe\"]//text()",
						"link": "//*[@class=\"wDYxhc NFQFxe\"]/div/a/@href"
					}
				},
				{
					"name": "people_also_searched",
					"xpath": "//*[@class=\"zVvuGd MRfBrb\"]/div",
					"xpaths": {
						"title": ".//a/@title",
						"link": ".//a/@href"
					}
				},
				{
					"name": "see_results_about",
					"xpath": "//*[@class=\"EfDVh wDYxhc NFQFxe\"]",
					"xpaths": {
						"title": ".//*[@class=\"RJn8N ellip tNxQIb ynAwRc\"]/text()",
						"link": ".//a/@href"
					}
				},
				{
					"name": "snippet_result",
					"xpath": "//*[@id=\"Odp5De\"]",
					"xpaths": {
						"text": ".//*[@class=\"hgKElc\"]//text()",
						"publisher": ".//*[@class=\"TbwUpd NJjxre\"]//*[@class=\"iUh30 qLRx3b tjvcx\"]/text()[1]",
						"title": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/text()",
						"link": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/ancestor::a/@href"
					}
				},
				{
					"name": "people_also_asked",
					"xpath": "//div[@class=\"iDjcJe IX9Lgd wwB5gf\"]",
					"xpaths": {
						"question": ".//span/text()"
					}
				},
				{
					"name": "top_stories",
					"xpath": "//a[@class=\"WlydOe\"]",
					"xpaths": {
						"publisher": ".//*[@class=\"CEMjEf NUnG9d\"]//text()",
						"title": ".//*[@class=\"mCBkyc tNxQIb ynAwRc nDgy9d\"]/text()",
						"link": ".//@href",
						"time": ".//*[@class=\"OSrXXb ZE0LJd\"]//text()"
					}
				},
				{
					"name": "related_searches",
					"xpath": "//*[@class=\"k8XOCe R0xfCb VCOFK s8bAkb\"]",
					"xpaths": {
						"title": ".//*[@class=\"s75CSd OhScic AB4Wff\"]//text()",
						"link": ".//@href"
					}
				},
				{
					"name": "maps_locations",
					"xpath": "//*[@class=\"yYlJEf VByer Q7PwXb VDgVie\"]/ancestor::div[@class=\"VkpGBb\"]",
					"xpaths": {
						"title": ".//*[@class=\"OSrXXb\"]/text()",
						"link": ".//*[@class=\"yYlJEf VByer Q7PwXb VDgVie\"]/@data-url",
						"properties": ".//*[@class=\"rllt__details\"]/div[position()>1]/text()",
						"rated": ".//*[@class=\"MvDXgc\"]//*[@aria-label]/@aria-label",
						"reviews_count": ".//*[@class=\"HypWnf YrbPuc\"]/text()"
					}
				},
				{
					"name": "videos",
					"xpath": "//*[@class=\"X5OiLe\"]",
					"xpaths": {
						"title": ".//*[@class=\"cHaqb\"]//text()",
						"host": ".//*[@class=\"pcJO7e\"]/cite/text()",
						"publisher": ".//*[@class=\"pcJO7e\"]/span/text()",
						"link": ".//*[@class=\"cHaqb\"]//text()/ancestor::a/@href",
						"time": ".//*[@class=\"hMJ0yc\"]/span/text()"
					}
				}
			]
		},
		"google.de/search?q=": {
			"name": "Google",
			"selectors": [
				{
					"name": "search_result",
					"xpath": "//*[@id=\"search\"]//*[@class=\"g Ww4FFb vt6azd tF2Cxc\"]",
					"xpaths": {
						"publisher": ".//div[@class=\"TbwUpd NJjxre\"]//*[@class=\"iUh30 qLRx3b tjvcx\" or @class=\"iUh30 tjvcx\"]/text()[1]",
						"title": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]//text()",
						"description": ".//*[@class=\"VwiC3b yXK7lf MUxGbd yDYNvb lyLwlc lEBKkf\"]/span[last()]//text()",
						"link": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/parent::a/@href",
						"date": ".//*[@class=\"MUxGbd wuQ4Ob WZ8Tjf\"]/span/text()"
					}
				},
				{
					"name": "featured_result",
					"xpath": "//*[@class=\"BYM4Nd\"]",
					"xpaths": {
						"publisher": ".//div[@class=\"TbwUpd NJjxre\"]//*[@class=\"iUh30 qLRx3b tjvcx\" or @class=\"iUh30 tjvcx\"]/text()[1]",
						"title": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/text()"
					}
				},
				{
					"name": "featured_result_links",
					"xpath": "//*[@class=\"BYM4Nd\"]//*[@class=\"usJj9c\"]",
					"xpaths": {
						"title": ".//h3//text()",
						"link": ".//a/@href",
						"description": ".//div//text()"
					}
				},
				{
					"name": "sidebar_result",
					"xpath": "//*[@class=\"I6TXqe\"]",
					"xpaths": {
						"title": ".//*[@class=\"qrShPb kno-ecr-pt PZPZlf HOpgu q8U8x hNKfZe\"]//text()",
						"link": "//*[@class=\"wDYxhc NFQFxe\"]/div/a/@href"
					}
				},
				{
					"name": "people_also_searched",
					"xpath": "//*[@class=\"zVvuGd MRfBrb\"]/div",
					"xpaths": {
						"title": ".//a/@title",
						"link": ".//a/@href"
					}
				},
				{
					"name": "see_results_about",
					"xpath": "//*[@class=\"EfDVh wDYxhc NFQFxe\"]",
					"xpaths": {
						"title": ".//*[@class=\"RJn8N ellip tNxQIb ynAwRc\"]/text()",
						"link": ".//a/@href"
					}
				},
				{
					"name": "snippet_result",
					"xpath": "//*[@id=\"Odp5De\"]",
					"xpaths": {
						"text": ".//*[@class=\"hgKElc\"]//text()",
						"publisher": ".//*[@class=\"TbwUpd NJjxre\"]//*[@class=\"iUh30 qLRx3b tjvcx\"]/text()[1]",
						"title": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/text()",
						"link": ".//*[@class=\"LC20lb MBeuO DKV0Md\"]/ancestor::a/@href"
					}
				},
				{
					"name": "people_also_asked",
					"xpath": "//div[@class=\"iDjcJe IX9Lgd wwB5gf\"]",
					"xpaths": {
						"question": ".//span/text()"
					}
				},
				{
					"name": "top_stories",
					"xpath": "//a[@class=\"WlydOe\"]",
					"xpaths": {
						"publisher": ".//*[@class=\"CEMjEf NUnG9d\"]//text()",
						"title": ".//*[@class=\"mCBkyc tNxQIb ynAwRc nDgy9d\"]/text()",
						"link": ".//@href",
						"time": ".//*[@class=\"OSrXXb ZE0LJd\"]//text()"
					}
				},
				{
					"name": "related_searches",
					"xpath": "//*[@class=\"k8XOCe R0xfCb VCOFK s8bAkb\"]",
					"xpaths": {
						"title": ".//*[@class=\"s75CSd OhScic AB4Wff\"]//text()",
						"link": ".//@href"
					}
				},
				{
					"name": "maps_locations",
					"xpath": "//*[@class=\"yYlJEf VByer Q7PwXb VDgVie\"]/ancestor::div[@class=\"VkpGBb\"]",
					"xpaths": {
						"title": ".//*[@class=\"OSrXXb\"]/text()",
						"link": ".//*[@class=\"yYlJEf VByer Q7PwXb VDgVie\"]/@data-url",
						"properties": ".//*[@class=\"rllt__details\"]/div[position()>1]/text()",
						"rated": ".//*[@class=\"MvDXgc\"]//*[@aria-label]/@aria-label",
						"reviews_count": ".//*[@class=\"HypWnf YrbPuc\"]/text()"
					}
				},
				{
					"name": "videos",
					"xpath": "//*[@class=\"X5OiLe\"]",
					"xpaths": {
						"title": ".//*[@class=\"cHaqb\"]//text()",
						"host": ".//*[@class=\"pcJO7e\"]/cite/text()",
						"publisher": ".//*[@class=\"pcJO7e\"]/span/text()",
						"link": ".//*[@class=\"cHaqb\"]//text()/ancestor::a/@href",
						"time": ".//*[@class=\"hMJ0yc\"]/span/text()"
					}
				}
			]
		},
		"duckduckgo.com/?q=": {
			"name": "DuckDuckGo",
			"selectors": [
				{
					"name": "search_result",
					"xpath": "//article",
					"xpaths": {
						"publisher": ".//*[@class=\"Wo6ZAEmESLNUuWBkbMxx\"]/text()[2]",
						"title": ".//*[@class=\"EKtkFWMYpwzMKOYr0GYm LQVY1Jpkk8nyJ6HBWKAk\"]/text()",
						"link": ".//*[@class=\"EKtkFWMYpwzMKOYr0GYm LQVY1Jpkk8nyJ6HBWKAk\"]/ancestor::a[@href]/@href",
						"description": ".//*[@class=\"OgdwYG6KE2qthn9XQWFC\"]//text()"
					}
				},
				{
					"name": "related_result",
					"xpath": ".//*[@class=\"related-searches__item\"]",
					"xpaths": {
						"title": ".//span//text()",
						"link": ".//a/@href"
					}
				},
				{
					"name": "sidebar_result",
					"xpath": "//*[@class=\"module__content js-about-module-content\"]",
					"xpaths": {
						"title": ".//*[@class=\"module__title__link\"]/text()",
						"link": ".//*[@class=\"module__official-url js-about-item-link\"]/@href",
						"description": ".//*[@class=\"module__text\"]//text()",
						"related_links": ".//*[@class=\"about-profiles__link js-about-profile-link\"]/@href"
					}
				},
				{
					"name": "recent_news",
					"xpath": "//*[@class=\"module--carousel__item has-image\"]",
					"xpaths": {
						"title": ".//a/@title",
						"link": ".//a/@href",
						"publisher": ".//*[@class=\"module--carousel__source result__url\"]/text()",
						"time": ".//*[@class=\"tile__time\"]/text()"
					}
				}
			]
		},
		"youtube.com/results?search_query=": {
			"name": "YouTube",
			"selectors": [
				{
					"name": "search_result",
					"xpath": "//ytd-video-renderer",
					"xpaths": {
						"title": ".//*[@id=\"video-title\"]/yt-formatted-string/text()",
						"views": ".//*[@id=\"metadata-line\"]/span[1]/text()",
						"time": ".//*[@id=\"metadata-line\"]/span[2]/text()",
						"channel_name": ".//*[@id=\"text\"]/a/text()",
						"channel_url": ".//*[@id=\"text\"]/a/@href",
						"description": ".//*[@id=\"dismissible\"]/div/div[3]/yt-formatted-string//text()",
						"badge": ".//*[@id=\"badges\"]/div//span/text()"
					}
				}
			]
		},
		"youtube.nl/results?search_query=": {
			"name": "YouTube",
			"selectors": [
				{
					"name": "search_result",
					"xpath": "//ytd-video-renderer",
					"xpaths": {
						"title": ".//*[@id=\"video-title\"]/yt-formatted-string/text()",
						"views": ".//*[@id=\"metadata-line\"]/span[1]/text()",
						"time": ".//*[@id=\"metadata-line\"]/span[2]/text()",
						"channel_name": ".//*[@id=\"text\"]/a/text()",
						"channel_url": ".//*[@id=\"text\"]/a/@href",
						"description": ".//*[@id=\"dismissible\"]/div/div[3]/yt-formatted-string//text()",
						"badge": ".//*[@id=\"badges\"]/div//span/text()"
					}
				}
			]
		},
		"youtube.be/results?search_query=": {
			"name": "YouTube",
			"selectors": [
				{
					"name": "search_result",
					"xpath": "//ytd-video-renderer",
					"xpaths": {
						"title": ".//*[@id=\"video-title\"]/yt-formatted-string/text()",
						"views": ".//*[@id=\"metadata-line\"]/span[1]/text()",
						"time": ".//*[@id=\"metadata-line\"]/span[2]/text()",
						"channel_name": ".//*[@id=\"text\"]/a/text()",
						"channel_url": ".//*[@id=\"text\"]/a/@href",
						"description": ".//*[@id=\"dismissible\"]/div/div[3]/yt-formatted-string//text()",
						"badge": ".//*[@id=\"badges\"]/div//span/text()"
					}
				}
			]
		},
		"youtube.de/results?search_query=": {
			"name": "YouTube",
			"selectors": [
				{
					"name": "search_result",
					"xpath": "//ytd-video-renderer",
					"xpaths": {
						"title": ".//*[@id=\"video-title\"]/yt-formatted-string/text()",
						"views": ".//*[@id=\"metadata-line\"]/span[1]/text()",
						"time": ".//*[@id=\"metadata-line\"]/span[2]/text()",
						"channel_name": ".//*[@id=\"text\"]/a/text()",
						"channel_url": ".//*[@id=\"text\"]/a/@href",
						"description": ".//*[@id=\"dismissible\"]/div/div[3]/yt-formatted-string//text()",
						"badge": ".//*[@id=\"badges\"]/div//span/text()"
					}
				}
			]
		},
		"twitter.com/search?q=": {
			"name": "Twitter",
			"selectors": [
				{
					"name": "search_result",
					"xpath": "//article[@data-testid]",
					"xpaths": {
						"name": ".//div[@class=\"css-1dbjc4n\"]//div[@class=\"css-1dbjc4n r-1awozwy r-18u37iz r-1wbh5a2 r-dnmrzs\"]//a//text()",
						"username": ".//div[@class=\"css-1dbjc4n\"]//*[contains(text(), \"@\")]/text()",
						"username_link": "//div[@class=\"css-1dbjc4n\"]//*[contains(text(), \"@\")]/ancestor::a/@href",
						"message": ".//div[@class=\"css-1dbjc4n\"]/div[@class=\"css-1dbjc4n\"][1]//text()"
					}
				},
				{
					"name": "people",
					"xpath": "//*[@class=\"css-18t94o4 css-1dbjc4n r-1ny4l3l r-ymttw5 r-1f1sjgu r-o7ynqc r-6416eg\"]",
					"xpaths": {
						"name": ".//span[contains(text(), \"@\")]/ancestor::a/../../../*[1]//span/text()",
						"link": ".//*[@class=\"css-901oao r-1nao33i r-xoduu5 r-18u37iz r-1q142lx r-1qd0xha r-a023e6 r-16dba41 r-rjixqe r-bcqeeo r-qvutc0\"]/ancestor::a/@href",
						"username": ".//span[contains(text(), \"@\")]/text()",
						"description": ".//div[@class=\"css-901oao r-1nao33i r-1qd0xha r-a023e6 r-16dba41 r-rjixqe r-bcqeeo r-1h8ys4a r-1jeg54m r-qvutc0\"]/span[@class=\"css-901oao css-16my406 r-poiln3 r-bcqeeo r-qvutc0\"]/text()"
					}
				}
			]
		},
        "bing.com/search?q=": {
            "name": "Bing",
            "selectors": [
                {
                    "name": "search_result",
                    "xpath": "//li[@class=\"b_algo\"]",
                    "xpaths": {
                        "title": ".//h2/a/text()",
                        "link": ".//h2/a/@href",
                        "description": ".//p/text()"
                    }
                }
            ]
        },
        "yahoo.com/search?q=": {
            "name": "Yahoo",
            "selectors": [
                {
                    "name": "search_result",
                    "xpath": "//div[@class=\"dd algo algo-sr srp algo-no-hover\"]",
                    "xpaths": {
                        "title": ".//h3/a/text()",
                        "link": ".//h3/a/@href",
                        "description": ".//p/text()"
                    }
                }
            ]
        }
	}
}