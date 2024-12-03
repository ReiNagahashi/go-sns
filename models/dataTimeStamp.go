package models

type DateTimeStamp struct{
	createdAt string
	updatedAt string
}

func NewDateTimeStamp(createdAt, updatedAt string) *DateTimeStamp{
	return &DateTimeStamp{createdAt, updatedAt}
}

func (dateTime DateTimeStamp) GetCreatedAt()string{
	return dateTime.createdAt
} 

func (dateTime *DateTimeStamp) SetCreatedAt(createdAt string){
	dateTime.createdAt = createdAt
}

func (dateTime DateTimeStamp) GetUpdatedAt()string{
	return dateTime.updatedAt
} 

func (dateTime *DateTimeStamp) SetUpdatedAt(updatedAt string){
	dateTime.updatedAt = updatedAt
}