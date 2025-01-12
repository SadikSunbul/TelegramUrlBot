package chart

import (
	"bytes"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func generateBarItems(values []int) []opts.BarData {
	items := make([]opts.BarData, 0)
	for _, v := range values {
		items = append(items, opts.BarData{Value: v})
	}
	return items
}

func CreateChart(xExsenDara []string, yExsenDara []int, title string) *bytes.Buffer {

	if len(xExsenDara) != len(yExsenDara) {
		// hatalı
		return nil
	}

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: title}),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(true)}),
	)
	bar.SetXAxis(xExsenDara).
		AddSeries("Değerler", generateBarItems(yExsenDara))

	buf := new(bytes.Buffer)
	bar.Render(buf)
	return buf
}
