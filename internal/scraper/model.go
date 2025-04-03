package scraper

type SearchInterface interface {
	SearchID(id string) (int, error)
	ActSearch(id string, bill []string) ([]PivotTable, error)
	Bills(link string) (Bills, error)
}

type Search struct {
	Token     string
	Session   string
	SearchURL string
	BillsURL  string
}

type Acts struct {
	StatementNumber string
	Date            string
	Amount          string
	Company         string
	Payer           string
	SelfPayer       string
	Bill            string
	Lot             string
	LinkBill        string
	IdPayer         string
	SiteID          string
}

type Bills struct {
	Date   string
	Amount string
}

type PivotTable struct {
	Acts  Acts
	Bills Bills
}
