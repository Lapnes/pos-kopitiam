package dto

type TaxConfigRequest struct {
	BranchID      string  `json:"branch_id"`
	Name          string  `json:"name" binding:"required"`
	TaxType       string  `json:"tax_type" binding:"required,oneof=tax service"`
	Percentage    float64 `json:"percentage" binding:"required,min=0,max=1"`
	IsActive      bool    `json:"is_active"`
	EffectiveFrom string  `json:"effective_from" binding:"required"` // YYYY-MM-DD
}