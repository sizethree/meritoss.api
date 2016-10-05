package services

import "fmt"
import "github.com/sizethree/miritos.api/models"
import "github.com/sizethree/miritos.api/context"

type UserManager struct {
	*context.Database
}

func (manager *UserManager) IsDuplicate(target *models.User) (bool, error) {
	var count int
	var existing models.User

	if target == nil {
		return true, fmt.Errorf("INVALID_USER")
	}

	err := manager.Model(existing).Where("email = ?", target.Email).Count(&count).Error
	return count >= 1, err
}

func (manager *UserManager) FindOrCreate(target *models.User) error {
	return manager.FirstOrCreate(target, *target).Error
}
