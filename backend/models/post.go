package models

import "encoding/json"

type Post struct {
	id          int
	title       string
	description string
	timeStamp   DateTimeStamp
}

type postJSON struct {
	ID          int           `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	TimeStamp   DateTimeStamp `json:"timeStamp"`
}

func NewPost(id int, title, description string, timeStamp DateTimeStamp) *Post {
	return &Post{id, title, description, timeStamp}
}

func (p Post) GetId() int {
	return p.id
}

func (p *Post) SetId(id int) {
	p.id = id
}

func (p Post) GetTitle() string {
	return p.title
}

func (p *Post) SetTitle(title string) {
	p.title = title
}

func (p Post) GetDescription() string {
	return p.description
}

func (p *Post) SetDescription(description string) {
	p.description = description
}

func (p Post) GetTimeStamp() DateTimeStamp {
	return p.timeStamp
}

func (p *Post) SetTimeStamp(timeStamp DateTimeStamp) {
	p.timeStamp = timeStamp
}

func (p *Post) GetFields() []interface{} {
	return []interface{}{p.title, p.description, p.timeStamp.GetCreatedAt(), p.timeStamp.GetUpdatedAt()}
}

func (p *Post) MarshalJSON() ([]byte, error) {
	return json.Marshal(&postJSON{
		ID:          p.id,
		Title:       p.title,
		Description: p.description,
		TimeStamp:   p.timeStamp,
	})
}

func (p *Post) UnmarshalJSON(data []byte) error {
	var pj postJSON
	if err := json.Unmarshal(data, &pj); err != nil {
		return err
	}
	p.id = pj.ID
	p.title = pj.Title
	p.description = pj.Description
	p.timeStamp = pj.TimeStamp
	
	return nil
}