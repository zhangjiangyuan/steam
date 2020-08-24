package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/go-gomail/gomail"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

var (
	low   float64 = 486
	tHour *time.Ticker
)

func main() {
	f := excelize.NewFile()
	// Create a new sheet.
	runtime := time.Now().Format("20060102_150405")

	err := os.MkdirAll("static", 0777)
	if err != nil {
		log.Fatalln(err)
	}

	f.Path = fmt.Sprintf("static%s%s.xlsx", string(os.PathSeparator), runtime)
	index := f.NewSheet("Sheet1")

	tSec := time.NewTicker(1 * time.Second)

	i := 1

	data := make(map[string]string)

	for {
		sec := <-tSec.C
		if sec.Minute() == 0 && sec.Second() == 0 {
			tSec.Stop()
			tHour = time.NewTicker(1 * time.Hour)
			break
		}
	}

	for {
		select {

		case <-tHour.C:
			t := time.Now().Format("2006-01-02 15:04:05")
			p := getPrice()
			data["time"] = t
			data["price"] = p
			fmt.Println(t + "\t" + p)

			price, err := strconv.ParseFloat(p, 64)

			if err != nil {
				fmt.Println(err)
			} else if price <= low {
				sendToMe(p)
			}
			saveXlsx(f, index, i, data)
			i++
		}

	}

}

func saveXlsx(f *excelize.File, index, i int, data map[string]string) {
	err := f.SetCellValue("Sheet1", fmt.Sprintf("A%d", i), data["time"])
	if err != nil {
		fmt.Println(err)
	}
	err = f.SetCellValue("Sheet1", fmt.Sprintf("B%d", i), data["price"])
	if err != nil {
		fmt.Println(err)
	}
	f.SetActiveSheet(index)
	if err := f.Save(); err != nil {
		fmt.Println(err)
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

func sendToMe(data string) {
	//定义收件人
	mailTo := []string{
		"978746873@qq.com",
	}
	//邮件主题为"Hello"
	subject := "淘宝steam点卡价格"
	// 邮件正文
	body := data

	err := sendMail(mailTo, subject, body)
	if err != nil {
		log.Println(err)
		fmt.Println("send fail")
	} else {
		fmt.Println("send successfully")

	}
}

func sendMail(mailTo []string, subject string, body string) error {
	//定义邮箱服务器连接信息，如果是阿里邮箱 pass填密码，qq邮箱填授权码
	mailConn := map[string]string{
		"user": "qq978746873@163.com",
		"pass": ".Gzsgcs0.",
		"host": "smtp.163.com",
		"port": "465",
	}
	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()
	m.SetHeader("From", "Zhangjiayuan"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“XD Game”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	m.SetHeader("To", mailTo...)                                 //发送给多个用户
	m.SetHeader("Subject", subject)                              //设置邮件主题
	m.SetBody("text/html", body)                                 //设置邮件正文

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	return err
}
