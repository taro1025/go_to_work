package main

import (
	"fmt"
    "log"
	"os"
	"time"
    "github.com/joho/godotenv"
	"github.com/sclevine/agouti"
)

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}
}

func main() {
	loadEnv()

	//ChromeDriverを使用するための記述
	driver := agouti.ChromeDriver()

	//WebDriverプロセスを開始する
	err := driver.Start()
	if err != nil {
		log.Fatal(err)
	}

	//WebDriverプロセスを停止する（main関数の最後で停止したいのでdeferで処理)
	defer driver.Stop()
	if err != nil {
		log.Fatal(err)
	}

	//NewPage()でDriverに対応したページを返す。（今回はChrome）
	page, err := driver.NewPage()

	page.Navigate("https://attendance.moneyforward.com/my_page")
	url, err := page.URL()
	if err != nil {
		log.Fatal(err)
	}
	if url != "https://attendance.moneyforward.com/my_page" { // ログインしていない場合はログインする
		page.FindByID("employee_session_form_office_account_name").Fill(os.Getenv("MONEY_COMPANY_ID"))
		page.FindByID("employee_session_form_account_name_or_email").Fill(os.Getenv("MONEY_EMAIL"))
		page.FindByID("employee_session_form_password").Fill(os.Getenv("MONEY_PASSWORD"))
		err = page.FindByName("commit").Click()
		if err != nil {
			log.Fatal(err)
		}
	}
	t := time.Now()
	beginningOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	const layout = "2006-01-02"
	bulkAttendancesUrl := "https://attendance.moneyforward.com/my_page/bulk_attendances/" + beginningOfMonth.Format(layout) + "/edit"

	err = page.Navigate(bulkAttendancesUrl)
	rows := page.All("table.attendance-table-contents > tbody > tr")
	rowsLength, err := rows.Count()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < rowsLength; i++ {
		attend := rows.At(i).All("td")

		// 平日のみ入力
		dayKind, _ := attend.At(1).Find("div").Text()
		if dayKind == "平日" {
			attend.At(3).Find("div > input").Fill("9:00")
			attend.At(4).Find("div > input").Fill("18:00")
		}
	}
	err = page.FindByName("commit").Click()
	if err != nil {
		log.Fatal(err)
	}
	page.Screenshot("complete.png")
	if err != nil {
		log.Fatal(err)
	}
}
