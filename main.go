package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sclevine/agouti"
	"log"
	"os"
	"time"
)

func main() {
	loadEnv()

	driver := startChromeDriver()
	defer driver.Stop()

	page, err := driver.NewPage() //Driverに対応したページを返す。（今回はChrome）

	bulkAttendancesUrl := thisMonthBulkAttendancesUrl()

	page.Navigate(bulkAttendancesUrl)
	url, err := page.URL()
	if err != nil {
		log.Fatal(err)
	}

	if url != bulkAttendancesUrl { // ログインしていない場合はログインする
		login(page)
		page.Navigate(bulkAttendancesUrl)
	}

	rows := page.All("table.attendance-table-contents > tbody > tr")
	rowsLength, err := rows.Count()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < rowsLength; i++ {
		attend := rows.At(i).All("td")
		dayKind, _ := attend.At(1).Find("div").Text()
		if dayKind == "平日" {
			attend.At(3).Find("div > input").Fill("9:00")
			attend.At(4).Find("div > input").Fill("18:00")
		}
	}
	time.Sleep(1 * time.Second)
	saveButton := page.Find("body > div > div.attendance-contents-inner > div > div > div > div.attendance-main-contents-inner > div > form > div.fixed-header-container > div > div > div.position-right > input[type=\"submit\"]:nth-child(2)")
	saveButton.MouseToElement()
	err = saveButton.Click()
	if err != nil {
		log.Fatal(err)
	}
	// すぐ消すの嫌なので3秒止める
	time.Sleep(3 * time.Second)
}


func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}
}

func startChromeDriver() *agouti.WebDriver{
	//ChromeDriverを使用するための記述
	driver := agouti.ChromeDriver()

	err := driver.Start()
	if err != nil {
		log.Fatal(err)
	}

	return driver
}

func thisMonthBulkAttendancesUrl() string{
	t := time.Now()
	beginningOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	const layout = "2006-01-02"
	return "https://attendance.moneyforward.com/my_page/bulk_attendances/" + beginningOfMonth.Format(layout) + "/edit"
}

func login(page *agouti.Page) {
	page.FindByID("employee_session_form_office_account_name").Fill(os.Getenv("MONEY_COMPANY_ID"))
	page.FindByID("employee_session_form_account_name_or_email").Fill(os.Getenv("MONEY_EMAIL"))
	page.FindByID("employee_session_form_password").Fill(os.Getenv("MONEY_PASSWORD"))
	err := page.FindByName("commit").Click()
	if err != nil {
		log.Fatal(err)
	}
}
