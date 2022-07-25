package auth

import (
	"context"
	"time"
	"url-shortner/config"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignupUser struct {
	Username string
	Password string
	Name     string
}
type AuthService interface {
	NormalLogin(ctx context.Context, email, password string) (jwt.MapClaims, error)
	GenerateToken(ctx context.Context, claims jwt.MapClaims) (string, error)
	RegisterUser(ctx context.Context, user RegisterInput) (jwt.MapClaims, error)
}

type authService struct {
	c  *mongo.Client
	db *mongo.Database
}

func (s *authService) NormalLogin(ctx context.Context, email, password string) (jwt.MapClaims, error) {
	if user, err := getUserByEmail(ctx, s, email); err != nil {
		return nil, err
	} else {
		// TODO: Add password checking logic.
		return jwt.MapClaims{
			"user_id": user.Id,
			"email":   user.Email,
			"name":    user.Name,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		}, nil
	}
}

func (s *authService) GenerateToken(ctx context.Context, claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.Get("JWT_SECRET").(string)))
	return t, err
}
func (s *authService) RegisterUser(ctx context.Context, user RegisterInput) (jwt.MapClaims, error) {
	return nil, nil
}

func NewAuthService(c *mongo.Client) AuthService {
	return &authService{
		c:  c,
		db: c.Database(config.Get("DB_NAME").(string)),
	}
}

// ========================================================
// Private Functions
// ========================================================
func getUserByEmail(ctx context.Context, s *authService, email string) (User, error) {
	usersCollection := s.db.Collection("users")
	opts := options.FindOne()
	query := bson.M{
		"email": email,
	}
	projections := bson.M{
		"_id":      1,
		"email":    1,
		"password": 1,
		"name":     1,
		"pet_name": 1,
	}
	opts.SetProjection(projections)
	var user User
	if err := usersCollection.FindOne(ctx, query, opts).Decode(&user); err != nil {
		return User{}, err
	}
	return user, nil
}
