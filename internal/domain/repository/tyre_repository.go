package repository

import "github.com/tms/tyre/internal/domain/entity"

type TyreRepository interface {
	Create(tyre *entity.TyreMaster) error
	GetByID(id uint) (*entity.TyreMaster, error)
	GetByBarcode(barcode string) (*entity.TyreMaster, error)
	GetBySerialNumber(sn string) (*entity.TyreMaster, error)
	Update(tyre *entity.TyreMaster) error
	Delete(id uint) error
	List(page, perPage int, companyID uint, status, brandID, sizeID, search string) ([]*entity.TyreMaster, int64, error)
	GetByUnitID(unitID uint) ([]*entity.TyreMaster, error)
	GetSpareTyres(companyID uint) ([]*entity.TyreMaster, error)
	Mount(tyreID uint, unitID uint, position string, rtd1, rtd2 *float64) error
	MountWithHM(tyreID uint, unitID uint, position string, rtd1, rtd2 *float64, mountHM float64) error
	UpdateLifetime(tyreID uint, mountHM, currentHM float64) error
	GetMountHM(tyreID uint) (float64, error)
	Dismount(tyreID uint, status string, rtd1, rtd2 *float64) error
	GetTyreHistory(tyreID uint) ([]*entity.TyreHistoryItem, error)

	// Status log
	CreateStatusLog(log *entity.TyreStatusLog) error
	GetStatusLogs(tyreID uint) ([]*entity.TyreStatusLog, error)
}
