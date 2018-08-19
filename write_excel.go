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
		if bucket.HasProductionData {
			addProjectName := true
			for _, test := range bucket.Tests.Data {
				row := sheet.AddRow()
				projectNameCell := row.AddCell()
				if addProjectName {
					projectNameCell.Value = bucket.ProjectName
					addProjectName = false
				}
				testNameCell := row.AddCell()
				testNameCell.Value = test.Name
				successRateCell := row.AddCell()
				successRateCell.Value = strconv.FormatFloat(test.TestMetrics.SuccessRate, 'f', 2, 64)
				avgRespTimeMsCell := row.AddCell()
				avgRespTimeMsCell.Value = strconv.FormatFloat(test.TestMetrics.AvgRespTimeMs, 'f', 2, 64)
				//Current Period
				fiftyThRespTimeMsCell := row.AddCell()
				fiftyThRespTimeMsCell.Value = strconv.FormatFloat(test.TestMetrics.ThisTimePeriod.RespTime50ThPercentile, 'f', 2, 64)
				nintyFifthThRespTimeMsCell := row.AddCell()
				nintyFifthThRespTimeMsCell.Value = strconv.FormatFloat(test.TestMetrics.ThisTimePeriod.RespTime95ThPercentile, 'f', 2, 64)
				nintyNinthThRespTimeMsCell := row.AddCell()
				nintyNinthThRespTimeMsCell.Value = strconv.FormatFloat(test.TestMetrics.ThisTimePeriod.RespTime99ThPercentile, 'f', 2, 64)
				// Last Period
				lastPeriodFiftyThRespTimeMsCell := row.AddCell()
				lastPeriodFiftyThRespTimeMsCell.Value = strconv.FormatFloat(test.TestMetrics.ChangesFromLastTimePeriod.RespTime50ThPercentile, 'f', 2, 64)
				lastPeriodNintyFifthThRespTimeMsCell := row.AddCell()
				lastPeriodNintyFifthThRespTimeMsCell.Value = strconv.FormatFloat(test.TestMetrics.ChangesFromLastTimePeriod.RespTime95ThPercentile, 'f', 2, 64)
				lastPeriodNintyNinthThRespTimeMsCell := row.AddCell()
				lastPeriodNintyNinthThRespTimeMsCell.Value = strconv.FormatFloat(test.TestMetrics.ChangesFromLastTimePeriod.RespTime99ThPercentile, 'f', 2, 64)
				// log.Printf("%+v sucess rate and %+v avg response time for %+v test \n\n", test.TestMetrics.SuccessRate, test.TestMetrics.AvgRespTimeMs, test.Name)
			}
		}
	}
	fileName := fmt.Sprintf("RunScopeStats_%v.xlsx", time.Now().Format("Jan_2_2006_at_3_04pm"))
	xlsxSaveErr := summaryXslsFile.Save(fileName)
	HandleError(xlsxSaveErr, "Unable to save Run scope summary sheet excel")
}
