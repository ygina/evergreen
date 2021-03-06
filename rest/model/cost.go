package model

import (
	"time"

	"github.com/evergreen-ci/evergreen/model/task"
	"github.com/pkg/errors"
)

// APIVersionCost is the model to be returned by the API whenever cost data is fetched by version id.
type APIVersionCost struct {
	VersionId    APIString     `json:"version_id"`
	SumTimeTaken time.Duration `json:"sum_time_taken"`
}

// BuildFromService converts from a service level task by loading the data
// into the appropriate fields of the APIVersionCost.
func (apiVersionCost *APIVersionCost) BuildFromService(h interface{}) error {
	switch v := h.(type) {
	case *task.VersionCost:
		apiVersionCost.VersionId = APIString(v.VersionId)
		apiVersionCost.SumTimeTaken = v.SumTimeTaken
	default:
		return errors.Errorf("incorrect type when fetching converting version cost type")
	}
	return nil
}

// ToService returns a service layer version cost using the data from APIVersionCost.
func (apiVersionCost *APIVersionCost) ToService() (interface{}, error) {
	return nil, errors.Errorf("ToService() is not implemented for APIVersionCost")
}

// APIDistroCost is the model to be returned by the API whenever cost data is fetched by distro id.
type APIDistroCost struct {
	DistroId     APIString     `json:"distro_id"`
	SumTimeTaken time.Duration `json:"sum_time_taken"`
}

// BuildFromService converts from a service level task by loading the data
// into the appropriate fields of the APIDistroCost.
func (apiDistroCost *APIDistroCost) BuildFromService(h interface{}) error {
	switch v := h.(type) {
	case *task.DistroCost:
		apiDistroCost.DistroId = APIString(v.DistroId)
		apiDistroCost.SumTimeTaken = v.SumTimeTaken
	default:
		return errors.Errorf("incorrect type when fetching converting distro cost type")
	}
	return nil
}

// ToService returns a service layer distro cost using the data from APIDistroCost.
func (apiDistroCost *APIDistroCost) ToService() (interface{}, error) {
	return nil, errors.Errorf("ToService() is not implemented for APIDistroCost")
}
