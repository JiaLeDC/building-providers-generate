package main

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/lucsky/cuid"
	"github.com/xuri/excelize/v2"
)

func processSheet(f *excelize.File, sheetName string, priceProfileID string, buildingWriter *csv.Writer, providerWriter *csv.Writer) error {

	// Read all rows from the sheet
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to get rows from sheet %s: %v", sheetName, err)
	}

	// Skip header row
	if len(rows) < 2 {
		return fmt.Errorf("sheet %s has no data rows", sheetName)
	}

	// Process each row
	for _, row := range rows[1:] {
		// Ensure row has enough columns
		if len(row) < 4 {
			continue
		}

		rawName := row[0]
		rawAddress := row[1]
		rawLatitudeStr := row[2]
		rawLongitudeStr := row[3]

		// Generate unique building ID
		buildingID := cuid.New()

		makeTime := time.Now()
		eightHoursAgo := makeTime.Add(8 * time.Hour)
		currentTime := eightHoursAgo.Format("2006-01-02 15:04:05")

		// Write building row
		if err := buildingWriter.Write([]string{
			buildingID, currentTime, currentTime, rawName, rawLatitudeStr, rawLongitudeStr, rawAddress, "NULL", "f",
		}); err != nil {
			return fmt.Errorf("failed to write building row: %v", err)
		}

		// Write provider DIA/BroadBand row
		if err := providerWriter.Write([]string{
			cuid.New(), currentTime, currentTime, "PhibeeTelecom", "f", "{GigabitEthernet}", "{30}", buildingID, priceProfileID,
		}); err != nil {
			return fmt.Errorf("failed to write provider row: %v", err)
		}
	}

	return nil
}

func main() {
	// Price Profile IDs
	priceProfileIDs := []string{
		"cm6ytex7r0000c4ka7rpd50a2",
	}

	city := "France"

	// List of Excel files to process
	excelFiles := []string{
		fmt.Sprintf("Buildings_Broadband.xlsx"),
		// Add more templates here and replace $city dynamically
	}

	// Prepare consolidated CSV files
	fileName_Building := fmt.Sprintf("%s_Consolidated_Buildings.csv", city)
	consolidatedBuilding, err := os.Create(fileName_Building)
	if err != nil {
		slog.Error("Failed to create consolidated buildings CSV", err)
		return
	}
	defer consolidatedBuilding.Close()
	buildingWriter := csv.NewWriter(consolidatedBuilding)
	defer buildingWriter.Flush()

	// Write building CSV header
	if err := buildingWriter.Write([]string{"id", "createTime", "updateTime", "name", "latitude", "longitude", "address", "googleMapsPlaceId", "systemManaged"}); err != nil {
		slog.Error("Failed to write building CSV header", err)
		return
	}

	fileName_Provider := fmt.Sprintf("%s_Consolidated_Providers.csv", city)
	consolidatedProvider, err := os.Create(fileName_Provider)
	if err != nil {
		slog.Error("Failed to create consolidated providers CSV", err)
		return
	}
	defer consolidatedProvider.Close()
	providerWriter := csv.NewWriter(consolidatedProvider)
	defer providerWriter.Flush()

	// Write provider CSV header
	if err := providerWriter.Write([]string{"id", "createTime", "updateTime", "provider", "systemManaged", "interfaceTypes", "ipv4PrefixLengths", "buildingId", "priceProfileId"}); err != nil {
		slog.Error("Failed to write provider CSV header", err)
		return
	}

	// Process each Excel file
	for _, excelFilePath := range excelFiles {
		slog.Info(fmt.Sprintf("Opening Excel file: %s", excelFilePath))

		// Open the Excel file
		f, err := excelize.OpenFile(excelFilePath)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to open Excel file %s", excelFilePath), err)
			continue
		}
		defer f.Close()

		// Get all sheet names
		sheetNames := f.GetSheetList()
		slog.Info(fmt.Sprintf("Sheets found in %s: %v", excelFilePath, sheetNames))

		// Process each sheet with a corresponding price profile ID
		for i, sheetName := range sheetNames {
			slog.Info(fmt.Sprintf("Processing sheet: %s", sheetName))

			// Check if we have enough price profile IDs
			if i >= len(priceProfileIDs) {
				slog.Warn(fmt.Sprintf("No more price profile IDs for sheet %s in file %s", sheetName, excelFilePath))
				break
			}

			if err := processSheet(f, sheetName, priceProfileIDs[i], buildingWriter, providerWriter); err != nil {
				slog.Error(fmt.Sprintf("Error processing sheet %s from %s", sheetName, excelFilePath), err)
				continue
			}
		}
	}

	slog.Info("Consolidated CSV files generated successfully.")
}
