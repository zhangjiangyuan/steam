package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

func main() {
	f := excelize.NewFile()
	// Create a new sheet.
	index := f.NewSheet("Sheet1")

	t1 := time.Tick(1 * time.Hour)
	i := 1
	for {

		select {

		case <-t1:
			t := time.Now().Format("2006-01-02 15:04:05")
			p := getPrice()
			fmt.Println(t + "\t" + p)
			err := f.SetCellValue("Sheet1", fmt.Sprintf("A%d",i), t)
			if err != nil {
				fmt.Println(err)
			}
			err = f.SetCellValue("Sheet1", fmt.Sprintf("B%d",i), p)
			if err != nil {
				fmt.Println(err)
			}
			f.SetActiveSheet(index)
			if err := f.SaveAs("Book1.xlsx"); err != nil {
				fmt.Println(err)
			}
			i++
		}

	}

}

func getPrice() string {
	reps, _ := http.Get("https://item.taobao.com/item.htm?spm=a230r.1.14.14.4dd53b520ODYJx&id=583610156835&ns=1&abbucket=11#detail")
	var body []byte
	body, err := ioutil.ReadAll(reps.Body)
	if err != nil {
		log.Fatalln(err)
	}
	reg, err := regexp.Compile("skuMap.*?\n")

	if err != nil {
		log.Fatalln(err)
	}
	ress := reg.FindAll(body, 1)
	i := -1
	for k, v := range ress[0] {
		if v == 123 {
			i = k
			break
		}
	}
	res := ress[0][i:]
	g := gjson.Parse(string(res))
	return g.Get(";30182:3319558;122216883:27447;.price").String()
}
