package api

import (
	"math"
	"time"
)

type Post interface {
	ID() int
	Deleted() bool
	Type() string
	By() string
	Time() time.Time
	Text() string
	Dead() bool
	Parent() int
	Poll() int
	Kids() []int
	URL() string
	Score() int
	Title() string
	Parts() []int
	Descendants() int
}

type post map[string]interface{}

func (p post) ID() int {
	id, _ := p["id"].(int)
	return id
}

func (p post) Deleted() bool {
	deleted, _ := p["deleted"].(bool)
	return deleted
}

func (p post) Type() string {
	pType, _ := p["type"].(string)
	return pType
}

func (p post) By() string {
	by, _ := p["by"].(string)
	return by
}

func (p post) Time() time.Time {
	rawTime, _ := p["time"].(int64)
	return time.Unix(rawTime, 0)
}

func (p post) Text() string {
	text, _ := p["text"].(string)
	return text
}

func (p post) Dead() bool {
	dead, _ := p["dead"].(bool)
	return dead
}

func (p post) Parent() int {
	parent, _ := p["parent"].(int)
	return parent
}

func (p post) Poll() int {
	poll, _ := p["poll"].(int)
	return poll
}

func (p post) Kids() []int {
	iKids, ok := p["kids"].([]interface{})
	if !ok {
		return []int{}
	}

	var kids []int

	for _, kid := range iKids {
		kids = append(kids, int(math.Abs(kid.(float64))))
	}

	return kids
}

func (p post) URL() string {
	url, _ := p["url"].(string)
	return url
}

func (p post) Score() int {
	score, _ := p["score"].(float64)
	return int(score)
}

func (p post) Title() string {
	title, _ := p["title"].(string)
	return title
}

func (p post) Parts() []int {
	iParts, ok := p["parts"].([]interface{})
	if !ok {
		return []int{}
	}

	var parts = make ([]int, len(iParts))

	for _, v := range iParts {
		parts = append(parts, v.(int))
	}
	return parts
}

func (p post) Descendants() int {
	descendants, _ := p["descendants"].(int)
	return descendants
}


