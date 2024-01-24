package router

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"mini_project/controller"
	docs "mini_project/docs"
	"mini_project/entity"
	"mini_project/services"
	gormutils "mini_project/utils/gorm"
	redis_utils "mini_project/utils/redis"
)

var (
	v1                                       = "/api/v1"
	db             *gorm.DB                  = gormutils.InitMySQL()
	rdb            *redis.Client             = redis_utils.InitRedis()
	userEntity     entity.UserEntity         = entity.NewUserEntity(db)
	userService    services.UserService      = services.NewUserService(userEntity)
	redisEntity    entity.RedisEntity        = entity.NewRedisEntity(rdb)
	jwtService     services.JWTService       = services.NewJWTService(redisEntity, userEntity)
	userController controller.UserController = controller.NewUserController(userService, jwtService)
)

func SetupRouter() *gin.Engine {
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic("Failed to load env file")
	}
	// Closing the db when the program stop
	defer gormutils.Close(db)
	defer redis_utils.Close(rdb)
	r := gin.Default()

	authRoutes := r.Group(v1 + "/auth")
	{
		authRoutes.POST("/login", userController.Login)
		authRoutes.POST("/register", userController.Register)
	}
	docs.SwaggerInfo.BasePath = v1
	swagger := r.Group(v1 + "/swagger")
	{
		swagger.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	err := r.Run()
	if err != nil {
		return nil
	}
	return r
}
