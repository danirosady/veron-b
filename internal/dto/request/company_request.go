package request

type CreateCompanyRequest struct {
	Name          string `json:"name" binding:"required,max=255"`
	Address       string `json:"address"`
	ContactPerson string `json:"contact_person" binding:"max=255"`
	Phone         string `json:"phone" binding:"max=50"`
	Email         string `json:"email" binding:"omitempty,email,max=255"`
}

type UpdateCompanyRequest struct {
	Name          string `json:"name" binding:"required,max=255"`
	Address       string `json:"address"`
	ContactPerson string `json:"contact_person" binding:"max=255"`
	Phone         string `json:"phone" binding:"max=50"`
	Email         string `json:"email" binding:"omitempty,email,max=255"`
	Status        string `json:"status" binding:"omitempty,oneof=active inactive"`
}
