package pdf

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type HeaderInfo struct {
	CompanyName string
	Address     string
	Logo        string
	Phone       string
	Email       string
}

type CustomerInfo struct {
	Name  string
	Phone string
	Email string
}
type OrderInfo struct {
	Date   time.Time
	Status string
}
type InvoiceInfo struct {
	OrderNo string
	Date    time.Time
}
type PaymentInfo struct {
	InvoiceNo string
	Date      time.Time
	Status    string
}
type OrderTemplate struct {
	pdf    *PDFBuilder
	Labels map[string]string
}

func GetOrderTemplate(lang string) *OrderTemplate {
	file, err := ioutil.ReadFile(filepath.Join("asset", "pdf_template.json"))
	translation := map[string]map[string]string{}
	err = json.Unmarshal([]byte(file), &translation)
	if err != nil {
		panic(err)
	}
	if v, ok := translation[lang]; ok {
		return &OrderTemplate{
			pdf:    NewPDFBuilder(),
			Labels: v,
		}
	}

	return &OrderTemplate{
		pdf:    NewPDFBuilder(),
		Labels: translation["EN"],
	}
}
func (temp *OrderTemplate) SetOrderNumber(orderNo string) {
	temp.pdf.SetFont("B", 12)
	orderNumber := fmt.Sprintf("%s: %s", temp.Labels["order_number"], orderNo)
	y := 44.0
	temp.pdf.RightCell(y, orderNumber, "#f9932d", "", Margin{}, Padding{Top: 2, Bottom: 2, Left: 2, Right: 2})
}

func (temp *OrderTemplate) SetInvoiceNumber(invoiceNo string) {
	temp.pdf.SetFont("B", 12)
	orderNumber := fmt.Sprintf("%s: %s", temp.Labels["invoice_number"], invoiceNo)
	y := 44.0
	temp.pdf.RightCell(y, orderNumber, "#f9932d", "", Margin{}, Padding{Top: 2, Bottom: 2, Left: 2, Right: 2})
}

func (temp *OrderTemplate) SetPaymentNumber(invoiceNo string) {
	temp.pdf.SetFont("B", 12)
	orderNumber := fmt.Sprintf("%s: %s", temp.Labels["payment_number"], invoiceNo)
	y := 44.0
	temp.pdf.RightCell(y, orderNumber, "#f9932d", "", Margin{}, Padding{Top: 2, Bottom: 2, Left: 2, Right: 2})

}
func (temp *OrderTemplate) SetCustomerInfo(cusInfo CustomerInfo) {
	y := 58.0
	temp.pdf.SetY(y)
	temp.pdf.SetFont("B", 10)
	temp.pdf.LeftCell(y, temp.Labels["order_to"]+":", "", "", Margin{}, Padding{Bottom: 2})
	temp.pdf.SetFont("", 10)
	temp.pdf.LeftCell(temp.pdf.GetY(), temp.Labels["customer_name"]+" - "+cusInfo.Name, "", "", Margin{}, Padding{Bottom: 2})
	temp.pdf.LeftCell(temp.pdf.GetY(), temp.Labels["customer_phone"]+" - "+cusInfo.Phone, "", "", Margin{}, Padding{Bottom: 2})
	temp.pdf.LeftCell(temp.pdf.GetY(), temp.Labels["customer_email"]+"- "+cusInfo.Email, "", "", Margin{}, Padding{Bottom: 2})
}

func (temp *OrderTemplate) SetOrderInfo(orderInfo OrderInfo) {
	y := 58.0
	temp.pdf.SetY(y)
	temp.pdf.SetFont("B", 10)
	temp.pdf.RightCell(y, temp.Labels["order_date"]+":", "", "", Margin{}, Padding{Bottom: 2})

	temp.pdf.SetFont("", 10)
	temp.pdf.RightCell(temp.pdf.GetY(), orderInfo.Date.Format("2006-01-02"), "", "", Margin{}, Padding{Bottom: 2})

	temp.pdf.SetFont("B", 10)
	temp.pdf.RightCell(temp.pdf.GetY(), temp.Labels["order_status"]+":", "", "", Margin{}, Padding{Bottom: 2})

	temp.pdf.SetFont("", 10)
	temp.pdf.RightCell(temp.pdf.GetY(), orderInfo.Status, "", "", Margin{}, Padding{Bottom: 2})
}

func (temp *OrderTemplate) SetInvoiceInfo(invoiceInfo InvoiceInfo) {
	y := 58.0
	temp.pdf.SetY(y)
	temp.pdf.SetFont("B", 10)
	temp.pdf.RightCell(y, temp.Labels["invoice_date"]+":", "", "", Margin{}, Padding{Bottom: 2})
	temp.pdf.SetFont("", 10)
	temp.pdf.RightCell(temp.pdf.GetY(), invoiceInfo.Date.Format("2006-01-02"), "", "", Margin{}, Padding{Bottom: 2})

	temp.pdf.SetFont("B", 10)
	temp.pdf.RightCell(temp.pdf.GetY(), temp.Labels["ref_order_no"]+":", "", "", Margin{}, Padding{Bottom: 2})

	temp.pdf.SetFont("", 10)
	temp.pdf.RightCell(temp.pdf.GetY(), invoiceInfo.OrderNo, "", "", Margin{}, Padding{Bottom: 2})
}

func (temp *OrderTemplate) SetPaymentInfo(paymentInfo PaymentInfo) {
	y := 58.0
	temp.pdf.SetY(y)
	temp.pdf.SetFont("B", 10)
	temp.pdf.RightCell(y, temp.Labels["payment_date"]+":", "", "", Margin{}, Padding{Bottom: 2})

	temp.pdf.SetFont("", 10)
	temp.pdf.RightCell(temp.pdf.GetY(), paymentInfo.Date.Format("2006-01-02"), "", "", Margin{}, Padding{Bottom: 2})

	temp.pdf.SetFont("B", 10)
	temp.pdf.RightCell(temp.pdf.GetY(), temp.Labels["payment_status"]+":", "", "", Margin{}, Padding{Bottom: 2})

	temp.pdf.SetFont("", 10)
	temp.pdf.RightCell(temp.pdf.GetY(), paymentInfo.Status, "", "", Margin{}, Padding{Bottom: 2})

	temp.pdf.SetFont("B", 10)
	temp.pdf.RightCell(temp.pdf.GetY(), temp.Labels["ref_invoice_no"]+":", "", "", Margin{}, Padding{Bottom: 2})

	temp.pdf.SetFont("", 10)
	temp.pdf.RightCell(temp.pdf.GetY(), paymentInfo.InvoiceNo, "", "", Margin{}, Padding{Bottom: 2})
}

func (temp *OrderTemplate) SetHeader(companyInfo *HeaderInfo) {
	//all page
	temp.pdf.SetHeaderFunc(func() {
		l, _, r, _ := temp.pdf.GetMargins()
		temp.pdf.DrawLine(20, "#f9932d", 0.3, Margin{
			Top:    0,
			Bottom: 0,
			Left:   int(l),
			Right:  int(r),
		})
		temp.pdf.Ln(10)
	})
	//first page
	y := 10.0 + 25.0

	if len(companyInfo.Logo) == 0 {
		temp.pdf.ImageOptions(filepath.Join("asset", "mac.png"), 10, 10, 25, 25, false, gofpdf.ImageOptions{
			ImageType:             "png",
			ReadDpi:               true,
			AllowNegativePosition: false,
		}, 0, "")
	} else {
		temp.pdf.LoadBase64Image(companyInfo.Logo, 10, 10, 25, 25, false, "", 0, "")
	}

	temp.pdf.SetFont("", 16)
	//company info
	temp.pdf.LeftCell(15, companyInfo.CompanyName, "", "", Margin{
		Left: 30,
	}, Padding{})
	temp.pdf.Ln(2)
	temp.pdf.SetFont("", 8)
	temp.pdf.LeftMultiCell(temp.pdf.GetY(), 0, companyInfo.Address,
		"", "", Margin{
			Left: 30,
		}, Padding{
			Top:    1,
			Bottom: 1})

	temp.pdf.Ln(2)
	if temp.pdf.GetY() > y {
		y = temp.pdf.GetY()
	}

	tel := companyInfo.Phone
	email := companyInfo.Email

	pageSizeW, _ := temp.pdf.GetPageSize()
	_, _, r, _ := temp.pdf.GetMargins()
	temp.pdf.SetFont("B", 8)
	temp.pdf.RightCellComplex(
		y, (pageSizeW - r),
		email,
		"#FFFFFF",
		"#f9932d",
		Margin{},
		Padding{Top: 1, Bottom: 1, Left: 2, Right: 2})
	sWidth := temp.pdf.GetStringWidth(email) + 2
	temp.pdf.SetFont("FA", 8)
	temp.pdf.RightCellComplex(
		y, (pageSizeW - r - sWidth), "\uf0e0",
		"#FFFFFF",
		"#f9932d",
		Margin{},
		Padding{Top: 1, Bottom: 1, Left: 2, Right: 2})

	sWidth += temp.pdf.GetStringWidth("\uf0e0") + 2
	temp.pdf.SetFont("B", 8)
	temp.pdf.RightCellComplex(
		y, (pageSizeW - r - sWidth),
		tel,
		"#FFFFFF",
		"#f9932d",
		Margin{},
		Padding{Top: 1, Bottom: 1, Left: 2, Right: 2})

	sWidth += temp.pdf.GetStringWidth(tel) + 2
	temp.pdf.SetFont("FA", 8)
	temp.pdf.RightCellComplex(
		y, (pageSizeW - r - sWidth), "\uf095",
		"#FFFFFF",
		"#f9932d",
		Margin{},
		Padding{Top: 1, Bottom: 1, Left: 2, Right: 2})
}
func (temp *OrderTemplate) GetOrderReceiptTableHeader() []string {
	return []string{
		temp.Labels["branch_name"],
		temp.Labels["storage_number"],
		temp.Labels["details"],
		temp.Labels["rent_period"],
		temp.Labels["start_to_end"],
		temp.Labels["total_amount"]}
}
func (temp *OrderTemplate) SetTable(column []string, data []string) {

	header := TableHeader{
		Color:           ("#ffffff"),
		BackgroundColor: ("#f9932d"),
		Data:            column,
	}
	body := TableBody{
		Color:           ("#000000"),
		BackgroundColor: ("#f5f5f5"),
		Data:            data,
	}
	colNum := len(header.Data)
	cellWidth := make([]float64, colNum)
	initY := temp.pdf.GetY() + 10
	if initY < 90 {
		initY = 90
	}
	temp.pdf.Table(initY, func(x, y, maxWidth float64) float64 {
		for i, str := range header.Data {
			wH := temp.pdf.GetStringWidth(str) + 1
			wB := temp.pdf.GetStringWidth(body.Data[i]) + 1
			if wH > wB {
				cellWidth[i] = wH
			} else {
				cellWidth[i] = wB
			}
			cellWidth[i] = maxWidth / float64(colNum)
		}
		temp.pdf.SetFont("B", 8)
		temp.pdf.SetColor(FillColor, header.BackgroundColor)
		temp.pdf.SetColor(TextColor, header.Color)
		temp.pdf.WrapTextTableCell(temp.pdf.GetY(), temp.pdf.GetX(), cellWidth, header.Data)
		return temp.pdf.GetY()
	}, func(x, y, maxWidth float64) float64 {
		temp.pdf.SetFont("", 8)
		temp.pdf.SetColor(FillColor, body.BackgroundColor)
		temp.pdf.SetColor(TextColor, body.Color)
		temp.pdf.WrapTextTableCell(temp.pdf.GetY(), temp.pdf.GetX(), cellWidth, body.Data)
		return temp.pdf.GetY()
	})
	temp.pdf.Ln(15)

	pageSizeW, _ := temp.pdf.GetPageSize()
	_, _, r, _ := temp.pdf.GetMargins()
	temp.pdf.SetFont("", 12)
	newY := temp.pdf.GetY()
	temp.pdf.SetColor(FillColor, body.BackgroundColor)
	temp.pdf.SetColor(TextColor, body.Color)

	totalAmount := data[len(data)-1]
	sWidth := temp.pdf.GetStringWidth(totalAmount)
	temp.pdf.RightCellComplex(newY, pageSizeW-r, data[len(data)-1], "", "", Margin{}, Padding{Right: 3})

	temp.pdf.SetFont("B", 14)
	temp.pdf.RightCellComplex(newY, (pageSizeW-r-sWidth)-2, column[len(data)-1], "#f9932d", "", Margin{}, Padding{Right: 15})
}

func (temp *OrderTemplate) SetDefaultFooter() {

	temp.pdf.SetFooterFunc(func() {
		bottom := 20.0
		pageSizeW, pageSizeH := temp.pdf.GetPageSize()

		l, _, r, _ := temp.pdf.GetMargins()
		temp.pdf.DrawLine(pageSizeH-bottom, "#f9932d", 0.3, Margin{
			Top:    0,
			Bottom: 0,
			Left:   int(l),
			Right:  int(r),
		})
		temp.pdf.Ln(2)
		temp.pdf.SetFont("", 8)
		_, fontHeight := temp.pdf.GetFontSize()

		temp.pdf.CellFormat(pageSizeW-l-r, fontHeight, temp.Labels["footer"], "0", 0, "CM", false, 0, "")
		temp.pdf.RightCell(temp.pdf.GetY(), strconv.Itoa(
			temp.pdf.PageCount()), "", "", Margin{}, Padding{Right: 5})
	})

}

func (temp *OrderTemplate) SetFooter(content string) {

	temp.pdf.SetFooterFunc(func() {
		bottom := 20.0
		pageSizeW, pageSizeH := temp.pdf.GetPageSize()

		l, _, r, _ := temp.pdf.GetMargins()
		temp.pdf.DrawLine(pageSizeH-bottom, "#f9932d", 0.3, Margin{
			Top:    0,
			Bottom: 0,
			Left:   int(l),
			Right:  int(r),
		})
		temp.pdf.Ln(2)
		temp.pdf.SetFont("", 8)
		_, fontHeight := temp.pdf.GetFontSize()

		temp.pdf.CellFormat(pageSizeW-l-r, fontHeight, content, "0", 0, "CM", false, 0, "")
		temp.pdf.RightCell(temp.pdf.GetY(), strconv.Itoa(
			temp.pdf.PageCount()), "", "", Margin{}, Padding{Right: 5})
	})

}

func (temp *OrderTemplate) SetRemarks(remarks string) {
	temp.pdf.SetFont("B", 10)
	temp.pdf.LeftCell(297-70, temp.Labels["remarks"]+":", "", "", Margin{}, Padding{})
	temp.pdf.SetFont("", 10)
	temp.pdf.LeftMultiCell(297-65, 0, remarks,
		"", "", Margin{}, Padding{})
}

func (temp *OrderTemplate) SetTandC(tandc *RichTextBlocks) {
	temp.pdf.AddPage()
	temp.pdf.SetFont("B", 16)
	_, fontHeight := temp.pdf.GetFontSize()
	pageSizeW, _ := temp.pdf.GetPageSize()
	l, t, r, _ := temp.pdf.GetMargins()
	temp.pdf.CellFormat(pageSizeW-l-r, fontHeight, temp.Labels["t_and_c"], "0", 0, "CM", false, 0, "")
	temp.pdf.Ln(12)
	// ht := temp.pdf.PointConvert(10)
	//write := func(str string) {
	//	temp.pdf.LeftMultiCell(temp.pdf.GetY(), 0, str,
	//		"", "", Margin{}, Padding{})
	//}
	// html := temp.pdf.HTMLBasicNew()
	// html.Write(ht, tandc)
	for _, block := range tandc.Blocks {
		align := "left"
		if block.Data != nil {
			align = block.Data.TextAlign
		}
		style := block.Type
		wordStyle := ""
		paddingLeft := 0
		for _, inline := range block.InlineStyleRanges {
			if inline.Style == "BOLD" {
				wordStyle = "B"
				break
			}
		}
		switch style {
		case "unordered-list-item":
			{
				paddingLeft = 3
				temp.pdf.SetFont(wordStyle, 8)
			}
		case "header-one":
			{
				temp.pdf.SetFont(wordStyle, 15)
			}
		case "header-two":
			{
				temp.pdf.SetFont(wordStyle, 14)
			}
		case "header-three":
			{
				temp.pdf.SetFont(wordStyle, 13)
			}
		case "header-four":
			{
				temp.pdf.SetFont(wordStyle, 12)
			}
		case "header-five":
			{
				temp.pdf.SetFont(wordStyle, 11)
			}
		case "header-six":
			{
				temp.pdf.SetFont(wordStyle, 10)
			}
		case "unstyled":
			{
				temp.pdf.SetFont(wordStyle, 8)
			}
		default:
			temp.pdf.SetFont("", 8)
		}

		if temp.pdf.GetY() > 260 {
			temp.pdf.AddPage()

			temp.pdf.SetY(30)
		}

		temp.pdf.SetMargins(20, t, 20)
		switch align {
		case "start":
			{
				temp.pdf.LeftMultiCell(temp.pdf.GetY(), 0, block.Text,
					"", "", Margin{
						Left: 0,
					}, Padding{
						Left:   paddingLeft,
						Top:    1,
						Bottom: 1})
			}
		case "left":
			{
				temp.pdf.LeftMultiCell(temp.pdf.GetY(), 0, block.Text,
					"", "", Margin{
						Left: 0,
					}, Padding{
						Left:   paddingLeft,
						Top:    1,
						Bottom: 1})
			}
		case "right":
			{
				temp.pdf.RightMultiCell(temp.pdf.GetY(), 0, block.Text,
					"", "", Margin{
						Left: 0,
					}, Padding{
						Left:   paddingLeft,
						Top:    1,
						Bottom: 1})
			}
		case "end":
			{
				temp.pdf.RightMultiCell(temp.pdf.GetY(), 0, block.Text,
					"", "", Margin{
						Left: 0,
					}, Padding{
						Left:   paddingLeft,
						Top:    1,
						Bottom: 1})
			}
		case "center":
			{
				temp.pdf.CenterMultiCell(temp.pdf.GetY(), 0, block.Text,
					"", "", Margin{
						Left: 0,
					}, Padding{
						Left:   paddingLeft,
						Top:    1,
						Bottom: 1})
			}
		}
		temp.pdf.SetMargins(l, t, r)
	}

}

func (temp *OrderTemplate) SaveAs(file string) error {

	return temp.pdf.OutputFileAndClose(file)
}
func (temp *OrderTemplate) Output(writer io.Writer) error {

	return temp.pdf.Output(writer)
}
