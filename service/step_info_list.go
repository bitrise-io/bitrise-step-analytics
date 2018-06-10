package service

import (
	"net/http"

	"github.com/slapec93/bitrise-step-analytics/models"
)

// NewStepInfosFromStepInfoDatas ...
func NewStepInfosFromStepInfoDatas(stepInfoDatas []models.StepInfoData) []StepInfo {
	stepInfos := []StepInfo{}
	for _, aStepInfo := range stepInfoDatas {
		stepInfo := StepInfo{
			StepName:   aStepInfo.StepName,
			Duration:   aStepInfo.Duration,
			IsCI:       aStepInfo.IsCI,
			LaunchDate: aStepInfo.LaunchDate,
		}
		stepInfos = append(stepInfos, stepInfo)
	}
	return stepInfos
}

// StepInfoListHandler ...
func StepInfoListHandler(w http.ResponseWriter, r *http.Request) error {
	stepInfos := models.ListStepInfos()

	return RespondWithSuccess(w, NewStepInfosFromStepInfoDatas(stepInfos))
}
