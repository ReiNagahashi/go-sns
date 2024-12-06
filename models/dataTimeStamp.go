package models

import "time"

type DateTimeStamp struct{
	createdAt time.Time
	updatedAt time.Time
}

func NewDateTimeStamp(createdAt, updatedAt time.Time) *DateTimeStamp{
	return &DateTimeStamp{createdAt, updatedAt}
}

func (dateTime DateTimeStamp) GetCreatedAt()time.Time{
	return dateTime.createdAt
} 

func (dateTime *DateTimeStamp) SetCreatedAt(createdAt time.Time){
	dateTime.createdAt = createdAt
}

func (dateTime DateTimeStamp) GetUpdatedAt()time.Time{
	return dateTime.updatedAt
} 

func (dateTime *DateTimeStamp) SetUpdatedAt(updatedAt time.Time){
	dateTime.updatedAt = updatedAt
}