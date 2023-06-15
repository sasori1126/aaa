package account

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/internal/actions"
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/storage"
	"encoding/json"
	"strconv"
	"time"
)

type UserRepoDb struct {
	client storage.Storage
}

func (d UserRepoDb) SearchUsers(limit, offset int, search string) ([]models.User, error) {
	return d.client.SearchUsers(limit, offset, search)
}

func (d UserRepoDb) UpdateUserByField(userId uint, data models.FindByField) error {
	return d.client.UpdateUserField(userId, data)
}

func (d UserRepoDb) GetUserAddressByID(addressID uint) (*models.UserAddress, error) {
	return d.client.GetUserAddressByField(models.FindByField{
		Field: "id",
		Value: addressID,
	})
}

func (d UserRepoDb) DelUserAddress(userAddress *models.UserAddress) error {
	return d.client.DelUserAddress(userAddress)
}

func (d UserRepoDb) SetUserDefaultAddress(data *models.UserAddress) error {
	return d.client.SetUserDefaultAddress(data)
}

func (d UserRepoDb) UpdateUserAddress(userAddress *models.UserAddress) error {
	return d.client.UpdateUserAddress(userAddress)
}

func (d UserRepoDb) GetUserAddressByField(fv models.FindByField) (*models.UserAddress, error) {
	return d.client.GetUserAddressByField(fv)
}

func (d UserRepoDb) GetUserAddresses(userId uint) ([]models.UserAddress, error) {
	return d.client.GetUserAddresses(userId)
}

func (d UserRepoDb) AddUserAddress(data models.UserAddress) (*models.UserAddress, error) {
	return d.client.AddUserAddress(data)
}

func (d UserRepoDb) GetUsers(limit, offset int) ([]models.User, error) {
	return d.client.GetUsers(limit, offset)
}

func (d UserRepoDb) CreateUser(u models.User) error {
	err := d.client.CreateUser(u)
	if err != nil {
		return err
	}
	return nil
}

func (d UserRepoDb) FindUserByEmail(email string) (*models.User, error) {
	return d.client.FindUserByField(models.FindByField{Field: "email", Value: email})
}

func (d UserRepoDb) GetUserById(id uint) (*models.User, error) {
	return d.client.FindUserByField(models.FindByField{Field: "id", Value: id})
}
func (d UserRepoDb) FindUserById(id uint) (*models.User, error) {
	//retrieve user cache
	hashedId, _ := models.EncodeHashId(id)
	f := func(v interface{}) (string, error) {
		model := v.(models.FindByField)
		u, err := d.client.FindUserByField(model)
		if err != nil {
			return "", err
		}
		data, err := json.Marshal(u)
		if err != nil {
			return "", err
		}

		return string(data), err
	}

	user, err := storage.Cache.Remember(hashedId, f, models.FindByField{Field: "id", Value: strconv.Itoa(int(id))})
	if err != nil {
		return nil, err
	}

	ms := &models.User{}
	err = json.Unmarshal([]byte(user), ms)
	if err != nil {
		return nil, err
	}

	return ms, nil
}

func (d UserRepoDb) UpdateUser(u *models.User, data map[string]interface{}) error {
	err := d.client.UpdateByField(u, data)
	if err != nil {
		return err
	}
	return nil
}

func (d UserRepoDb) CreateAuth(u *models.User, token *actions.TokenDetail) error {
	at := time.Unix(token.AtExpires, 0)
	rt := time.Unix(token.RtExpires, 0)
	now := time.Now()

	hashedId, _ := models.EncodeHashId(u.ID)

	rdb := storage.Cache // Redis Cache
	err := rdb.SetWithTime(token.AccessUuid, hashedId, at.Sub(now))
	if err != nil {
		configs.Logger.Error(err)
		return err
	}
	err = rdb.SetWithTime(token.RefreshUuid, hashedId, rt.Sub(now))
	if err != nil {
		configs.Logger.Error(err)
		return err
	}
	return nil
}

func NewUserRepoDb(db storage.Storage) UserRepo {
	return UserRepoDb{client: db}
}
