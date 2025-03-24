package app

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"ord_crm/config"
	"ord_crm/internal/application"
	"ord_crm/internal/domain/entity"
	"ord_crm/internal/infrastructure/api"
	"ord_crm/internal/infrastructure/file"
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

	e := file.NewExcelRepo()
	list, err := e.ExcelParse()
	if err != nil {
		fmt.Println("Error parsing file excel...:", err)
		return err
	}

	searchRepo := api.NewSearchRepo(entity.Search{
		Token:     token.Token,
		Session:   token.Session,
		SearchURL: token.SearchURL,
		BillsURL:  token.BillsURL,
	})
	searchService := application.NewSearchService(searchRepo)

	companyRepo := api.NewCompanyRepo(token.BaseURL, token.Session)
	companyService := application.NewCompanyService(companyRepo)

	var allData [][]entity.PivotTable

	for _, c := range list {

		// Для того что бы забирать акты только по мегатрифику .

		id, err := searchService.SearchID(c.ID)
		if err != nil {
			fmt.Println("Error search inn client")
		}

		billArray, err := companyService.GetClientLot(id)
		if err != nil {
			fmt.Println("error get array lot")
			return err
		}

		data, err := searchService.ActSearch(c.ID, billArray)
		if err != nil {
			fmt.Println("error get data from CRM")
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
