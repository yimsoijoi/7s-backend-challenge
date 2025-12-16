package mongo

import (
	"context"

	"github.com/yimsoijoi/7s-backend-challenge/internal/domain"
	"github.com/yimsoijoi/7s-backend-challenge/internal/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ColUser = "users"
)

type UserRepository struct {
	col *mongo.Collection
}

func NewUserRepository(db *mongo.Database) ports.UserRepository {
	return &UserRepository{col: db.Collection(ColUser)}
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	doc, err := toDocument(u)
	if err != nil {
		return err
	}

	_, err = r.col.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	u.ID = doc.ID.Hex()

	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	var u domain.User
	err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&u)
	return &u, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.col.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	return &u, err
}

func (r *UserRepository) FindAll(
	ctx context.Context,
	// Pagination is not used for now.
	// The number of users is expected to be small.
	// We rely on an index on `createdAt` to keep this query fast.
	// Add pagination if the user list becomes large.
) ([]*domain.User, error) {

	filter := bson.M{}

	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetProjection(bson.M{
			"password": 0,
		})

	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var users []*domain.User
	for cur.Next(ctx) {
		var u domain.User
		if err := cur.Decode(&u); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, u *domain.User) error {
	oid, _ := primitive.ObjectIDFromHex(u.ID)
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"name": u.Name, "email": u.Email}})
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	return r.col.CountDocuments(ctx, bson.M{})
}
