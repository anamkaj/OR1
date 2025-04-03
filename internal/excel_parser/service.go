package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"ord_crm/internal/scraper"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ExcelRepo struct{}

type RowData struct {
	Key      int
	Value    string
	Acts     int64
	FullName string
	Count    int
	Amount   string
	Company  string
}

type RowNum struct {
	NumActYa int64
	Row      int
}

type IgnoreEntry struct {
	Name  string
	ID    string
	Error string
	Link  string
	Row   int
}

var sheet = "Sheet1"

var ignor = make(map[int]IgnoreEntry)

func NewExcelRepo() ExcelRepository {
	return &ExcelRepo{}
}

func (e *ExcelRepo) ExcelParse() ([]Excel, error) {
	path := "./../storage/input/a.xlsx"

	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	var client []Excel
	column := "R"
	d := "B"

	counter := make(map[string]int)
	for row := 3; row < 300; row++ {

		login, err := f.GetCellValue(sheet, fmt.Sprintf("%s%d", d, row))
		if err != nil {
			log.Println("Error reading row")
			return nil, err
		}

		// остановка в конце документа где начинаются пустые строки.
		if login == "" {
			break
		}

		number, err := f.GetCellValue(sheet, fmt.Sprintf("%s%d", column, row))
		if err != nil {
			log.Println("Error reading row")
			return nil, err
		}
		// пропуск итерации если ИНН пустой
		if number == "" {
			continue
		}

		// [10,20,772736109269] == ["772736109269"]
		// поиск повторов
		exists := false
		for _, r := range client {
			if r.ID == number {
				exists = true
				break
			}

		}
		// подсчет строк в документе по одному клиенту
		counter[login]++

		if !exists {
			client = append(client, Excel{
				Login: login,
				ID:    number,
			})

		}

	}

	for i := range client {
		if count, found := counter[client[i].Login]; found {
			client[i].Count = count
		}
	}

	// for _, r := range client {
	// 	fmt.Printf("Добавлено: %s (%s) | Занимает: %d строк\n", r.ID, r.Login, r.Count)
	// }

	return client, err
}

func (e *ExcelRepo) ExcelImport(data [][]scraper.PivotTable) error {

	path := "./../storage/input/a.xlsx"
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)

		}
	}()

	rowData, err := scanDocument(f)
	if err != nil {
		return err
	}

	item := newSearchClient(data, rowData)
	if len(item) == 0 {
		fmt.Println("Array is empty_: ", item)

	}

	checkCompany(f, rowData, item)

	// где то сдесь можно сделать пометку разолкации  -- как отметить сумму в докуменете или вывести в консоль.

	groupItem := make(map[string][]scraper.PivotTable)
	for _, v := range item {
		groupItem[v.Acts.IdPayer] = append(groupItem[v.Acts.IdPayer], v)
	}

	for idPayer := range groupItem {
		sort.Slice(groupItem[idPayer], func(i, j int) bool {
			dateI, err := time.Parse("02.01.2006", groupItem[idPayer][i].Bills.Date)
			if err != nil {
				log.Printf("Error parsing Bills.Date %s: %v", groupItem[idPayer][i].Bills.Date, err)
				return false
			}
			dateJ, err := time.Parse("02.01.2006", groupItem[idPayer][j].Bills.Date)
			if err != nil {
				log.Printf("Error parsing Bills.Date %s: %v", groupItem[idPayer][j].Bills.Date, err)
				return false
			}
			return dateI.Before(dateJ)
		})
	}

	for idPayer, item := range groupItem {

		rw, err := updateNumber(rowData, idPayer)
		if err != nil {
			fmt.Println("error update number act ya")
		}

		switch {
		case len(rw) == 0:
			ignor[len(ignor)] = IgnoreEntry{
				Name:  item[0].Acts.Company,
				ID:    idPayer,
				Error: "No rows defined in excel table for this IdPayer",
				Link:  "",
			}
			continue

		case len(rw) < len(item):

			fmt.Println(len(rw), len(item))
			// если актов несколоко , а строчек в докуменете меньше , дублируется последняя ячейка и вставляется акты
			if len(rw) > 0 {
				lastRow := rw[len(rw)-1].Row

				if err := f.DuplicateRowTo(sheet, lastRow, lastRow+1); err != nil {
					fmt.Println("error duplicate row from table:", err)
				}

				fmt.Printf("Add new line to:%s-%d\n", item[0].Acts.Company, lastRow+1)

				rowData, err = scanDocument(f)
				if err != nil {
					return err
				}

				rw = nil

				rw, err = updateNumber(rowData, idPayer)
				if err != nil {
					fmt.Println("error update number act YA")
				}

			}

			for i, element := range item {
				if len(item) == len(rw) {
					row := rw[i].Row
					err := insertData(f, element, row)
					if err != nil {
						return fmt.Errorf("error inserting data for payer %s: %w", element.Acts.Company, err)
					}
				}

			}

		case len(item) == 1:
			element := item[0]
			for i := 0; i < len(rw); i++ {
				row := rw[i].Row
				err := insertData(f, element, row)
				if err != nil {
					return fmt.Errorf("error inserting data for payer %s: %w", element.Acts.Company, err)
				}
			}

		case len(rw) == len(item):
			for i, element := range item {
				row := rw[i].Row
				err := insertData(f, element, row)
				if err != nil {
					return fmt.Errorf("error inserting data for payer %s: %w", element.Acts.Company, err)
				}
			}

		default:
			indexMap := make(map[int64]int)
			nextIndex := 0

			for i := 0; i < len(rw); i++ {

				element := item[i%len(item)]
				row := rw[i].Row

				if s, ok := indexMap[rw[i].NumActYa]; ok {
					err := insertData(f, item[s], row)
					if err != nil {
						return fmt.Errorf("error inserting data for payer %s: %w", element.Acts.Company, err)
					}
				} else {
					index := nextIndex % len(item)
					indexMap[rw[i].NumActYa] = index

					err := insertData(f, item[index], row)
					if err != nil {
						return fmt.Errorf("error inserting data for payer %s: %w", element.Acts.Company, err)
					}

					nextIndex++
				}

			}
		}

	}

	for _, el := range ignor {
		// fmt.Println("Данные клиенты не будут добавленны в файл. Возможные ошибки: нет актов, количество строк больше чем актов, требуется добавить ячеек.")

		fmt.Printf("Error client %s : %s\n", el.Name, el.Error)

	}

	if err := f.SaveAs(path); err != nil {
		fmt.Println("Error saving file:", err)
		return err
	}

	return nil
}

func updateNumber(rowData []RowData, idPayer string) ([]RowNum, error) {

	var rw []RowNum
	for _, t := range rowData {
		if t.Value == idPayer {
			rw = append(rw, RowNum{
				NumActYa: t.Acts,
				Row:      t.Key,
			})
		}
	}

	sort.Slice(rw, func(i, j int) bool {
		return rw[i].NumActYa < rw[j].NumActYa
	})
	return rw, nil
}

func scanDocument(f *excelize.File) ([]RowData, error) {
	var rowData []RowData

	for rw := 3; rw < 1000; rw++ {
		number, err := f.GetCellValue(sheet, fmt.Sprintf("R%d", rw))
		if err != nil {
			log.Println("Error reading number")
			return nil, err
		}

		act, err := f.GetCellValue(sheet, fmt.Sprintf("G%d", rw))
		if err != nil {
			log.Println("Error reading act")
			return nil, err
		}

		amount, err := f.GetCellValue(sheet, fmt.Sprintf("AK%d", rw))
		if err != nil {
			log.Println("Error reading amount")
			return nil, err
		}
		company, err := f.GetCellValue(sheet, fmt.Sprintf("K%d", rw))
		if err != nil {
			fmt.Println("error parse comp0any name")
		}

		if number != "" {
			pref := strings.TrimPrefix(act, "YB-")
			num, err := strconv.ParseInt(pref, 10, 64)
			if err != nil {
				log.Printf("Error parsing number %s: %v", pref, err)
				continue
			}

			rowData = append(rowData, RowData{
				Key:      rw,
				Value:    number,
				FullName: act,
				Acts:     num,
				Count:    0,
				Amount:   amount,
				Company:  company,
			})

		}

		for i := range rowData {
			count := 0
			for j := range rowData {
				if rowData[j].Value == rowData[i].Value {
					count++
				}
			}
			rowData[i].Count = count
		}

	}
	return rowData, nil
}

func newSearchClient(data [][]scraper.PivotTable, cells []RowData) []scraper.PivotTable {

	amountMap := make(map[string][]struct {
		Amount string
		Row    int
	})

	for _, cell := range cells {
		amountMap[cell.Value] = append(amountMap[cell.Value], struct {
			Amount string
			Row    int
		}{
			Amount: cell.Amount,
			Row:    cell.Key,
		})
	}

	items := make([]scraper.PivotTable, 0)

	for i, r := range data {
		for j, v := range r {
			if _, ok := amountMap[v.Acts.IdPayer]; ok {
				str := v.Bills.Amount
				result := strings.TrimSuffix(str, "руб.")
				data[i][j].Bills.Amount = result
				items = append(items, data[i][j])
			}

		}
	}

	return items

}

func checkCompany(f *excelize.File, rowData []RowData, items []scraper.PivotTable) {
	rowMap := make(map[string][]RowData)
	for _, s := range rowData {
		rowMap[s.Value] = append(rowMap[s.Value], s)
	}

	for _, item := range items {
		r := rowMap[item.Acts.IdPayer]
		for _, s := range r {
			company, err := f.GetCellValue(sheet, fmt.Sprintf("K%d", s.Key))
			if err != nil {
				fmt.Println("error parse comp0any name")
			}

			cleanCompany := cleanCompanyName(company)
			cleanSelfPayer := cleanCompanyName(item.Acts.SelfPayer)

			if !strings.Contains(cleanCompany, cleanSelfPayer) && !strings.Contains(cleanSelfPayer, cleanCompany) {
				ignor[s.Key] = IgnoreEntry{
					Name:  item.Acts.Company,
					ID:    item.Acts.IdPayer,
					Error: fmt.Sprintf("the company do not match pay: %d row", s.Key),
					Link:  "",
				}
			}

		}
	}

}

func cleanCompanyName(company string) string {
	company = strings.ToLower(company)
	company = strings.TrimSpace(company)
	company = strings.ReplaceAll(company, "\u00A0", " ") // Неразрывные пробелы
	company = strings.ReplaceAll(company, `"`, "")       // Кавычки
	company = strings.ReplaceAll(company, "–", "-")      // Тире
	company = strings.ReplaceAll(company, " ", "")       // Все пробелы
	company = strings.ReplaceAll(company, "-", "")       // Все дефисы
	company = strings.Replace(company, "общество с ограниченной ответственностью", "", 1)
	return company
}

func insertData(f *excelize.File, v scraper.PivotTable, row int) error {

	data := map[string]string{
		"Z":  v.Acts.Bill,
		"AA": "contract",
		"AB": "distribution",
		"AC": "org-distribution",
		"AD": billsTime(v.Bills.Date),
		"AE": v.Bills.Amount,
		"AF": v.Acts.StatementNumber,
		"AG": v.Acts.Amount,
		"AH": v.Acts.Date,
	}

	for col, value := range data {
		cell := fmt.Sprintf("%s%d", col, row)

		val, err := f.GetCellValue(sheet, cell)
		if err != nil {
			return err
		}

		if val == "" {
			fmt.Printf("Add new data to: %s-%d\n", v.Acts.Company, row)
			if err := f.SetCellValue(sheet, cell, value); err != nil {
				return err
			}
		}

	}

	return nil
}

func billsTime(d string) string {

	layoutInput := "02.01.2006"
	layoutOutput := "2006-01-02"

	parsedTime, err := time.Parse(layoutInput, d)
	if err != nil {
		fmt.Println(err)
		return "error parsing time"
	}

	return parsedTime.Format(layoutOutput)

}
