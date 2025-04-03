package scraper

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"
)

type SearchApi struct {
	Search Search
}

func NewSearchRepo(s Search) *SearchApi {
	return &SearchApi{
		Search: s,
	}
}

func (с *SearchApi) SearchID(id string) (int, error) {

	params := url.Values{}
	params.Add("order", "id_desc")
	params.Add("q[payer_major_requisite_value_contains]", id)

	url := fmt.Sprintf("%s?%s", с.Search.SearchURL, params.Encode())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("error request to CRM")
		return 0, err
	}

	req.Header.Add("Cookie", с.Search.Session)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error request")
		return 0, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Err reading body response")
		return 0, err

	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)

	}

	var company string

	doc.Find("tbody").Each(func(index int, tbody *goquery.Selection) {
		tbody.Find("tr").Each(func(index int, tr *goquery.Selection) {
			c, ok := tr.Find("td.company").Find("a").Attr("href")
			if !ok {
				fmt.Println("Error ? not find company link ", err)
			} else {
				company = c
			}

		})
	})

	parts := strings.Split(company, "/")
	str := parts[len(parts)-1]

	if str == "" {
		return 0, nil
	}

	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	return num, nil

}

func (c *SearchApi) ActSearch(id string, bills []string) ([]PivotTable, error) {

	params := url.Values{}
	params.Add("utf8", "✓")
	params.Add("commit", "Отобрать")
	params.Add("order", "id_desc")
	params.Add("q[payer_major_requisite_value_contains]", id)

	url := fmt.Sprintf("%s?%s", c.Search.SearchURL, params.Encode())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("error request to CRM")
		return nil, err
	}

	req.Header.Add("Cookie", c.Search.Session)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error request")
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Err reading body response")
		return nil, err

	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	var data []PivotTable
	act, err := parseAct(doc, id, bills)
	if err != nil {
		fmt.Println("Error parsing act info")
		return nil, err
	}

	for _, bill := range act {
		b, err := c.Bills(strings.TrimSpace(bill.LinkBill))
		if err != nil {
			fmt.Println("Error parsing act info")
			return nil, err
		}

		data = append(data, PivotTable{
			Acts:  bill,
			Bills: b,
		})

	}

	return data, nil
}

func (c *SearchApi) Bills(link string) (Bills, error) {
	url := fmt.Sprintf("%s%s", c.Search.BillsURL, link)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("error request to Bills")

	}

	req.Header.Add("Cookie", c.Search.Session)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error request")

	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Err reading body response")

	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	bills, err := parseBills(doc)
	if err != nil {
		fmt.Println("Error parsing act info")
	}

	return bills, nil
}

func parseBills(n *goquery.Document) (Bills, error) {

	var row Bills

	n.Find("tbody").Each(func(index int, tbody *goquery.Selection) {
		tbody.Find("tr").Each(func(index int, tr *goquery.Selection) {

			thText := tr.Find("th").Text()
			if strings.TrimSpace(thText) == "Сумма счета" {
				row.Amount = tr.Find("td").Text()

			}
			if strings.TrimSpace(thText) == "Дата создания" {
				row.Date = tr.Find("td").Text()

			}

		})
	})

	return row, nil

}

func parseAct(n *goquery.Document, id string, bills []string) ([]Acts, error) {
	var act []Acts

	var billsUpper []string
	for _, s := range bills {
		billsUpper = append(billsUpper, strings.ToUpper(s))
	}

	n.Find("tbody").Each(func(index int, tbody *goquery.Selection) {
		tbody.Find("tr").Each(func(index int, tr *goquery.Selection) {

			bill := strings.ToUpper(strings.TrimSpace(tr.Find("td.bill").Text()))
			// если есть номера счетов , собираем только те которые есть в массиве bills
			if slices.Contains(billsUpper, bill) {
				date := tr.Find("td.date").Text()
				result, err := timeParse(date)
				if err != nil {
					fmt.Println("Error parsing date for", bill, ":", err)
				}

				if result {
					row := Acts{
						StatementNumber: tr.Find("td.statement_number").Text(),
						Date:            tr.Find("td.date").Text(),
						Amount:          tr.Find("td.amount").Text(),
						Company:         tr.Find("td.company").Text(),
						Payer:           tr.Find("td.payer").Text(),
						SelfPayer:       tr.Find("td.self_payer").Text(),
						Bill:            bill,
						LinkBill:        tr.Find("td.bill").Find("a").AttrOr("href", ""),
						Lot:             tr.Find("td.lot").Text(),
						IdPayer:         id,
					}

					c, ok := tr.Find("td.company").Find("a").Attr("href")
					if !ok {
						fmt.Println("Error ? not find company link ")
					} else {
						parts := strings.Split(c, "/")
						row.SiteID = parts[len(parts)-1]

					}

					act = append(act, row)
				}

			}

		})
	})

	return act, nil

}

// что делать с теми актами которые позже выставленны на 2 месяца ??? Мы забираем только те которые выставленны за 1 месяц !!!

// дата не должна выходить за пределы месяца старого иначе он берет новаые акты
// АД4193	31.03.2025
// 	АД3076	28.02.2025    	01.02.2025 - 28.02.2025

func timeParse(date string) (bool, error) {
	layout := "02.01.2006" // dd.mm.yyyy
	parseTime, err := time.Parse(layout, date)
	if err != nil {
		fmt.Println("Ошибка разбора даты", err)
		return false, err
	}
	parseTimeUTC := parseTime.UTC()
	now := time.Now().UTC()
	startOfPrevMonth := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.UTC) // 01.02
	endOfPrevMonth := startOfPrevMonth.AddDate(0, -1, 0)

	if parseTimeUTC.Before(startOfPrevMonth) && parseTimeUTC.After(endOfPrevMonth) {
		return true, nil
	}

	return false, nil
}
