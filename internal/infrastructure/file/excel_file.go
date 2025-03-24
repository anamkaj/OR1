package file

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"ord_crm/internal/domain/entity"
	"ord_crm/internal/domain/repository"
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
}

type RowNum struct {
	Num int64
	Row int
}

var sheet = "Sheet1"

func NewExcelRepo() repository.ExcelRepository {
	return &ExcelRepo{}
}

func (e *ExcelRepo) ExcelParse() ([]entity.Excel, error) {
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

	var client []entity.Excel

	column := "R"
	d := "B"

	counter := make(map[string]int)

	for row := 3; row < 1000; row++ {

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
			client = append(client, entity.Excel{
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

func (e *ExcelRepo) ExcelImport(data [][]entity.PivotTable) error {
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
		fmt.Println("Array is empty_: ", rowData)
	}

	// где то сдесь можно сделать пометку разолкации

	// проверка фирм куда поступил платеж медиаплощадь vs адверта

	groupItem := make(map[string][]entity.PivotTable)
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

	ignor := make(map[int]struct {
		name  string
		id    string
		error string
		link  string
	})

	for idPayer, item := range groupItem {
		var count int
		for _, record := range rowData {
			if record.Value == idPayer {
				count = record.Count
				break
			}
		}

		rw, err := updateNumber(rowData, idPayer)
		if err != nil {
			fmt.Println("error update number act ya")
		}

		switch {
		case len(rw) == 0:
			ignor[len(ignor)] = struct {
				name  string
				id    string
				error string
				link  string
			}{
				name:  item[0].Acts.Company,
				id:    idPayer,
				error: "No rows defined in ex for this IdPayer",
				link:  "",
			}
			continue

		case count < len(item):
			// если актов несколоко , а строчек в докуменете меньше , дублируется последняя ячейка и вставляется акты
			if len(rw) > 0 {
				lastRow := rw[len(rw)-1].Row

				if err := f.DuplicateRowTo(sheet, lastRow, lastRow+1); err != nil {
					fmt.Println("Ошибка при дублировании строки:", err)
				}

				fmt.Printf("Add new line to:%s-%d\n", item[0].Acts.Company, lastRow+1)

				rowData, err = scanDocument(f)
				if err != nil {
					return err
				}

				rw = nil

				rw, err = updateNumber(rowData, idPayer)
				if err != nil {
					fmt.Println("error update number act ya")
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
			for i := 0; i < count; i++ {
				row := rw[i].Row
				err := insertData(f, element, row)
				if err != nil {
					return fmt.Errorf("error inserting data for payer %s: %w", element.Acts.Company, err)
				}
			}

		case count == len(item):
			for i, element := range item {
				row := rw[i].Row
				err := insertData(f, element, row)
				if err != nil {
					return fmt.Errorf("error inserting data for payer %s: %w", element.Acts.Company, err)
				}
			}

		default:
			num := 0
			for i := 0; i < count; i++ {

				if i >= len(item) {
					break
				}

				element := item[i]
				row := rw[i].Row

				if num == int(rw[i].Num) {
					el := item[i-1]
					err := insertData(f, el, row)
					if err != nil {
						return fmt.Errorf("error inserting data for payer %s: %w", element.Acts.Company, err)
					}

				}

				if i+1 < len(rw) {
					nextRow := rw[i+1].Row
					err := insertData(f, element, nextRow)
					if err != nil {
						return fmt.Errorf("error inserting data for payer %s: %w", element.Acts.Company, err)
					}
				}

				num = int(rw[i].Num)

			}

		}

	}

	for _, el := range ignor {
		// fmt.Println("Данные клиенты не будут добавленны в файл. Возможные ошибки: нет актов, количество строк больше чем актов, требуется добавить ячеек.")

		fmt.Printf("Error adding client %s : %s\n", el.name, el.error)

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
				Num: t.Acts,
				Row: t.Key,
			})
		}
	}

	sort.Slice(rw, func(i, j int) bool {
		return rw[i].Num < rw[j].Num
	})
	return rw, nil
}

func scanDocument(f *excelize.File) ([]RowData, error) {
	// Когда добавляем ячейки нужен пересчет документов

	var rowData []RowData

	for rw := 3; rw < 1000; rw++ {
		number, err := f.GetCellValue(sheet, fmt.Sprintf("R%d", rw))
		if err != nil {
			log.Println("Error reading row")
			return nil, err
		}

		act, err := f.GetCellValue(sheet, fmt.Sprintf("G%d", rw))
		if err != nil {
			log.Println("Error reading row")
			return nil, err
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

func newSearchClient(data [][]entity.PivotTable, cells []RowData) []entity.PivotTable {

	valueMap := make(map[string]struct{})
	for _, cell := range cells {
		valueMap[cell.Value] = struct{}{}
	}

	var items []entity.PivotTable
	for i, r := range data {
		for j, v := range r {
			if _, exists := valueMap[v.Acts.IdPayer]; exists {
				str := v.Bills.Amount
				result := strings.TrimSuffix(str, "руб.")
				data[i][j].Bills.Amount = result
				items = append(items, data[i][j])
			}
		}
	}

	return items

}

func insertData(f *excelize.File, v entity.PivotTable, row int) error {

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
		fmt.Println("error parsing date bill")
		return ""
	}

	return parsedTime.Format(layoutOutput)

}
