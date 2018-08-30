package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

func WriteStatsToExcel(buckets *Buckets, restMVPFile, criticalAppsFile *string) {
	sheetName := "Run Scope Stats"
	summaryXslsFile := xlsx.NewFile()
	runScopeStatsSheet, sheetErr := summaryXslsFile.AddSheet(sheetName)
	HandleError(sheetErr, "Unable to create a sheet in excel")
	//Add Header
	addHeaderForMainSheet(runScopeStatsSheet)
	criticalAppsBucketChann := make(chan bool, 1)
	criticalAppsBucketExists := false
	restBucketChann := make(chan bool, 1)
	restBucketExists := false
	//For Every Bucket
	for _, bucket := range buckets.Data {
		if bucket.HasProductionData {
			if strings.EqualFold(RestBucketId, bucket.BucketKey) {
				restBucketExists = true
				restSheet, restSheetErr := summaryXslsFile.AddSheet("REST")
				HandleError(restSheetErr, "Unable to create sheet for Rest")
				go writeRestDataSheet(restSheet, bucket, restBucketChann, restMVPFile)
			}
			if strings.EqualFold(CriticalAppsBucketId, bucket.BucketKey) {
				criticalAppsBucketExists = true
				criticalAppsSheet, criticalAppsSheetErr := summaryXslsFile.AddSheet("Critical Apps")
				HandleError(criticalAppsSheetErr, "Unable to create sheet for Critical Apps")
				go writeCriticalAppsDataSheet(criticalAppsSheet, bucket, criticalAppsBucketChann, criticalAppsFile)
			}
			addProjectName := true
			for _, test := range bucket.Tests.Data {
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
	if !restBucketExists {
		restBucketChann <- true
	}
	if !criticalAppsBucketExists {
		criticalAppsBucketChann <- true
	}
	select {
	case <-restBucketChann:
	case <-criticalAppsBucketChann:
	}
	fileName := fmt.Sprintf("RunScopeStats_%v.xlsx", time.Now().Format("Jan_2_2006_at_3_04pm"))
	xlsxSaveErr := summaryXslsFile.Save(fileName)
	HandleError(xlsxSaveErr, "Unable to save Run scope summary sheet excel")
}

func msToSeconds(ms float64) float64 {
	return float64(ms / 1000)
}

func writeRestDataSheet(restSheet *xlsx.Sheet, restBucket Bucket, chann chan bool, restMVPFile *string) {
	restMvpDataEntries := ReadRestMvpData(restMVPFile)
	testEntries := make(map[string]Test)
	for _, test := range restBucket.Tests.Data {
		testEntries[test.Id] = test
	}
	addHeaderForRestSheet(restSheet)
	for _, restMvpData := range restMvpDataEntries {
		row := restSheet.AddRow()
		apiCell := row.AddCell()
		apiCell.Value = restMvpData.TestName
		if test, ok := testEntries[restMvpData.TestId]; ok {

			responseRateCell := row.AddCell()
			responseRateCell.Value = floatToString(msToSeconds(test.TestMetrics.AvgRespTimeMs))

			availabilityCell := row.AddCell()
			availabilityCell.Value = floatToString(test.TestMetrics.SuccessRate)

			stabilityCell := row.AddCell()
			if test.TestMetrics.AvgRespTimeMs == 0.0 {
				stabilityCell.Value = "0.0"
			} else {
				stabilityCell.Value = floatToString(float64(restMvpData.MVPVal) / float64(msToSeconds(test.TestMetrics.AvgRespTimeMs)) * 100)
			}

		} else {
			responseRateCell := row.AddCell()
			responseRateCell.Value = "This test is missing in Run Scope"
		}
	}
	chann <- true
}

func floatToString(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}

func writeCriticalAppsDataSheet(criticalAppsSheet *xlsx.Sheet, criticalAppsBucket Bucket, chann chan bool, criticalAppsFile *string) {
	criticalAppsDataEntries := ReadCriticalAppsData(criticalAppsFile)
	testEntries := make(map[string]Test)
	for _, test := range criticalAppsBucket.Tests.Data {
		testEntries[test.Id] = test
	}
	addHeaderForCriticalAppsSheet(criticalAppsSheet)
	for _, criticalAppsData := range criticalAppsDataEntries {
		row := criticalAppsSheet.AddRow()

		appCell := row.AddCell()
		appCell.Value = criticalAppsData.TestName

		if test, ok := testEntries[criticalAppsData.TestId]; ok {
			monthlyAvgCell := row.AddCell()
			monthlyAvgCell.Value = floatToString(msToSeconds(test.TestMetrics.AvgRespTimeMs))

			availabilityCell := row.AddCell()
			availabilityCell.Value = floatToString(test.TestMetrics.SuccessRate)
		} else {
			monthlyAvgCell := row.AddCell()
			monthlyAvgCell.Value = "This test is missing in Run Scope"
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
