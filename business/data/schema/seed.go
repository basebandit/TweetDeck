package schema

import "github.com/jmoiron/sqlx"

//Seed runs the queries that populate the database with initial data.The queries are ran in a
//if any fail transaction and rolled back
func Seed(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(seeds); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}

const seeds = `
 -- Create a test user with password "gophers"
 INSERT INTO users (id,firstname,lastname,email,password_hash,active,created_at,updated_at) VALUES
 ('5cf37266-3473-4006-984f-9325122678b7','test','user1','testuser1@gmail.com','$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a',TRUE,'2020-09-04 00:00:00','2020-09-04 00:00:00'),
 ('45b5fbd3-755f-4379-8f07-a58d4a30fa2f','test','user2','testuser2@gmail.com','$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW',TRUE,'2020-09-04 00:00:00','2020-09-04 00:00:00') ON CONFLICT DO NOTHING;
 
 INSERT INTO avatars (id,username,user_id,created_at,updated_at) VALUES 
 ('98b6d4b8-f04b-4c79-8c2e-a0aef46854b7','DKJnr3','45b5fbd3-755f-4379-8f07-a58d4a30fa2f','2020-09-04 02:00:00','2020-09-04 02:30:00'),
 ('85f6fb09-eb05-4874-ae39-82d1a30fe0d7','FelistusQ','45b5fbd3-755f-4379-8f07-a58d4a30fa2f','2020-09-04 02:00:00','2020-09-04 03:00:00'),
 ('a235be9e-ab5d-44e6-a987-fa1c749264c7','jean_wangari','5cf37266-3473-4006-984f-9325122678b7','2020-09-04 03:00:00','2020-09-04 04:30:00') 
 ON CONFLICT DO NOTHING;

 -- Create a test avatar with the default(empty) user_id field
 INSERT INTO avatars(id,username,created_at,updated_at) VALUES
 ('6ba7b810-9dad-11d1-80b4-00c04fd430c8','TPS_Ke','2020-09-04 08:00:00','2020-09-04 09:00:00')
 ON CONFLICT DO NOTHING;


 INSERT INTO profiles (id,avatar_id,bio,profile_image_url,"name",twitter_id,join_date,followers,"following",tweets,likes,last_tweet_time,created_at) VALUES 
 ('a2b0639f-2cc6-44b8-b97b-15d69dbb511e','98b6d4b8-f04b-4c79-8c2e-a0aef46854b7','Gíkúyú ní wendo','https://pbs.twimg.com/profile_images/1293103745648201729/CMwG39AN_normal.jpg','DK Jnr.','1267766034179776512','Tue Jun 02 10:32:52 +0000',238,626,486,87,'2020-09-04 04:10:10','2020-09-04 00:00:00'),
 ('72f8b983-3eb4-48db-9ed0-e45cc6bd716b','85f6fb09-eb05-4874-ae39-82d1a30fe0d7','A free spirit. I stand for justice. Proud to be black.','https://pbs.twimg.com/profile_images/1277896115342528512/uNVpTeIW_normal.jpg','Felistus Waithira','1268059887906557953','Wed Jun 03 06:00:31 +0000 2020',1025,970,647,738,'2020-09-04 02:20:10','2020-09-04 00:00:00'),
 ('6ba7b814-9dad-11d1-80b4-00c04fd430c8','a235be9e-ab5d-44e6-a987-fa1c749264c7','Cuppycake Living large','https://pbs.twimg.com/profile_images/1288401307204804608/0s5DK5ej_normal.jpg','Jean Wangari','1267757177999101953','Tue Jun 02 09:59:17 +0000 2020',1291,1344,1130,3070,'2020-08-23 15:34:42','2020-09-04 05:10:10')
 ON CONFLICT DO NOTHING;

 -- Create other profiles of already existing avatars.
 INSERT INTO profiles (id,avatar_id,bio,profile_image_url,"name",twitter_id,join_date,followers,"following",tweets,likes,last_tweet_time,created_at) VALUES 
 ('c918083b-ecba-4af5-a19b-b3231a7bd09e','98b6d4b8-f04b-4c79-8c2e-a0aef46854b7','Gíkúyú ní wendo','https://pbs.twimg.com/profile_images/1293103745648201729/CMwG39AN_normal.jpg','DK Jnr.','1267766034179776512','Tue Jun 02 10:32:52 +0000',308,826,686,207,'2020-09-07 02:15:10','2020-09-07 02:15:00'),

 ('d4c50057-992e-4396-a8ac-fefbbeee208c','98b6d4b8-f04b-4c79-8c2e-a0aef46854b7','Gíkúyú ní wendo','https://pbs.twimg.com/profile_images/1293103745648201729/CMwG39AN_normal.jpg','DK Jnr.','1267766034179776512','Tue Jun 02 10:32:52 +0000',498,996,1006,327,'2020-09-10 03:10:10','2020-09-10 03:15:00'),

 ('ff6b6efd-87c3-435e-a493-167f32bbc98b','85f6fb09-eb05-4874-ae39-82d1a30fe0d7','A free spirit. I stand for justice. Proud to be black.','https://pbs.twimg.com/profile_images/1277896115342528512/uNVpTeIW_normal.jpg','Felistus Waithira','1268059887906557953','Wed Jun 03 06:00:31 +0000 2020',1026,976,648,740,'2020-09-08 09:20:10','2020-09-08 10:00:00'),

 ('9f4c8bd4-30ba-42bd-9c6c-adb9fc83d3bd','a235be9e-ab5d-44e6-a987-fa1c749264c7','Cuppycake Living large','https://pbs.twimg.com/profile_images/1288401307204804608/0s5DK5ej_normal.jpg','Jean Wangari','1267757177999101953','Tue Jun 02 09:59:17 +0000 2020',1291,1344,1130,3070,'2020-08-23 15:34:42','2020-09-04 05:10:10')
 ON CONFLICT DO NOTHING;
`

//DeleteAll runs the drop table queries. The queries are ran in a
//transaction and rolled back if any fail.
func DeleteAll(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(deleteAll); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}

//queries to clean up the database between tests.
const deleteAll = `
 DELETE FROM users;
 DELETE FROM avatars;
 DELETE FROM profiles;
`
