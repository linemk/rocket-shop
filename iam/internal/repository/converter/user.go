package converter

import (
	internalModel "github.com/linemk/rocket-shop/iam/internal/model"
	repoModel "github.com/linemk/rocket-shop/iam/internal/repository/model"
)

// ToInternalUser конвертирует repository User в internal User
func ToInternalUser(user *repoModel.User) *internalModel.User {
	if user == nil {
		return nil
	}

	notificationMethods := make([]internalModel.NotificationMethod, len(user.NotificationMethods))
	for i, method := range user.NotificationMethods {
		notificationMethods[i] = internalModel.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		}
	}

	return &internalModel.User{
		UserUUID:            user.UserUUID,
		Login:               user.Login,
		PasswordHash:        user.PasswordHash,
		Email:               user.Email,
		NotificationMethods: notificationMethods,
	}
}

// ToRepoUser конвертирует internal User в repository User
func ToRepoUser(user *internalModel.User) *repoModel.User {
	if user == nil {
		return nil
	}

	notificationMethods := make([]repoModel.NotificationMethod, len(user.NotificationMethods))
	for i, method := range user.NotificationMethods {
		notificationMethods[i] = repoModel.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		}
	}

	return &repoModel.User{
		UserUUID:            user.UserUUID,
		Login:               user.Login,
		PasswordHash:        user.PasswordHash,
		Email:               user.Email,
		NotificationMethods: notificationMethods,
	}
}
