package api

import (
	"fmt"
	"log"
	"net/http"
	"shivamaravanthe/HosangadiReports/database"
	"strconv"
	"strings"
)

func Gst(w http.ResponseWriter, r *http.Request) {
	type data struct {
		CPstr    string `gorm:"column:cost_price"`
		Qtystr   string `gorm:"column:sales_prod_qty"`
		SPstr    string `gorm:"column:sales_prod_sp"`
		GSTstr   string `gorm:"column:gst_value"`
		Salesref string `gorm:"column:sales_ref"`
	}
	SEMsales := 0.0
	gstSSMMap := map[int]float64{0: 0.0, 5: 0.0, 12: 0.0, 18: 0.0, 28: 0.0}
	gstConsantMap := map[int]float64{0: 0.0, 5: 0.4761905, 12: 0.10714286, 18: 0.15254237, 28: 0.21875}
	// SELECT cost_price,sales_prod_qty,sales_prod_sp,gst_value FROM somanath2023.sales_sp sp
	// INNER JOIN somanath2023.sales ON sp.sales_ref = sales.sales_ref
	// WHERE sp.sales_ref LIKE 'SEM_%'
	sqlData := []data{}
	if err := database.DB.Table("sales_sp sp").
		Select("sp.sales_ref as sales_ref,cost_price,sales_prod_qty,sales_prod_sp,gst_value").
		Where("sp.sales_ref LIKE 'SEM_%'").Joins("INNER JOIN sales ON sp.sales_ref = sales.sales_ref").
		Scan(&sqlData).Error; err != nil {
		log.Printf("Failed to fetch from database %v", err)
		return
	}
	gst := 0
	cp := 0.0
	sp := 0.0
	qty := 0.0
	intStringMap := map[string]int{}
	floatStringMap := map[string]float64{}
	for _, each := range sqlData {
		cpArray := strings.Split(each.CPstr, ":")
		spArray := strings.Split(each.SPstr, ":")
		qtyArray := strings.Split(each.Qtystr, ":")
		gstArray := strings.Split(each.GSTstr, ":")

		gstArray = gstArray[1 : len(gstArray)-1]
		qtyArray = qtyArray[1 : len(qtyArray)-1]
		spArray = spArray[1 : len(spArray)-1]
		cpArray = cpArray[1 : len(cpArray)-1]
		var err error
		for i, item := range cpArray {
			if _, ok := intStringMap[gstArray[i]]; !ok {
				gst, err = strconv.Atoi(gstArray[i])
				if err != nil {
					gst, _ = strconv.Atoi(strings.Split(gstArray[i], ".")[0])
				}
				intStringMap[gstArray[i]] = gst
			} else {
				gst = intStringMap[gstArray[i]]
			}
			if _, ok := floatStringMap[item]; !ok {
				cp, _ = strconv.ParseFloat(item, 64)
				floatStringMap[item] = cp
			} else {
				cp = floatStringMap[item]
			}
			if _, ok := floatStringMap[qtyArray[i]]; !ok {
				qty, _ = strconv.ParseFloat(qtyArray[i], 64)
				floatStringMap[qtyArray[i]] = qty
			} else {
				qty = floatStringMap[qtyArray[i]]
			}
			if _, ok := floatStringMap[spArray[i]]; !ok {
				sp, _ = strconv.ParseFloat(spArray[i], 64)
				floatStringMap[spArray[i]] = sp
			} else {
				sp = floatStringMap[spArray[i]]
			}
			temp := (sp - cp) * qty
			gstSSMMap[gst] += temp
			SEMsales += sp * qty
		}
	}

	paidGSTMap := map[int]float64{0: 0.0, 5: 0.0, 12: 0.0, 18: 0.0, 28: 0.0}
	totalPaid := 0.0

	for i, value := range gstSSMMap {
		totalPaid += value * gstConsantMap[i]
		paidGSTMap[i] = value * gstConsantMap[i]
	}
	fmt.Print("GST | Profit\n")
	fmt.Printf("0%s  %.2f\n", "%", gstSSMMap[0])
	fmt.Printf("5%s  %.2f\n", "%", gstSSMMap[5])
	fmt.Printf("12%s %.2f\n", "%", gstSSMMap[12])
	fmt.Printf("18%s %.2f\n", "%", gstSSMMap[18])
	fmt.Printf("28%s %.2f\n\n", "%", gstSSMMap[28])
	fmt.Printf("Total Profit %.2f\n\n", gstSSMMap[28]+gstSSMMap[18]+gstSSMMap[12]+gstSSMMap[5]+gstSSMMap[0])
	fmt.Print("GST SSM would have paid\n")
	fmt.Printf("0%s %.2f\n", "%", paidGSTMap[0])
	fmt.Printf("5%s %.2f\n", "%", paidGSTMap[5])
	fmt.Printf("12%s %.2f\n", "%", paidGSTMap[12])
	fmt.Printf("18%s %.2f\n", "%", paidGSTMap[18])
	fmt.Printf("28%s %.2f\n", "%", paidGSTMap[28])
	fmt.Printf("Total GST SSM PAID %.2f\n", totalPaid)
	fmt.Printf("Total GST SEM PAID %#.2f\n\n", SEMsales*0.01)
	fmt.Printf("Total SAVED %#.2f\n", totalPaid-(SEMsales*0.01))
}
