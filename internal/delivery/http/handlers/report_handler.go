package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/tms/tyre/internal/delivery/http/response"
	"github.com/tms/tyre/internal/usecase"
)

type ReportHandler struct {
	reportUseCase *usecase.ReportUseCase
}

func NewReportHandler(reportUseCase *usecase.ReportUseCase) *ReportHandler {
	return &ReportHandler{reportUseCase: reportUseCase}
}

// GetReplacementReport returns a replacement report
// GET /api/v1/reports/replacement
func (h *ReportHandler) GetReplacementReport(c *gin.Context) {
	filters := buildReportFilters(c)

	report, err := h.reportUseCase.GetReplacementReport(c.Request.Context(), filters)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data laporan replacement")
		return
	}

	response.Success(c, http.StatusOK, "Success", report)
}

// GetInventoryReport returns an inventory report
// GET /api/v1/reports/inventory
func (h *ReportHandler) GetInventoryReport(c *gin.Context) {
	filters := buildReportFilters(c)

	report, err := h.reportUseCase.GetInventoryReport(c.Request.Context(), filters)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data laporan inventory")
		return
	}

	response.Success(c, http.StatusOK, "Success", report)
}

// GetScheduleReport returns a tyre schedule report
// GET /api/v1/reports/schedule
func (h *ReportHandler) GetScheduleReport(c *gin.Context) {
	filters := buildReportFilters(c)

	report, err := h.reportUseCase.GetScheduleReport(c.Request.Context(), filters)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data laporan schedule")
		return
	}

	response.Success(c, http.StatusOK, "Success", report)
}

// GetReplacementReportExport exports replacement report as CSV
// GET /api/v1/reports/replacement/export
func (h *ReportHandler) GetReplacementReportExport(c *gin.Context) {
	filters := buildReportFilters(c)

	report, err := h.reportUseCase.GetReplacementReport(c.Request.Context(), filters)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data laporan replacement")
		return
	}

	headers := []string{
		"Date", "Unit Code", "Plate Number", "Company", "Position",
		"Action", "Old Tyre SN", "New Tyre SN", "Brand", "Size",
		"Pattern", "Old RTD", "New RTD", "HM", "Driver Name",
		"Operator Name", "Remarks",
	}

	h.writeCSV(c, "replacement_report.csv", headers, func(csvWriter *csv.Writer) {
		for _, row := range report {
			record := []string{
				row.Date,
				row.UnitCode,
				row.PlateNumber,
				row.CompanyName,
				row.Position,
				row.Action,
				row.OldTyreSN,
				row.NewTyreSN,
				row.Brand,
				row.Size,
				row.Pattern,
				fmt.Sprintf("%.2f", row.OldRTD),
				fmt.Sprintf("%.2f", row.NewRTD),
				fmt.Sprintf("%.2f", row.HM),
				row.DriverName,
				row.OperatorName,
				row.Remarks,
			}
			csvWriter.Write(record)
		}
	})
}

// GetInventoryReportExport exports inventory report as CSV
// GET /api/v1/reports/inventory/export
func (h *ReportHandler) GetInventoryReportExport(c *gin.Context) {
	filters := buildReportFilters(c)

	report, err := h.reportUseCase.GetInventoryReport(c.Request.Context(), filters)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data laporan inventory")
		return
	}

	headers := []string{
		"Barcode", "Serial Number", "Brand", "Size", "Pattern",
		"Status", "OTD", "RTD", "Percent Worn (%)", "Unit Code",
		"Plate Number", "Company", "Last Position", "Last Mount Date",
		"Total Depth", "Current Depth",
	}

	h.writeCSV(c, "inventory_report.csv", headers, func(csvWriter *csv.Writer) {
		for _, row := range report {
			record := []string{
				row.Barcode,
				row.SerialNumber,
				row.Brand,
				row.Size,
				row.Pattern,
				row.Status,
				fmt.Sprintf("%.2f", row.OTD),
				fmt.Sprintf("%.2f", row.RTD),
				fmt.Sprintf("%.2f", row.PercentWorn),
				row.UnitCode,
				row.PlateNumber,
				row.CompanyName,
				row.LastPosition,
				row.LastMountDate,
				fmt.Sprintf("%.2f", row.TotalDepth),
				fmt.Sprintf("%.2f", row.CurrentDepth),
			}
			csvWriter.Write(record)
		}
	})
}

// GetScheduleReportExport exports schedule report as CSV
// GET /api/v1/reports/schedule/export
func (h *ReportHandler) GetScheduleReportExport(c *gin.Context) {
	filters := buildReportFilters(c)

	report, err := h.reportUseCase.GetScheduleReport(c.Request.Context(), filters)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data laporan schedule")
		return
	}

	headers := []string{
		"Serial Number", "Barcode", "Brand", "Size", "Pattern",
		"Unit Code", "Plate Number", "Position", "Current RTD",
		"Estimated End RTD", "Days Remaining", "Recommended Action",
	}

	h.writeCSV(c, "schedule_report.csv", headers, func(csvWriter *csv.Writer) {
		for _, row := range report {
			record := []string{
				row.SerialNumber,
				row.Barcode,
				row.Brand,
				row.Size,
				row.Pattern,
				row.UnitCode,
				row.PlateNumber,
				row.Position,
				fmt.Sprintf("%.2f", row.CurrentRTD),
				fmt.Sprintf("%.2f", row.EstimatedEnd),
				fmt.Sprintf("%.1f", row.DaysRemaining),
				row.RecommendedAction,
			}
			csvWriter.Write(record)
		}
	})
}

// writeCSV writes a CSV response with the given filename, headers, and data.
func (h *ReportHandler) writeCSV(c *gin.Context, filename string, headers []string, writeRows func(*csv.Writer)) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write headers
	if err := writer.Write(headers); err != nil {
		return
	}

	// Write data rows
	writeRows(writer)
}

// buildReportFilters builds filter map from query parameters.
func buildReportFilters(c *gin.Context) map[string]interface{} {
	filters := make(map[string]interface{})

	if companyID, _ := strconv.ParseUint(c.Query("company_id"), 10, 32); companyID > 0 {
		filters["company_id"] = uint(companyID)
	}
	if projectID, _ := strconv.ParseUint(c.Query("project_id"), 10, 32); projectID > 0 {
		filters["project_id"] = uint(projectID)
	}
	if unitID, _ := strconv.ParseUint(c.Query("unit_id"), 10, 32); unitID > 0 {
		filters["unit_id"] = uint(unitID)
	}
	if brandID, _ := strconv.ParseUint(c.Query("brand_id"), 10, 32); brandID > 0 {
		filters["brand_id"] = uint(brandID)
	}
	if sizeID, _ := strconv.ParseUint(c.Query("size_id"), 10, 32); sizeID > 0 {
		filters["size_id"] = uint(sizeID)
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		filters["date_from"] = dateFrom
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		filters["date_to"] = dateTo
	}
	if action := c.Query("action"); action != "" {
		filters["action"] = action
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	return filters
}

// RegisterRoutes registers report routes on the given router group.
func (h *ReportHandler) RegisterRoutes(rg *gin.RouterGroup) {
	reports := rg.Group("/reports")
	{
		reports.GET("/replacement", h.GetReplacementReport)
		reports.GET("/replacement/export", h.GetReplacementReportExport)
		reports.GET("/inventory", h.GetInventoryReport)
		reports.GET("/inventory/export", h.GetInventoryReportExport)
		reports.GET("/schedule", h.GetScheduleReport)
		reports.GET("/schedule/export", h.GetScheduleReportExport)
	}
}
