SELECT u."id", u."username", u."password_hash"
	FROM "user" u
	WHERE u."id" = $1