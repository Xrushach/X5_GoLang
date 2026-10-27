package analytics

import "honnef.co/go/tools/analysis/report"

type Sale struct {
	Product          string
	Category         string
	Quantity         int
	UnitPriceKopecks int64
	Cancelled        bool
}

type Report struct {
	RevenueKopecks    int64
	UnitsSold         int
	RevenueByCategory map[string]int64
	TopProduct        string
	InvalidSales      int
}

func AnalyzeSales(sales []Sale) Report {
	// TODO: реализуйте расчёт отчёта по правилам задания.
	rep := Report{RevenueByCategory: make(map[string]int64)}
	topQ := 0
	for _, sale := range sales {
		if sale.Cancelled {
			continue
		}
		if sale.Product == "" || sale.Quantity <= 0 || sale.UnitPriceKopecks <= 0 || sale.Category == "" {
			rep.InvalidSales++
			continue
		}
		morzha := int64(sale.Quantity) * sale.UnitPriceKopecks

		rep.RevenueKopecks += morzha
		rep.UnitsSold += sale.Quantity
		rep.RevenueByCategory[sale.Category] += morzha

		if sale.Quantity > topQ || (sale.Quantity == topQ && sale.Product < report.TopProduct) {
			topQ = sale.Quantity
			report.TopProduct = sale.Product
		}
	}
	return rep
}
