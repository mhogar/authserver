// Auto generated. DO NOT EDIT.

package scripts

// GetCreateMigrationTableScript gets the CreateMigrationTable script
func GetCreateMigrationTableScript() string {
	return `
CREATE TABLE IF NOT EXISTS public."migration" (
    "timestamp" varchar(14) NOT NULL,
    CONSTRAINT migration_pk PRIMARY KEY ("timestamp")
);
`
}

// GetCreateUserTableScript gets the CreateUserTable script
func GetCreateUserTableScript() string {
	return `
CREATE TABLE  IF NOT EXISTS public."user" (
	id uuid NOT NULL,
	username varchar(30) NOT NULL,
	password_hash bytea NOT NULL,
	CONSTRAINT user_pk PRIMARY KEY (id),
	CONSTRAINT user_username_un UNIQUE (username)
);
`
}

// GetDeleteMigrationByTimestampScript gets the DeleteMigrationByTimestamp script
func GetDeleteMigrationByTimestampScript() string {
	return `
DELETE FROM "migration"
   WHERE "timestamp" = $1
`
}

// GetGetLatestTimestampScript gets the GetLatestTimestamp script
func GetGetLatestTimestampScript() string {
	return `
SELECT m."timestamp" FROM "migration" m
    ORDER BY m."timestamp" DESC
    LIMIT 1
`
}

// GetGetMigrationByTimestampScript gets the GetMigrationByTimestamp script
func GetGetMigrationByTimestampScript() string {
	return `
SELECT m."timestamp" 
    FROM "migration" m 
    WHERE m."timestamp" = $1
`
}

// GetGetUserByIDScript gets the GetUserByID script
func GetGetUserByIDScript() string {
	return `
SELECT u."id", u."username", u."password_hash"
	FROM "user" u
	WHERE u."id" = $1
`
}

// GetSaveMigrationScript gets the SaveMigration script
func GetSaveMigrationScript() string {
	return `
INSERT INTO "migration" ("timestamp") 
    VALUES ($1)
`
}

// GetSaveUserScript gets the SaveUser script
func GetSaveUserScript() string {
	return `
INSERT INTO "user" ("id", "username", "password_hash")
	VALUES ($1, $2, $3)
`
}
