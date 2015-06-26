package mp

import (
    "fmt"
    "io/ioutil"
    "log"
    "regexp"
    "time"
    "strconv"
    "net/http"
    "math"

    "github.com/qiniu/iconv"
    "github.com/a2n/alu"
)

type MovableProperty struct {
    Amount int
    PageNo uint8
    Items map[string]*Item
}

func NewService() *MovableProperty {
    return &MovableProperty {
	PageNo: 1,
	Items: make(map[string]*Item, 0),
    }
}

func (mp *MovableProperty) Get() {
    html := mp.getPage()
    reAmount := regexp.MustCompile(`共(\d+)筆`)
    amountStr := reAmount.FindAllStringSubmatch(html, -1)
    amount, err := strconv.Atoi(amountStr[0][1])
    if err != nil {
	    log.Printf("%s has error, %s", alu.Caller(), err.Error())
	    return
    }
    mp.Amount = amount
    fmt.Printf("%d items.\n", mp.Amount)

    for mp.PageNo <= uint8(math.Ceil(float64(mp.Amount) / 50.0)) {
	mp.page(html)
	fmt.Printf("pageNo: %d, items: %d, amount: %d\n", mp.PageNo, len(mp.Items), mp.Amount)
	mp.PageNo++
	html = mp.getPage()
    }
}

func (mp *MovableProperty) getPage() string {
    urlStr := fmt.Sprintf("http://www.tpkonsale.moj.gov.tw/sale/sale1QryM.asp?pageno=%d", mp.PageNo)
    resp, err := http.DefaultClient.Get(urlStr)
    if err != nil {
	log.Printf("%s has error, %s", alu.Caller(), err.Error())
	return ""
    }

    b, err := ioutil.ReadAll(resp.Body)
    if err != nil {
	log.Printf("%s has error, %s", alu.Caller(), err.Error())
	return ""
    }
    resp.Body.Close()

    cd, err := iconv.Open("utf-8//ignore", "big5")
    if err != nil {
	    log.Printf("%s has error, %s", alu.Caller(), err.Error())
	    return ""
    }

    utf8 := cd.ConvString(string(b))

    return utf8
}

type Item struct {
    Branch string
    DateTime time.Time
    Kind string
    Attachments []string
}

func (mp *MovableProperty) page(html string) {
    reNoLoc := regexp.MustCompile(`(\d{3}-\d{2}-\d{8})<BR> \((.{2})`)
    noLocSlice := reNoLoc.FindAllStringSubmatch(html, -1)

    reDate := regexp.MustCompile(`(\d{4})/(\d)/(\d+) (.{2}) (\d{2}):(\d{2}):(\d{2})`)
    dateSlice := reDate.FindAllStringSubmatch(html, -1)

    reKind := regexp.MustCompile(`<td width="12%">(.{2})`)
    kindSlice := reKind.FindAllStringSubmatch(html, -1)

    reAttachment := regexp.MustCompile(`/TPKOnSale/((\d+)\w+\.\w+)`)
    attachmentSlice := reAttachment.FindAllStringSubmatch(html, -1)
    lastIndex := 0

    for k, v := range noLocSlice {
	    // date
	    y, _ := strconv.Atoi(dateSlice[k][1])
	    m, _ := strconv.Atoi(dateSlice[k][2])
	    d, _ := strconv.Atoi(dateSlice[k][3])
	    hh, _ := strconv.Atoi(dateSlice[k][5])
	    mm, _ := strconv.Atoi(dateSlice[k][6])
	    date := time.Date(y, time.Month(m), d, hh, mm, 0, 0, time.Local)

	    // attachments
	    as := make([]string, 0)
	    for {
		    as = append(as, attachmentSlice[lastIndex][1])
		    pin, _ := strconv.Atoi(attachmentSlice[lastIndex][2])
		    if lastIndex == len(attachmentSlice) - 1 {
			    break
		    }
		    nextPin, _ := strconv.Atoi(attachmentSlice[lastIndex + 1][2])
		    lastIndex++
		    if pin != nextPin {
			    break
		    }
	    }

	    i := Item {
		    Branch: v[2],
		    DateTime: date,
		    Kind: kindSlice[k][1],
		    Attachments: as,
	    }
	    mp.Items[v[1]] = &i

	    /*
	    fmt.Printf("No:\t%s\n", v[1])
	    fmt.Printf("Branch:\t%s\n", i.Branch)
	    fmt.Printf("Date:\t%s\n", i.DateTime.String())
	    fmt.Printf("Kind:\t%s\n", i.Kind)
	    fmt.Printf("Files:\t%s\n", i.Attachments)
	    fmt.Printf("\n\n")
	    */
    }
}
