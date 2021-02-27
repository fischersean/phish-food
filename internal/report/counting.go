package report

import (
	"fmt"
	"github.com/fischersean/phish-food/internal/reddit"
	"github.com/fischersean/phish-food/internal/stocks"
	"math"
	"regexp"
	"sort"
)

const (
	MAX_COMMENT_DEPTH = 8
	MAX_REPLY_LOOK    = 100
)

var (
	TickerBlackList = []string{
		"USD",
		"YOLO",
		"DD",
		"RH",
		"HOLD",
		"IPO",
		"ITM",
		"OTM",
		"IMO",
		"CEO",
		"CFO",
		"CTO",
		"FTC",
		"FCC",
		"FEC",
		"MOON",
		"CAN",
		"EPS",
		"ATH",
		"ATL",
		"FOR",
		"GDP",
		"EDIT",
	}
	TickerRestrictedList = map[string]string{
		"A": "",
		"X": "",
		"U": "",
	}
)

const (
	MinTotalScore = 10.0
)

func reFromTicker(s stocks.Stock) (*regexp.Regexp, error) {

	if _, present := TickerRestrictedList[s.Symbol]; !present {
		return regexp.Compile(fmt.Sprintf(`\b%s\b`, s.Symbol))
	}

	return regexp.Compile(fmt.Sprintf(`\$%s\b`, s.Symbol))
}

func commentCountReferences(c reddit.Comment, re *regexp.Regexp) (count CountResult) {

	if c.Depth > MAX_COMMENT_DEPTH {
		return count
	}

	if re.FindString(c.Body) != "" {
		count.CommentScore = c.Score
		count.CommentMentions = 1
	}

	for i, v := range c.Replies {
		if i > MAX_REPLY_LOOK {
			break
		}
		count.Add(commentCountReferences(v, re))
	}

	return count
}

func postCountReferences(p reddit.Post, re *regexp.Regexp) (count CountResult) {

	if re.FindString(fmt.Sprintf("%s %s", p.Title, p.Body)) != "" {
		count.PostScore = p.Score
		count.PostMentions += 1
	}

	for i, c := range p.Comments {
		if i > MAX_REPLY_LOOK {
			break
		}
		count.Add(commentCountReferences(c, re))
	}

	return count
}

func countTotalScore(c CountResult) float64 {

	mentionScore := float64(c.PostMentions) * float64(c.CommentMentions)
	upScore := math.Log(float64(c.PostScore) + 0.1*float64(c.CommentScore) + 1.0)
	return mentionScore * upScore
}

func CountStockReferences(tickers []stocks.Stock, posts []reddit.Post) ([]StockReport, error) {

	var reports = make([]StockReport, len(tickers))
	reportCount := 0

	tickBlackListMap := make(map[string]string)

	for _, v := range TickerBlackList {
		tickBlackListMap[v] = ""
	}

	for _, stk := range tickers {
		r := StockReport{Stock: stk}
		if _, present := tickBlackListMap[stk.Symbol]; present {
			continue
		}
		re, err := reFromTicker(stk)
		if err != nil {
			return nil, err
		}

		for _, p := range posts {
			r.Count.Add(postCountReferences(p, re))
		}

		r.Count.TotalScore = countTotalScore(r.Count)
		if r.Count.TotalScore > 0 {
			reports[reportCount] = r
			reportCount += 1
		}
	}

	if reportCount == 0 {
		return nil, nil
	}
	return reports[0 : reportCount-1], nil
}

// SortStockReports sorts the stockReports slice from greatest to least
func SortStockReports(stockReports []StockReport) {

	lessFunc := func(i, j int) bool {
		return stockReports[i].Count.TotalScore > stockReports[j].Count.TotalScore
	}

	sort.Slice(stockReports, lessFunc)
}
