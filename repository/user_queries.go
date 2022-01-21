package repository

var userSQL = `SELECT p.id, coalesce(c.code, '') as code, coalesce(c.name, '') as name, c.uf, c.zip_code, p.description, CAST(p.is_instructor AS UNSIGNED) as IsInstructor,
								p.username, p.password, p.first_name, p.last_name, 	COALESCE(p.profile_picture, '')
								FROM profile p LEFT JOIN city c
                        on p.city_code = c.code`

var instructorsByCitySQL = `SELECT p.id, coalesce(c.code, '') as code, coalesce(c.name, '') as name, c.uf, c.zip_code, p.description, CAST(p.is_instructor AS UNSIGNED) as IsInstructor,
						p.username, p.password, p.first_name, p.last_name, 	COALESCE(p.profile_picture, ''),
						(SELECT count(*) from connections cn where cn.contact_id = p.id and cn.owner_id = ?) as isConnection
						FROM profile p LEFT JOIN city c
				on p.city_code = c.code`

var createUserSQL = `INSERT into profile(city_code, description, is_instructor, username, password, first_name, last_name, profile_picture, fcm_token) 
						values(?, ?, ?, ?, ?, ?, ?, ?, ? )
`

var SQL_UPDATE_USER = `UPDATE profile set first_name = ?, last_name = ?, description = ?, profile_picture = ?, password = ?, city_code = ? where id = ? `

var SQL_UPDATE_FCM = `UPDATE profile set fcm_token = ? where username = ?`

var SQL_LOAD_INSTRUCTOR_SCHEDULE = `SELECT si.week_day, si.duration, si.hour, si.minutes from schedule_item si inner join schedule sch on si.schedule_id = sch.id where sch.instructor_id = ?`
