package repository

import (
	"database/sql"
	"fitgoapi/database"
	"fitgoapi/model"
	"fmt"
	"math"
)

type IUserRepository interface {
	SaveUser(user model.User) (int64, model.User, error)
	FetchUserByEmail(email string) (model.User, error)
	FetchUserByEmailAndPassword(email string, password string) (model.User, error)
	LoadInstructorsByCityCode(userId int64, cityCode string, offset int, size int) (model.PageableUser, error)
	FetchUserByID(userID int64) (model.User, error)
	VerifyEmail(email string) (bool, error)
	MakeConnection(userID int64, contactID int64) error
	LoadUserConnections(userID int64, page int, size int) (model.PageableUserPlain, error)
	UpdateUser(userID int64, updatedUser model.User) error
	UpdateFcmToken(username string, fcmToken string) error
}

type UserRepository struct{}

// Implementando IUserRepository
func (r UserRepository) SaveUser(user model.User) (int64, model.User, error) {
	db := database.Connection

	var isInstructor int

	if user.IsInstructor {
		isInstructor = 1
	} else {
		isInstructor = 0
	}

	result, err := db.Exec(
		createUserSQL,
		user.City.Code,
		user.Description,
		isInstructor,
		user.Username,
		user.Password,
		user.FirstName,
		user.LastName,
		user.ProfilePicture,
		user.FcmTokenFirebase,
	)

	if err == nil {

		userID, _ := result.LastInsertId()

		if err == nil {
			SkillRepository{}.CreateSkillsForUser(userID, user.Skills)
		}

	} else {
		fmt.Printf("%s", err.Error())
	}

	count, err := result.RowsAffected()

	var createdUser model.User

	if count > 0 && err == nil {
		createdUser, err = r.FetchUserByEmailAndPassword(user.Username, user.Password)
	}

	return count, createdUser, err
}

func (r UserRepository) FetchUserByEmailAndPassword(email string, password string) (model.User, error) {
	db := database.Connection

	user, err := loadUserFromResultSet(db.QueryRow(userSQL+` where p.username = ? and p.password = ?`, email, password))

	if err == nil {
		user.Skills, _ = skillRepository.FetchUserSkills(user.ID)
	}
	return user, err
}

func (r UserRepository) FetchUserByEmail(email string) (model.User, error) {
	db := database.Connection

	user, err := loadUserFromResultSet(db.QueryRow(userSQL+` where p.username = ?`, email))

	if err == nil {
		user.Skills, _ = skillRepository.FetchUserSkills(user.ID)
	}

	return user, err
}

func (r UserRepository) FetchUserByID(userID int64) (model.User, error) {
	db := database.Connection

	user, err := loadUserFromResultSet(db.QueryRow(userSQL+` where p.id = ?`, userID))

	if err == nil {
		user.Skills, _ = skillRepository.FetchUserSkills(user.ID)
	}

	return user, err
}

func (r UserRepository) LoadInstructorsByCityCode(userId int64, cityCode string, offset int, size int) (model.PageableUser, error) {
	pageable := model.PageableUser{
		Users: []model.User{},
	}
	var errorFound error

	db := database.Connection

	err := db.QueryRow("Select count(*) from ("+userSQL+" where p.is_instructor = 1 and c.code = ?) as aux", cityCode).Scan(&pageable.Total)

	if err != nil {
		errorFound = err
	} else {
		if pageable.Total >= offset {

			users, err := loadUserFromResultRows(db.Query(instructorsByCitySQL+` where p.is_instructor = 1 and c.code = ? limit ? offset ?`, userId, cityCode, size, offset))

			if err != nil {
				errorFound = err
			} else {
				for index := range users {
					users[index].Skills, _ = skillRepository.FetchUserSkills(users[index].ID)
				}
			}

			pageable.Users = users
		}
	}

	return pageable, errorFound
}

// VerifyEmail returns true if email count is < 1
func (r UserRepository) VerifyEmail(email string) (bool, error) {
	db := database.Connection

	var emailCount int

	err := db.QueryRow("Select count(*) from profile p where p.username = ?", email).Scan(&emailCount)

	if err != nil {
		return true, err
	}
	return emailCount < 1, err
}

func (r UserRepository) MakeConnection(userID int64, contactID int64) error {
	db := database.Connection

	_, err := db.Exec(SQL_MAKE_CONNECTION, userID, contactID)

	return err
}

func (r UserRepository) LoadUserConnections(userID int64, page int, size int) (model.PageableUserPlain, error) {

	db := database.Connection

	pageable := model.PageableUserPlain{
		Users: []model.SimpleUser{},
	}

	err := db.QueryRow("Select count(*) from (select * from connections c where c.owner_id = ? ) as aux", userID).Scan(&pageable.Total)

	offset := page*size - (size)

	rows, err := db.Query(SQL_LOAD_USER_CONNECTIONS, userID, size, offset)

	fmt.Printf("page %v size %v offset %v", page, size, offset)

	totalPages := float64(pageable.Total) / float64(size)

	pageable.TotalPages = int(math.Ceil(totalPages))

	if totalPages < 1 {
		pageable.TotalPages = 1
	}

	if page > pageable.TotalPages {
		page = pageable.TotalPages
	}

	pageable.Page = page

	if err != nil {
		return model.PageableUserPlain{}, err
	}

	for rows.Next() {

		var user model.SimpleUser

		isInstructor := 0

		err = rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Description,
			&isInstructor,
			&user.ProfilePicture,
		)

		user.IsInstructor = isInstructor != 0

		pageable.Users = append(pageable.Users, user)
	}

	return pageable, err
}

func (r UserRepository) UpdateUser(userID int64, updatedUser model.User) error {
	db := database.Connection

	_, err := db.Exec(
		SQL_UPDATE_USER,
		updatedUser.FirstName,
		updatedUser.LastName,
		updatedUser.Description,
		updatedUser.ProfilePicture,
		updatedUser.Password,
		updatedUser.City.Code,
		userID,
	)

	return err
}

func (r UserRepository) UpdateFcmToken(username string, fcmToken string) error {
	db := database.Connection

	_, err := db.Exec(SQL_UPDATE_FCM, fcmToken, username)

	return err
}

func loadUserFromResultRows(rows *sql.Rows, errQuery error) ([]model.User, error) {
	var users []model.User
	var errorFounded error

	if errQuery != nil {
		errorFounded = errQuery
	} else {

		for rows.Next() {
			user := model.User{
				City:   model.City{},
				Skills: []model.Skill{},
			}

			isInstructor := 0
			isConnection := 0

			err := rows.Scan(
				&user.ID,
				&user.City.Code,
				&user.City.Name,
				&user.City.UF,
				&user.City.ZipCode,
				&user.Description,
				&isInstructor,
				&user.Username,
				&user.Password,
				&user.FirstName,
				&user.LastName,
				&user.ProfilePicture,
				&isConnection,
			)
			if err != nil {
				errorFounded = err
			}

			user.IsInstructor = isInstructor != 0
			user.IsConnection = isConnection != 0

			users = append(users, user)
		}
	}
	return users, errorFounded
}

func loadUserFromResultSet(row *sql.Row) (model.User, error) {
	user := model.User{
		City:   model.City{},
		Skills: []model.Skill{},
	}

	isInstructor := 0

	err := row.Scan(
		&user.ID,
		&user.City.Code,
		&user.City.Name,
		&user.City.UF,
		&user.City.ZipCode,
		&user.Description,
		&isInstructor,
		&user.Username,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.ProfilePicture,
	)

	if err != nil {

	}

	user.IsInstructor = isInstructor != 0
	return user, err
}
