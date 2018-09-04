package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/tealeg/xlsx"
)

func WriteStatsToExcel(buckets *Buckets, restMVPData map[string]RestMVPData, criticalAppsData map[string]CriticalAppsData) {
	sheetName := "Run Scope Stats"
	summaryXslsFile := xlsx.NewFile()
	runScopeStatsSheet, sheetErr := summaryXslsFile.AddSheet(sheetName)
	HandleError(sheetErr, "Unable to create a sheet in excel")
	//Add Header
	addHeaderForMainSheet(runScopeStatsSheet)
	criticalAppsChann := make(chan bool)
	restChann := make(chan bool)
	//For Every Bucket
	for _, bucket := range buckets.Data {
		if bucket.HasProductionData {
			addProjectName := true
			for _, test := range bucket.Tests.Data {

				if restMvpRecord, ok := restMVPData[test.Id]; ok {
					restMvpRecord.TestData = test
					restMVPData[test.Id] = restMvpRecord
				}

				if criticalAppsRecord, ok := criticalAppsData[test.Id]; ok {
					criticalAppsRecord.TestData = test
					criticalAppsData[test.Id] = criticalAppsRecord
				}

				row := runScopeStatsSheet.AddRow()
				projectNameCell := row.AddCell()
				if addProjectName {
					projectNameCell.Value = bucket.ProjectName
					addProjectName = false
				}
				testNameCell := row.AddCell()
				testNameCell.Value = test.Name
				successRateCell := row.AddCell()
				successRateCell.Value = floatToString(test.TestMetrics.SuccessRate)
				avgRespTimeSecondsCell := row.AddCell()
				avgRespTimeSecondsCell.Value = floatToString(msToSeconds(test.TestMetrics.AvgRespTimeMs))
				//Current Period
				fiftyThRespTimeSecCell := row.AddCell()
				fiftyThRespTimeSecCell.Value = floatToString(msToSeconds(test.TestMetrics.ThisTimePeriod.RespTime50ThPercentile))
				nintyFifthThRespTimeSecCell := row.AddCell()
				nintyFifthThRespTimeSecCell.Value = floatToString(msToSeconds(test.TestMetrics.ThisTimePeriod.RespTime95ThPercentile))
				nintyNinthThRespTimeSecCell := row.AddCell()
				nintyNinthThRespTimeSecCell.Value = floatToString(msToSeconds(test.TestMetrics.ThisTimePeriod.RespTime99ThPercentile))
				// Last Period
				lastPeriodFiftyThRespTimeSecCell := row.AddCell()
				lastPeriodFiftyThRespTimeSecCell.Value = floatToString(msToSeconds(test.TestMetrics.ChangesFromLastTimePeriod.RespTime50ThPercentile))
				lastPeriodNintyFifthThRespTimeSecCell := row.AddCell()
				lastPeriodNintyFifthThRespTimeSecCell.Value = floatToString(msToSeconds(test.TestMetrics.ChangesFromLastTimePeriod.RespTime95ThPercentile))
				lastPeriodNintyNinthThRespTimeSecCell := row.AddCell()
				lastPeriodNintyNinthThRespTimeSecCell.Value = floatToString(msToSeconds(test.TestMetrics.ChangesFromLastTimePeriod.RespTime99ThPercentile))
				// log.Printf("%+v sucess rate and %+v avg response time for %+v test \n\n", test.TestMetrics.SuccessRate, test.TestMetrics.AvgRespTimeMs, test.Name)
			}
		}
	}

	restSheet, restSheetErr := summaryXslsFile.AddSheet("REST")
	HandleError(restSheetErr, "Unable to create sheet for Rest")
	go writeRestDataSheet(restSheet, restMVPData, restChann)

	criticalAppsSheet, criticalAppsSheetErr := summaryXslsFile.AddSheet("Critical Apps")
	HandleError(criticalAppsSheetErr, "Unable to create sheet for Critical Apps")
	go writeCriticalAppsDataSheet(criticalAppsSheet, criticalAppsData, criticalAppsChann)

	<-restChann
	<-criticalAppsChann

	fileName := fmt.Sprintf("RunScopeStats_%v.xlsx", time.Now().Format("Jan_2_2006_at_3_04pm"))
	xlsxSaveErr := summaryXslsFile.Save(fileName)
	HandleError(xlsxSaveErr, "Unable to save Run scope summary sheet excel")
}

func msToSeconds(ms float64) float64 {
	return float64(ms / 1000)
}

func writeRestDataSheet(restSheet *xlsx.Sheet, restMVPDataEntriesMap map[string]RestMVPData, chann chan bool) {
	addHeaderForRestSheet(restSheet)
	for i := 1; i <= len(restMVPDataEntriesMap); i++ {
		for _, restMvpDataEntry := range restMVPDataEntriesMap {
			if restMvpDataEntry.Sno == i {
				row := restSheet.AddRow()
				apiCell := row.AddCell()
				apiCell.Value = restMvpDataEntry.TestName
				if len(restMvpDataEntry.TestData.Id) != 0 {
					responseRateCell := row.AddCell()
					responseRateCell.Value = floatToString(msToSeconds(restMvpDataEntry.TestData.TestMetrics.AvgRespTimeMs))

					availabilityCell := row.AddCell()
					availabilityCell.Value = floatToString(restMvpDataEntry.TestData.TestMetrics.SuccessRate)

					stabilityCell := row.AddCell()
					if restMvpDataEntry.TestData.TestMetrics.AvgRespTimeMs == 0.0 {
						stabilityCell.Value = "0.0"
					} else {
						stabilityCell.Value = floatToString(float64(restMvpDataEntry.MVPVal) / float64(msToSeconds(restMvpDataEntry.TestData.TestMetrics.AvgRespTimeMs)) * 100)
					}
				} else {
					responseRateCell := row.AddCell()
					responseRateCell.Value = "This test is missing in Run Scope. Is Test id in youe metadata excel matching run scope test id ?"
				}
				break
			}
		}
	}
	chann <- true
}

func floatToString(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}

func writeCriticalAppsDataSheet(criticalAppsSheet *xlsx.Sheet, criticalAppsDataEntriesMap map[string]CriticalAppsData, chann chan bool) {
	addHeaderForCriticalAppsSheet(criticalAppsSheet)
	for i := 1; i <= len(criticalAppsDataEntriesMap); i++ {
		for _, criticalAppsDataEntry := range criticalAppsDataEntriesMap {
			if criticalAppsDataEntry.Sno == i {
				row := criticalAppsSheet.AddRow()
				appCell := row.AddCell()
				appCell.Value = criticalAppsDataEntry.TestName
				if len(criticalAppsDataEntry.TestData.Id) != 0 {
					monthlyAvgCell := row.AddCell()
					monthlyAvgCell.Value = floatToString(msToSeconds(criticalAppsDataEntry.TestData.TestMetrics.AvgRespTimeMs))

					availabilityCell := row.AddCell()
					availabilityCell.Value = floatToString(criticalAppsDataEntry.TestData.TestMetrics.SuccessRate)
				} else {
					monthlyAvgCell := row.AddCell()
					monthlyAvgCell.Value = "This test is missing in Run Scope. Is Test id in youe metadata excel matching run scope test id ?"
				}
				break
			}
		}
	}
	chann <- true
}

func addHeaderForMainSheet(sheet *xlsx.Sheet) {
	headerRow := sheet.AddRow()
	for _, header := range XLSXHeaders {
		cell := headerRow.AddCell()
		cell.Value = header
	}
}
func addHeaderForRestSheet(sheet *xlsx.Sheet) {
	headerRow := sheet.AddRow()
	for _, header := range RESTHeaders {
		cell := headerRow.AddCell()
		cell.Value = header
	}
}
func addHeaderForCriticalAppsSheet(sheet *xlsx.Sheet) {
	headerRow := sheet.AddRow()
	for _, header := range CriticalAppsHeaders {
		cell := headerRow.AddCell()
		cell.Value = header
	}
}
