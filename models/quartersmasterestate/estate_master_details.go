// Package quartersmasterestate contains structs and queries for Estate Master Details API.
//
// Path: /var/www/html/go_projects/HRMODULE/Final_Mergecode/Meivan/models/quartersmasterestate
// --- Creator's Info ---
//
// Creator: Ramya M R
// Created On: 19-01-2026
//
// Last Modified By: Ramya M R
//
// Last Modified Date: 11-02-2026
package quartersmasterestate

import (
	"database/sql"
	"fmt"
)

// SQL Query for Estate Master Details
var MyQueryEstateMasterDetails = (`
SELECT
    qm.id AS quarters_id,
    qm.quartersnumber,
    qm.floor_id,
    fm.floor_name,
    qm.building_id,
    bm.building_name,
    bm.campus_id,
    c.campuscode AS campus_code,
    qm.displayname,
    qm.street,
    qm.plintharea,
    qm.licencefee,
    qm.swdcharges,
    qm.servicecharges,
    qm.ebcharges,
    qm.garagecharges,
    qm.cautiondeposit,
    qm.quartersstatus,
    qm.address,
    CASE 
        WHEN qm.is_servant_quarters = 1 THEN 'Yes'
        WHEN qm.is_servant_quarters = 0 THEN 'No'
        ELSE NULL
    END AS is_servant_quarters,
    qm.servant_quartersno,
    qm.effectivefrom
FROM humanresources.quartersmaster qm
JOIN humanresources.buildingmaster bm 
    ON bm.id = qm.building_id
LEFT JOIN humanresources.campus c 
    ON c.id = bm.campus_id
JOIN humanresources.quarterscategory qc 
    ON qc.id = bm.quarters_category
LEFT JOIN humanresources.floormaster fm 
    ON fm.id = qm.floor_id
WHERE qc.id = $1
  AND ( $2::int8[] IS NULL OR qm.building_id = ANY($2::int8[]) )
  AND ( $3::int8[] IS NULL OR qm.id = ANY($3::int8[]) )
ORDER BY qm.displayname;
`)

// Struct for Estate Master Details
type EstateMasterDetailsStruct struct {
	QuartersID        *int     `json:"Quarters_Id"`
	QuartersNumber    *string  `json:"Quarters_Number"`
	FloorID           *int     `json:"Floor_Id"`
	FloorName         *string  `json:"Floor_Name"`
	BuildingID        *int     `json:"Building_Id"`
	BuildingName      *string  `json:"Building_Name"`
	CampusID          *int     `json:"Campus_Id"`
	CampusCode        *string  `json:"Campus_Code"`
	DisplayName       *string  `json:"Display_Name"`
	Street            *string  `json:"Street"`
	PlinthArea        *float64 `json:"Plinth_Area"`
	LicenceFee        *float64 `json:"Licence_Fee"`
	SWDCharges        *float64 `json:"SWD_Charges"`
	ServiceCharges    *float64 `json:"Service_Charges"`
	EBCharges         *float64 `json:"EB_Charges"`
	GarageCharges     *float64 `json:"Garage_Charges"`
	CautionDeposit    *float64 `json:"Caution_Deposit"`
	QuartersStatus    *string  `json:"Quarters_Status"`
	Address           *string  `json:"Address"`
	IsServantQuarters *string  `json:"Is_Servant_Quarters"`
	ServantQuartersNo *string  `json:"Servant_Quarters_No"`
	EffectiveFrom     *string  `json:"Effective_From"`
}

// RetrieveEstateMasterDetails retrieves estate master details data
func RetrieveEstateMasterDetails(rows *sql.Rows) ([]EstateMasterDetailsStruct, error) {
	var list []EstateMasterDetailsStruct

	for rows.Next() {
		var s EstateMasterDetailsStruct

		err := rows.Scan(
			&s.QuartersID,
			&s.QuartersNumber,
			&s.FloorID,
			&s.FloorName,
			&s.BuildingID,
			&s.BuildingName,
			&s.CampusID,
			&s.CampusCode,
			&s.DisplayName,
			&s.Street,
			&s.PlinthArea,
			&s.LicenceFee,
			&s.SWDCharges,
			&s.ServiceCharges,
			&s.EBCharges,
			&s.GarageCharges,
			&s.CautionDeposit,
			&s.QuartersStatus,
			&s.Address,
			&s.IsServantQuarters,
			&s.ServantQuartersNo,
			&s.EffectiveFrom,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning estate master details: %v", err)
		}

		list = append(list, s)
	}

	return list, nil
}
