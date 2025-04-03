package crm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CompanyApi struct {
	baseURL string
	session string
}

func NewCompanyRepo(baseURL string, session string) CompanyRepository {
	return &CompanyApi{
		baseURL: baseURL, 
		session: session,
	}
}

func (r *CompanyApi) GetClientList(user_id string, page int) ([]CompanyList, error) {
	url := fmt.Sprintf("%s/clients?user_id=%s&page=%d", r.baseURL, user_id, page)

	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req.Header.Add("Cookie", r.session)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error resp: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reade body response: %w", err)
	}

	if string(body) == "[]" {
		fmt.Println("Empty array client, finish page")
		return nil, fmt.Errorf("Empty array client, finish page: %s", string(body))
	}

	switch string(body[0]) {
	case "{":
		var data CompanyPage
		if err := json.Unmarshal(body, &data); err != nil {
			return nil, fmt.Errorf("error serialization JSON: %w", err)
		}

		return data.Data, nil

	case "[":
		var data []CompanyList
		if err := json.Unmarshal(body, &data); err != nil {
			return nil, fmt.Errorf("error serialization JSON: %w", err)
		}
		return data, nil

	default:
		return nil, fmt.Errorf("unexpected response : %s ", err)

	}

}

func (r *CompanyApi) GetClientLot(company int) ([]string, error) {
	url := fmt.Sprintf("%s/additional_lots_info?order=end_date_desc&company_id=%d", r.baseURL, company)

	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req.Header.Add("Cookie", r.session)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error reade body response: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reade body response: %w", err)
	}

	var data []LotList
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("error serialization JSON: %w", err)
	}

	var filteredData []string
	for _, d := range data {
		if d.ServiceTitle == "МЕГАТРАФИК (продление)" || d.ServiceTitle == "МЕГАТРАФИК (первая продажа)" {
			filteredData = append(filteredData, d.BillNumber)
		}

	}

	return filteredData, nil
}

func (r *CompanyApi) GetClientBills(company int) ([]BillsList, error) {
	url := fmt.Sprintf("%s/additional_bills_info?client_id=%d", r.baseURL, company)

	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req.Header.Add("Cookie", r.session)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error resp: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reade body response: %w", err)
	}

	var data []BillsList
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("error serialization JSON: %w", err)
	}
	fmt.Println("Bills successfully")

	return data, nil
}
