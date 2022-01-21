package model

import "time"

type ScheduleItemDTO struct {
	WeekDay  string `json:"weekDay"`
	Hour     int8   `json:"hour"`
	Minutes  int8   `json:"minutes"`
	Duration int8   `json:"duration"`
}

type DecomposedSchedule struct {
	WeekDay  int8
	Hour     int8
	Minutes  int8
	Duration int8
}

type IncidentDTO struct {
	WeekDay   string    `json:"weekDay"`
	Hour      int8      `json:"hour"`
	Minutes   int8      `json:"minutes"`
	Duration  int8      `json:"duration"`
	DayChange time.Time `json:"dayChange"`
	Motive    string    `json:"motive"`
}
