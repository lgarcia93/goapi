package repository

import (
	sql "database/sql"
	"fitgoapi/database"
	"fitgoapi/model"
	"fmt"
	"math"
	"time"
)

type IScheduleRepository interface {
	FindAllForUserPending(userID int64) ([]model.Schedule, error)
	LoadUserSchedule(userID int64, week int, page int, size int) (model.PageableSchedule, error)
	LoadUserScheduleByDay(userID int64, week int, dayOfWeek int, page int, size int) (model.PageableFlatScheduleItem, error)
	AcceptSchedule(userID int64, scheduleID int64, accept bool) (int64, error)
	LoadFullInstructorSchedule(instructorID int64) ([]model.SimpleScheduleItem, error)
	ValidateSchedule(instructorID int64, schedule []model.DecomposedSchedule, when time.Time) ([]model.ScheduleItem, error)
	ScheduleClass(userID int64, instructorID int64, schedule []model.DecomposedSchedule) (model.Schedule, error)
	LoadSchedule(scheduleID int64, userID int64) (model.Schedule, error)
	LoadScheduleItemByIdAndUser(userID int64, itemID int64) (model.FlatScheduleItem, error)
	FindIncidentForUserAndTime(userID int64, date time.Time) (model.Incident, error)
	CancelScheduleItem(userID int64, scheduleItem model.FlatScheduleItem, motive string, dayChange time.Time) (model.Incident, error)
	ChangeScheduleItem(userID int64, itemID int64, dayChange time.Time, weekDay int8, hour int8, minutes int8, duration int8, motive string) (model.Incident, error)
	LoadIncidentById(itemID int64) (model.Incident, error)
	LoadIncidentByScheduleItemAndDate(itemID int64, dayChange time.Time) (model.Incident, error)
	LoadIncidentByIdAndUser(itemID int64, userID int64) (model.Incident, error)
	LoadPendingCancellationsByUser(userID int64) ([]model.Incident, error)
	UpdateIncident(incident model.Incident) (model.Incident, error)
}

type ScheduleRepository struct {
}

func (s ScheduleRepository) LoadSchedule(scheduleID int64, userID int64) (model.Schedule, error) {
	db := database.Connection
	var schedule model.Schedule

	schedules, err := loadScheduleFromResultRows(db.Query(SQL_SCHEDULE_BY_ID, userID, userID, scheduleID))

	if len(schedules) > 0 {
		schedule = schedules[0]
	}

	return schedule, err
}

func (s ScheduleRepository) FindAllForUserPending(userID int64) ([]model.Schedule, error) {
	db := database.Connection

	results, err := db.Query(
		fmt.Sprintf(
			`%s where (sch.instructor_id = ? or sch.student_id = ?) AND sch.accepted = 0 `,
			scheduleSQL,
		),
		userID,
		userID,
	)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	var schedules []model.Schedule

	for results.Next() {

		var schedule model.Schedule

		var studentID int64
		var instructorID int64
		var accepted int
		var skillID int64

		err = results.Scan(&schedule.ID, &instructorID, &studentID, &skillID, &accepted)

		user, _ := UserRepository{}.FetchUserByID(studentID)

		schedule.User = model.SimpleUserFromUser(user)

		schedule.Accepted = false
		schedule.Skill, err = skillRepository.FetchSkillByID(skillID)

		if err != nil {
			panic(err.Error())
		}

		schedules = append(schedules, schedule)
	}

	return schedules, err
}

func (s ScheduleRepository) LoadUserSchedule(userID int64, week int, page int, size int) (model.PageableSchedule, error) {
	pageable := model.PageableSchedule{
		Content: []model.Schedule{},
		Page:    page,
	}

	db := database.Connection

	startDay, endDay, err := loadPeriodByWeekOfYear(week)

	err = db.QueryRow("Select count(distinct sid) as total from ("+SQL_SCHEDULE_BY_USER+") as outer_aux",
		userID, userID, userID, userID, startDay, endDay, userID, userID, startDay, endDay).Scan(&pageable.Total)

	if err != nil {
		return model.PageableSchedule{}, err
	} else {
		offset := page * size

		pageable.TotalPages = int(math.Ceil(float64(pageable.Total) / float64(size)))

		if pageable.Total >= offset {

			schedules, err := loadScheduleFromResultRows(db.Query(SQL_SCHEDULE_BY_USER+` limit ? offset ?`,
				userID, userID, userID, userID, startDay, endDay, userID, userID, startDay, endDay, size, offset))
			if err != nil {
				return model.PageableSchedule{}, err
			}

			pageable.Content = schedules
		}
	}

	return pageable, nil
}

func (s ScheduleRepository) LoadUserScheduleByDay(userID int64, week int, dayOfWeek int, page int, size int) (model.PageableFlatScheduleItem, error) {
	pageable := model.PageableFlatScheduleItem{
		Content: []model.FlatScheduleItem{},
		Page:    page,
	}

	db := database.Connection

	startDay, endDay, err := loadPeriodByWeekOfYear(week)

	err = db.QueryRow("Select count(distinct SCHEDULE_ID) as total from ("+SQL_SCHEDULE_BY_USER_AND_DAY+") as aux",
		userID, userID, userID, userID, dayOfWeek, startDay, endDay, userID, userID, dayOfWeek, startDay, endDay).Scan(&pageable.Total)

	if err != nil {
		return model.PageableFlatScheduleItem{}, err
	} else {
		offset := page * size

		pageable.TotalPages = int(math.Ceil(float64(pageable.Total) / float64(size)))

		if pageable.Total >= offset {

			schedules, err := loadScheduleItemFromResultRows(db.Query(SQL_SCHEDULE_BY_USER_AND_DAY+` limit ? offset ?`,
				userID, userID, userID, userID, dayOfWeek, startDay, endDay, userID, userID, dayOfWeek, startDay, endDay, size, offset))

			if err != nil {
				return model.PageableFlatScheduleItem{}, err
			}

			pageable.Content = schedules
		}
	}

	return pageable, nil
}

func (s ScheduleRepository) AcceptSchedule(userID int64, scheduleID int64, accept bool) (int64, error) {
	db := database.Connection

	var accepted int

	if accept {
		accepted = 1
	} else {
		accepted = 0
	}

	result, err := db.Exec("UPDATE schedule SET accepted = ? where (student_id = ? or instructor_id = ?) and id = ?", accepted, userID, userID, scheduleID)

	var rowsAffected int64

	if result != nil {
		rowsAffected, err = result.RowsAffected()
	}

	return rowsAffected, err
}

func (s ScheduleRepository) LoadFullInstructorSchedule(instructorID int64) ([]model.SimpleScheduleItem, error) {
	db := database.Connection

	var scheduleItems []model.SimpleScheduleItem

	rows, err := db.Query(SQL_LOAD_INSTRUCTOR_SCHEDULE, instructorID)

	if err != nil {
		return scheduleItems, err
	}

	for rows.Next() {
		var scheduleItem model.SimpleScheduleItem

		err = rows.Scan(
			&scheduleItem.WeekDay,
			&scheduleItem.Duration,
			&scheduleItem.Hour,
			&scheduleItem.Minutes,
		)

		scheduleItems = append(scheduleItems, scheduleItem)
	}

	return scheduleItems, err
}

func (s ScheduleRepository) ValidateSchedule(userID int64, schedule []model.DecomposedSchedule, when time.Time) ([]model.ScheduleItem, error) {

	var querySQL = ` SELECT si.id, si.week_day, si.hour, si.minutes, si.duration 
					 from schedule_item si inner join schedule s on si.schedule_id = s.id 
					 where (s.instructor_id = ? or s.student_id = ?) AND `

	querySQL = querySQL + decomposingEventsIntoRows(schedule)
	var isNotTodayDate = !(when.Truncate(24 * time.Hour).Equal(time.Now().Truncate(24 * time.Hour)))

	if isNotTodayDate {
		querySQL = querySQL +
			` union all
	  SELECT si.id, si.week_day, si.hour, si.minutes, si.duration 
	  from incident inc
	   inner join schedule_item si on si.id =  inc.schedule_item_id
	   inner join schedule s on s.id = si.schedule_id
	   inner join profile p on p.id = inc.requested_by
	   where (s.instructor_id = ? or s.student_id = ?)
	   AND date(inc.day_change) = date(?) AND 	
	` + decomposingIncidentIntoRows(schedule)
	}

	fmt.Println(querySQL)

	db := database.Connection

	var scheduleItems []model.ScheduleItem

	var rows *sql.Rows
	var err error

	if isNotTodayDate {
		rows, err = db.Query(querySQL, userID, userID, when, userID, userID)
	} else {
		rows, err = db.Query(querySQL, userID, userID)
	}

	if err != nil {
		return scheduleItems, err
	}

	for rows.Next() {
		var scheduleItem model.ScheduleItem

		err = rows.Scan(
			&scheduleItem.ID,
			&scheduleItem.WeekDay,
			&scheduleItem.Hour,
			&scheduleItem.Minutes,
			&scheduleItem.Duration,
		)

		scheduleItems = append(scheduleItems, scheduleItem)
	}

	return scheduleItems, err
}

func (s ScheduleRepository) ScheduleClass(userID int64, instructorID int64, events []model.DecomposedSchedule) (model.Schedule, error) {
	var errorFound error
	var insertedId int64
	var scheduleId int64
	var ids []int64

	db := database.Connection

	result, err := insertScheduleDatabase(db, userID, instructorID)
	errorFound = err

	if err == nil {
		scheduleId = result

		for _, event := range events {
			insertedId, err = insertScheduleItemDatabase(db, scheduleId, event)
			if err != nil {
				errorFound = err
				break
			}
			ids = append(ids, insertedId)
		}

		if errorFound != nil && scheduleId > 0 {
			go rollbackTransaction(scheduleId, ids)
			var schedule model.Schedule
			return schedule, errorFound
		}
	}

	return s.LoadSchedule(scheduleId, userID)
}

func (s ScheduleRepository) LoadScheduleItemByIdAndUser(userID int64, itemID int64) (model.FlatScheduleItem, error) {
	db := database.Connection

	results, err := db.Query(SQL_SCHEDULE_ITEM_BY_ID_AND_USER, userID, userID, itemID, userID, userID)

	if err != nil {
		return model.FlatScheduleItem{}, err
	}

	for results.Next() {

		var schedule = model.FlatScheduleItem{
			User: model.SimpleUser{},
		}
		err = results.Scan(
			&schedule.ID,
			&schedule.WeekDay,
			&schedule.Hour,
			&schedule.Minutes,
			&schedule.Duration,
			&schedule.ScheduleID,
			&schedule.User.ID,
			&schedule.User.FirstName,
			&schedule.User.LastName,
			&schedule.User.ProfilePicture,
		)

		if err != nil {
			return model.FlatScheduleItem{}, err
		} else {
			return schedule, err
		}
	}

	return model.FlatScheduleItem{}, fmt.Errorf("nenhum registro encontrado para o id %d", itemID)
}

func (s ScheduleRepository) FindIncidentForUserAndTime(userID int64, date time.Time) (model.Incident, error) {
	db := database.Connection

	results, err := db.Query(SQL_LOAD_INCIDENT_BY_USER_AND_TIME, userID, date, date)

	if err != nil {
		return model.Incident{}, err
	}

	for results.Next() {

		incident := model.Incident{
			RequestedBy: model.SimpleUser{},
		}

		err = results.Scan(
			&incident.ID,
			&incident.ScheduleItem,
			&incident.Accepted,
			&incident.Answered,
			&incident.Created,
			&incident.DayOfChange,
			&incident.WeekDay,
			&incident.Hour,
			&incident.Minutes,
			&incident.Type,
			&incident.Motive,
			&incident.RequestedBy.ID,
			&incident.RequestedBy.FirstName,
			&incident.RequestedBy.LastName,
			&incident.RequestedBy.ProfilePicture,
		)

		if err != nil {
			return model.Incident{}, err
		} else {
			return incident, err
		}
	}

	return model.Incident{}, nil
}

func (s ScheduleRepository) LoadIncidentById(itemID int64) (model.Incident, error) {
	db := database.Connection

	results, err := db.Query(SQL_LOAD_INCIDENT, itemID)

	if err != nil {
		return model.Incident{}, err
	}

	for results.Next() {

		incident := model.Incident{
			RequestedBy: model.SimpleUser{},
		}

		err = results.Scan(
			&incident.ID,
			&incident.ScheduleItem,
			&incident.Accepted,
			&incident.Answered,
			&incident.Created,
			&incident.DayOfChange,
			&incident.WeekDay,
			&incident.Hour,
			&incident.Minutes,
			&incident.Type,
			&incident.Motive,
			&incident.RequestedBy.ID,
			&incident.RequestedBy.FirstName,
			&incident.RequestedBy.LastName,
			&incident.RequestedBy.ProfilePicture,
		)

		if err != nil {
			return model.Incident{}, err
		} else {
			return incident, err
		}
	}

	return model.Incident{}, fmt.Errorf("nenhum registro encontrado para o id %d", itemID)
}

func (s ScheduleRepository) LoadIncidentByScheduleItemAndDate(itemID int64, dayChange time.Time) (model.Incident, error) {

	start, end := loadPeriodByCurrentDate(dayChange)
	db := database.Connection

	fmt.Printf(SQL_LOAD_INCIDENT_BY_ITEM_TIME_AND_TYPE)
	results, err := db.Query(SQL_LOAD_INCIDENT_BY_ITEM_TIME_AND_TYPE, itemID, start, end)

	if err != nil {
		return model.Incident{}, err
	}

	for results.Next() {

		incident := model.Incident{
			RequestedBy: model.SimpleUser{},
		}

		err = results.Scan(
			&incident.ID,
			&incident.ScheduleItem,
			&incident.Accepted,
			&incident.Answered,
			&incident.Created,
			&incident.DayOfChange,
			&incident.WeekDay,
			&incident.Hour,
			&incident.Minutes,
			&incident.Type,
			&incident.Motive,
			&incident.RequestedBy.ID,
			&incident.RequestedBy.FirstName,
			&incident.RequestedBy.LastName,
			&incident.RequestedBy.ProfilePicture,
		)

		if err != nil {
			return model.Incident{}, err
		} else {
			return incident, err
		}
	}

	return model.Incident{}, fmt.Errorf("nenhum registro encontrado para o id %d", itemID)
}

func (s ScheduleRepository) CancelScheduleItem(userID int64, scheduleItem model.FlatScheduleItem, motive string, dayChange time.Time) (model.Incident, error) {
	return persistIncident(userID, scheduleItem.ID, dayChange, model.Cancellation, scheduleItem.WeekDay, scheduleItem.Hour, scheduleItem.Minutes, scheduleItem.Duration, motive, s)
}

func (s ScheduleRepository) ChangeScheduleItem(userID int64, itemID int64, dayChange time.Time, weekDay int8, hour int8, minutes int8, duration int8, motive string) (model.Incident, error) {
	return persistIncident(userID, itemID, dayChange, model.Change, weekDay, hour, minutes, duration, motive, s)
}

func (s ScheduleRepository) LoadPendingCancellationsByUser(userID int64) ([]model.Incident, error) {
	db := database.Connection

	results, err := db.Query(SQL_LOAD_PENDING_INCIDENTS, userID, userID, model.Cancellation)

	if err != nil {
		return []model.Incident{}, err
	}

	var incidents []model.Incident

	for results.Next() {

		incident := model.Incident{
			RequestedBy: model.SimpleUser{},
		}

		err = results.Scan(
			&incident.ID,
			&incident.ScheduleItem,
			&incident.Accepted,
			&incident.Answered,
			&incident.Created,
			&incident.DayOfChange,
			&incident.WeekDay,
			&incident.Hour,
			&incident.Minutes,
			&incident.Type,
			&incident.Motive,
			&incident.RequestedBy.ID,
			&incident.RequestedBy.FirstName,
			&incident.RequestedBy.LastName,
			&incident.RequestedBy.ProfilePicture,
		)

		if err != nil {
			return []model.Incident{}, err
		} else {
			incidents = append(incidents, incident)
		}
	}

	return incidents, err
}

func (s ScheduleRepository) LoadIncidentByIdAndUser(itemID int64, userID int64) (model.Incident, error) {
	db := database.Connection

	results, err := db.Query(SQL_LOAD_INCIDENT_BY_USER_AND_ID, userID, userID, itemID)

	if err != nil {
		return model.Incident{}, err
	}

	for results.Next() {

		incident := model.Incident{
			RequestedBy: model.SimpleUser{},
		}

		err = results.Scan(
			&incident.ID,
			&incident.ScheduleItem,
			&incident.Accepted,
			&incident.Answered,
			&incident.Created,
			&incident.DayOfChange,
			&incident.WeekDay,
			&incident.Hour,
			&incident.Minutes,
			&incident.Type,
			&incident.Motive,
			&incident.RequestedBy.ID,
			&incident.RequestedBy.FirstName,
			&incident.RequestedBy.LastName,
			&incident.RequestedBy.ProfilePicture,
		)

		if err != nil {
			return model.Incident{}, err
		} else {
			return incident, err
		}
	}

	return model.Incident{}, fmt.Errorf("nenhum registro encontrado para o id %d", itemID)
}

func (s ScheduleRepository) UpdateIncident(incident model.Incident) (model.Incident, error) {
	db := database.Connection

	stmt, err := db.Prepare("UPDATE incident SET accepted = ?, answered = true where id = ? ")
	if err != nil {
		return model.Incident{}, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(incident.Accepted, incident.ID)
	if err != nil {
		return model.Incident{}, err
	}

	return incident, err
}

func persistIncident(userID int64, itemID int64, dayChange time.Time, incidentType model.IncidentType, weekDay int8, hour int8, minutes int8, duration int8, motive string, repository ScheduleRepository) (model.Incident, error) {
	db := database.Connection

	//TODO: Ver mudar coluna type, string vai pesar banco
	stmt, err := db.Prepare("INSERT INTO incident(schedule_item_id, requested_by, accepted, created, day_change, type, week_day, hour, minutes, duration, motive) VALUE (?, ?, false, now(), ?, ?, ? ,?, ?, ?, ?)")
	if err != nil {
		return model.Incident{}, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(itemID, userID, dayChange, incidentType, weekDay, hour, minutes, duration, motive)
	if err != nil {
		return model.Incident{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return model.Incident{}, err
	}

	return repository.LoadIncidentById(id)
}

func insertScheduleDatabase(db *sql.DB, userID int64, instructorID int64) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO schedule(INSTRUCTOR_ID, STUDENT_ID, SKILL_ID, ACCEPTED, UPDATED) VALUE (?, ?, NULL, FALSE, now())")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(instructorID, userID)
	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

func insertScheduleItemDatabase(db *sql.DB, scheduleID int64, event model.DecomposedSchedule) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO schedule_item(week_day, duration, hour, minutes, schedule_id) VALUE (?, ?, ?, ?, ?)")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(event.WeekDay, event.Duration, event.Hour, event.Minutes, scheduleID)
	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

func rollbackTransaction(scheduleId int64, ids []int64) {
	db := database.Connection

	for _, id := range ids {
		_, _ = db.Exec("DELETE from schedule_item where id = ? ", id)
	}
	_, _ = db.Exec("DELETE from schedule where id = ? ", scheduleId)
}

func loadScheduleFromResultRows(rows *sql.Rows, errQuery error) ([]model.Schedule, error) {
	schedules := []model.Schedule{}
	var errorFound error

	if errQuery != nil {
		errorFound = errQuery
	} else {
		var mapRaws = map[int64]model.Schedule{}

		for rows.Next() {
			scheduleItem := model.ScheduleItem{}

			schedule := model.Schedule{
				User:     model.SimpleUser{},
				Accepted: true,
			}

			isChange := false

			err := rows.Scan(
				&schedule.ID,
				&schedule.Updated,
				&scheduleItem.ID,
				&scheduleItem.WeekDay,
				&scheduleItem.Hour,
				&scheduleItem.Minutes,
				&scheduleItem.Duration,
				&schedule.User.ID,
				&schedule.User.FirstName,
				&schedule.User.LastName,
				&schedule.User.ProfilePicture,
				&isChange,
			)
			if err != nil {
				errorFound = err
			}

			item, ok := mapRaws[schedule.ID]
			if ok {
				item.ScheduleItems = append(schedule.ScheduleItems, scheduleItem)
				mapRaws[schedule.ID] = item
			} else {
				schedule.ScheduleItems = append(schedule.ScheduleItems, scheduleItem)
				mapRaws[schedule.ID] = schedule
			}
		}

		var keys []int64
		for k := range mapRaws {
			keys = append(keys, k)
		}

		for _, k := range keys {
			schedules = append(schedules, mapRaws[k])
		}

	}
	return schedules, errorFound
}

func loadScheduleItemFromResultRows(rows *sql.Rows, errQuery error) ([]model.FlatScheduleItem, error) {
	items := []model.FlatScheduleItem{}
	var errorFound error

	if errQuery != nil {
		errorFound = errQuery
	} else {
		for rows.Next() {
			scheduleItem := model.FlatScheduleItem{}

			isChange := false

			var updated time.Time
			err := rows.Scan(
				&scheduleItem.ScheduleID,
				&updated,
				&scheduleItem.ID,
				&scheduleItem.WeekDay,
				&scheduleItem.Hour,
				&scheduleItem.Minutes,
				&scheduleItem.Duration,
				&scheduleItem.User.ID,
				&scheduleItem.User.FirstName,
				&scheduleItem.User.LastName,
				&scheduleItem.User.ProfilePicture,
				&isChange,
			)
			if err != nil {
				errorFound = err
			}

			items = append(items, scheduleItem)
		}

	}
	return items, errorFound
}

func decomposingIncidentIntoRows(events []model.DecomposedSchedule) string {
	var query = " ("

	for index, event := range events {
		query = query + fmt.Sprintf(` 
			(inc.week_day = %d and (
            (if((inc.minutes + inc.duration) > 60, inc.hour + 1, inc.hour) = (%d)
             AND if((inc.minutes + inc.duration) > 60, 60 - (inc.duration - inc.minutes), inc.minutes + inc.duration) > %d)
             OR
            ( inc.hour = %d and  (inc.minutes + inc.duration) > %d )
      	)) `, event.WeekDay, event.Hour, event.Minutes,
			event.Hour, event.Minutes)

		if index < (len(events) - 1) {
			query = query + " OR "
		}
	}

	query = query + ") "

	return query
}

func decomposingEventsIntoRows(events []model.DecomposedSchedule) string {
	var query = " ("

	for index, event := range events {
		query = query + fmt.Sprintf(` 
			(week_day = %d and (
            (if((minutes + duration) > 60, hour + 1, hour) = (%d)
             AND if((minutes + duration) > 60, 60 - (duration - minutes), minutes + duration) > %d)
             OR
            ( hour = %d and  (minutes + duration) > %d )
      	)) `, event.WeekDay, event.Hour, event.Minutes,
			event.Hour, event.Minutes)

		if index < (len(events) - 1) {
			query = query + " OR "
		}
	}

	query = query + ") "

	return query
}

func loadPeriodByWeekOfYear(week int) (time.Time, time.Time, error) {

	if week > 53 {
		return time.Time{}, time.Time{}, fmt.Errorf("numero de semana inválido .: %d", week)
	}
	var today = time.Now()
	_, thisWeek := today.ISOWeek()

	toDayInTheWeekPassed := today.AddDate(0, 0, 7*(week-thisWeek))
	start, end := loadPeriodByCurrentDate(toDayInTheWeekPassed)
	return start, end, nil
}

func loadPeriodByCurrentDate(currentDate time.Time) (time.Time, time.Time) {
	var startDate time.Time
	var endDate time.Time

	if currentDate.Weekday() == time.Sunday {
		startDate = currentDate
		endDate = currentDate.AddDate(0, 0, 6)
	} else {
		startDate = currentDate.AddDate(0, 0, int(currentDate.Weekday())*-1)
		endDate = currentDate.AddDate(0, 0, int(6-currentDate.Weekday()))
	}
	return startDate, endDate
}
