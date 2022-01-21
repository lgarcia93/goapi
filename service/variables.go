package service

import r "fitgoapi/repository"

var skillRepository r.ISkillRepository = r.SkillRepository{}
var cityRepository r.ICityRepository = r.CityRepository{}
var scheduleRepository r.IScheduleRepository = r.ScheduleRepository{}
var userRepository r.IUserRepository = r.UserRepository{}
