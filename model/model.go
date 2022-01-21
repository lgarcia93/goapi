package model

import (
	"time"
)

// City struct
type City struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	ZipCode string `json:"zipCode"`
	UF      string `json:"uf"`
}

// Skill struct
type Skill struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// User struct
type User struct {
	ID               int64  `json:"id"`
	Username         string `json:"username"`
	Password         string
	FirstName        string  `json:"firstName"`
	LastName         string  `json:"lastName"`
	Description      string  `json:"description"`
	ProfilePicture   string  `json:"profilePicture"`
	FcmTokenFirebase string  `json:"fcmTokenFirebase"`
	City             City    `json:"city"`
	Skills           []Skill `json:"skills"`
	IsInstructor     bool    `json:"isInstructor"`
	IsConnection     bool    `json:"isConnection"`
}

type ScheduleItem struct {
	ID         int64 `json:"id"`
	ScheduleID int64 `json:"scheduleId"`
	WeekDay    int8  `json:"weekDay"`
	Hour       int8  `json:"hour"`
	Minutes    int8  `json:"minutes"`
	Duration   int8  `json:"duration"`
}

type FlatScheduleItem struct {
	ID         int64      `json:"id"`
	ScheduleID int64      `json:"scheduleId"`
	User       SimpleUser `json:"user"`
	WeekDay    int8       `json:"weekDay"`
	Hour       int8       `json:"hour"`
	Minutes    int8       `json:"minutes"`
	Duration   int8       `json:"duration"`
}

type SimpleScheduleItem struct {
	WeekDay  int8 `json:"weekDay"`
	Hour     int8 `json:"hour"`
	Minutes  int8 `json:"minutes"`
	Duration int8 `json:"duration"`
}

type SimpleUser struct {
	ID             int64  `json:"id"`
	Username       string `json:"username"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Description    string `json:"description"`
	ProfilePicture string `json:"profilePicture"`
	IsInstructor   bool   `json:"isInstructor"`
}

func SimpleUserFromUser(user User) SimpleUser {
	return SimpleUser{
		ID:             user.ID,
		Description:    user.Description,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		ProfilePicture: user.ProfilePicture,
		IsInstructor:   user.IsInstructor,
		Username:       user.Username,
	}
}

type Schedule struct {
	ID            int64          `json:"id"`
	User          SimpleUser     `json:"user"`
	Skill         Skill          `json:"skill"`
	ScheduleItems []ScheduleItem `json:"scheduleItems"`
	Updated       time.Time      `json:"updated"`
	Accepted      bool           `json:"accepted"`
}

type Incident struct {
	ID           int64        `json:"id"`
	ScheduleItem int64        `json:"scheduleItemId"`
	RequestedBy  SimpleUser   `json:"requestedBy"`
	DayOfChange  time.Time    `json:"dayOfChange"`
	WeekDay      int8         `json:"weekDay"`
	Hour         int8         `json:"hour"`
	Minutes      int8         `json:"minutes"`
	Duration     int8         `json:"duration"` //TODO: rever se tira ou não
	Created      time.Time    `json:"created"`
	Accepted     bool         `json:"accepted"`
	Type         IncidentType `json:"type"`
	Answered     bool         `json:"answered"`
	Motive       string       `json:"motive"`
}

type IncidentType string

const (
	Cancellation IncidentType = "cancellation"
	Change                    = "change"
)

type PageableUser struct {
	Users      []User `json:"content"`
	Total      int    `json:"total"`
	TotalPages int    `json:"totalPages"`
	Page       int    `json:"current"`
}

type PageableUserPlain struct {
	Users      []SimpleUser `json:"content"`
	Total      int          `json:"total"`
	TotalPages int          `json:"totalPages"`
	Page       int          `json:"current"`
}

type PageableSchedule struct {
	Content    []Schedule `json:"content"`
	Total      int        `json:"total"`
	TotalPages int        `json:"totalPages"`
	Page       int        `json:"current"`
}

type PageableFlatScheduleItem struct {
	Content    []FlatScheduleItem `json:"content"`
	Total      int                `json:"total"`
	TotalPages int                `json:"totalPages"`
	Page       int                `json:"current"`
}
