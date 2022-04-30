package pdf

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"

	color "dmglab.com/mac-crm/pkg/util/color"
)

type Align struct {
	TextAlign string `json:"text-align"`
}
type InlineStyleRanges struct {
	Length int    `json:"length"`
	Offset int    `json:"offset"`
	Style  string `json:"style"`
}
type Block struct {
	Data *Align `json:"data,omitempty"`
	// Depth             int           `json:"depth"`
	// EntityRanges      []interface{} `json:"entityRanges"`
	InlineStyleRanges []InlineStyleRanges `json:"inlineStyleRanges"`
	Key               string              `json:"key"`
	Text              string              `json:"text"`
	Type              string              `json:"type"`
}
type RichTextBlocks struct {
	Blocks []Block `json:"blocks"`
	// EntityMap struct {
	// 	Num0 struct {
	// 		Data struct {
	// 			Map struct {
	// 				Data struct {
	// 					TargetOption string `json:"targetOption"`
	// 					Title        string `json:"title"`
	// 					URL          string `json:"url"`
	// 				} `json:"data"`
	// 				Mutability string `json:"mutability"`
	// 				Type       string `json:"type"`
	// 			} `json:"_map"`
	// 			TargetOption string `json:"targetOption"`
	// 			Title        string `json:"title"`
	// 			URL          string `json:"url"`
	// 		} `json:"data"`
	// 		Mutability string `json:"mutability"`
	// 		Type       string `json:"type"`
	// 	} `json:"0"`
	// 	Num1 struct {
	// 		Data struct {
	// 			Map struct {
	// 				Data struct {
	// 					TargetOption string `json:"targetOption"`
	// 					Title        string `json:"title"`
	// 					URL          string `json:"url"`
	// 				} `json:"data"`
	// 				Mutability string `json:"mutability"`
	// 				Type       string `json:"type"`
	// 			} `json:"_map"`
	// 			TargetOption string `json:"targetOption"`
	// 			Title        string `json:"title"`
	// 			URL          string `json:"url"`
	// 		} `json:"data"`
	// 		Mutability string `json:"mutability"`
	// 		Type       string `json:"type"`
	// 	} `json:"1"`
	// 	Num2 struct {
	// 		Data struct {
	// 			Map struct {
	// 				Data struct {
	// 					TargetOption string `json:"targetOption"`
	// 					Title        string `json:"title"`
	// 					URL          string `json:"url"`
	// 				} `json:"data"`
	// 				Mutability string `json:"mutability"`
	// 				Type       string `json:"type"`
	// 			} `json:"_map"`
	// 			TargetOption string `json:"targetOption"`
	// 			Title        string `json:"title"`
	// 			URL          string `json:"url"`
	// 		} `json:"data"`
	// 		Mutability string `json:"mutability"`
	// 		Type       string `json:"type"`
	// 	} `json:"2"`
	// 	Num3 struct {
	// 		Data struct {
	// 			Map struct {
	// 				Data struct {
	// 					TargetOption string `json:"targetOption"`
	// 					Title        string `json:"title"`
	// 					URL          string `json:"url"`
	// 				} `json:"data"`
	// 				Mutability string `json:"mutability"`
	// 				Type       string `json:"type"`
	// 			} `json:"_map"`
	// 			TargetOption string `json:"targetOption"`
	// 			Title        string `json:"title"`
	// 			URL          string `json:"url"`
	// 		} `json:"data"`
	// 		Mutability string `json:"mutability"`
	// 		Type       string `json:"type"`
	// 	} `json:"3"`
	// 	Num4 struct {
	// 		Data struct {
	// 			Map struct {
	// 				Data struct {
	// 					TargetOption string `json:"targetOption"`
	// 					Title        string `json:"title"`
	// 					URL          string `json:"url"`
	// 				} `json:"data"`
	// 				Mutability string `json:"mutability"`
	// 				Type       string `json:"type"`
	// 			} `json:"_map"`
	// 			TargetOption string `json:"targetOption"`
	// 			Title        string `json:"title"`
	// 			URL          string `json:"url"`
	// 		} `json:"data"`
	// 		Mutability string `json:"mutability"`
	// 		Type       string `json:"type"`
	// 	} `json:"4"`
	// } `json:"entityMap"`
}
type PDFBuilder struct {
	Language string
	Fonts    map[string]string
	*gofpdf.Fpdf
}
type Margin struct {
	Top    int
	Bottom int
	Left   int
	Right  int
}
type Padding struct {
	Top    int
	Bottom int
	Left   int
	Right  int
}

//var font = "NotoSansTC-Regular"

func NewPDFBuilder() *PDFBuilder {

	font := map[string]string{}
	pdf := gofpdf.New("P", "mm", "A4", "") //210 x 297 = A4
	font["R"] = "TaipeiSansTCBeta-Light"
	pdf.AddUTF8Font("TaipeiSansTCBeta-Light", "", filepath.Join("asset", "TaipeiSansTCBeta-Light.ttf"))
	font["B"] = "TaipeiSansTCBeta-Bold"
	pdf.AddUTF8Font("TaipeiSansTCBeta-Bold", "", filepath.Join("asset", "TaipeiSansTCBeta-Bold.ttf"))

	font["FA"] = "Font Awesome 6 Free-Solid-900"
	pdf.AddUTF8Font("Font Awesome 6 Free-Solid-900", "", filepath.Join("asset", "Font Awesome 6 Free-Solid-900.ttf"))

	pdf.AddPage()
	// pdf.SetAutoPageBreak(true, 30)
	pdf.SetMargins(10, 40, 10)
	return &PDFBuilder{
		Language: "",
		Fonts:    font,
		Fpdf:     pdf,
	}
}

func (pdf *PDFBuilder) SetFont(style string, size float64) {
	if style == "B" {
		pdf.Fpdf.SetFont(pdf.Fonts[style], "", size)
	} else if style == "FA" {
		pdf.Fpdf.SetFont(pdf.Fonts[style], "", size)
	} else {
		pdf.Fpdf.SetFont(pdf.Fonts["R"], "", size)
	}

}

func (pdf *PDFBuilder) LoadBase64Image(imageStr string, x, y, w, h float64, flow bool, tp string, link int, linkStr string) error {

	indicator := "base64,"
	if len(imageStr) == 0 || len(imageStr) < len(indicator) {
		return nil
	}
	imgName := uuid.New().String()
	index := strings.Index(imageStr, indicator)

	str := imageStr[index+len(indicator):]
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return err
	}
	i, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return err
	}
	buff := new(bytes.Buffer)
	err = png.Encode(buff, i)
	if err != nil {
		return err
	}
	// reader := strings.NewReader(base64.StdEncoding.EncodeToString([]byte(cmp.SettingValue)))
	_ = pdf.RegisterImageOptionsReader(imgName+".png", gofpdf.ImageOptions{ImageType: "png"}, buff)

	pdf.Image(imgName+".png", x, y, w, h, flow, tp, link, linkStr)
	return nil
}
func (pdf *PDFBuilder) DrawLine(initY float64, colorCode string, lineHeight float64, margin Margin) {
	pageSizeW, _ := pdf.GetPageSize()
	pdf.SetY(initY)
	pdf.SetX(float64(margin.Left))
	rgb, _ := color.Hex2RGB(color.Hex(colorCode))
	pdf.SetFillColor(rgb.Red, rgb.Green, rgb.Blue)
	pdf.CellFormat(pageSizeW-float64(margin.Left)-float64(margin.Right), lineHeight, "", "0", 0, "mm", true, 0, "")
}
func (pdf *PDFBuilder) RightCellComplex(initY, initX float64, str string, textColor, backgroundColor string, margin Margin, padding Padding) {
	pdf.SetY(initY)
	//pageSizeW, _ := pdf.GetPageSize()
	//_, _, r, _ := pdf.GetMargins()
	_, fontHeight := pdf.GetFontSize()
	strWidth := pdf.GetStringWidth(str)
	pdf.SetX(initX - strWidth - float64(margin.Right) - float64(padding.Left) - float64(padding.Right))
	lineHeight := fontHeight + float64(padding.Top) + float64(padding.Bottom)
	width := strWidth + float64(padding.Left) + float64(padding.Right)
	fill := false
	if len(backgroundColor) > 0 {
		fill = true
		rgb, err := color.Hex2RGB(color.Hex(backgroundColor))
		if err == nil {
			pdf.SetFillColor(rgb.Red, rgb.Green, rgb.Blue)
		}
	}
	if len(textColor) > 0 {
		rgb, err := color.Hex2RGB(color.Hex(textColor))
		if err == nil {
			pdf.SetTextColor(rgb.Red, rgb.Green, rgb.Blue)
		} else {
			pdf.SetTextColor(0, 0, 0)
		}
	} else {
		pdf.SetTextColor(0, 0, 0)
	}
	if padding.Left > 0 {
		pdf.CellFormat(width, lineHeight, "", "0", 0, "L", fill, 0, "")
		pdf.SetX(initX - strWidth - float64(margin.Right) - float64(padding.Right))
		pdf.CellFormat(width, lineHeight, str, "0", 0, "L", false, 0, "")
	} else {
		pdf.CellFormat(width, lineHeight, str, "0", 0, "L", fill, 0, "")
	}
	pdf.SetY(initY + lineHeight)
}
func (pdf *PDFBuilder) RightCell(initY float64, str string, textColor, backgroundColor string, margin Margin, padding Padding) {
	pdf.SetY(initY)
	pageSizeW, _ := pdf.GetPageSize()
	_, _, r, _ := pdf.GetMargins()
	_, fontHeight := pdf.GetFontSize()
	strWidth := pdf.GetStringWidth(str)
	pdf.SetX(pageSizeW - r - strWidth - float64(margin.Right) - float64(padding.Left) - float64(padding.Right))
	lineHeight := fontHeight + float64(padding.Top) + float64(padding.Bottom)
	width := strWidth + float64(padding.Left) + float64(padding.Right)
	fill := false
	if len(backgroundColor) > 0 {
		fill = true
		rgb, err := color.Hex2RGB(color.Hex(backgroundColor))
		if err == nil {
			pdf.SetFillColor(rgb.Red, rgb.Green, rgb.Blue)
		}
	}
	if len(textColor) > 0 {
		rgb, err := color.Hex2RGB(color.Hex(textColor))
		if err == nil {
			pdf.SetTextColor(rgb.Red, rgb.Green, rgb.Blue)
		} else {
			pdf.SetTextColor(0, 0, 0)
		}
	} else {
		pdf.SetTextColor(0, 0, 0)
	}
	if padding.Left > 0 {
		pdf.CellFormat(width, lineHeight, "", "0", 0, "L", fill, 0, "")
		pdf.SetX(pageSizeW - r - strWidth - float64(margin.Right) - float64(padding.Right))
		pdf.CellFormat(width, lineHeight, str, "0", 0, "L", false, 0, "")
	} else {
		pdf.CellFormat(width, lineHeight, str, "0", 0, "L", fill, 0, "")
	}
	pdf.SetY(initY + lineHeight)
}

func (pdf *PDFBuilder) LeftCell(initY float64, str string, textColor, backgroundColor string, margin Margin, padding Padding) {
	pdf.SetY(initY)
	_, fontHeight := pdf.GetFontSize()
	strWidth := pdf.GetStringWidth(str) + 2
	l, _, _, _ := pdf.GetMargins()
	x := float64(margin.Left) + float64(padding.Left) + l
	pdf.SetX(x)
	lineHeight := fontHeight + float64(padding.Top) + float64(padding.Bottom)
	width := strWidth + float64(padding.Left) + float64(padding.Right)
	fill := false
	if len(backgroundColor) > 0 {
		fill = true
		rgb, err := color.Hex2RGB(color.Hex(backgroundColor))
		if err == nil {
			pdf.SetFillColor(rgb.Red, rgb.Green, rgb.Blue)
		}
	}
	if len(textColor) > 0 {
		rgb, err := color.Hex2RGB(color.Hex(textColor))
		if err == nil {
			pdf.SetTextColor(rgb.Red, rgb.Green, rgb.Blue)
		} else {
			pdf.SetTextColor(0, 0, 0)
		}
	} else {
		pdf.SetTextColor(0, 0, 0)
	}
	if padding.Left > 0 {
		pdf.CellFormat(width, lineHeight, "", "0", 0, "L", fill, 0, "")
		pdf.SetX(x)
		pdf.CellFormat(width, lineHeight, str, "0", 0, "L", false, 0, "")
	} else {
		pdf.CellFormat(width, lineHeight, str, "0", 0, "L", fill, 0, "")
	}
	pdf.SetY(initY + lineHeight)
}

func (pdf *PDFBuilder) RightMultiCell(initY, maxWidth float64, str string, textColor, backgroundColor string, margin Margin, padding Padding) {
	pdf.SetY(initY)
	pageSizeW, _ := pdf.GetPageSize()
	_, fontHeight := pdf.GetFontSize()
	_, _, r, _ := pdf.GetMargins()
	strWidth := pdf.GetStringWidth(str)
	x := pageSizeW - strWidth - float64(margin.Right) - float64(padding.Left) - float64(padding.Right) - r
	pdf.SetX(x)

	lineHeight := fontHeight + float64(padding.Top) + float64(padding.Bottom)
	// width := strWidth + float64(padding.Left) + float64(padding.Right)
	containerWidth := maxWidth
	if containerWidth == 0 {
		containerWidth = (pageSizeW - r - x)
	}
	totalLine := pdf.LineCount(str, containerWidth, Padding{Left: padding.Left, Right: padding.Right})

	fill := false
	if len(backgroundColor) > 0 {
		fill = true
		rgb, err := color.Hex2RGB(color.Hex(backgroundColor))
		if err == nil {
			pdf.SetFillColor(rgb.Red, rgb.Green, rgb.Blue)
		}
	}
	if len(textColor) > 0 {
		rgb, err := color.Hex2RGB(color.Hex(textColor))
		if err == nil {
			pdf.SetTextColor(rgb.Red, rgb.Green, rgb.Blue)
		} else {
			pdf.SetTextColor(0, 0, 0)
		}
	} else {
		pdf.SetTextColor(0, 0, 0)
	}
	if padding.Left > 0 {
		pdf.MultiCell(containerWidth, lineHeight, "", "0", "L", fill)
		pdf.SetX(pageSizeW - strWidth - float64(margin.Right) - float64(padding.Right))
		pdf.MultiCell(containerWidth, lineHeight, str, "0", "L", false)
	} else {
		pdf.MultiCell(containerWidth, lineHeight, str, "0", "L", fill)
	}
	pdf.SetY(initY + (lineHeight * float64(totalLine)))
}
func (pdf *PDFBuilder) LineCount(content string, containerWidth float64, padding Padding) int {
	strWidth := pdf.GetStringWidth(content) + 2
	totalLine := 1
	contentWidth := strWidth + float64(padding.Left) + float64(padding.Right)
	if contentWidth > containerWidth {
		lineCount := (contentWidth / containerWidth)
		totalLine = int(contentWidth / containerWidth)
		if lineCount-float64(totalLine) > 0 {
			totalLine++
		}
	}
	scanner := bufio.NewScanner(strings.NewReader(content))
	scanner.Split(bufio.ScanLines)
	count := 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		// fmt.Fprintln(os.Stderr, "reading input:", err)
	}
	return totalLine + count
}

func (pdf *PDFBuilder) CenterMultiCell(initY, maxWidth float64, str string, textColor, backgroundColor string, margin Margin, padding Padding) {
	pdf.SetY(initY)
	pageSizeW, _ := pdf.GetPageSize()
	l, _, r, _ := pdf.GetMargins()
	_, fontHeight := pdf.GetFontSize()
	// strWidth := pdf.GetStringWidth(str) + 2
	x := l + float64(margin.Left) + float64(padding.Left)
	pdf.SetX(x)
	lineHeight := fontHeight + float64(padding.Top) + float64(padding.Bottom)
	// width := strWidth + float64(padding.Left) + float64(padding.Right)
	containerWidth := maxWidth
	if containerWidth == 0 {
		containerWidth = (pageSizeW - r - x)
	}
	totalLine := pdf.LineCount(str, containerWidth, Padding{Left: padding.Left, Right: padding.Right})
	fill := false
	if len(backgroundColor) > 0 {
		fill = true
		rgb, err := color.Hex2RGB(color.Hex(backgroundColor))
		if err == nil {
			pdf.SetFillColor(rgb.Red, rgb.Green, rgb.Blue)
		}
	}
	if len(textColor) > 0 {
		rgb, err := color.Hex2RGB(color.Hex(textColor))
		if err == nil {
			pdf.SetTextColor(rgb.Red, rgb.Green, rgb.Blue)
		} else {
			pdf.SetTextColor(0, 0, 0)
		}
	} else {
		pdf.SetTextColor(0, 0, 0)
	}
	if padding.Left > 0 {
		pdf.MultiCell(containerWidth, lineHeight, "", "0", "L", fill)
		pdf.SetX(float64(margin.Left) + float64(padding.Left))
		pdf.MultiCell(containerWidth, lineHeight, str, "0", "C", false)
	} else {
		pdf.MultiCell(containerWidth, lineHeight, str, "0", "C", fill)
	}
	pdf.SetY(initY + (lineHeight * float64(totalLine)))
}
func (pdf *PDFBuilder) LeftMultiCell(initY, maxWidth float64, str string, textColor, backgroundColor string, margin Margin, padding Padding) {
	pdf.SetY(initY)
	pageSizeW, _ := pdf.GetPageSize()
	l, _, r, _ := pdf.GetMargins()
	_, fontHeight := pdf.GetFontSize()
	// strWidth := pdf.GetStringWidth(str) + 2
	x := l + float64(margin.Left) + float64(padding.Left)
	pdf.SetX(x)
	lineHeight := fontHeight + float64(padding.Top) + float64(padding.Bottom)
	// width := strWidth + float64(padding.Left) + float64(padding.Right)
	containerWidth := maxWidth
	if containerWidth == 0 {
		containerWidth = (pageSizeW - r - x)
	}
	totalLine := pdf.LineCount(str, containerWidth, Padding{Left: padding.Left, Right: padding.Right})
	fill := false
	if len(backgroundColor) > 0 {
		fill = true
		rgb, err := color.Hex2RGB(color.Hex(backgroundColor))
		if err == nil {
			pdf.SetFillColor(rgb.Red, rgb.Green, rgb.Blue)
		}
	}
	if len(textColor) > 0 {
		rgb, err := color.Hex2RGB(color.Hex(textColor))
		if err == nil {
			pdf.SetTextColor(rgb.Red, rgb.Green, rgb.Blue)
		} else {
			pdf.SetTextColor(0, 0, 0)
		}
	} else {
		pdf.SetTextColor(0, 0, 0)
	}
	if padding.Left > 0 {
		pdf.MultiCell(containerWidth, lineHeight, "", "0", "L", fill)
		pdf.SetX(float64(margin.Left) + float64(padding.Left))
		pdf.MultiCell(containerWidth, lineHeight, str, "0", "L", false)
	} else {
		pdf.MultiCell(containerWidth, lineHeight, str, "0", "L", fill)
	}
	pdf.SetY(initY + (lineHeight * float64(totalLine)))
}

func (pdf *PDFBuilder) GetColWidth(colNum int) float64 {
	pageSizeW, _ := pdf.GetPageSize()
	l, _, r, _ := pdf.GetMargins()
	containerWidth := (pageSizeW - l - r)
	return containerWidth / float64(colNum)
}

func (pdf *PDFBuilder) Set2Column(initY float64, col1 func(x1, y1, maxWidth float64) float64, col2 func(x2, y2, maxWidth float64) float64) {
	pdf.SetY(initY)
	l, _, _, _ := pdf.GetMargins()
	x := l
	maxW := pdf.GetColWidth(2)
	col1EndY := col1(x, initY, maxW)
	pdf.SetY(initY)
	col2EndY := col2(x, initY, maxW)
	if col1EndY > col2EndY {
		pdf.SetY(col1EndY)
	} else {
		pdf.SetY(col2EndY)
	}
}

type TableHeader struct {
	Color           string
	BackgroundColor string
	Data            []string
}
type TableBody struct {
	Color           string
	BackgroundColor string
	Data            []string
}
type ColorType string

const FillColor ColorType = "FillColor"
const TextColor ColorType = "TextColor"

func (pdf *PDFBuilder) SetColor(colorType ColorType, colorCode string) {
	switch colorType {
	case FillColor:
		{
			rgb, err := color.Hex2RGB(color.Hex(colorCode))
			if err == nil {
				pdf.SetFillColor(rgb.Red, rgb.Green, rgb.Blue)
			} else {
				pdf.SetFillColor(0, 0, 0)
			}
		}
	case TextColor:
		{
			rgb, err := color.Hex2RGB(color.Hex(colorCode))
			if err == nil {
				pdf.SetTextColor(rgb.Red, rgb.Green, rgb.Blue)
			} else {
				pdf.SetTextColor(0, 0, 0)
			}
		}
	}

}

func (pdf *PDFBuilder) Table(initY float64, tableHead func(x, y, maxWidth float64) float64, tableBody func(x, y, maxWidth float64) float64) {
	pageSizeW, _ := pdf.GetPageSize()
	l, _, r, _ := pdf.GetMargins()
	containerWidth := (pageSizeW - l - r)
	pdf.SetY(initY)
	y := initY
	initX := l
	pdf.SetX(initX)
	//header
	endY := tableHead(initX, y, containerWidth)

	pdf.Line(l, y, pageSizeW-r, y)
	endY = tableBody(initX, endY, containerWidth)
	// pdf.SetX(initX)
	y = endY
	pdf.Line(l, y, pageSizeW-r, y)
	// pdf.CellFormat(containerWidth, 0.1, "", "T", 0, "", false, 0, "")
}

func (pdf *PDFBuilder) WrapTextTableCell(rowY, rowX float64, cellWidth []float64, data []string) {
	pdf.SetY(rowY)
	pdf.SetX(rowX)
	maxH := 0.0
	for i, _ := range data {
		cellHeight := pdf.MaxHeight(cellWidth[i], data) + 3
		if cellHeight > maxH {
			maxH = cellHeight
		}
		pdf.CellFormat(cellWidth[i], cellHeight, "", "LR", 0, "C", true, 0, "")
	}

	_, lineHeight := pdf.GetFontSize()
	pdf.SetY(rowY)
	pdf.SetX(rowX)

	for i, str := range data {
		totalLine := 0
		for _, str := range data {
			totalLine = pdf.LineCount(str, cellWidth[i], Padding{})
		}

		pdf.SetY(rowY + maxH/2 - (lineHeight*float64(totalLine))/2)
		pdf.SetX(rowX)
		pdf.MultiCell(cellWidth[i], lineHeight, str, "", "C", false)
		rowX += cellWidth[i]
	}

	pdf.SetY(rowY + maxH)
}

func (pdf *PDFBuilder) MaxHeight(size float64, data []string) float64 {
	numberOfLine := 0

	_, lineHeight := pdf.GetFontSize()
	for _, str := range data {
		totalLine := pdf.LineCount(str, size, Padding{})
		if totalLine > numberOfLine {
			numberOfLine = totalLine
		}
	}
	return lineHeight * float64(numberOfLine)
}
