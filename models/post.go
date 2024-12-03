package models

type Post struct{
	id int
	title string
	description string
	timeStamp DateTimeStamp
}

func NewPost(id int, title, description string, timeStamp DateTimeStamp) *Post{
	return &Post{id, title, description, timeStamp}
}

func (p Post) GetId()int{
	return p.id
}

func (p *Post) SetId(id int){
	p.id = id
}

func (p Post) GetTitle() string{
	return p.title
}

func (p *Post) SetTitle(title string){
	p.title = title
}

func (p Post) GetDescription()string{
	return p.description
}

func (p *Post) SetDescription(description string){
	p.description = description
}

func (p Post) GetTimeStamp() DateTimeStamp{
	return p.timeStamp
}

func (p *Post) SetTimeStamp(timeStamp DateTimeStamp){
	p.timeStamp = timeStamp
}


