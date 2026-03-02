package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthClient 认证服务客户端
type AuthClient struct {
	conn   *grpc.ClientConn
	client AuthServiceClient
}

// NewAuthClient 创建认证服务客户端
func NewAuthClient(addr string) (*AuthClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return &AuthClient{
		conn:   conn,
		client: NewAuthServiceClient(conn),
	}, nil
}

// Close 关闭连接
func (c *AuthClient) Close() error {
	return c.conn.Close()
}

// Register 用户注册
func (c *AuthClient) Register(ctx context.Context, username, email, password string) (*LoginResponse, error) {
	return c.client.Register(ctx, &RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	})
}

// Login 用户登录
func (c *AuthClient) Login(ctx context.Context, username, password string) (*LoginResponse, error) {
	return c.client.Login(ctx, &LoginRequest{
		Username: username,
		Password: password,
	})
}

// RefreshToken 刷新Token
func (c *AuthClient) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	return c.client.RefreshToken(ctx, &RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
}

// ValidateToken 验证Token
func (c *AuthClient) ValidateToken(ctx context.Context, accessToken string) (*ValidateTokenResponse, error) {
	return c.client.ValidateToken(ctx, &ValidateTokenRequest{
		AccessToken: accessToken,
	})
}

// GetCurrentUser 获取当前用户信息
func (c *AuthClient) GetCurrentUser(ctx context.Context, userID uint64) (*UserInfo, error) {
	return c.client.GetCurrentUser(ctx, &GetCurrentUserRequest{
		UserId: userID,
	})
}
