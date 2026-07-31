package response

import "github.com/tms/tyre/internal/domain/entity"

type MasterBrandResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func ToMasterBrandResponse(e *entity.MasterBrand) *MasterBrandResponse {
	if e == nil {
		return nil
	}
	return &MasterBrandResponse{
		ID:     e.ID,
		Name:   e.Name,
		Status: e.Status,
	}
}

func ToMasterBrandResponses(entities []*entity.MasterBrand) []*MasterBrandResponse {
	responses := make([]*MasterBrandResponse, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, ToMasterBrandResponse(e))
	}
	return responses
}

type MasterSizeResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func ToMasterSizeResponse(e *entity.MasterSize) *MasterSizeResponse {
	if e == nil {
		return nil
	}
	return &MasterSizeResponse{
		ID:     e.ID,
		Name:   e.Name,
		Status: e.Status,
	}
}

func ToMasterSizeResponses(entities []*entity.MasterSize) []*MasterSizeResponse {
	responses := make([]*MasterSizeResponse, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, ToMasterSizeResponse(e))
	}
	return responses
}

type MasterTypeResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func ToMasterTypeResponse(e *entity.MasterType) *MasterTypeResponse {
	if e == nil {
		return nil
	}
	return &MasterTypeResponse{
		ID:     e.ID,
		Name:   e.Name,
		Status: e.Status,
	}
}

func ToMasterTypeResponses(entities []*entity.MasterType) []*MasterTypeResponse {
	responses := make([]*MasterTypeResponse, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, ToMasterTypeResponse(e))
	}
	return responses
}

type MasterPatternResponse struct {
	ID      uint                 `json:"id"`
	BrandID uint                 `json:"brand_id"`
	Name    string               `json:"name"`
	Status  string               `json:"status"`
	Brand   *MasterBrandResponse `json:"brand,omitempty"`
}

func ToMasterPatternResponse(e *entity.MasterPattern) *MasterPatternResponse {
	if e == nil {
		return nil
	}
	return &MasterPatternResponse{
		ID:      e.ID,
		BrandID: e.BrandID,
		Name:    e.Name,
		Status:  e.Status,
		Brand:   ToMasterBrandResponse(e.Brand),
	}
}

func ToMasterPatternResponses(entities []*entity.MasterPattern) []*MasterPatternResponse {
	responses := make([]*MasterPatternResponse, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, ToMasterPatternResponse(e))
	}
	return responses
}

type MasterReasonResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func ToMasterReasonResponse(e *entity.MasterReason) *MasterReasonResponse {
	if e == nil {
		return nil
	}
	return &MasterReasonResponse{
		ID:     e.ID,
		Name:   e.Name,
		Status: e.Status,
	}
}

func ToMasterReasonResponses(entities []*entity.MasterReason) []*MasterReasonResponse {
	responses := make([]*MasterReasonResponse, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, ToMasterReasonResponse(e))
	}
	return responses
}

type MasterActionResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func ToMasterActionResponse(e *entity.MasterAction) *MasterActionResponse {
	if e == nil {
		return nil
	}
	return &MasterActionResponse{
		ID:     e.ID,
		Name:   e.Name,
		Status: e.Status,
	}
}

func ToMasterActionResponses(entities []*entity.MasterAction) []*MasterActionResponse {
	responses := make([]*MasterActionResponse, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, ToMasterActionResponse(e))
	}
	return responses
}

type MasterRemarkResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func ToMasterRemarkResponse(e *entity.MasterRemark) *MasterRemarkResponse {
	if e == nil {
		return nil
	}
	return &MasterRemarkResponse{
		ID:     e.ID,
		Name:   e.Name,
		Status: e.Status,
	}
}

func ToMasterRemarkResponses(entities []*entity.MasterRemark) []*MasterRemarkResponse {
	responses := make([]*MasterRemarkResponse, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, ToMasterRemarkResponse(e))
	}
	return responses
}

type UnitTypeConfigResponse struct {
	ID             string                 `json:"id"`
	UnitType       string                 `json:"unit_type"`
	DisplayName    string                 `json:"display_name"`
	MaxPosition    int                    `json:"max_position"`
	PositionConfig entity.PositionConfigs `json:"position_config"`
}

func ToUnitTypeConfigResponse(e *entity.UnitTypeConfig) *UnitTypeConfigResponse {
	if e == nil {
		return nil
	}
	return &UnitTypeConfigResponse{
		ID:             "",
		UnitType:       e.UnitType,
		DisplayName:    e.DisplayName,
		MaxPosition:    e.MaxPosition,
		PositionConfig: e.PositionConfig,
	}
}

func ToUnitTypeConfigResponses(entities []*entity.UnitTypeConfig) []*UnitTypeConfigResponse {
	responses := make([]*UnitTypeConfigResponse, 0, len(entities))
	for _, e := range entities {
		responses = append(responses, ToUnitTypeConfigResponse(e))
	}
	return responses
}
