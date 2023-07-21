This application consists of 5 parts

1.) 'dbInit' is for database and tables initialisation (in our case creates myAppDatabase database and promotions table).

2.) 'persistance' part is for interacting with database (get connection and crud operations with promotion table).

3.) in 'resources' part config data and promotions.csv file is keeping.

4.) 'filerader' is for reading data for huge csv file and save it to db. 
You can just run main.go file inside and it will read and save huge data every 30 minutes.

5.)  'api' is simple REST API to interact by postman or curl command with promotion data, to fetch or save data. 
      examnple : GET  http://localhost:9090/promotions/2aaf8f7b-872e-4a14-afbe-1a9b4799dee3
      If you run main.go file it will initiate and run gin web server.