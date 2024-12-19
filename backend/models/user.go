package models

import "encoding/json"

type User struct {
	id        int
	name      string
	email     string
	timeStamp DateTimeStamp
}

type userJSON struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	TimeStamp DateTimeStamp `json:"timeStamp"`
}

func NewUser(id int, name, email string, timeStamp DateTimeStamp) *User {
	return &User{id, name, email, timeStamp}
}

func (u User) GetId() int {
	return u.id
}

func (u *User) SetId(id int) {
	u.id = id
}

func (u User) GetName() string {
	return u.name
}

func (u *User) SetName(name string) {
	u.name = name
}

func (u User) GetEmail() string {
	return u.email
}

func (u *User) SetEmail(email string) {
	u.email = email
}

func (u User) GetTimeStamp() DateTimeStamp {
	return u.timeStamp
}

func (u *User) SetTimeStamp(timeStamp DateTimeStamp) {
	u.timeStamp = timeStamp
}

func (u *User) GetFields() []interface{} {
	return []interface{}{u.name, u.email, u.timeStamp.GetCreatedAt(), u.timeStamp.GetUpdatedAt()}
}

func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(&userJSON{
		ID:        u.id,
		Name:      u.name,
		Email:     u.email,
		TimeStamp: u.timeStamp,
	})
}

func (u *User) UnmarshalJSON(data []byte) error {
	var uj userJSON
	if err := json.Unmarshal(data, &uj); err != nil {
		return err
	}
	u.id = uj.ID
	u.name = uj.Name
	u.email = uj.Email
	u.timeStamp = uj.TimeStamp

	return nil
}
