// Package databasequartersmasterestate handles DB access for Estate Master Details API.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/database/quartersmasterestate
// --- Creator's Info ---
//
// Creator: Ramya M R
//
// Created On: 19-01-2026
//
// Last Modified By:
//
// Last Modified Date:

package databasequartersmasterestate

import (
	credentials "Hrmodule/dbconfig"
	modelsquartersmasterestate "Hrmodule/models/quartersmasterestate"
	"fmt"
	"strings"

	"github.com/lib/pq" // Required for handling Go slices in Postgres
)

// Helper function to parse IDs from various possible input types (string, float64, or slice)
func parseIDList(input interface{}) []int64 {
	var ids []int64

	if input == nil {
		return nil
	}

	switch val := input.(type) {
	case float64:
		// Case: Single number from JSON
		ids = append(ids, int64(val))

	case string:
		// Case: Comma-separated string like "00, 02"
		strParts := strings.Split(val, ",")
		for _, s := range strParts {
			s = strings.TrimSpace(s)
			if s == "" {
				continue
			}
			var id int64
			// Use Sscanf to handle the numeric conversion
			if _, err := fmt.Sscanf(s, "%d", &id); err == nil {
				ids = append(ids, id)
			}
		}

	case []interface{}:
		// Case: JSON array like [1, 2, 3]
		for _, item := range val {
			if num, ok := item.(float64); ok {
				ids = append(ids, int64(num))
			} else if s, ok := item.(string); ok {
				var id int64
				if _, err := fmt.Sscanf(strings.TrimSpace(s), "%d", &id); err == nil {
					ids = append(ids, id)
				}
			}
		}
	}

	if len(ids) == 0 {
		return nil
	}
	return ids
}

// GetEstateMasterDetailsFromDB fetches estate master details data from DB
func GetEstateMasterDetailsFromDB(
	decryptedData map[string]interface{},
) ([]modelsquartersmasterestate.EstateMasterDetailsStruct, int, error) {

	// 1. Establish Database Connection
	db := credentials.GetDB()

	// 2. Extract & validate category_id (REQUIRED)
	categoryVal, ok := decryptedData["category_id"]
	if !ok || categoryVal == nil {
		return nil, 0, fmt.Errorf("category_id is required")
	}

	var categoryID int64
	// Handle category_id being sent as a string or number
	if cv, ok := categoryVal.(float64); ok {
		categoryID = int64(cv)
	} else if cv, ok := categoryVal.(string); ok {
		fmt.Sscanf(cv, "%d", &categoryID)
	}

	// 3. Robustly parse building_ids and quarters_ids
	// This handles strings "00, 02", arrays [0, 2], or single numbers
	buildingIDs := parseIDList(decryptedData["building_ids"])
	quartersIDs := parseIDList(decryptedData["quarters_ids"])

	// 4. Prepare parameters for SQL ANY() using pq.Array
	var bParam, qParam interface{}

	if len(buildingIDs) > 0 {
		bParam = pq.Array(buildingIDs)
	} else {
		bParam = nil // SQL will handle this via: $2 IS NULL
	}

	if len(quartersIDs) > 0 {
		qParam = pq.Array(quartersIDs)
	} else {
		qParam = nil // SQL will handle this via: $3 IS NULL
	}

	// 5. Execute query
	rows, err := db.Query(
		modelsquartersmasterestate.MyQueryEstateMasterDetails,
		categoryID, // $1
		bParam,     // $2
		qParam,     // $3
	)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	// 6. Fetch results
	data, err := modelsquartersmasterestate.RetrieveEstateMasterDetails(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("retrieving result failed: %v", err)
	}

	return data, len(data), nil
}
