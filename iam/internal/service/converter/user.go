package converter

import (
	commonv1 "github.com/linemk/rocket-shop/shared/pkg/proto/common/v1"

	"github.com/linemk/rocket-shop/iam/internal/model"
)

func UserToProto(user *model.User) *commonv1.User {
	if user == nil {
		return nil
	}

	notificationMethods := make([]*commonv1.NotificationMethod, len(user.NotificationMethods))
	for i, method := range user.NotificationMethods {
		notificationMethods[i] = &commonv1.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		}
	}

	return &commonv1.User{
		UserUuid:            user.UserUUID,
		Login:               user.Login,
		Email:               user.Email,
		NotificationMethods: notificationMethods,
	}
}

func UserFromProto(user *commonv1.User) *model.User {
	if user == nil {
		return nil
	}

	notificationMethods := make([]model.NotificationMethod, len(user.NotificationMethods))
	for i, method := range user.NotificationMethods {
		notificationMethods[i] = model.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		}
	}

	return &model.User{
		UserUUID:            user.UserUuid,
		Login:               user.Login,
		Email:               user.Email,
		NotificationMethods: notificationMethods,
	}
}
