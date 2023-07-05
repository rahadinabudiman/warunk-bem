package repository

import (
	"context"
	"fmt"
	"time"
	"warunk-bem/domain"
	"warunk-bem/dtos"
	"warunk-bem/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	DB         mongo.Database
	Collection mongo.Collection
}

const (
	timeFormat     = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
	collectionName = "user"
)

func NewUserRepository(DB mongo.Database) domain.UserRepository {
	return &userRepository{DB, DB.Collection(collectionName)}
}

func (m *userRepository) InsertOne(ctx context.Context, req *domain.User) (*domain.User, error) {
	var (
		err error
	)

	_, err = m.Collection.InsertOne(ctx, req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (m *userRepository) FindOne(ctx context.Context, id string) (*domain.User, error) {
	var (
		user domain.User
		err  error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &user, err
	}

	err = m.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(&user)
	if err != nil {
		return &user, err
	}

	return &user, nil
}

func (m *userRepository) GetAllWithPage(ctx context.Context, rp int64, p int64, filter interface{}, setsort interface{}) ([]domain.User, int64, error) {

	var (
		user []domain.User
		skip int64
		opts *options.FindOptions
	)

	skip = (p * rp) - rp
	if setsort != nil {
		opts = options.MergeFindOptions(
			options.Find().SetLimit(rp),
			options.Find().SetSkip(skip),
			options.Find().SetSort(setsort),
		)
	} else {
		opts = options.MergeFindOptions(
			options.Find().SetLimit(rp),
			options.Find().SetSkip(skip),
		)
	}

	cursor, err := m.Collection.Find(
		ctx,
		filter,
		opts,
	)

	if err != nil {
		return nil, 0, err
	}
	if cursor == nil {
		return nil, 0, fmt.Errorf("nil cursor value")
	}
	err = cursor.All(ctx, &user)
	if err != nil {
		return nil, 0, err
	}

	count, err := m.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return user, 0, err
	}

	return user, count, err
}

func (m *userRepository) UpdateOne(ctx context.Context, user *domain.User, id string) (*domain.User, error) {
	var (
		err error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	filter := bson.M{"_id": idHex}
	update := bson.M{"$set": bson.M{
		"name":            user.Name,
		"email":           user.Email,
		"username":        user.Username,
		"password":        user.Password,
		"updated_at":      time.Now(),
		"role":            user.Role,
		"verified":        user.Verified,
		"loginverif":      user.LoginVerif,
		"verification":    user.VerificationCode,
		"activation_code": user.ActivationCode,
	}}

	_, err = m.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return user, err
	}

	err = m.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (m *userRepository) GetByCredential(ctx context.Context, req *dtos.LoginUserRequest) (*domain.User, error) {
	var (
		user domain.User
		err  error
	)

	credential := bson.M{
		"email":    req.Email,
		"password": req.Password,
	}

	err = m.Collection.FindOne(ctx, credential).Decode(&user)
	if err != nil {
		return &user, err
	}

	return &user, nil
}
func (m *userRepository) DeleteOne(ctx context.Context, id string) error {
	var (
		err error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = m.Collection.DeleteOne(ctx, bson.M{"_id": idHex})
	if err != nil {
		return err
	}

	return nil
}

func (m *userRepository) FindUsername(ctx context.Context, username string) (*domain.User, error) {
	var (
		user domain.User
		err  error
	)

	err = m.Collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return &user, err
	}

	return &user, nil
}

func (m *userRepository) FindEmail(ctx context.Context, email string) (*domain.User, error) {
	var (
		user domain.User
		err  error
	)

	err = m.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return &user, err
	}

	return &user, nil
}

func (m *userRepository) FindVerificationCode(ctx context.Context, verification int) (*domain.User, error) {
	var (
		user domain.User
		err  error
	)

	err = m.Collection.FindOne(ctx, bson.M{"verification": verification}).Decode(&user)
	if err != nil {
		return &user, err
	}

	return &user, nil
}

func (m *userRepository) FindActivationCode(ctx context.Context, activation int) (*domain.User, error) {
	var (
		user domain.User
		err  error
	)
	err = m.Collection.FindOne(ctx, bson.M{"activation_code": activation}).Decode(&user)
	if err != nil {
		return &user, err
	}

	return &user, nil
}
