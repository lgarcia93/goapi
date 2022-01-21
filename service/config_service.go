package service

// IConfigService interface
type IConfigService interface {
	LoadConfigs() (map[string]interface{}, error)
}

// // Login ... struct
type ConfigService struct {
}

// LoadConfigs ...
func (service ConfigService) LoadConfigs() (map[string]interface{}, error) {
	skills, err := skillRepository.FetchSkills()
	cities, err := cityRepository.FetchCities()

	return map[string]interface{}{"skills": skills, "cities": cities}, err
}
