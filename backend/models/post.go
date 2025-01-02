package models

import "encoding/json"

type Post struct {
	id          	int
	title       	string
	description 	string
	timeStamp   	DateTimeStamp
	submitted_by 	int
}

type postJSON struct {
	ID          	int           	`json:"id"`
	Title       	string        	`json:"title"`
	Description 	string        	`json:"description"`
	TimeStamp   	DateTimeStamp 	`json:"timeStamp"`
	Submitted_by 	int 			`json:"submitted_by"`
}

func NewPost(id, submitted_by int, title, description string, timeStamp DateTimeStamp) *Post {
	return &Post{id, title, description, timeStamp, submitted_by}
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


func (p *Post) GetSubmitted_by() int{
	return p.submitted_by
}

func (p *Post) SetSubmitted_by(submitted_by int){
	p.submitted_by = submitted_by
}


func (p *Post) GetFields() []interface{} {
	return []interface{}{p.title, p.description, p.submitted_by, p.timeStamp.GetCreatedAt(), p.timeStamp.GetUpdatedAt()}
}

func (p *Post) MarshalJSON() ([]byte, error) {
	return json.Marshal(&postJSON{
		ID:          p.id,
		Title:       p.title,
		Description: p.description,
		Submitted_by: p.submitted_by,
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
	p.submitted_by = pj.Submitted_by
	p.timeStamp = pj.TimeStamp
	
	return nil
}