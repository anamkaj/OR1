package models

type CompanyPage struct {
	Count int           `json:"count"`
	Data  []CompanyList `json:"data"`
}

type CompanyList struct {
	CompaniesMany       bool        `json:"companies_many"`
	CompaniesIds        []int       `json:"companies_ids"`
	CompanyID           int         `json:"company_id"`
	SourceID            int         `json:"source_id"`
	ClientID            int         `json:"client_id"`
	ClientTitle         string      `json:"client_title"`
	TotalByMonthVisits  int         `json:"total_by_month_visits"`
	URL                 string      `json:"url"`
	CompanyTitle        string      `json:"company_title"`
	ManagerAssignedAt   interface{} `json:"manager_assigned_at"`
	EditorAssignedAt    string      `json:"editor_assigned_at"`
	Manager             interface{} `json:"manager"`
	ManagerID           interface{} `json:"manager_id"`
	Editor              string      `json:"editor"`
	EditorID            int         `json:"editor_id"`
	CompanySource       string      `json:"company_source"`
	CompanyCriteria     string      `json:"company_criteria"`
	Region              string      `json:"region"`
	ProductCount        int         `json:"product_count"`
	Packet              string      `json:"packet"`
	PacketEndDate       string      `json:"packet_end_date"`
	LastOperationDate   string      `json:"last_operation_date"`
	LastOperationResult string      `json:"last_operation_result"`
	LastOperationGoal   string      `json:"last_operation_goal"`
	LastOperationStatus string      `json:"last_operation_status"`
	NextOperationDate   string      `json:"next_operation_date"`
	NextOperationTime   string      `json:"next_operation_time"`
	LotEndDate          string      `json:"lot_end_date"`
	BannerLotEndDate    interface{} `json:"banner_lot_end_date"`
}

type LotList struct {
	ID               int         `json:"id"`
	BillID           int         `json:"bill_id"`
	AccountNumber    string      `json:"account_number"`
	ShiftLotID       interface{} `json:"shift_lot_id"`
	BillNumber       string      `json:"bill_number"`
	StartDate        string      `json:"start_date"`
	EndDate          string      `json:"end_date"`
	ServiceTitle     string      `json:"service_title"`
	DateColor        string      `json:"date_color"`
	SummaryPrice     string      `json:"summary_price"`
	Debt             string      `json:"debt"`
	Paid             string      `json:"paid"`
	Payer            string      `json:"payer"`
	SelfPayer        string      `json:"self_payer"`
	Company          string      `json:"company"`
	CompanySourceID  int         `json:"company_source_id"`
	User             string      `json:"user"`
	LastPaymentDate  string      `json:"last_payment_date"`
	Cluster          string      `json:"cluster"`
	Amount           int         `json:"amount"`
	OneUnitCost      interface{} `json:"one_unit_cost"`
	AccountingCenter string      `json:"accounting_center"`
	CanUpdate        bool        `json:"can_update"`
	CanUpdateBill    bool        `json:"can_update_bill"`
	CanDelete        bool        `json:"can_delete"`
	ShowBatchForm    bool        `json:"show_batch_form"`
	BannerName       interface{} `json:"banner_name"`
	BannerID         interface{} `json:"banner_id"`
	CopyBannerShape  bool        `json:"copy_banner_shape"`
	CanPrint         bool        `json:"can_print"`
}

type BillsList struct {
	ID            int    `json:"id"`
	AccountNumber string `json:"account_number"`
	Lots          []struct {
		ID            int    `json:"id"`
		AccountNumber string `json:"account_number"`
	} `json:"lots"`
	CreatedAt string `json:"created_at"`
	SelfPayer struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	} `json:"self_payer"`
	Payer struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	} `json:"payer"`
	AccountingCenter struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	} `json:"accounting_center"`
	SummaryPrice struct {
		Formatted string `json:"formatted"`
		Number    string `json:"number"`
	} `json:"summary_price"`
	SummaryPriceCurrency struct {
		Formatted string `json:"formatted"`
		Number    string `json:"number"`
	} `json:"summary_price_currency"`
	LastPaymentDate string `json:"last_payment_date"`
	Paid            struct {
		Formatted string `json:"formatted"`
		Number    string `json:"number"`
	} `json:"paid"`
	PaidCurrency struct {
		Formatted string `json:"formatted"`
		Number    string `json:"number"`
	} `json:"paid_currency"`
	PaidColor         string `json:"paid_color"`
	CanRead           bool   `json:"can_read"`
	CanCreate         bool   `json:"can_create"`
	CanUpdate         bool   `json:"can_update"`
	CanPrint          bool   `json:"can_print"`
	CanSendMonetaMail bool   `json:"can_send_moneta_mail"`
	CanDeleteLots     bool   `json:"can_delete_lots"`
}
