package main

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/tealeg/xlsx"
)

func ReadRestMvpData(excelFileName *string, restMvpDataEntriesChan chan map[string]RestMVPData) {
	restMvpXslx, err := xlsx.OpenFile(*excelFileName)
	HandleError(err, "Unable to open rest mvp input excel file")
	sheet := restMvpXslx.Sheets[0] // Only conside first sheet
	restMVPDataEnties := make(map[string]RestMVPData)
	for rowId, row := range sheet.Rows {
		if rowId == 0 {
			//Ignoring Header
			continue
		}
		if len(row.Cells) != 4 {
			panic(errors.New("Required number of cells missing in input Rest MVP Xlsx"))
		}
		restMVPData := RestMVPData{}
		for i, cell := range row.Cells {
			text := cell.String()
			if len(text) == 0 {
				break
			}
			switch i {
			case 0:
				restMVPData.Sno, err = strconv.Atoi(text)
				HandleError(err, "Serial No values supplied in input excel is not valid")
			case 1:
				restMVPData.TestId = text
			case 2:
				restMVPData.TestName = text
			case 3:
				restMVPData.MVPVal, err = strconv.Atoi(text)
				HandleError(err, "Rest MVP values supplied in input excel is not valid")
			}
		}
		if len(restMVPData.TestId) != 0 {
			restMVPDataEnties[restMVPData.TestId] = restMVPData
		}
	}
	restMvpDataEntriesChan <- restMVPDataEnties
}

func ReadCriticalAppsData(criticalAppsFile *string, criticalAppsDataEntriesChan chan map[string]CriticalAppsData) {
	criticalAppsDataXlsx, err := xlsx.OpenFile(*criticalAppsFile)
	HandleError(err, "Unable to open rest mvp input excel file")
	sheet := criticalAppsDataXlsx.Sheets[0] // Only consider first sheet
	criticalAppsEntries := make(map[string]CriticalAppsData)
	for rowId, row := range sheet.Rows {
		if rowId == 0 {
			//Ignoring Header
			continue
		}
		if len(row.Cells) != 3 {
			panic(errors.New("Required number of cells missing in input Critical Apps Xlsx"))
		}
		criticalAppsData := CriticalAppsData{}
		for i, cell := range row.Cells {
			text := cell.String()
			if len(text) == 0 {
				break
			}
			switch i {
			case 0:
				criticalAppsData.Sno, err = strconv.Atoi(text)
				HandleError(err, "Serial No values supplied in input excel is not valid")
			case 1:
				criticalAppsData.TestId = text
			case 2:
				criticalAppsData.TestName = text
			}
		}
		if len(criticalAppsData.TestId) != 0 {
			criticalAppsEntries[criticalAppsData.TestId] = criticalAppsData
		}
	}
	criticalAppsDataEntriesChan <- criticalAppsEntries
}
