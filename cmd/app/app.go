package app

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"ord_crm/internal/config"
	"ord_crm/internal/crm_parser"
	"ord_crm/internal/excel_parser"
	"ord_crm/internal/scraper"
)

type AppModule struct {
	DB *redis.Client
}

func NewApp(db *redis.Client) *AppModule {
	return &AppModule{
		DB: db,
	}
}

func (a *AppModule) App() error {
	token, err := config.GetToken()
	if err != nil {
		fmt.Println("Error reading env....")
		return err
	}

	e := excel.NewExcelRepo()
	list, err := e.ExcelParse()
	if err != nil {
		fmt.Println("Error parsing file excel...:", err)
		return err
	}

	searchRepo := scraper.NewSearchRepo(scraper.Search{
		Token:     token.Token,
		Session:   token.Session,
		SearchURL: token.SearchURL,
		BillsURL:  token.BillsURL,
	})

	companyRepo := crm.NewCompanyRepo(token.BaseURL, token.Session)

	var allData [][]scraper.PivotTable

	for _, c := range list {

		// Для того что бы забирать акты только по мегатрифику .

		fmt.Printf("ID: %s\n", c.ID)

		id, err := searchRepo.SearchID(c.ID)
		if err != nil {
			fmt.Println("Error search inn client", err)
			return err
		}

		if id == 0 {
			fmt.Printf("ID is empty: %s\n", c.ID)
			continue
		}

		billArray, err := companyRepo.GetClientLot(id)
		if err != nil {
			fmt.Println("error get array lot", err)
			return err
		}

		data, err := searchRepo.ActSearch(c.ID, billArray)
		if err != nil {
			fmt.Println("error get data from CRM", err)
			return err
		}

		allData = append(allData, data)

	}
	err = e.ExcelImport(allData)
	if err != nil {
		fmt.Println("error insert data in excel document")
	}

	return nil
}
