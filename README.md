RUNNING BACKEND
===============
There are two binaries one for admin and one for the app. The admin is used to set up the generating public and private ssh keys for use in generating jwt tokens.

STEP 1
======
Ensure you have docker and docker-compose installed.
Make sure that no program is running on port 5432
      sudo service postgresql stop
Then there is a docker-compose file for our postgresql container image. 
Run it like so
      docker-compose -f docker-compose.test.yml up postgres

STEP 2
======
Run the admin binary with the following command

./admin keygen ~/.avatarlysis

STEP 3
======
There is an env.sh file that sets the environment variables.
Import it in your environment like so:

source env.sh

STEP 4
=======

Run the app binary

./app

STEP 5
=======
Run the frontend.To run it you will use docker to load the frontend image like so 

  docker load -i frontend.tar

After loading it in memory.Run the image like so

  docker run -it -p 8080:80 --rm avatarlysis:web-dev

Then in your browser, go to http://localhost:8080

   username: testuser1@gmail.com
   password: gophers