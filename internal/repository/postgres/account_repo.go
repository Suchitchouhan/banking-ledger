package postgres

import (
	"fmt"
	"time"

	"banking-ledger/internal/domain"

	"banking-ledger/internal/repository/models"

	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func mapDomainToModel(account *domain.Account) *models.Account {
	return &models.Account{
		ID:        account.ID,
		Name:      account.Name,
		Balance:   account.Balance,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
}

func mapModelToDomain(model *models.Account) *domain.Account {
	return &domain.Account{
		ID:        model.ID,
		Name:      model.Name,
		Balance:   model.Balance,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

// Inserts a new account into the database
func (r *AccountRepository) Create(account *domain.Account) error {
	model := mapDomainToModel(account)
	result := r.db.Create(model)
	if result.Error != nil {
		return fmt.Errorf("failed to create account: %v", result.Error)
	}

	return nil
}

// retrieves an account by its ID
func (r *AccountRepository) GetByID(id string) (*domain.Account, error) {
	var model models.Account
	result := r.db.First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to retrieve account: %v", result.Error)
	}
	return mapModelToDomain(&model), nil
}

// Updates the balance of an account
func (r *AccountRepository) UpdateBalance(id string, newBalance float64) error {
	result := r.db.Model(&models.Account{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"balance":    newBalance,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update account balance: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("account not found")
	}

	return nil
}

// retrieves all accounts
func (r *AccountRepository) List() ([]*domain.Account, error) {
	var models []models.Account
	result := r.db.Order("created_at DESC").Find(&models)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list accounts: %v", result.Error)
	}

	accounts := make([]*domain.Account, len(models))
	for i, model := range models {
		modelCopy := model
		accounts[i] = mapModelToDomain(&modelCopy)
	}

	return accounts, nil
}
