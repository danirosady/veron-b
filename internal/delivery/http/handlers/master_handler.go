package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/tms/tyre/internal/delivery/http/response"
	"github.com/tms/tyre/internal/domain/entity"
	"github.com/tms/tyre/internal/usecase"
)

type MasterHandler struct {
	masterUseCase *usecase.MasterUseCase
}

func NewMasterHandler(masterUseCase *usecase.MasterUseCase) *MasterHandler {
	return &MasterHandler{masterUseCase: masterUseCase}
}

// ListBrands returns all active tyre brands
// GET /api/v1/master/brands
func (h *MasterHandler) ListBrands(c *gin.Context) {
	brands, err := h.masterUseCase.ListBrands(c.Request.Context())
	if err != nil {
		response.InternalError(c, "Gagal mengambil data brand")
		return
	}
	response.Success(c, http.StatusOK, "Success", brands)
}

// ListSizes returns all active tyre sizes
// GET /api/v1/master/sizes
func (h *MasterHandler) ListSizes(c *gin.Context) {
	sizes, err := h.masterUseCase.ListSizes(c.Request.Context())
	if err != nil {
		response.InternalError(c, "Gagal mengambil data size")
		return
	}
	response.Success(c, http.StatusOK, "Success", sizes)
}

// ListTypes returns all active tyre types
// GET /api/v1/master/types
func (h *MasterHandler) ListTypes(c *gin.Context) {
	types, err := h.masterUseCase.ListTypes(c.Request.Context())
	if err != nil {
		response.InternalError(c, "Gagal mengambil data type")
		return
	}
	response.Success(c, http.StatusOK, "Success", types)
}

// ListPatterns returns active patterns, optionally filtered by brand
// GET /api/v1/master/patterns?brand_id=X
func (h *MasterHandler) ListPatterns(c *gin.Context) {
	brandID, _ := strconv.ParseUint(c.DefaultQuery("brand_id", "0"), 10, 32)

	patterns, err := h.masterUseCase.ListPatterns(c.Request.Context(), uint(brandID))
	if err != nil {
		response.InternalError(c, "Gagal mengambil data pattern")
		return
	}
	response.Success(c, http.StatusOK, "Success", patterns)
}

// ListReasons returns all active replacement reasons
// GET /api/v1/master/reasons
func (h *MasterHandler) ListReasons(c *gin.Context) {
	reasons, err := h.masterUseCase.ListReasons(c.Request.Context())
	if err != nil {
		response.InternalError(c, "Gagal mengambil data reason")
		return
	}
	response.Success(c, http.StatusOK, "Success", reasons)
}

// ListActions returns all active replacement actions
// GET /api/v1/master/actions
func (h *MasterHandler) ListActions(c *gin.Context) {
	actions, err := h.masterUseCase.ListActions(c.Request.Context())
	if err != nil {
		response.InternalError(c, "Gagal mengambil data action")
		return
	}
	response.Success(c, http.StatusOK, "Success", actions)
}

// ListRemarks returns all active remarks
// GET /api/v1/master/remarks
func (h *MasterHandler) ListRemarks(c *gin.Context) {
	remarks, err := h.masterUseCase.ListRemarks(c.Request.Context())
	if err != nil {
		response.InternalError(c, "Gagal mengambil data remark")
		return
	}
	response.Success(c, http.StatusOK, "Success", remarks)
}

// ListUnitTypeConfigs returns all unit type configurations
// GET /api/v1/master/unit-types
func (h *MasterHandler) ListUnitTypeConfigs(c *gin.Context) {
	configs, err := h.masterUseCase.ListUnitTypeConfigs(c.Request.Context())
	if err != nil {
		response.InternalError(c, "Gagal mengambil data unit type config")
		return
	}
	response.Success(c, http.StatusOK, "Success", configs)
}

// GetUnitTypeConfig returns the parsed position configuration for a unit type
// GET /api/v1/master/unit-types/:unit_type/config
func (h *MasterHandler) GetUnitTypeConfig(c *gin.Context) {
	unitType := c.Param("unit_type")
	if unitType == "" {
		response.BadRequest(c, "Unit type wajib diisi", nil)
		return
	}

	positions, err := h.masterUseCase.GetUnitTypeConfig(c.Request.Context(), unitType)
	if err != nil {
		if errors.Is(err, usecase.ErrUnitTypeConfigNotFound) {
			response.NotFound(c, "Unit type config tidak ditemukan")
			return
		}
		response.InternalError(c, "Gagal mengambil data unit type config")
		return
	}

	response.Success(c, http.StatusOK, "Success", positions)
}

// CreateBrand creates a new tyre brand
// POST /api/v1/master/brands
func (h *MasterHandler) CreateBrand(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	brand, err := h.masterUseCase.CreateBrand(c.Request.Context(), req.Name)
	if err != nil {
		response.InternalError(c, "Gagal membuat brand")
		return
	}
	response.Success(c, http.StatusCreated, "Brand berhasil dibuat", brand)
}

// UpdateBrand updates a tyre brand
// PUT /api/v1/master/brands/:id
func (h *MasterHandler) UpdateBrand(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	brand, err := h.masterUseCase.UpdateBrand(c.Request.Context(), uint(id), req.Name)
	if err != nil {
		response.InternalError(c, "Gagal mengupdate brand")
		return
	}
	response.Success(c, http.StatusOK, "Brand berhasil diupdate", brand)
}

// DeleteBrand deletes a tyre brand
// DELETE /api/v1/master/brands/:id
func (h *MasterHandler) DeleteBrand(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}
	if err := h.masterUseCase.DeleteBrand(c.Request.Context(), uint(id)); err != nil {
		response.InternalError(c, "Gagal menghapus brand")
		return
	}
	response.Success(c, http.StatusOK, "Brand berhasil dihapus", nil)
}

// CreateSize creates a new tyre size
// POST /api/v1/master/sizes
func (h *MasterHandler) CreateSize(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	size, err := h.masterUseCase.CreateSize(c.Request.Context(), req.Name)
	if err != nil {
		response.InternalError(c, "Gagal membuat size")
		return
	}
	response.Success(c, http.StatusCreated, "Size berhasil dibuat", size)
}

// UpdateSize updates a tyre size
// PUT /api/v1/master/sizes/:id
func (h *MasterHandler) UpdateSize(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	size, err := h.masterUseCase.UpdateSize(c.Request.Context(), uint(id), req.Name)
	if err != nil {
		response.InternalError(c, "Gagal mengupdate size")
		return
	}
	response.Success(c, http.StatusOK, "Size berhasil diupdate", size)
}

// DeleteSize deletes a tyre size
// DELETE /api/v1/master/sizes/:id
func (h *MasterHandler) DeleteSize(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}
	if err := h.masterUseCase.DeleteSize(c.Request.Context(), uint(id)); err != nil {
		response.InternalError(c, "Gagal menghapus size")
		return
	}
	response.Success(c, http.StatusOK, "Size berhasil dihapus", nil)
}

// CreateType creates a new tyre type
// POST /api/v1/master/types
func (h *MasterHandler) CreateType(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	t, err := h.masterUseCase.CreateType(c.Request.Context(), req.Name)
	if err != nil {
		response.InternalError(c, "Gagal membuat type")
		return
	}
	response.Success(c, http.StatusCreated, "Type berhasil dibuat", t)
}

// UpdateType updates a tyre type
// PUT /api/v1/master/types/:id
func (h *MasterHandler) UpdateType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	t, err := h.masterUseCase.UpdateType(c.Request.Context(), uint(id), req.Name)
	if err != nil {
		response.InternalError(c, "Gagal mengupdate type")
		return
	}
	response.Success(c, http.StatusOK, "Type berhasil diupdate", t)
}

// DeleteType deletes a tyre type
// DELETE /api/v1/master/types/:id
func (h *MasterHandler) DeleteType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}
	if err := h.masterUseCase.DeleteType(c.Request.Context(), uint(id)); err != nil {
		response.InternalError(c, "Gagal menghapus type")
		return
	}
	response.Success(c, http.StatusOK, "Type berhasil dihapus", nil)
}

// CreateReason creates a new replacement reason
// POST /api/v1/master/reasons
func (h *MasterHandler) CreateReason(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	reason, err := h.masterUseCase.CreateReason(c.Request.Context(), req.Name)
	if err != nil {
		response.InternalError(c, "Gagal membuat reason")
		return
	}
	response.Success(c, http.StatusCreated, "Reason berhasil dibuat", reason)
}

// UpdateReason updates a replacement reason
// PUT /api/v1/master/reasons/:id
func (h *MasterHandler) UpdateReason(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	reason, err := h.masterUseCase.UpdateReason(c.Request.Context(), uint(id), req.Name)
	if err != nil {
		response.InternalError(c, "Gagal mengupdate reason")
		return
	}
	response.Success(c, http.StatusOK, "Reason berhasil diupdate", reason)
}

// DeleteReason deletes a replacement reason
// DELETE /api/v1/master/reasons/:id
func (h *MasterHandler) DeleteReason(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}
	if err := h.masterUseCase.DeleteReason(c.Request.Context(), uint(id)); err != nil {
		response.InternalError(c, "Gagal menghapus reason")
		return
	}
	response.Success(c, http.StatusOK, "Reason berhasil dihapus", nil)
}

// CreateAction creates a new replacement action
// POST /api/v1/master/actions
func (h *MasterHandler) CreateAction(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	action, err := h.masterUseCase.CreateAction(c.Request.Context(), req.Name)
	if err != nil {
		response.InternalError(c, "Gagal membuat action")
		return
	}
	response.Success(c, http.StatusCreated, "Action berhasil dibuat", action)
}

// UpdateAction updates a replacement action
// PUT /api/v1/master/actions/:id
func (h *MasterHandler) UpdateAction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	action, err := h.masterUseCase.UpdateAction(c.Request.Context(), uint(id), req.Name)
	if err != nil {
		response.InternalError(c, "Gagal mengupdate action")
		return
	}
	response.Success(c, http.StatusOK, "Action berhasil diupdate", action)
}

// DeleteAction deletes a replacement action
// DELETE /api/v1/master/actions/:id
func (h *MasterHandler) DeleteAction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}
	if err := h.masterUseCase.DeleteAction(c.Request.Context(), uint(id)); err != nil {
		response.InternalError(c, "Gagal menghapus action")
		return
	}
	response.Success(c, http.StatusOK, "Action berhasil dihapus", nil)
}

// CreateRemark creates a new remark
// POST /api/v1/master/remarks
func (h *MasterHandler) CreateRemark(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	remark, err := h.masterUseCase.CreateRemark(c.Request.Context(), req.Name)
	if err != nil {
		response.InternalError(c, "Gagal membuat remark")
		return
	}
	response.Success(c, http.StatusCreated, "Remark berhasil dibuat", remark)
}

// UpdateRemark updates a remark
// PUT /api/v1/master/remarks/:id
func (h *MasterHandler) UpdateRemark(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	remark, err := h.masterUseCase.UpdateRemark(c.Request.Context(), uint(id), req.Name)
	if err != nil {
		response.InternalError(c, "Gagal mengupdate remark")
		return
	}
	response.Success(c, http.StatusOK, "Remark berhasil diupdate", remark)
}

// DeleteRemark deletes a remark
// DELETE /api/v1/master/remarks/:id
func (h *MasterHandler) DeleteRemark(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}
	if err := h.masterUseCase.DeleteRemark(c.Request.Context(), uint(id)); err != nil {
		response.InternalError(c, "Gagal menghapus remark")
		return
	}
	response.Success(c, http.StatusOK, "Remark berhasil dihapus", nil)
}

// CreatePattern creates a new tyre pattern
// POST /api/v1/master/patterns
func (h *MasterHandler) CreatePattern(c *gin.Context) {
	var req struct {
		BrandID uint   `json:"brand_id"`
		Name    string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	pattern, err := h.masterUseCase.CreatePattern(c.Request.Context(), req.BrandID, req.Name)
	if err != nil {
		response.InternalError(c, "Gagal membuat pattern")
		return
	}
	response.Success(c, http.StatusCreated, "Pattern berhasil dibuat", pattern)
}

// UpdatePattern updates a tyre pattern
// PUT /api/v1/master/patterns/:id
func (h *MasterHandler) UpdatePattern(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	pattern, err := h.masterUseCase.UpdatePattern(c.Request.Context(), uint(id), req.Name)
	if err != nil {
		response.InternalError(c, "Gagal mengupdate pattern")
		return
	}
	response.Success(c, http.StatusOK, "Pattern berhasil diupdate", pattern)
}

// DeletePattern deletes a tyre pattern
// DELETE /api/v1/master/patterns/:id
func (h *MasterHandler) DeletePattern(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}
	if err := h.masterUseCase.DeletePattern(c.Request.Context(), uint(id)); err != nil {
		response.InternalError(c, "Gagal menghapus pattern")
		return
	}
	response.Success(c, http.StatusOK, "Pattern berhasil dihapus", nil)
}

// CreateUnitType creates a new unit type configuration
// POST /api/v1/master/unit-types
func (h *MasterHandler) CreateUnitType(c *gin.Context) {
	var req struct {
		UnitType       string                `json:"unit_type" binding:"required"`
		DisplayName    string                `json:"display_name" binding:"required"`
		MaxPosition    int                   `json:"max_position" binding:"required,min=1"`
		PositionConfig json.RawMessage       `json:"position_config"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	var positionConfigs entity.PositionConfigs
	if req.PositionConfig != nil {
		if err := json.Unmarshal(req.PositionConfig, &positionConfigs); err != nil {
			response.BadRequest(c, "Format position_config tidak valid", nil)
			return
		}
	}

	config, err := h.masterUseCase.CreateUnitTypeConfig(c.Request.Context(), req.UnitType, req.DisplayName, req.MaxPosition, positionConfigs)
	if err != nil {
		response.InternalError(c, "Gagal membuat unit type")
		return
	}
	response.Success(c, http.StatusCreated, "Unit type berhasil dibuat", config)
}

// UpdateUnitType updates a unit type configuration
// PUT /api/v1/master/unit-types/:id
func (h *MasterHandler) UpdateUnitType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}
	var req struct {
		DisplayName    string                `json:"display_name" binding:"required"`
		MaxPosition    int                   `json:"max_position" binding:"required,min=1"`
		PositionConfig json.RawMessage       `json:"position_config"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	var positionConfigs entity.PositionConfigs
	if req.PositionConfig != nil {
		if err := json.Unmarshal(req.PositionConfig, &positionConfigs); err != nil {
			response.BadRequest(c, "Format position_config tidak valid", nil)
			return
		}
	}

	config, err := h.masterUseCase.UpdateUnitTypeConfig(c.Request.Context(), uint(id), req.DisplayName, req.MaxPosition, positionConfigs)
	if err != nil {
		response.InternalError(c, "Gagal mengupdate unit type")
		return
	}
	response.Success(c, http.StatusOK, "Unit type berhasil diupdate", config)
}

// DeleteUnitType deletes a unit type configuration
// DELETE /api/v1/master/unit-types/:id
func (h *MasterHandler) DeleteUnitType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}
	if err := h.masterUseCase.DeleteUnitTypeConfig(c.Request.Context(), uint(id)); err != nil {
		response.InternalError(c, "Gagal menghapus unit type")
		return
	}
	response.Success(c, http.StatusOK, "Unit type berhasil dihapus", nil)
}

// RegisterRoutes registers master data routes on the given router group
func (h *MasterHandler) RegisterRoutes(rg *gin.RouterGroup) {
	master := rg.Group("/master")
	{
		master.GET("/brands", h.ListBrands)
		master.POST("/brands", h.CreateBrand)
		master.PUT("/brands/:id", h.UpdateBrand)
		master.DELETE("/brands/:id", h.DeleteBrand)

		master.GET("/sizes", h.ListSizes)
		master.POST("/sizes", h.CreateSize)
		master.PUT("/sizes/:id", h.UpdateSize)
		master.DELETE("/sizes/:id", h.DeleteSize)

		master.GET("/types", h.ListTypes)
		master.POST("/types", h.CreateType)
		master.PUT("/types/:id", h.UpdateType)
		master.DELETE("/types/:id", h.DeleteType)

		master.GET("/patterns", h.ListPatterns)
		master.POST("/patterns", h.CreatePattern)
		master.PUT("/patterns/:id", h.UpdatePattern)
		master.DELETE("/patterns/:id", h.DeletePattern)

		master.GET("/reasons", h.ListReasons)
		master.POST("/reasons", h.CreateReason)
		master.PUT("/reasons/:id", h.UpdateReason)
		master.DELETE("/reasons/:id", h.DeleteReason)

		master.GET("/actions", h.ListActions)
		master.POST("/actions", h.CreateAction)
		master.PUT("/actions/:id", h.UpdateAction)
		master.DELETE("/actions/:id", h.DeleteAction)

		master.GET("/remarks", h.ListRemarks)
		master.POST("/remarks", h.CreateRemark)
		master.PUT("/remarks/:id", h.UpdateRemark)
		master.DELETE("/remarks/:id", h.DeleteRemark)

		master.GET("/unit-types", h.ListUnitTypeConfigs)
		master.POST("/unit-types", h.CreateUnitType)
		master.PUT("/unit-types/:id", h.UpdateUnitType)
		master.DELETE("/unit-types/:id", h.DeleteUnitType)
		master.GET("/unit-types/:unit_type/config", h.GetUnitTypeConfig)
	}
}