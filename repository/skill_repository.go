package repository

import (
	database "fitgoapi/database"
	"fitgoapi/model"
)

// ISkillRepository ...
type ISkillRepository interface {
	FetchUserSkills(userID int64) ([]model.Skill, error)
	FetchSkillByID(skillID int64) (model.Skill, error)
	FetchSkills() ([]model.Skill, error)
	CreateSkillsForUser(userID int64, skills []model.Skill) error
}

// SkillRepository ...
type SkillRepository struct {
}

// FetchSkillByID ...
func (r SkillRepository) FetchSkillByID(skillID int64) (model.Skill, error) {
	db := database.Connection

	row := db.QueryRow("SELECT s.id, s.name from skill s where s.id = ?", skillID)

	var skill model.Skill

	err := row.Scan(&skill.ID, &skill.Name)

	return skill, err
}

// FetchUserSkills ...
func (r SkillRepository) FetchUserSkills(userID int64) ([]model.Skill, error) {

	db := database.Connection

	var errorFound error

	skills := []model.Skill{}

	results, err := db.Query("SELECT s.id, s.name from skill s left join profile_skill ps on s.id = ps.skill_id where ps.user_id = ?", userID)
	errorFound = err

	if err == nil {
		for results.Next() {
			var skill model.Skill
			err = results.Scan(&skill.ID, &skill.Name)

			if err == nil {
				skills = append(skills, skill)
			} else {
				errorFound = err
			}
		}
	}

	return skills, errorFound
}

func (r SkillRepository) FetchSkills() ([]model.Skill, error) {

	db := database.Connection

	skills := []model.Skill{}

	results, err := db.Query("SELECT id, name FROM skill")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {

		var skill model.Skill
		err = results.Scan(&skill.ID, &skill.Name)
		if err != nil {
			panic(err.Error())
		}
		skills = append(skills, skill)
	}

	return skills, err
}

func (r SkillRepository) CreateSkillsForUser(userID int64, skills []model.Skill) error {
	db := database.Connection

	var err error

	for _, skill := range skills {
		_, err = db.Exec("INSERT into profile_skill(skill_id, user_id) values(?, ?)", skill.ID, userID)
	}

	return err
}
