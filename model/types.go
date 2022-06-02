package model

type (
	MaxEffectiveTimeCalculator struct {
		OperationalProcess string `json:"operational_process" validate:"required"`
		Seconds            int    `json:"seconds" validate:"required"`
		Warehouse          string `json:"warehouse" validate:"required"`
		Attribute          string `json:"attribute" validate:"required"`
		DateFrom           string `json:"date_from" validate:"required"`
	}

	MaxEffectiveTimeDTO struct {
		ProcessName           string `json:"process_name" validate:"required,max=25,strregex=alpha_num"`
		FacilityID            string `json:"facility" validate:"required,max=25,strregex=alpha_num"`
		ProcessAttributeValue string `json:"attribute" validate:"required,max=25,strregex=charaters"`
		Seconds               int    `json:"seconds" validate:"required,numeric"`
		ValidFrom             string `json:"valid_from" validate:"required"`
		ValidTo               string `json:"valid_to,omitempty"`
	}
)
