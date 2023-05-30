# Digitale Polarisatie
## Project structure
- api 
    - Contains the API code
- docker
    - Contains docker container configuration and docker compose scripts
- ext
    - Contains the browser extension code

## API
### Installation
- Python 3.8.10 was used for the development of the API. (So it's recommended to use the same version, however higher versions might work as well)
- MongoDB 5.0.9 was used as database

By running the dev.cmd (Windows) or dev.sh (Linux) script, this will automatically install the required packages in order for the application to start. In order to install the specific version of packages that were used during the creation of this project, then you can run the following command to install the packages:
    
    ```bash
    pip install -r requirements.txt
    ```

This will install the packages that are specified in the requirements.txt


### Configuration
In order to configure the API, you can either configure the application through the .env file or through setting the environment variables via command line.

For example, to set the environment variable `API_PORT` to `8080`, you can run the following command:

    ```bash
    export API_PORT=8080
    ```

## EXT-BMS

### Installation
- NodeJS v16.13.2 was used for the development of the project
- PNPM v7.8.0 (Yarn or NPM can be used)

By running the dev.cmd (Windows) this will install the packages automatically and execute the build script for Google Chrome (this was used mostly during development). The dev.cmd script is mainly used for development and not autically building for production.

In order to build for production, you can run the following command.
Chrome:

    ```bash
    pnpm run app:chrome
    ```

Firefox:
    
    ```bash
    pnpm run app:firefox
    ```

Edge:

    ```bash
    pnpm run app:edge
    ```

This will build the application for Google Chrome in a *dist* folder. Each of the command build in the same dist folder. Therefore the commands have to be run each time for each version of browser.

There is a potential bug that has occurred multiple times with running the build process for Firefox. A solution for this problem is to do the following:

- Run the command `pnpm run app:chrome`
- Copy the contents of the file at location `public/firefox_manifest.json` into the file `dist/manifest.json`
- Then install the extension in Firefox browser

## Browser
### Installation

To install the extension in the browser, you can watch the videos linked below.
- https://youtu.be/5pWS-PM-8C8 (Chrome)
- https://youtu.be/kb_Mb9JNwPI (Firefox)

The videos are also provided in the `docs` folder.


## Data
This chapter is dedicated to explain what data is gathered from all the search engines. This data is extracted through the Extractor class within the "content.ts" file. To extend or modify how the data is extracted, this can be done in the "config.x.x.x.json" file. 

The extractor uses XPath to extract data instead of CSS selector, reason for this choice is that because the XPath can be used to traverse backwards (such as getting the parents of the object) and CSS selector can't unless it's mixed with JavaScript. 

### Google News
| Feature | Description |
| --- | --- |
| Publisher   | The publisher of the article       |
| Title       | The title of the article           |
| Description | The description of the article     |
| Time        | The time the article was published |
| Link        | The URL/Link of the article        |

### Google Videos
| Feature | Description |
| --- | --- |
| Publisher   | The publisher of the video       |
| Title       | The title of the video           |
| Link        | The URL/Link of the video        |
| Description | The description of the video     |
| Subtitle    | The subtitle of the video        |

### Google Search

#### Results
This is the general search result.

| Feature | Description |
| --- | --- |
| Publisher   | The publisher of the article       |
| Title       | The title of the article           |
| Description | The description of the article     |
| Link        | The URL/Link of the article        |
| Date        | The date the article was published |

#### Featured Result
| Feature | Description |
| --- | --- |
| Publisher   | The publisher of the article       |
| Title       | The title of the article           |

#### Featured Result Links
This is shown at the top when there are multiple links within a website/page that are related to the search query.

| Feature | Description |
| --- | --- |
| Title       | The title of the link           |
| Link        | The URL/Link of the link        |
| Description | The description of the link     |

#### Sidebar Results
This is result that is shown on the right side, usually containing information about an orgnization or a person.

| Feature | Description |
| --- | --- |
| Title       | The title of the link           |
| Link        | The URL/Link of the link        |

#### People Also Searched
| Feature | Description |
| --- | --- |
| Title       | The title of the link           |
| Link        | The URL/Link of the link        |

#### See Results About
| Feature | Description |
| --- | --- |
| Title       | The title of the link           |
| Link        | The URL/Link of the link        |

#### Snippet Results
This is the result that is given by google that usually contains a bolded text or a highlighted text from a page, that is related to the query.

| Feature | Description |
| --- | --- |
| Text       | The text of the snippet           |
| Publisher  | The publisher of the snippet      |
| Title      | The title of the snippet          |
| Link       | The URL/Link of the snippet       |

#### People Also Ask
| Feature | Description |
| --- | --- |
| Question       | The question of the query |

#### Top Stories
This is typically shown to the user when there's is a related news to the search query.

| Feature | Description |
| --- | --- |
| Publisher   | The publisher of the article       |
| Title       | The title of the article           |
| Link        | The URL/Link of the article        |
| Time        | The time the article was published |

#### Related Searches
| Feature | Description |
| --- | --- |
| Title       | The title of the link           |
| Link        | The URL/Link of the link        |

#### Map Locations
| Feature | Description |
| --- | --- |
| Title       | The title of the link                       |
| Link        | The URL/Link of the link                    |
| Properties  | The details of the location                 |
| Rated       | The rating of the location (out of 5 stars) |
| Reviews Count | The amount of reviews for the location    |


### DuckDuckGo

#### Search Results
| Feature | Description |
| --- | --- |
| Publisher   | The publisher of the article       |
| Title       | The title of the article           |
| Link        | The URL/Link of the article        |
| Description | The description of the article     |

#### Related Results
| Feature | Description |
| --- | --- |
| Title       | The title of the link           |
| Link        | The URL/Link of the link        |

#### Sidebar Results
| Feature | Description |
| --- | --- |
| Title       | The title of the link           |
| Link        | The URL/Link of the link        |
| Description | The description of the link     |
| Related Links | The related links of the link |

#### Recent News
| Feature | Description |
| --- | --- |
| Title       | The title of the link           |
| Link        | The URL/Link of the link        |
| Publisher   | The publisher of the link       |
| Description | The description of the link     |

### YouTube
| Feature | Description |
| --- | --- |
| Title        | The title of the video           |
| Views        | The amount of views the video has |
| Time         | The time the video was published |
| Channel Name | The name of the channel        |
| Channel URL  | The URL/Link of the channel    |
| Description  | The description of the video   |
| Badge        | The badge, usually seen like "new" or "hd" |

### Twitter

#### Search Results
| Feature | Description |
| --- | --- |
| Name          | The name of the user           |
| Username      | The username of the user       |
| Username Link | The URL/Link of the user       |
| Message       | The message of the user        |

#### People
| Feature | Description |
| --- | --- |
| Name          | The name of the user           |
| Link          | The URL/Link of the user       |
| Username      | The username of the user       |
| Description   | The description of the user    |

