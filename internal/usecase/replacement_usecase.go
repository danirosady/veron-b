package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/tms/tyre/internal/domain/entity"
	"github.com/tms/tyre/internal/domain/repository"
	"github.com/tms/tyre/internal/dto/request"
)

// Common errors used by the replacement use case
var (
	ErrReplacementNotFound      = errors.New("replacement not found")
	ErrReplacementUnitRequired  = errors.New("unit is required")
	ErrReplacementTyreRequired = errors.New("new tyre is required for mount action")
	ErrReplacementTyreNotSpare = errors.New("new tyre must be in spare status for mount action")
	ErrReplacementTyreNotMounted = errors.New("old tyre must be in mounted status for dismount action")
	ErrReplacementTyreNotFound  = errors.New("tyre not found")
	ErrReplacementUnitNotFound  = errors.New("unit not found")
	ErrReplacementDriverNotFound = errors.New("driver not found")
	ErrReplacementInvalidAction = errors.New("invalid replacement action")
	ErrReplacementSwapBothRequired = errors.New("swap action requires both old_tyre_id and new_tyre_id")
)

// ReplacementUseCase handles business logic for tyre replacement operations.
type ReplacementUseCase struct {
	replacementRepo repository.ReplacementRepository
	tyreRepo       repository.TyreRepository
	unitRepo       repository.UnitRepository
	driverRepo     repository.DriverRepository
}

// NewReplacementUseCase creates a new ReplacementUseCase instance.
func NewReplacementUseCase(
	replacementRepo repository.ReplacementRepository,
	tyreRepo repository.TyreRepository,
	unitRepo repository.UnitRepository,
	driverRepo repository.DriverRepository,
) *ReplacementUseCase {
	return &ReplacementUseCase{
		replacementRepo: replacementRepo,
		tyreRepo:       tyreRepo,
		unitRepo:       unitRepo,
		driverRepo:     driverRepo,
	}
}

// Create creates a new replacement record with tyre status transitions.
func (uc *ReplacementUseCase) Create(ctx context.Context, req *request.CreateReplacementRequest, operatorID uint) (*entity.Replacement, error) {
	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	// Validate unit exists
	unit, err := uc.unitRepo.GetByID(*req.UnitID)
	if err != nil {
		return nil, err
	}
	if unit == nil {
		return nil, ErrReplacementUnitNotFound
	}

	// Validate driver exists if provided
	if req.DriverID != nil && *req.DriverID != 0 {
		driver, err := uc.driverRepo.GetByID(*req.DriverID)
		if err != nil {
			return nil, err
		}
		if driver == nil {
			return nil, ErrReplacementDriverNotFound
		}
	}

	// Build replacement details from the request
	details := make([]entity.ReplacementDetail, 0, len(req.Details))
	for _, d := range req.Details {
		detail := entity.ReplacementDetail{
			Position:        d.Position,
			Action:         d.Action,
			OldTyreID:    d.OldTyreID,
			NewTyreID:    d.NewTyreID,
			OldTyreTread1: d.OldTyreTread1,
			OldTyreTread2: d.OldTyreTread2,
			NewTyreTread1: d.NewTyreTread1,
			NewTyreTread2: d.NewTyreTread2,
			NewTyreStatus: d.NewTyreStatus,
			Remark:       d.Remark,
		}

		// Validate and process each detail based on action
		switch d.Action {
		case "mount":
			if d.NewTyreID == nil {
				return nil, ErrReplacementTyreRequired
			}
			newTyre, err := uc.tyreRepo.GetByID(*d.NewTyreID)
			if err != nil {
				return nil, err
			}
			if newTyre == nil {
				return nil, ErrReplacementTyreNotFound
			}
			if newTyre.Status != string(entity.TyreStatusSpare) && newTyre.Status != string(entity.TyreStatusDismounted) {
				return nil, ErrReplacementTyreNotSpare
			}
			// Populate new tyre info
			detail.NewTyreSerialNum = newTyre.SerialNumber
			detail.NewTyrePattern = getPatternName(newTyre)
			detail.NewTyreSize = getSizeName(newTyre)

		case "dismount":
			if d.OldTyreID != nil {
				oldTyre, err := uc.tyreRepo.GetByID(*d.OldTyreID)
				if err != nil {
					return nil, err
				}
				if oldTyre == nil {
					return nil, ErrReplacementTyreNotFound
				}
				if oldTyre.Status != string(entity.TyreStatusMounted) {
					return nil, ErrReplacementTyreNotMounted
				}
				detail.OldTyreSerialNum = oldTyre.SerialNumber
				detail.OldTyrePattern = getPatternName(oldTyre)
				detail.OldTyreSize = getSizeName(oldTyre)
				detail.OldTyreStatus = oldTyre.Status
				if oldTyre.RTD1 != nil {
					detail.OldTyreTread1 = oldTyre.RTD1
				}
				if oldTyre.RTD2 != nil {
					detail.OldTyreTread2 = oldTyre.RTD2
				}
			}

		case "swap":
			if d.OldTyreID == nil || d.NewTyreID == nil {
				return nil, ErrReplacementSwapBothRequired
			}
			oldTyre, err := uc.tyreRepo.GetByID(*d.OldTyreID)
			if err != nil {
				return nil, err
			}
			if oldTyre == nil {
				return nil, ErrReplacementTyreNotFound
			}
			if oldTyre.Status != string(entity.TyreStatusMounted) {
				return nil, ErrReplacementTyreNotMounted
			}
			newTyre, err := uc.tyreRepo.GetByID(*d.NewTyreID)
			if err != nil {
				return nil, err
			}
			if newTyre == nil {
				return nil, ErrReplacementTyreNotFound
			}
			if newTyre.Status != string(entity.TyreStatusSpare) && newTyre.Status != string(entity.TyreStatusDismounted) {
				return nil, ErrReplacementTyreNotSpare
			}
			detail.OldTyreSerialNum = oldTyre.SerialNumber
			detail.OldTyrePattern = getPatternName(oldTyre)
			detail.OldTyreSize = getSizeName(oldTyre)
			detail.OldTyreStatus = oldTyre.Status
			if oldTyre.RTD1 != nil {
				detail.OldTyreTread1 = oldTyre.RTD1
			}
			if oldTyre.RTD2 != nil {
				detail.OldTyreTread2 = oldTyre.RTD2
			}
			detail.NewTyreSerialNum = newTyre.SerialNumber
			detail.NewTyrePattern = getPatternName(newTyre)
			detail.NewTyreSize = getSizeName(newTyre)
		}

		details = append(details, detail)
	}

	// Create the replacement record
	replacement := &entity.Replacement{
		CompanyID:     unit.CompanyID,
		ProjectID:     unit.ProjectID,
		UnitID:        *req.UnitID,
		DriverID:       derefUint(req.DriverID),
		Date:          date,
		HMUpdate:      derefFloat(req.HMUpdate),
		CurrentLifeHM: derefFloat(req.CurrentLifeHM),
		HMPlan:        derefFloat(req.HMPlan),
		Remarks:       req.Remarks,
		CreatedBy:     operatorID,
		Details:       details,
	}

	result, err := uc.replacementRepo.Create(replacement)
	if err != nil {
		return nil, err
	}

	// Apply tyre status changes after creation
	for _, d := range details {
		if d.Action == "mount" && d.NewTyreID != nil {
			// Mount: update tyre to mounted
			newTyre, _ := uc.tyreRepo.GetByID(*d.NewTyreID)
			if newTyre != nil {
				newTyre.Status = string(entity.TyreStatusMounted)
				newTyre.UnitID = req.UnitID
				pos := d.Position
				newTyre.MountedPosition = &pos
				uc.tyreRepo.Update(newTyre)
			}
		} else if (d.Action == "dismount" || d.Action == "swap") && d.OldTyreID != nil {
			// Dismount/Swap: update old tyre to spare or scrap
			oldTyre, _ := uc.tyreRepo.GetByID(*d.OldTyreID)
			if oldTyre != nil {
				newStatus := string(entity.TyreStatusSpare)
				if d.NewTyreStatus != "" {
					newStatus = d.NewTyreStatus
				} else {
					// Determine status based on RTD
					if oldTyre.RTD < 2.0 {
						newStatus = string(entity.TyreStatusScrap)
					}
				}
				oldTyre.Status = newStatus
				oldTyre.UnitID = nil
				oldTyre.MountedPosition = nil
				uc.tyreRepo.Update(oldTyre)
			}
		}
	}

	// Reload with relations
	return uc.replacementRepo.GetByID(result.ID)
}

// GetByID returns a single replacement by ID.
func (uc *ReplacementUseCase) GetByID(ctx context.Context, id uint) (*entity.Replacement, error) {
	replacement, err := uc.replacementRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if replacement == nil {
		return nil, ErrReplacementNotFound
	}
	return replacement, nil
}

// List returns a paginated list of replacements.
func (uc *ReplacementUseCase) List(ctx context.Context, page, perPage int, filters map[string]interface{}) ([]*entity.Replacement, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	var companyID, projectID, unitID uint
	var dateFrom, dateTo *time.Time

	if filters != nil {
		if v, ok := filters["company_id"].(uint); ok {
			companyID = v
		}
		if v, ok := filters["project_id"].(uint); ok {
			projectID = v
		}
		if v, ok := filters["unit_id"].(uint); ok {
			unitID = v
		}
		if v, ok := filters["date_from"].(string); ok && v != "" {
			t, err := time.Parse("2006-01-02", v)
			if err == nil {
				dateFrom = &t
			}
		}
		if v, ok := filters["date_to"].(string); ok && v != "" {
			t, err := time.Parse("2006-01-02", v)
			if err == nil {
				dateTo = &t
			}
		}
	}

	return uc.replacementRepo.List(page, perPage, companyID, projectID, unitID, dateFrom, dateTo)
}

// Update updates an existing replacement record.
func (uc *ReplacementUseCase) Update(ctx context.Context, id uint, req *request.UpdateReplacementRequest, operatorID uint) (*entity.Replacement, error) {
	replacement, err := uc.replacementRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if replacement == nil {
		return nil, ErrReplacementNotFound
	}

	if req.Date != "" {
		date, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return nil, errors.New("invalid date format, use YYYY-MM-DD")
		}
		replacement.Date = date
	}

	replacement.HMUpdate = req.HMUpdate
	replacement.CurrentLifeHM = req.CurrentLifeHM
	replacement.HMPlan = req.HMPlan
	replacement.Remarks = req.Remarks

	if req.DriverID != nil && *req.DriverID != 0 {
		replacement.DriverID = *req.DriverID
	}

	// Note: updating details is complex (would need to revert/apply tyre status changes)
	// For now, update main replacement fields only
	if err := uc.replacementRepo.Update(replacement); err != nil {
		return nil, err
	}

	return uc.replacementRepo.GetByID(id)
}

// Delete removes a replacement record.
func (uc *ReplacementUseCase) Delete(ctx context.Context, id uint) error {
	replacement, err := uc.replacementRepo.GetByID(id)
	if err != nil {
		return err
	}
	if replacement == nil {
		return ErrReplacementNotFound
	}
	return uc.replacementRepo.Delete(id)
}

// GetLastByUnitID returns the most recent replacement for a given unit.
func (uc *ReplacementUseCase) GetLastByUnitID(unitID uint) (*entity.Replacement, error) {
	return uc.replacementRepo.GetLastReplacementByUnitID(unitID)
}

// getPatternName returns the pattern name from a tyre.
func getPatternName(tyre *entity.TyreMaster) string {
	if tyre.Pattern != nil {
		return tyre.Pattern.Name
	}
	return ""
}

// getSizeName returns the size name from a tyre.
func getSizeName(tyre *entity.TyreMaster) string {
	if tyre.Size != nil {
		return tyre.Size.Name
	}
	return ""
}

func derefFloat(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

func derefUint(v *uint) uint {
	if v == nil {
		return 0
	}
	return *v
}
