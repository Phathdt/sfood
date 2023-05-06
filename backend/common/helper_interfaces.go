package common

import (
	"context"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type GINComponent interface {
	GetPort() int
	GetRouter() *gin.Engine
}

type GormComponent interface {
	GetDB() *gorm.DB
}

type GrpcServer interface {
	SetRegisterHdl(hdl func(*grpc.Server))
}

type ClerkComponent interface {
	GetClient() clerk.Client
}

type JWTProvider interface {
	IssueToken(ctx context.Context, id, sub string) (token string, expSecs int, err error)
	ParseToken(ctx context.Context, tokenString string) (claims *jwt.RegisteredClaims, err error)
}

type Config interface {
	GetGRPCPort() int
	GetGRPCServerAddress() string
}
