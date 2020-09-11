CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  firstname TEXT,
  lastname TEXT,
  roles TEXT[],
  active BOOLEAN DEFAULT TRUE,
  email TEXT UNIQUE,
  password_hash TEXT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS avatars (
  id UUID PRIMARY KEY,
  username TEXT UNIQUE,
  user_id UUID,
  active BOOLEAN,
 created_at TIMESTAMP,
 updated_at TIMESTAMP,

 FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT
);

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
);

INSERT INTO users (id,firstname,lastname,email,password_hash,roles,created_at,updated_at) VALUES 
('5cf37266-3473-4006-984f-9325122678b7','test','user1','testuser1@gmail.com','$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a','{ADMIN}','2019-08-24 00:00:00', '2019-08-24 00:00:00'),
('45b5fbd3-755f-4379-8f07-a58d4a30fa2f','test','user2','testuser2@gmail.com','$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW','{USER}','2019-08-24 00:00:00', '2019-08-24 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO avatars(id,username,user_id,active,created_at,updated_at) VALUES 
('a2b0639f-2cc6-44b8-b97b-15d69dbb511e','DKSnr4','45b5fbd3-755f-4379-8f07-a58d4a30fa2f',TRUE,'2019-01-01 00:00:01.000001+00','2019-01-01 00:00:01.000001+00'),
('72f8b983-3eb4-48db-9ed0-e45cc6bd716b','FelistusJ','45b5fbd3-755f-4379-8f07-a58d4a30fa2f',TRUE,'2019-01-01 00:00:02.000001+00','2019-01-01 00:00:02.000001+00'),
('84b8ff3e-85ec-4929-b045-b2e2d72eb4a7','jean_waithera','5cf37266-3473-4006-984f-9325122678b7',TRUE,'2019-01-05 00:00:02.000001+00','2019-01-05 00:00:02.000001+00')
ON CONFLICT DO NOTHING;

INSERT INTO profiles(id,avatar_id,"name",bio,followers,tweets,"following",likes,profile_image_url,join_date,last_tweet_time,twitter_id,created_at,updated_at)
VALUES
('98b6d4b8-f04b-4c79-8c2e-a0aef46854b7','a2b0639f-2cc6-44b8-b97b-15d69dbb511e','DK Jnr.','Gikuyu ni wendo',244,480,633,602,'https://pbs.twimg.com/profile_images/1293103745648201729/CMwG39AN.jpg','2020-06-02 10:32:52','2020-08-23 15:34:42','1267766034179776512','2019-09-01 00:00:03.000001+00','2019-09-01 01:00:03.000001+00'),
('85f6fb09-eb05-4874-ae39-82d1a30fe0d7','a2b0639f-2cc6-44b8-b97b-15d69dbb511e','DK Jnr.','Gikuyu ni wendo',240,530,653,632,'https://pbs.twimg.com/profile_images/1293103745648201729/CMwG39AN.jpg','2020-06-02 10:32:52','2020-08-23 15:34:42','1267766034179776512','2019-09-02 00:00:03.000001+00','2019-09-02 00:00:03.000001+00'),
('a235be9e-ab5d-44e6-a987-fa1c749264c7','72f8b983-3eb4-48db-9ed0-e45cc6bd716b','Felistus Waithira','CA free spirit.I stand for justice. Proud to be black.',1036,648,961,743,'https://pbs.twimg.com/profile_images/1277896115342528512/uNVpTeIW.jpg','2020-06-03 06:00:31','2020-08-19 08:55:58','1268059887906557953','2019-09-02 00:00:03.000001+00','2019-09-02 00:00:03.000001+00'),
('6bd8a31c-6a58-46f0-8727-8b27b2360a90','84b8ff3e-85ec-4929-b045-b2e2d72eb4a7','Jean Wangari','Cuppycake\n\nLiving large',1284,1105,1361,3062,'https://pbs.twimg.com/profile_images/1288401307204804608/0s5DK5ej.jpg','2020-06-02 09:59:17','2020-08-22 18:55:03','1267757177999101953','2019-09-02 00:00:03.000001+00','2019-09-02 00:00:03.000001+00')
ON CONFLICT DO NOTHING;