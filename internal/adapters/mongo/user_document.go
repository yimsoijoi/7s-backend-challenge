package mongo

import (
	"time"

	"github.com/yimsoijoi/7s-backend-challenge/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	CreatedAt time.Time          `bson:"created_at"`
}

func toDocument(u *domain.User) (*userDocument, error) {
	oid, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		oid = primitive.NewObjectID()
	}

	return &userDocument{
		ID:        oid,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
	}, nil
}

func toDomain(d *userDocument) *domain.User {
	return &domain.User{
		ID:        d.ID.Hex(),
		Name:      d.Name,
		Email:     d.Email,
		Password:  d.Password,
		CreatedAt: d.CreatedAt,
	}
}
