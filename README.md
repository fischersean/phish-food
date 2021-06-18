# PhishFood
[![Deploy](https://github.com/fischersean/phish-food/actions/workflows/deploy.yml/badge.svg)](https://github.com/fischersean/phish-food/actions/workflows/deploy.yml) [![CodeQL](https://github.com/fischersean/phish-food/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/fischersean/phish-food/actions/workflows/codeql-analysis.yml)

*Like the Ben & Jerry's Ice Cream*

## About
PhishFood is the source code for the ETL pipeline and associated cloud infrastructure required to produce the TheKettle database. This project aims to provide quality, reliable data that an end user can have confidence in. At it's core, the pipeline and database attempts to collect and summarize what is "Hot" on Reddit's most popular trading subreddits. 

Below is an example database entry (NoSQL version):

```
{
    "id": "wallstreetbets_20210408",
    "hour": 18,
    "data": [
        {
            "Stock": {
                "Symbol": "GME",
                "FullName": "GameStop Corporation Common Stock",
                "Exchange": "NYSE"
            },
            "Count": {
                "PostScore": 16852,
                "CommentScore": 3306,
                "TotalScore": 975.1710719570256,
                "PostMentions": 2,
                "CommentMentions": 50
            }
        }
    ]
}
```

Currently there are 3 supported subreddits:
- stocks
- wallstreetbets
- investing

## Why?
It was widely reported during the GameStop hype that hedge funds were setting up or buying applications to scrape Reddit for the latest trending stock data. I thought it would be helpful to a retail trader to have access to the same data.
