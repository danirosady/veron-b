package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/tms/tyre/internal/domain/entity"
	"github.com/tms/tyre/internal/domain/repository"
	"github.com/tms/tyre/internal/dto/request"
	"gorm.io/gorm"
)

// Common errors used by the tyre use case
var (
	ErrTyreNotFound          = errors.New("tyre not found")
	ErrTyreBarcodeRequired   = errors.New("barcode is required")
	ErrTyreBarcodeTooLong    = errors.New("barcode must not exceed 100 characters")
	ErrTyreBarcodeExists     = errors.New("barcode already exists")
	ErrTyreSerialRequired    = errors.New("serial number is required")
	ErrTyreSerialTooLong     = errors.New("serial number must not exceed 100 characters")
	ErrTyreSerialExists      = errors.New("serial number already exists")
	ErrTyreSizeNotFound      = errors.New("tyre size not found")
	ErrTyreBrandNotFound     = errors.New("tyre brand not found")
	ErrTyrePatternNotFound   = errors.New("tyre pattern not found")
	ErrTyreRTDExceedsOTD     = errors.New("rtd must be less than or equal to otd")
	ErrTyreRTDInvalid        = errors.New("rtd_1 and rtd_2 must be between 0 and otd")
	ErrTyrePSIOutOfRange     = errors.New("psi must be between 50 and 200")
	ErrTyreCompanyIDRequired = errors.New("company id is required")
)

// TyreUseCase handles business logic for tyre master records.
type TyreUseCase struct {
	tyreRepo    repository.TyreRepository
	companyRepo repository.CompanyRepository
	masterRepo  repository.MasterRepository
}

// NewTyreUseCase creates a new TyreUseCase instance.
func NewTyreUseCase(
	tyreRepo repository.TyreRepository,
	companyRepo repository.CompanyRepository,
	masterRepo repository.MasterRepository,
) *TyreUseCase {
	return &TyreUseCase{
		tyreRepo:    tyreRepo,
		companyRepo: companyRepo,
		masterRepo:  masterRepo,
	}
}

// List returns a paginated list of tyres.
func (uc *TyreUseCase) List(ctx context.Context, page, perPage int, companyID uint, status, brandID, sizeID string) ([]*entity.TyreMaster, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	return uc.tyreRepo.List(page, perPage, companyID, status, brandID, sizeID)
}

// GetByID returns a tyre by ID.
func (uc *TyreUseCase) GetByID(ctx context.Context, id uint) (*entity.TyreMaster, error) {
	tyre, err := uc.tyreRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if tyre == nil {
		return nil, ErrTyreNotFound
	}
	return tyre, nil
}

// GetByBarcode returns a tyre by barcode.
func (uc *TyreUseCase) GetByBarcode(ctx context.Context, barcode string) (*entity.TyreMaster, error) {
	if strings.TrimSpace(barcode) == "" {
		return nil, ErrTyreBarcodeRequired
	}
	tyre, err := uc.tyreRepo.GetByBarcode(barcode)
	if err != nil {
		return nil, err
	}
	if tyre == nil {
		return nil, ErrTyreNotFound
	}
	return tyre, nil
}

// GetSpareTyres returns all spare (and dismounted) tyres for a company.
func (uc *TyreUseCase) GetSpareTyres(ctx context.Context, companyID uint) ([]*entity.TyreMaster, error) {
	if companyID == 0 {
		return nil, ErrTyreCompanyIDRequired
	}
	return uc.tyreRepo.GetSpareTyres(companyID)
}

// Create creates a new tyre master record.
func (uc *TyreUseCase) Create(ctx context.Context, companyID uint, req *request.CreateTyreRequest) (*entity.TyreMaster, error) {
	if companyID == 0 {
		return nil, ErrTyreCompanyIDRequired
	}

	barcode := strings.TrimSpace(req.Barcode)
	if barcode == "" {
		return nil, ErrTyreBarcodeRequired
	}
	if len(barcode) > 100 {
		return nil, ErrTyreBarcodeTooLong
	}

	serial := strings.TrimSpace(req.SerialNumber)
	if serial == "" {
		return nil, ErrTyreSerialRequired
	}
	if len(serial) > 100 {
		return nil, ErrTyreSerialTooLong
	}

	company, err := uc.companyRepo.GetByID(companyID)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, ErrCompanyNotFound
	}

	if _, err := uc.masterRepo.GetSizeByID(req.SizeID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTyreSizeNotFound
		}
		return nil, err
	}
	if _, err := uc.masterRepo.GetBrandByID(req.BrandID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTyreBrandNotFound
		}
		return nil, err
	}
	if _, err := uc.masterRepo.GetPatternByID(req.PatternID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTyrePatternNotFound
		}
		return nil, err
	}

	if req.OTD < 0 {
		return nil, ErrTyreRTDExceedsOTD
	}

	if req.RTD > req.OTD {
		return nil, ErrTyreRTDExceedsOTD
	}

	if req.RTD1 != nil {
		if *req.RTD1 < 0 || *req.RTD1 > req.OTD {
			return nil, ErrTyreRTDInvalid
		}
	}
	if req.RTD2 != nil {
		if *req.RTD2 < 0 || *req.RTD2 > req.OTD {
			return nil, ErrTyreRTDInvalid
		}
	}

	if req.PSI != nil {
		if *req.PSI < 50 || *req.PSI > 200 {
			return nil, ErrTyrePSIOutOfRange
		}
	}

	rtd := req.RTD
	if rtd == 0 {
		rtd = req.OTD
	}
	if req.RTD1 != nil && req.RTD2 != nil {
		rtd = (*req.RTD1 + *req.RTD2) / 2.0
	}

	tyreType := strings.TrimSpace(req.Type)
	if tyreType == "" {
		tyreType = "Radial"
	}

	existing, err := uc.tyreRepo.GetByBarcode(barcode)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrTyreBarcodeExists
	}

	existingSN, err := uc.tyreRepo.GetBySerialNumber(serial)
	if err != nil {
		return nil, err
	}
	if existingSN != nil {
		return nil, ErrTyreSerialExists
	}

	tyre := &entity.TyreMaster{
		CompanyID:  companyID,
		Barcode:    barcode,
		SerialNumber:    serial,
		DOTCode:         req.DOTCode,
		Type:            tyreType,
		SizeID:          req.SizeID,
		BrandID:         req.BrandID,
		PatternID:       req.PatternID,
		OTD:             req.OTD,
		RTD:             rtd,
		RTD1:            req.RTD1,
		RTD2:            req.RTD2,
		Lifetime:        req.Lifetime,
		PSI:             req.PSI,
		Status:          string(entity.TyreStatusSpare),
		Remarks:         req.Remarks,
	}

	if err := uc.tyreRepo.Create(tyre); err != nil {
		return nil, err
	}

	return tyre, nil
}

// Update updates an existing tyre master record.
func (uc *TyreUseCase) Update(ctx context.Context, id uint, req *request.UpdateTyreRequest) (*entity.TyreMaster, error) {
	tyre, err := uc.tyreRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if tyre == nil {
		return nil, ErrTyreNotFound
	}

	if req.RTD < 0 || req.RTD > tyre.OTD {
		return nil, ErrTyreRTDExceedsOTD
	}

	if req.RTD1 != nil {
		if *req.RTD1 < 0 || *req.RTD1 > tyre.OTD {
			return nil, ErrTyreRTDInvalid
		}
	}
	if req.RTD2 != nil {
		if *req.RTD2 < 0 || *req.RTD2 > tyre.OTD {
			return nil, ErrTyreRTDInvalid
		}
	}

	if req.PSI != nil {
		if *req.PSI < 50 || *req.PSI > 200 {
			return nil, ErrTyrePSIOutOfRange
		}
	}

	rtd := req.RTD
	if req.RTD1 != nil && req.RTD2 != nil {
		rtd = (*req.RTD1 + *req.RTD2) / 2.0
	}

	tyre.DOTCode = req.DOTCode
	if strings.TrimSpace(req.Type) != "" {
		tyre.Type = req.Type
	}
	tyre.RTD = rtd
	tyre.RTD1 = req.RTD1
	tyre.RTD2 = req.RTD2
	tyre.PSI = req.PSI
	tyre.Remarks = req.Remarks
	if req.Status != "" {
		tyre.Status = req.Status
	}

	if err := uc.tyreRepo.Update(tyre); err != nil {
		return nil, err
	}

	return tyre, nil
}

// Delete removes a tyre by ID.
func (uc *TyreUseCase) Delete(ctx context.Context, id uint) error {
	tyre, err := uc.tyreRepo.GetByID(id)
	if err != nil {
		return err
	}
	if tyre == nil {
		return ErrTyreNotFound
	}
	return uc.tyreRepo.Delete(id)
}
