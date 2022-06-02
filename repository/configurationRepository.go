package repository

import (
	"strings"

	"github.com/mercadolibre/seed_configuration-service/configuration"
	"github.com/mercadolibre/seed_configuration-service/model"
)

type CalculatorRepository struct {
	ConfigurationDbClient IDBClient
}

const queryString = "SELECT operational_process,seconds,warehouse,attribute, date_from FROM `max_effective_time_historic` maxEffTime WHERE maxEffTime.`operational_process` = ? ORDER BY warehouse, `attribute`, date_from ASC"

func NewCalculatorRepository() *CalculatorRepository {
	return &CalculatorRepository{
		ConfigurationDbClient: &DBClient{connectionString: configuration.GetConnectionStringCalculator()},
	}
}

func (lr *CalculatorRepository) Get() (int, error) {
	db := lr.ConfigurationDbClient.openDB()
	defer lr.ConfigurationDbClient.closeDB(db)

	quantity := 0
	result := db.QueryRow("SELECT COUNT(1) FROM max_effective_time_historic")
	err := result.Scan(&quantity)

	return quantity, err
}

func (lr *CalculatorRepository) GetAll(process string) ([]model.MaxEffectiveTimeCalculator, error) {
	db := lr.ConfigurationDbClient.openDB()
	defer lr.ConfigurationDbClient.closeDB(db)

	response := []model.MaxEffectiveTimeCalculator{}
	rows, err := db.Query(queryString, strings.ToLower(process))
	if err != nil {
		return response, nil
	}

	defer rows.Close()

	for rows.Next() {
		var aux model.MaxEffectiveTimeCalculator
		err = rows.Scan(
			&aux.OperationalProcess,
			&aux.Seconds,
			&aux.Warehouse,
			&aux.Attribute,
			&aux.DateFrom)
		if err != nil {
			return nil, err
		}
		response = append(response, aux)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return response, err
}

func (lr *CalculatorRepository) Upsert() error {
	return nil
}
