package model

import (
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	ID          string `json:"id"`
	CompanyName string `json:"companyName"`
	FoundedYear string `json:"foundedYear"`
	Location    string `json:"location"`
	Address     string `json:"address"`
}

func NewCompanyService(db *gorm.DB) *CompanyService {
	return &CompanyService{
		DB: db,
	}
}

type CompanyService struct {
	DB *gorm.DB
}

func (d *CompanyService) CreateCompany(companyInput NewCompany) (*Company, error) {
	company := &Company{
		CompanyName: companyInput.CompanyName,
		FoundedYear: companyInput.FoundedYear,
		Location:    companyInput.Location,
		Address:     companyInput.Address,
	}

	result := d.DB.Create(company)
	if result.Error != nil {
		return nil, result.Error
	}

	return company, nil
}
func (d *CompanyService) GetAllCompanies() ([]*Company, error) {
	companies := make([]*Company, 0)
	result := d.DB.Find(&companies)
	if result.Error != nil {
		return nil, result.Error
	}
	return companies, nil
}
