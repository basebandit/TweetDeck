package schema

import (
	"github.com/dimiro1/darwin"
	"github.com/jmoiron/sqlx"
)
//Migrate runs queries that update the db schema up to date with the migrations
//defined in this package
func Migrate(db *sqlx.DB)error{
 driver := darwin.NewGenericDriver(db.DB,darwin.PostgresDialect{})
 d := darwin.New(driver,migrations,nil)
 return d.Migrate()
}

var migrations = []darwin.Migration{
	{
		Version: 1,
		Description:"Add users",
		Script:`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			firstname TEXT,
			lastname TEXT,
			email TEXT UNIQUE,
			password_hash TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		);`,
	},
	{
		Version: 2,
		Description:"Add avatars",
		Script:`
		CREATE TABLE IF NOT EXISTS avatars (
			id UUID PRIMARY KEY,
			username TEXT,
			user_id UUID,
			active BOOLEAN,
		 created_at TIMESTAMP,
		 updated_at TIMESTAMP,
		
		 FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT
		);`,
	},
	{
		Version:3,
		Description: "Add profiles",
		Script:`
		CREATE TABLE IF NOT EXISTS profiles (
			id UUID PRIMARY KEY,
			avatar_id UUID,
		 followers INT,
			"following" INT,
		 tweets INT,
		 likes INT,
		 join_date TEXT,
		 profile_image_url TEXT,
		 bio TEXT,
		 "name" TEXT,
		 twitter_id TEXT,
		 last_tweet_time TEXT,
		 created_at TIMESTAMP,
		 updated_at TIMESTAMP,
		
		 FOREIGN KEY(avatar_id) REFERENCES avatars(id) ON DELETE RESTRICT
		);`,
	},
}