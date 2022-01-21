package service

import (
	"fitgoapi/model"
	"fitgoapi/utils"
	"fmt"
	"time"
)

// IScheduleService ...
type IScheduleService interface {
	LoadUserSchedule(userID int64, week int, page int, size int) (model.PageableSchedule, error)
	LoadUserScheduleByDay(userID int64, week int, dayOfWeek int, page int, size int) (model.PageableFlatScheduleItem, error)
	FindAllUserPendingSchedule(userID int64) ([]model.Schedule, error)
	AcceptOrRejectSchedule(userID int64, scheduleID int64, accept bool) (int64, error)
	LoadFullInstructorSchedule(instructorID int64) (map[string][]model.SimpleScheduleItem, error)
	ScheduleClass(userID int64, instructorID int64, items []model.ScheduleItemDTO) (model.Schedule, error)
	CancelScheduleEvent(userID int64, itemID int64, week int) (model.Incident, error)
	ChangeScheduleEvent(userID int64, itemID int64, incident model.IncidentDTO) (model.Incident, error)
	LoadIncidentByPeriod(itemID int64, dayOfChange time.Time) (model.Incident, error)
	LoadScheduleItemByIdAndUser(itemID int64, userID int64) (model.FlatScheduleItem, error)
	LoadPendingIncidents(userID int64, incidentType model.IncidentType) ([]model.Incident, error)
	AcceptOrRejectChange(userID int64, changeID64 int64, answer bool) error
}

// ScheduleService ...
type ScheduleService struct {
}

// LoadUserSchedule ...
func (service ScheduleService) LoadUserSchedule(userID int64, week int, offset int, size int) (model.PageableSchedule, error) {
	return scheduleRepository.LoadUserSchedule(userID, week, offset, size)
}

func (service ScheduleService) LoadUserScheduleByDay(userID int64, week int, dayOfWeek int, offset int, size int) (model.PageableFlatScheduleItem, error) {
	return scheduleRepository.LoadUserScheduleByDay(userID, week, dayOfWeek, offset, size)
}

// FindAllUserPendingSchedule ...
func (service ScheduleService) FindAllUserPendingSchedule(userID int64) ([]model.Schedule, error) {
	return scheduleRepository.FindAllForUserPending(userID)
}

// AcceptOrRejectSchedule ...
func (service ScheduleService) AcceptOrRejectSchedule(userID int64, scheduleID int64, accept bool) (int64, error) {
	return scheduleRepository.AcceptSchedule(userID, scheduleID, accept)
}

func (service ScheduleService) LoadFullInstructorSchedule(instructorID int64) (map[string][]model.SimpleScheduleItem, error) {
	scheduleItems, err := scheduleRepository.LoadFullInstructorSchedule(instructorID)

	itemsMap := map[string][]model.SimpleScheduleItem{}

	if err != nil {
		return itemsMap, err
	}

	for index, _ := range scheduleItems {

		weekDay := utils.WeekDayFromNumber(scheduleItems[index].WeekDay)

		if itemsMap[weekDay] == nil {
			itemsMap[weekDay] = []model.SimpleScheduleItem{}
		}

		itemsMap[weekDay] = append(itemsMap[weekDay], scheduleItems[index])
	}

	return itemsMap, err
}

func (service ScheduleService) ScheduleClass(userID int64, instructorID int64, items []model.ScheduleItemDTO) (model.Schedule, error) {
	var schedule = model.Schedule{}

	var weekDaySchedule, err = decomposeScheduleDays(items)

	if err != nil {
		return schedule, err
	}

	if len(weekDaySchedule) == 0 {
		return schedule, fmt.Errorf("lista de horários vazia.: ")
	}

	var events, errValidation = scheduleRepository.ValidateSchedule(instructorID, weekDaySchedule, time.Now())

	if errValidation != nil {
		return schedule, errValidation
	}

	if len(events) > 0 {
		return schedule, generateErrorFromEvents(events)
	}

	schedule, errorSchedule := scheduleRepository.ScheduleClass(userID, instructorID, weekDaySchedule)

	return schedule, errorSchedule
}

func (service ScheduleService) LoadScheduleItemByIdAndUser(itemID int64, userID int64) (model.FlatScheduleItem, error) {
	scheduleItem, err := scheduleRepository.LoadScheduleItemByIdAndUser(userID, itemID)

	if err != nil {
		return model.FlatScheduleItem{}, err
	}

	if scheduleItem.ID == 0 {
		return model.FlatScheduleItem{}, fmt.Errorf("nenhum evento localizado para o id %d", itemID)
	}

	return scheduleItem, err
}

func (service ScheduleService) CancelScheduleEvent(userID int64, itemID int64, week int) (model.Incident, error) {

	scheduleItem, err := service.LoadScheduleItemByIdAndUser(itemID, userID)

	if err != nil {
		return model.Incident{}, err
	}

	var toValidate = time.Now()
	_, thisWeek := toValidate.ISOWeek()

	toValidate = toValidate.AddDate(0, 0, (week-thisWeek)*7)
	var weekDayToCompare = toValidate.Weekday()

	var when = toValidate.AddDate(0, 0, int(scheduleItem.WeekDay-int8(weekDayToCompare)))
	if time.Now().After(when) {
		return model.Incident{}, fmt.Errorf("semana de referência inválida.: %d", week)
	}

	incident, _ := service.LoadIncidentByPeriod(itemID, when)

	if incident.ID > 0 {
		return model.Incident{}, fmt.Errorf("já existe um pedido de cancelamento para este dia .: %s", when.String())
	}

	incidentRegistered, err := scheduleRepository.CancelScheduleItem(userID, scheduleItem, "", when)

	return incidentRegistered, err
}

func (service ScheduleService) ChangeScheduleEvent(userID int64, itemID int64, incident model.IncidentDTO) (model.Incident, error) {

	scheduleItem, err := service.LoadScheduleItemByIdAndUser(itemID, userID)

	if err != nil {
		return model.Incident{}, err
	}

	weekNumber, err := utils.NumberFromWeekDay(incident.WeekDay)
	if err != nil {
		return model.Incident{}, err
	}

	var dto = model.DecomposedSchedule{
		WeekDay:  weekNumber,
		Hour:     incident.Hour,
		Minutes:  incident.Minutes,
		Duration: incident.Duration,
	}

	changeForTheWeek, err := scheduleRepository.LoadIncidentByScheduleItemAndDate(itemID, incident.DayChange)
	if err != nil {
		return model.Incident{}, err
	}

	if changeForTheWeek.ID > 0 {
		return model.Incident{}, fmt.Errorf("já existe uma alteração de horários para esta semana .: %d/%d %d:%d",
			incident.DayChange.Day(), incident.DayChange.Month(), incident.Hour, incident.Minutes)
	}

	items, err := scheduleRepository.ValidateSchedule(scheduleItem.User.ID, []model.DecomposedSchedule{dto}, incident.DayChange)

	if err != nil {
		return model.Incident{}, err
	}

	if len(items) > 0 {
		return model.Incident{}, fmt.Errorf("o usuário já está com o horário comprometido .: %d/%d %d:%d",
			incident.DayChange.Day(), incident.DayChange.Month(), incident.Hour, incident.Minutes)
	}

	incidentRegistered, err := scheduleRepository.ChangeScheduleItem(userID, itemID, incident.DayChange, weekNumber,
		incident.Hour, incident.Minutes, scheduleItem.Duration, incident.Motive)

	return incidentRegistered, err
}

func (service ScheduleService) LoadIncidentByPeriod(itemID int64, dayOfChange time.Time) (model.Incident, error) {
	return scheduleRepository.LoadIncidentByScheduleItemAndDate(itemID, dayOfChange)
}

func (service ScheduleService) LoadPendingIncidents(userID int64, incidentType model.IncidentType) ([]model.Incident, error) {

	if incidentType == model.Cancellation {
		return scheduleRepository.LoadPendingCancellationsByUser(userID)
	} else {
		return scheduleRepository.LoadPendingCancellationsByUser(userID)
	}
}

func (service ScheduleService) AcceptOrRejectChange(userID int64, changeID64 int64, answer bool) error {
	incident, err := scheduleRepository.LoadIncidentByIdAndUser(changeID64, userID)

	if err != nil {
		return err
	}

	if incident.RequestedBy.ID == userID {
		return fmt.Errorf("você não pode aceitar ou rejeitar uma requisição feita por você mesmo.: %d", changeID64)
	}

	if incident.Answered {
		return fmt.Errorf("requisição já respondida.: %d", changeID64)
	}

	incident.Accepted = answer

	_, err = scheduleRepository.UpdateIncident(incident)
	return err
}

func generateErrorFromEvents(events []model.ScheduleItem) error {
	var occurrences = ""
	for _, occurrence := range events {
		occurrences = occurrences + fmt.Sprintf("{ Dia: %s, Hora: %d, Minutes: %d }",
			utils.WeekDayFromNumber(occurrence.WeekDay), occurrence.Hour, occurrence.Minutes)
	}
	return fmt.Errorf("este instrutor já possui horários agendados para os dias: %s", occurrences)
}

func decomposeScheduleDays(items []model.ScheduleItemDTO) ([]model.DecomposedSchedule, error) {
	var errorFound error
	var weekDaySchedule []model.DecomposedSchedule
	for _, item := range items {
		dayNumber, err := utils.NumberFromWeekDay(item.WeekDay)

		if err != nil {
			errorFound = err
			break
		}
		event := model.DecomposedSchedule{WeekDay: dayNumber, Hour: item.Hour, Minutes: item.Minutes, Duration: item.Duration}

		weekDaySchedule = append(weekDaySchedule, event)
	}
	return weekDaySchedule, errorFound
}
