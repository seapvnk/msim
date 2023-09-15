package system

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
	"msim/app/shared"
)

type SystemEntity struct {
	ID    uuid.UUID
	Key   string
	Value string
	Type  string
}

// Return value as int.
func (s *SystemEntity) AsInt() int {
	if strings.ToLower(s.Type) != "int" {
		return 0
	}

	num, err := strconv.Atoi(s.Value)
	if err != nil {
		return 0
	}

	return num
}

// Return value as float.
func (s *SystemEntity) AsFloat() float64 {
	if strings.ToLower(s.Type) != "float" {
		return 0
	}

	num, err := strconv.ParseFloat(s.Value, 64)
	if err != nil {
		return 0
	}

	return num
}

// Return value as boolean.
func (s *SystemEntity) AsBool() bool {
	if strings.ToLower(s.Type) != "boolean" {
		return false
	}

	boolean, err := strconv.ParseBool(s.Value)
	if err != nil {
		return false
	}

	return boolean
}

type SystemService struct {
	systemRepository *SystemRepository
}

type SystemEnvDTO struct {
	Key   string
	Value string
	Type  string
}

func (service *SystemService) Create(s *SystemEnvDTO) (*SystemEntity, *shared.Exception) {
	entity := &SystemEntity{ID: uuid.New(), Key: s.Key, Value: s.Value, Type: s.Type}
	result, err := service.systemRepository.Create(entity)

	if err != nil {
		return nil, shared.DefaultException(shared.ALREADY_CREATED_EX, "env")
	}

	return result, nil
}

type SystemKeyDTO struct {
	Key string
}

// Get a system variable by key.
func (service *SystemService) GetByKey(dto *SystemKeyDTO) (*SystemEntity, *shared.Exception) {
	result, err := service.systemRepository.GetByKey(dto.Key)

	if err != nil {
		return nil, shared.DefaultException(shared.NOT_FOUND_EX, "env")
	}

	return result, nil
}

// Retrieves all system variables.
func (service *SystemService) GetAll() ([]*SystemEntity, *shared.Exception) {
	result, err := service.systemRepository.GetAll()

	if err != nil {
		return nil, shared.DefaultException(shared.INTERNAL_EX, "system")
	}

	return result, nil
}

type SystemKeyUpdateDTO struct {
	Key   string
	Value string
}

// Edit a system variable.
func (service *SystemService) UpdateValueByKey(dto *SystemKeyUpdateDTO) (*SystemEntity, *shared.Exception) {
	result, err := service.systemRepository.UpdateValueByKey(dto.Key, dto.Value)

	if err != nil {
		msg := "system variable not found"
		return nil, shared.DefaultException(shared.INTERNAL_EX, msg)
	}

	return result, nil
}
