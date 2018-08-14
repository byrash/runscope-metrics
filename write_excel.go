package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/tealeg/xlsx"
)

func WriteStatsToExcel(buckets *Buckets) {
	sheetName := "Run Scope Stats"
	summaryXslsFile := xlsx.NewFile()
	sheet, sheetErr := summaryXslsFile.AddSheet(sheetName)
	HandleError(sheetErr, "Unable to create a sheet in excel")
	//Add Header
	headerRow := sheet.AddRow()
	for _, header := range XLSXHeaders {
		cell := headerRow.AddCell()
		cell.Value = header
	}
	//For Every Bucket
	for _, bucket := range buckets.Data {
		for _, test := range bucket.Tests.Data {
			row := sheet.AddRow()
			projectNameCell := row.AddCell()
			projectNameCell.Value = bucket.ProjectName
			testNameCell := row.AddCell()
			testNameCell.Value = test.Name
			successRateCell := row.AddCell()
			successRateCell.Value = strconv.FormatFloat(test.TestMetrics.SuccessRate, 'f', 2, 64)
			avgRespTimeMsCell := row.AddCell()
			avgRespTimeMsCell.Value = strconv.FormatFloat(test.TestMetrics.AvgRespTimeMs, 'f', 2, 64)
			// log.Printf("%+v sucess rate and %+v avg response time for %+v test \n\n", test.TestMetrics.SuccessRate, test.TestMetrics.AvgRespTimeMs, test.Name)
		}
	}
	fileName := fmt.Sprintf("RunScopeStats_%v.xlsx", time.Now().Format("Jan_2_2006_at_3_04pm"))
	xlsxSaveErr := summaryXslsFile.Save(fileName)
	HandleError(xlsxSaveErr, "Unable to save Run scope summary sheet excel")
}
