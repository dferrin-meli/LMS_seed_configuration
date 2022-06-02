package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mercadolibre/seed_configuration-service/model"
	"github.com/mercadolibre/seed_configuration-service/repository"
)

var client = http.Client{Timeout: 5 * time.Second}

const (
	token_fury          = "cb70d19679786e7c308dbe3eed1901fe0aac06308ec46b271aa237abcff602ab"
	scope_calculator    = "https://jdespo-dev_lms-calculator.furyapps.io/labour/maxeffectivetime/get/process/%s"
	scope_configuration = "https://production_lms-configuration-service.furyapps.io/labour/processes/max-effective-time"
)

func main() {
	// processes := []string{"batchsorter", "checkin", "cycle_count", "hu_assembly", "inbound", "inbound_audit", "inbound_audit_movable",
	// 	"inbound_pick", "packing", "packing_withdrawals", "picking", "picking_transfer", "picking_withdrawals",
	// 	"putaway", "putwallin", "receiving", "sales_dispatch", "stock_audit"}

	repositoryConfiguration := repository.NewCalculatorRepository()

	metsCalculator, err := repositoryConfiguration.GetAll("stock_audit")
	if err != nil {
		fmt.Printf("error get quantity repository. %s", err.Error())
	}

	//consulta a la bd antigua
	// metCalculator, err := getMaxEffectiveTimeCalculator("batchsorter")
	// if err != nil {
	// 	return
	// }

	var metsToConfiguration []model.MaxEffectiveTimeDTO
	for _, antiguo := range metsCalculator {
		if len(metsToConfiguration) == 0 {
			metToConfiguration := model.MaxEffectiveTimeDTO{
				ProcessName:           antiguo.OperationalProcess,
				FacilityID:            antiguo.Warehouse,
				ProcessAttributeValue: antiguo.Attribute,
				Seconds:               antiguo.Seconds,
				ValidFrom:             antiguo.DateFrom,
			}

			metsToConfiguration = append(metsToConfiguration, metToConfiguration)
			continue
		}
		for i, nuevo := range metsToConfiguration {
			if antiguo.OperationalProcess == nuevo.ProcessName && antiguo.Warehouse == nuevo.FacilityID && antiguo.Attribute == nuevo.ProcessAttributeValue && nuevo.ValidTo == "" {
				if antiguo.Seconds != nuevo.Seconds {
					metToConfiguration := model.MaxEffectiveTimeDTO{
						ProcessName:           antiguo.OperationalProcess,
						FacilityID:            antiguo.Warehouse,
						ProcessAttributeValue: antiguo.Attribute,
						Seconds:               antiguo.Seconds,
						ValidFrom:             antiguo.DateFrom,
					}
					dateFromAntiguo, _ := formatStringToDate(antiguo.DateFrom)
					metsToConfiguration[i].ValidTo = formatDateToString(dateFromAntiguo.AddDate(0, 0, -1))
					metsToConfiguration = append(metsToConfiguration, metToConfiguration)
					break
				}
			} else if i+1 == len(metsToConfiguration) {
				metToConfiguration := model.MaxEffectiveTimeDTO{
					ProcessName:           antiguo.OperationalProcess,
					FacilityID:            antiguo.Warehouse,
					ProcessAttributeValue: antiguo.Attribute,
					Seconds:               antiguo.Seconds,
					ValidFrom:             antiguo.DateFrom,
				}
				metsToConfiguration = append(metsToConfiguration, metToConfiguration)
			}
		}
	}

	for _, met := range metsToConfiguration {
		_, err := postMaxEffectiveTimeConfiguration(met)
		if err != nil {
			fmt.Printf("\nerror post met to configuration. %s", err.Error())
			fmt.Println(met)
		}
	}
}

func getMaxEffectiveTimeCalculator(process string) ([]model.MaxEffectiveTimeCalculator, error) {
	result := []model.MaxEffectiveTimeCalculator{}
	req, _ := http.NewRequest("GET", fmt.Sprintf(scope_calculator, process), nil)
	req.Header.Set("x-auth-token", token_fury)
	res, err := client.Do(req)

	if err != nil {
		fmt.Printf("\n Error get met from calculator. Error: %s", err.Error())
		return result, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Printf("\nError status code. Status: %s", res.Status)
		return result, errors.New("error. response is not OK")
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(bodyBytes, &result)

	return result, nil
}

func postMaxEffectiveTimeConfiguration(met model.MaxEffectiveTimeDTO) (model.MaxEffectiveTimeDTO, error) {
	result := model.MaxEffectiveTimeDTO{}
	json_data, _ := json.Marshal(met)
	req, _ := http.NewRequest("POST", scope_configuration, bytes.NewBuffer(json_data))
	req.Header.Set("x-auth-token", token_fury)
	res, err := client.Do(req)

	if err != nil {
		fmt.Printf("\n Error post met to configuration. Error: %s", err.Error())
		return result, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Printf("\nError post met. status code. Status: %s", res.Status)
		return result, errors.New("error. response is not OK")
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(bodyBytes, &result)

	return result, nil
}

func formatDateToString(date time.Time) string {
	isoDate, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))
	return isoDate.UTC().Format(time.RFC3339)
}

func formatStringToDate(date string) (time.Time, error) {
	return time.Parse(time.RFC3339, date)
}
