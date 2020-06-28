// Auto generated. DO NOT EDIT.

package scripts

// StaticScriptRepository is an implementation of the sql script repository interface that fetches statically defined scripts.
type StaticScriptRepository struct {}

// CreateMigrationTableScript gets the CreateMigrationTable script
func (StaticScriptRepository) CreateMigrationTableScript() string {
	return `
CREATE TABLE IF NOT EXISTS public."migration" (
    "timestamp" varchar(14) NOT NULL,
    CONSTRAINT migration_pk PRIMARY KEY ("timestamp")
);
`
}

// CreateUserTableScript gets the CreateUserTable script
func (StaticScriptRepository) CreateUserTableScript() string {
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

// DeleteMigrationByTimestampScript gets the DeleteMigrationByTimestamp script
func (StaticScriptRepository) DeleteMigrationByTimestampScript() string {
	return `
DELETE FROM "migration"
   WHERE "timestamp" = $1
`
}

// DeleteUserScript gets the DeleteUser script
func (StaticScriptRepository) DeleteUserScript() string {
	return `
DELETE FROM "user" u
    WHERE u."id" = $1
`
}

// GetLatestTimestampScript gets the GetLatestTimestamp script
func (StaticScriptRepository) GetLatestTimestampScript() string {
	return `
SELECT m."timestamp" FROM "migration" m
    ORDER BY m."timestamp" DESC
    LIMIT 1
`
}

// GetMigrationByTimestampScript gets the GetMigrationByTimestamp script
func (StaticScriptRepository) GetMigrationByTimestampScript() string {
	return `
SELECT m."timestamp" 
    FROM "migration" m 
    WHERE m."timestamp" = $1
`
}

// GetUserByIDScript gets the GetUserByID script
func (StaticScriptRepository) GetUserByIDScript() string {
	return `
SELECT u."id", u."username", u."password_hash"
	FROM "user" u
	WHERE u."id" = $1
`
}

// GetUserByUsernameScript gets the GetUserByUsername script
func (StaticScriptRepository) GetUserByUsernameScript() string {
	return `
SELECT u."id", u."username", u."password_hash"
	FROM "user" u
	WHERE u."username" = $1
`
}

// SaveMigrationScript gets the SaveMigration script
func (StaticScriptRepository) SaveMigrationScript() string {
	return `
INSERT INTO "migration" ("timestamp") 
    VALUES ($1)
`
}

// SaveUserScript gets the SaveUser script
func (StaticScriptRepository) SaveUserScript() string {
	return `
INSERT INTO "user" ("id", "username", "password_hash")
	VALUES ($1, $2, $3)
`
}

// UpdateUserScript gets the UpdateUser script
func (StaticScriptRepository) UpdateUserScript() string {
	return `
UPDATE "user" SET
    "username" = $2,
    "password_hash" = $3
WHERE "id" = $1
`
}
