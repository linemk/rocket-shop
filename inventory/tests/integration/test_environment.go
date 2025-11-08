//go:build integration

package integration

import (
	"github.com/linemk/rocket-shop/platform/pkg/testcontainers/app"
	"github.com/linemk/rocket-shop/platform/pkg/testcontainers/mongo"
	"github.com/linemk/rocket-shop/platform/pkg/testcontainers/network"
)

// TestEnvironment — структура для хранения ресурсов тестового окружения
type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.Container
	App     *app.Container
}
