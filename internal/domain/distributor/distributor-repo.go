package distributor

import (
	"axis/ecommerce-backend/internal/models"
	"axis/ecommerce-backend/internal/storage"
)

type RepoDistributorDb struct {
	client storage.Storage
}

func (d RepoDistributorDb) UpdateDistributor(data *models.Distributor) error {
	return d.client.UpdateDistributor(data)
}

func (d RepoDistributorDb) DeleteDistributor(id uint) error {
	return d.client.DeleteDistributor(id)
}

func (d RepoDistributorDb) GetDistributor(id uint) (*models.Distributor, error) {
	distributor, err := d.client.FindDistributorByField(models.FindByField{Field: "id", Value: id})
	if err != nil {
		return nil, err
	}

	return distributor, nil
}

func (d RepoDistributorDb) GetDistributors(limit, offset int) ([]models.Distributor, error) {
	return d.client.GetDistributors(limit, offset)
}

func (d RepoDistributorDb) CreateDistributor(data models.Distributor) error {
	err := d.client.CreateDistributor(data)
	if err != nil {
		return err
	}
	return nil
}

func NewRepoDistributorDb(db storage.Storage) RepoDistributor {
	return RepoDistributorDb{client: db}
}
