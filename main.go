package main

import (
	"github.com/kevinburke/twilio-go"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"gopkg.in/go-toast/toast.v1"
)

import "fmt"

const AccountSID = "ACd3f9deffef1a2df67e0246a3b9310bb7"
const AuthToken = "aee938906be7f3845443a96c915c7c85"
const ProxyURL = "https://raw.githubusercontent.com/clarketm/proxy-list/master/proxy-list-raw.txt"

var Proxies []string

type LinkForCPU struct {
	url  string
	name string
}

var links = []LinkForCPU{
	{
		url:  "https://www.bestbuy.com/site/amd-ryzen-9-5900x-4th-gen-12-core-24-threads-unlocked-desktop-processor-without-cooler/6438942.p?skuId=6438942",
		name: "5900x",
	},
	//{
	//	url:  "https://www.amazon.com/AMD-Ryzen-5900X-24-Thread-Processor/dp/B08164VTWH/",
	//	name: "5900x",
	//},
	{
		url:  "https://www.bhphotovideo.com/c/product/1598373-REG/amd_100_100000061wof_ryzen_9_5900x_3_7.html?SID=trd-us-9843011037592730000",
		name: "5900x",
	},
	{
		url:  "https://www.newegg.com/amd-ryzen-9-5900x/p/N82E16819113664",
		name: "5900x",
	},
	{
		url:  "https://www.newegg.com/Product/ComboDealDetails?ItemList=Combo.4207305&quicklink=true",
		name: "5900x",
	},
	{
		url:  "https://www.newegg.com/Product/ComboDealDetails?ItemList=Combo.4208036",
		name: "5900x",
	},
	{
		url:  "https://www.amd.com/en/direct-buy/5450881500/us?add-to-cart=true",
		name: "5900x",
	},
	{
		url:  "https://www.newegg.com/Product/ComboDealDetails?ItemList=Combo.4207319&quicklink=true",
		name: "5900x",
	},
	{
		url:  "https://www.bestbuy.com/site/amd-ryzen-9-5950x-4th-gen-16-core-32-threads-unlocked-desktop-processor-without-cooler/6438941.p?skuId=6438941",
		name: "5950x",
	},
	{
		url:  "https://www.bhphotovideo.com/c/product/1598372-REG/amd_100_100000059wof_ryzen_9_5950x_3_4.html?SID=trd-us-1012364251640402200",
		name: "5950x",
	},
	{
		url:  "https://www.tigerdirect.com/applications/SearchTools/item-details.asp?EdpNo=6945711&Sku=42259884&SRCCODE=3WCJ&utm_source=cj&utm_content=8808717&utm_term=12646018&cjevent=0454c97a343611eb81f90cd70a24060d",
		name: "5950x",
	},
	{
		url:  "https://www.newegg.com/amd-ryzen-9-5950x/p/N82E16819113663",
		name: "5950x",
	},
	//{
	//	url:  "https://www.amazon.com/AMD-Ryzen-5950X-32-Thread-Processor/dp/B0815Y8J9N/ref=cm_cr_arp_d_product_top?ie=UTF8",
	//	name: "5950x",
	//},
	{
		url:  "https://www.amd.com/en/direct-buy/5450881400/us?add-to-cart=true",
		name: "5950x",
	},
}

var outOfStockStrings = []string{
	"currently unavailable",
	"out of stock",
	"sold out",
	"currently sold out",
	"notify when available",
}

var numChecks = 0

func getProxies() {
	fmt.Printf("Getting proxies from %s...\n", ProxyURL)
	res, err := http.Get(ProxyURL)
	if err != nil {
		log.Fatalln(err)
	}

	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	proxiesStr := string(bs)
	Proxies = strings.Split(proxiesStr, "\n")

	fmt.Printf("%d proxies found.\n", len(Proxies))
}

func main() {
	getProxies()

	go func() {
		for {
			time.Sleep(time.Minute * 120)
			println("Refreshing proxies...")
			getProxies()
		}
	}()

	var wg sync.WaitGroup
	for _, link := range links {
		wg.Add(1)
		go func(cpu LinkForCPU) {
			for {
				_ = checkLink(cpu)
				numChecks += 1
				fmt.Printf("Checked %d times\n", numChecks)
				u, _ := url.Parse(cpu.url)

				rand.Seed(time.Now().UnixNano())
				n := rand.Intn(8000-4000+1) + 4000

				if strings.Contains(u.Host, "amazon") {
					time.Sleep(time.Duration(n*2) * time.Millisecond)
				} else {
					time.Sleep(time.Duration(n) * time.Millisecond)
				}

			}
			wg.Done()
		}(link)

		time.Sleep(620 * time.Millisecond) // Wait half a second before starting another request so we don't have a bunch at once
	}
	wg.Wait()
}

func checkLink(cpu LinkForCPU) bool {
	u, err := url.Parse(cpu.url)
	if err != nil {
		log.Fatalf("Error parsing URL: %s", err)
	}

	req, _ := http.NewRequest("GET", cpu.url, nil)
	//req.Header.Set("Host", u.Host)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,la;q=0.8")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("cookie", "ubid-main=132-2823154-0272167; lc-main=en_US; s_fid=46D497D5777EF32F-193FFBD5745F7772; i18n-prefs=USD; x-main=\"NAYYynlQpGj9qKt6ACnLnLV2aT@ZxxXAxZSyoD9c4YW@RP@mGjxoe8hH@X0fTufN\"; at-main=Atza|IwEBICMV-QXExhKbFGNGUrxYnGC4Wws2rVhJRxCAh-QrTjgNvCg0MpGVLBYN1VDQc7d5W_NWiJptwbF4qci2mIsbjNVYMyO3OatH2ndfVgFkXRG_zB7UhCRyxuiWiwTtK-ipkykbAohMlW5ld78nxvarvskEz-UBC4khWwxNjD-uVD31jQN3GueXFS6CX3Uo8uxaybLAwpXlOsb9pt6wPl1ftN-H; sess-at-main=\"yhahS9IJ+rO4s2tDZjquTD8aofEF4RtFqZkF0A+Ffdg=\"; sst-main=Sst1|PQHkzImiJukp5rQkXlAA40wzCc6fMttmdWa3Ozjw95X_kUhAP1XdS8SQHO_FBQupcguVG7bX-xg_qud2FG3oWJ2E6iCi8cl_cCSEb0WtPrDmCD7L3mNEcsf2v1u1H0s8I8PoHwrsb_NKmcFuo36bKfYfBYA2aAp0ywo4AOmySnsTxdMX6Zx3uPsGDseJB7zR0xQSoR5ShdiQp_XIEsip1NXJX0679rsxU0uj-nEEuwbeF8Eljd4SUXqss045tU8v_-W_-7Jxi6nBe0FTJMWTIgsjTOApERnrGaLSa2UORHRDLxM; s_vnum=2038505867425%26vn%3D1; session-id-apay=147-9208917-0592937; session-id=134-2551160-6521529; regStatus=pre-register; aws-target-visitor-id=1606863222601-264985; aws-target-data=%7B%22support%22%3A%221%22%7D; s_dslv=1606871195959; s_vn=1638399215415%26vn%3D2; s_nr=1606871195961-Repeat; session-token=IlPIJxqkhjCNV/tukkCwg/KWJcvr16hlESv6KQf+0iOnKhD7BTxXxmER0+DmUswCKh+BDJ8jUYKwzg4XDCggJZ1qrx19PhF/XSyZsD+mlWDNzrIm3J+IIYuHZmQrlhRJ/US2JueNkFiOX4FCy12dorNVzSgGBbUIn6HtnVjALcu5/hvqag4r2gVIyIzBTcGrgLThGgP3nzjrWjh6FWsqAPSjrxXnHXfy; session-id-time=2082787201l; csm-hit=tb:9GC9Y92H4DTP9JQRYA02+s-0P0VZ6YVZCQNB0BJ64X1|1606896309366&t:1606896309366&adb:adblk_yes")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		println("can't do req")
		log.Println(err)
	}

	defer res.Body.Close()
	if res.StatusCode >= 300 {
		println("status code >= 300; waiting 30 seconds to retry")
		log.Println(res.StatusCode)
		time.Sleep(30 * time.Second)
	}

	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		println("error reading byte stream")
		sendText(fmt.Sprintf("error reading byte stream: %s", err))
		log.Fatal(err)
	}

	bodystr := string(bs)

	return checkIfInStock(bodystr, cpu, u)
}

func checkIfInStock(body string, cpu LinkForCPU, url *url.URL) bool {
	for _, str := range outOfStockStrings {
		if strings.Contains(
			strings.ToLower(body),
			strings.ToLower(str),
		) {
			fmt.Printf("%s not in stock at %s!\n", cpu.name, url.Host)
			return false
		}

	}

	println(body)

	inStockAlert(cpu, url)
	return true
}

func inStockAlert(cpu LinkForCPU, url *url.URL) {
	sendTextMessage(cpu, url)
	fmt.Printf("%s is in stock at %s!\n", cpu.name, url.Host)

	t := toast.Notification{
		AppID:               "GimmeMyCPU",
		Title:               "Stock Alert!",
		Message:             fmt.Sprintf("%s is in stock at %s!", cpu.name, url),
		Icon:                "",
		ActivationType:      "",
		ActivationArguments: cpu.url,
		Actions:             nil,
		Audio:               toast.LoopingAlarm9,
		Loop:                true,
	}
	_ = t.Push()
}

func sendText(msg string) {
	client := twilio.NewClient(
		AccountSID,
		AuthToken,
		nil,
	)

	client.Messages.SendMessage(
		"15615108136",
		"17723418776",
		msg,
		nil,
	)
}

func sendTextMessage(cpu LinkForCPU, url *url.URL) {
	sendText(fmt.Sprintf("%s is in stock!\n%s", cpu.name, cpu.url))
}
