package mongo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yimsoijoi/7s-backend-challenge/internal/adapters/mongo"
	"github.com/yimsoijoi/7s-backend-challenge/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestUserRepository_Create(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := mongo.NewUserRepository(mt.DB)

		mt.AddMockResponses(
			mtest.CreateSuccessResponse(),
		)

		user := &domain.User{
			Name:  "John",
			Email: "john@test.com",
		}

		err := repo.Create(context.Background(), user)

		assert.NoError(t, err)
		assert.NotEmpty(t, user.ID) // ✔ THIS is correct
	})
}

func TestUserRepository_Count(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := mongo.NewUserRepository(mt.DB)

		namespace := mt.DB.Name() + "." + mongo.ColUser

		mt.AddMockResponses(
			mtest.CreateCursorResponse(
				1,
				namespace,
				mtest.FirstBatch,
				bson.D{{Key: "n", Value: int64(5)}},
			),
			mtest.CreateCursorResponse(
				0,
				namespace,
				mtest.NextBatch,
			),
		)

		count, err := repo.Count(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})
}

func TestUserRepository_Delete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := mongo.NewUserRepository(mt.DB)
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := repo.Delete(context.Background(), primitive.NewObjectID().Hex())
		assert.NoError(t, err)
	})
}

func TestUserRepository_Update(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := mongo.NewUserRepository(mt.DB)

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		user := &domain.User{
			ID:    primitive.NewObjectID().Hex(),
			Name:  "Updated",
			Email: "updated@test.com",
		}

		err := repo.Update(context.Background(), user)
		assert.NoError(t, err)
	})
}

func TestUserRepository_FindAll(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("multiple users", func(mt *mtest.T) {
		repo := mongo.NewUserRepository(mt.DB)
		namespace := mt.DB.Name() + "." + mongo.ColUser
		mt.AddMockResponses(mtest.CreateCursorResponse(
			1,
			namespace,
			mtest.FirstBatch,
			bson.D{{Key: "name", Value: "User1"}},
			bson.D{{Key: "name", Value: "User2"}},
		))

		users, err := repo.FindAll(context.Background())

		assert.NoError(t, err)
		assert.Len(t, users, 2)
	})
}

func TestUserRepository_FindByEmail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("found", func(mt *mtest.T) {
		repo := mongo.NewUserRepository(mt.DB)
		namespace := mt.DB.Name() + "." + mongo.ColUser

		mt.AddMockResponses(mtest.CreateCursorResponse(
			1,
			namespace,
			mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: primitive.NewObjectID()},
				{Key: "email", Value: "john@test.com"},
				{Key: "name", Value: "John"},
			},
		))

		user, err := repo.FindByEmail(context.Background(), "john@test.com")

		assert.NoError(t, err)
		assert.Equal(t, "john@test.com", user.Email)
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("found", func(mt *mtest.T) {
		repo := mongo.NewUserRepository(mt.DB)
		namespace := mt.DB.Name() + "." + mongo.ColUser
		oid := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateCursorResponse(
			1,
			namespace,
			mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: oid},
				{Key: "email", Value: "john@test.com"},
				{Key: "name", Value: "John"},
			},
		))

		user, err := repo.FindByID(context.Background(), oid.Hex())

		assert.NoError(t, err)
		assert.Equal(t, "john@test.com", user.Email)
	})
}
