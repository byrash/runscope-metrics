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
				avgRespTimeSecondsCell := row.AddCell()
				avgRespTimeSecondsCell.Value = strconv.FormatFloat(msToSeconds(test.TestMetrics.AvgRespTimeMs), 'f', 2, 64)
				//Current Period
				fiftyThRespTimeSecCell := row.AddCell()
				fiftyThRespTimeSecCell.Value = strconv.FormatFloat(msToSeconds(test.TestMetrics.ThisTimePeriod.RespTime50ThPercentile), 'f', 2, 64)
				nintyFifthThRespTimeSecCell := row.AddCell()
				nintyFifthThRespTimeSecCell.Value = strconv.FormatFloat(msToSeconds(test.TestMetrics.ThisTimePeriod.RespTime95ThPercentile), 'f', 2, 64)
				nintyNinthThRespTimeSecCell := row.AddCell()
				nintyNinthThRespTimeSecCell.Value = strconv.FormatFloat(msToSeconds(test.TestMetrics.ThisTimePeriod.RespTime99ThPercentile), 'f', 2, 64)
				// Last Period
				lastPeriodFiftyThRespTimeSecCell := row.AddCell()
				lastPeriodFiftyThRespTimeSecCell.Value = strconv.FormatFloat(msToSeconds(test.TestMetrics.ChangesFromLastTimePeriod.RespTime50ThPercentile), 'f', 2, 64)
				lastPeriodNintyFifthThRespTimeSecCell := row.AddCell()
				lastPeriodNintyFifthThRespTimeSecCell.Value = strconv.FormatFloat(msToSeconds(test.TestMetrics.ChangesFromLastTimePeriod.RespTime95ThPercentile), 'f', 2, 64)
				lastPeriodNintyNinthThRespTimeSecCell := row.AddCell()
				lastPeriodNintyNinthThRespTimeSecCell.Value = strconv.FormatFloat(msToSeconds(test.TestMetrics.ChangesFromLastTimePeriod.RespTime99ThPercentile), 'f', 2, 64)
				// log.Printf("%+v sucess rate and %+v avg response time for %+v test \n\n", test.TestMetrics.SuccessRate, test.TestMetrics.AvgRespTimeMs, test.Name)
			}
		}
	}
	fileName := fmt.Sprintf("RunScopeStats_%v.xlsx", time.Now().Format("Jan_2_2006_at_3_04pm"))
	xlsxSaveErr := summaryXslsFile.Save(fileName)
	HandleError(xlsxSaveErr, "Unable to save Run scope summary sheet excel")
}

func msToSeconds(ms float64) float64 {
	return float64(ms / 1000)
}
