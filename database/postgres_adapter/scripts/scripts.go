// Auto generated. DO NOT EDIT.

package scripts

// GetCreateMigrationTableScript gets the CreateMigrationTable script
func GetCreateMigrationTableScript() string {
	return `
CREATE TABLE IF NOT EXISTS public.migration (
    "timestamp" varchar(14) NOT NULL,
    CONSTRAINT migration_pk PRIMARY KEY ("timestamp")
);
`
}

// GetDeleteMigrationByTimestampScript gets the DeleteMigrationByTimestamp script
func GetDeleteMigrationByTimestampScript() string {
	return `
DELETE FROM migration
    WHERE "timestamp" = $1
`
}

// GetGetLatestTimestampScript gets the GetLatestTimestamp script
func GetGetLatestTimestampScript() string {
	return `
SELECT m."timestamp" FROM migration m
    ORDER BY m."timestamp" DESC
    LIMIT 1
`
}

// GetGetMigrationByTimestampScript gets the GetMigrationByTimestamp script
func GetGetMigrationByTimestampScript() string {
	return `
SELECT m."timestamp" 
    FROM migration m 
    WHERE m."timestamp" = $1
`
}

// GetSaveMigrationScript gets the SaveMigration script
func GetSaveMigrationScript() string {
	return `
INSERT INTO migration ("timestamp") 
    VALUES ($1)
`
}
