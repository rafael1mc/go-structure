package shared

import (
	"gomodel/internal/auth"
	"gomodel/internal/institution"
	"gomodel/internal/member"
	"gomodel/internal/ping"
	"gomodel/internal/shared/database"
	"gomodel/internal/shared/env"
	"gomodel/internal/shared/jwt"
	"gomodel/internal/shared/logger"
	"gomodel/internal/shared/redis"
	"gomodel/internal/shared/timeprovider"
	"gomodel/internal/user"

	"go.uber.org/dig"
)

func BuildSharedContainer() *dig.Container {
	container := dig.New()

	// Shared
	container.Provide(timeprovider.NewTimeProviderImpl)
	container.Provide(env.NewEnv)
	container.Provide(database.NewDatabase)
	container.Provide(redis.NewRedisWrapper)
	container.Provide(logger.NewLogHandler)
	container.Provide(logger.NewLogger)
	container.Provide(jwt.NewJwtWrapper)

	// Internal
	container.Provide(ping.NewPing)
	container.Provide(user.NewUser)
	container.Provide(auth.NewAuth)
	container.Provide(institution.NewInstitution)
	container.Provide(member.NewMember)

	return container
}
