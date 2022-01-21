package repository

var scheduleSQL = `SELECT sch.id, sch.instructor_id, sch.student_id, sch.skill_id, sch.accepted from schedule sch`

var SQL_SCHEDULE_BY_USER = `

SELECT SCHEDULE_ID as sid,
	   SCHEDULE_UPDATE,
       ITEM_ID,
       WEEKDAY,
       HOUR,
       MINUTES,
       DURATION,
       USER_ID,
       USER_NAME,
       USER_LASTNAME,
       USER_PIC,
	   IS_CHANGE	
from(
	SELECT s.id as SCHEDULE_ID,
		   s.updated as SCHEDULE_UPDATE,
		   si.id as ITEM_ID,
		   si.week_day as WEEKDAY,
		   si.hour as HOUR,
		   si.minutes as MINUTES,
		   si.duration as DURATION,
		   p.id as USER_ID,
		   COALESCE(p.first_name, '') as USER_NAME,
		   COALESCE(p.last_name, '') as USER_LASTNAME,
		   COALESCE(p.profile_picture, '') as USER_PIC,
	   	   false as IS_CHANGE	
	from schedule as s
			 inner join schedule_item si on s.id = si.schedule_id
			 inner join profile p on (
				(p.id = s.instructor_id AND p.id <> ?)
				 or
				(p.id = s.student_id AND p.id <> ?)
			 )
	where s.accepted = 1 
	and (s.student_id = ? OR s.instructor_id = ?)
	and si.id not in (
		Select inc.id from incident inc 
					  inner join schedule_item item on item.id =  inc.schedule_item_id
				 where date(inc.day_of_change) between date(?) and date(?) 
				 and inc.accepted = true
		)
UNION ALL
SELECT s.id as SCHEDULE_ID,
	   inc.created as SCHEDULE_UPDATE,
       si.id as ITEM_ID,
       inc.week_day as WEEKDAY,
       inc.hour as HOUR,
       inc.minutes as MINUTES,
       inc.duration as DURATION,
       p.id as USER_ID,
       COALESCE(p.first_name, '') as USER_NAME,
       COALESCE(p.last_name, '') as USER_LASTNAME,
       COALESCE(p.profile_picture, '') as USER_PIC,
	   true as IS_CHANGE
from incident inc
inner join schedule_item si on si.id =  inc.schedule_item_id
inner join schedule s on s.id = si.schedule_id
inner join profile p on p.id = inc.requested_by
where
(s.instructor_id = ? or s.student_id = ?)
and date(inc.day_of_change) between date(?) and date(?)
and inc.accepted = true
and inc.type = 'change'
) as aux
`

var SQL_SCHEDULE_BY_USER_AND_DAY = ` SELECT SCHEDULE_ID,
	   SCHEDULE_UPDATE,
       ITEM_ID,
       WEEKDAY,
       HOUR,
       MINUTES,
       DURATION,
       USER_ID,
       USER_NAME,
       USER_LASTNAME,
       USER_PIC,
	   IS_CHANGE	
from(
	SELECT s.id as SCHEDULE_ID,
		   s.updated as SCHEDULE_UPDATE,
		   si.id as ITEM_ID,
		   si.week_day as WEEKDAY,
		   si.hour as HOUR,
		   si.minutes as MINUTES,
		   si.duration as DURATION,
		   p.id as USER_ID,
		   COALESCE(p.first_name, '') as USER_NAME,
		   COALESCE(p.last_name, '') as USER_LASTNAME,
		   COALESCE(p.profile_picture, '') as USER_PIC,
	   	   false as IS_CHANGE	
	from schedule as s
			 inner join schedule_item si on s.id = si.schedule_id
			 inner join profile p on (
				(p.id = s.instructor_id AND p.id <> ?)
				 or
				(p.id = s.student_id AND p.id <> ?)
			 )
	where s.accepted = 1 
	and (s.student_id = ? OR s.instructor_id = ?)
	and si.week_day = ? 
	and si.id not in (
		Select inc.id from incident inc 
					  inner join schedule_item item on item.id =  inc.schedule_item_id
				 where date(inc.day_of_change) between date(?) and date(?) 
				 and inc.accepted = true
		)
UNION ALL
SELECT s.id as SCHEDULE_ID,
	   inc.created as SCHEDULE_UPDATE,
       si.id as ITEM_ID,
       inc.week_day as WEEKDAY,
       inc.hour as HOUR,
       inc.minutes as MINUTES,
       inc.duration as DURATION,
       p.id as USER_ID,
       COALESCE(p.first_name, '') as USER_NAME,
       COALESCE(p.last_name, '') as USER_LASTNAME,
       COALESCE(p.profile_picture, '') as USER_PIC,
	   true as IS_CHANGE
from incident inc
inner join schedule_item si on si.id =  inc.schedule_item_id
inner join schedule s on s.id = si.schedule_id
inner join profile p on p.id = inc.requested_by
where
(s.instructor_id = ? or s.student_id = ?)
and si.week_day = ? 
and date(inc.day_of_change) between date(?) and date(?)
and inc.accepted = true
and inc.type = 'change'
) as aux `

var SQL_MAKE_CONNECTION = `INSERT INTO connections(owner_id, contact_id) values(?, ?)`

var SQL_LOAD_USER_CONNECTIONS = `SELECT p.id, p.first_name, p.last_name, p.description, CAST(p.is_instructor AS UNSIGNED), coalesce(p.profile_picture, '')
from profile p inner join connections cn on p.id = cn.contact_id where cn.owner_id = ? limit ?  offset ?`

var SQL_SCHEDULE_BY_ID = `
	SELECT s.id as sid,
       s.updated as supdate,
       si.id as ssid,
       si.week_day,
       si.hour,
       si.minutes,
       si.duration,
       p.id as psid,
       COALESCE(p.first_name, '') as pifirstname,
       COALESCE(p.last_name, '') as pilastname,
       COALESCE(p.profile_picture, '') as piprofile_pic
from schedule as s
         inner join schedule_item si on s.id = si.schedule_id
         inner join profile p on (
            (p.id = s.instructor_id AND p.id <> ?)
             or
            (p.id = s.student_id AND p.id <> ?)
         )
where s.id = ?
`

var SQL_SCHEDULE_ITEM_BY_ID_AND_USER = `
SELECT si.id as ssid,
       si.week_day,
       si.hour,
       si.minutes,
       si.duration,
       s.id as sid,
	   p.id,
       COALESCE(p.first_name, '') as pifirstname,
       COALESCE(p.last_name, '') as pilastname,
       COALESCE(p.profile_picture, '') as piprofile_pic
from schedule as s
	 inner join schedule_item si on s.id = si.schedule_id
	 inner join profile p on (
		(p.id = s.instructor_id AND p.id <> ?)
		 or
		(p.id = s.student_id AND p.id <> ?)
	 )
where si.id = ?
AND (s.instructor_id = ? or s.student_id = ?)
`
var _SQL_INCIDENT = `
SELECT inc.id,
       inc.schedule_item_id,
       inc.accepted,
       inc.answered,
       inc.created,
       inc.day_change,
       inc.week_day,
       inc.hour,
       inc.minutes,
       inc.type,
       inc.motive,
	   p.id,	
       COALESCE(p.first_name, '') as pifirstname,
       COALESCE(p.last_name, '') as pilastname,
       COALESCE(p.profile_picture, '') as piprofile_pic
from incident inc
inner join schedule_item si on si.id =  inc.schedule_item_id
inner join schedule s on s.id = si.schedule_id
inner join profile p on p.id = inc.requested_by
`

var SQL_LOAD_INCIDENT = _SQL_INCIDENT + ` where inc.id = ? `

var SQL_SEARCH_INCIDENTS_BY_START_AND_END_DATE = _SQL_INCIDENT + ` where(inc.day_change between ? and ?)`

var SQL_LOAD_INCIDENT_BY_USER_AND_TIME = _SQL_INCIDENT + `
	where
	(s.instructor_id = ? or s.student_id = ?) and 
	date(inc.day_change) = date(?)
`

var SQL_LOAD_INCIDENT_BY_ITEM_TIME_AND_TYPE = _SQL_INCIDENT + `
	where
    si.id = ? and
	date(inc.day_change) between date(?) and date(?)
	AND inc.accepted = true
`

var SQL_LOAD_PENDING_INCIDENTS = _SQL_INCIDENT + `
	where
	(s.instructor_id = ? or s.student_id = ?) and 
	inc.answered = false and
    date(inc.day_change) >= now()
	AND inc.type = ?
	ORDER BY inc.day_change asc
`

var SQL_LOAD_INCIDENT_BY_USER_AND_ID = _SQL_INCIDENT + `
	where
	(s.instructor_id = ? or s.student_id = ?) and
    inc.id = ?
`
