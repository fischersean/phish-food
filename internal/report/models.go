package report

import (
	"github.com/fischersean/phish-food/internal/stocks"
)

type StockReport struct {
	Stock stocks.Stock
	Count CountResult
}

type CountResult struct {
	PostScore       int
	CommentScore    int
	TotalScore      float64
	PostMentions    int
	CommentMentions int
}

func (c *CountResult) Add(c2 CountResult) {

	c.PostScore += c2.PostScore
	c.CommentScore += c2.CommentScore
	c.PostMentions += c2.PostMentions
	c.CommentMentions += c2.CommentMentions

}
