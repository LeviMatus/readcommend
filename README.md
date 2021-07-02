# Readcommend

Readcommend is a book recommendation web app for the true book aficionados and disavowed human-size bookworms. It allows to search for book recommendations with best ratings, based on different search criteria.

# Instructions

The front-end single-page app has already been developed using Node/TypeScript/React (see `/app`) and a PostgreSQL database with sample data is also provided (see `/migrate.*`). Your mission - if you accept it - is to implement a back-end microservice that will fetch the data from the database and serve it to the front-end app.

- In the `service` directory, write a back-end microservice in the language of your choice (with a preference for Go, if you know it) that listens on `http://localhost:5000`.
- Write multiple REST endpoints, all under the `/api/v1` route, as specified in the `/open-api.yaml` swagger file.
- The most important endpoint, `/books`, must return book search results in order of descending ratings (from 5.0 to 1.0 stars) and filtered according to zero, one or multiple user selected criteria: author(s), genre(s), min/max pages, start/end publication date (the "era"). A maximum number of results can also be specified.
- It's OK to use libraries for http handling/routing and SQL (ie: query builders), but try to refrain from relying heavily on end-to-end frameworks (ie: Django) and ORMs that handle everything and leave little room to showcase your coding skills! ;)
- Write some documentation (ie: at the end of this file) to explain how to deploy and run your service.
- Write some unit tests to demonstrate how you test your code.
- Keep your code simple, clean and well-organized.
- If you use Git during development (and we recommend you do!), please ensure your repo is configured as private to prevent future candidates from finding it.
- When you are done, please zip your entire project (excluding the `.git` hidden folder if any) and send the archive to us for review.
- Don't hesitate to come back to us with any questions along the way. We prefer that you ask questions, rather than assuming and misinterpreting requirements.
- You have no time limit to complete this exercise, but the more time you take, the higher our expectations in terms of quality and completeness.
- You will be evaluated mainly based on how well you respect the above instructions. However, we understand that you may have a life (some people do), so if you don't have the time to respect all instructions, simply do your best and focus on what you deem most important.

# Development environment

## Docker Desktop

Make sure you have the latest version of Docker Desktop installed, with sufficient memory allocated to it, otherwise you might run into errors such as:

```
app_1         | Killed
app_1         | npm ERR! code ELIFECYCLE
app_1         | npm ERR! errno 137.
```

If that happens, first try running the command again, but if it doesn't help, try increasing the amount of memory allocated to Docker in Preferences > Resources.

## Starting front-end app and database

In this repo's root dir, run this command to start the front-end app (on port 8080) and PostgreSQL database (on port 5432):

```bash
$ docker-compose up --build
```

(later you can press Ctrl+C to stop this docker composition when you no longer need it)

Wait for everything to build and start properly.

## Creating and seeding database tables

In another terminal window, run this command to create and seed the PostgreSQL database:

```bash
$ ./migrate.sh
```

## Connecting to database

During development, you can connect to and experiment with the PostgreSQL database by running this command:

```bash
$ ./psql.sh
```

To exit the PostgreSQL session, type `\q` and press `ENTER`.

## Accessing front-end app

Point your browser to http://localhost:8080

Be patient, the first time it might take up to 1 or 2 minutes for parcel to build and serve the front-end app.

You should see the front-end app appear, with all components displaying error messages because the back-end service does not exist yet.

# Deploying and running back-end microservice

The Readcommend API is built in a composable manner. Right now it uses a backend Postgres
database, and a REST API (built in Chi) to serve data. The way it is constructed, any backend database
can be dropped in with minimal refactoring due to the use of repository interfaces. Similarly, the API
logic could be replaced by other request handlers, such as gRPC, without extensive code surgery.

Make receipes are provided for backend actions. To use the Makefile, its recommended changing
your working directory to [service/](service).

## Testing
### Unit Tests
Unit tests are provided for the backend service. To run this, issue `make test`.

### Benchmarks
Benchmark tests are provided for the api using minimal mock data. To run these, issue `make benchmark`.

### Test Coverage
To see what percentage of the code has been covered by test cases, issue `make coverage`

### Check linting
To see if the code passes common Golang linting tools, issue `make lint`

## Building

To build the backend API, run `make build`. This will generate a binary that can be executed. Alternatively,
you may elect to install the binary to your path. If you'd like to do this, you can run `make install`. It will
place to binary in your `$GOPATH/bin`.

The output of `make build` is `readcommend`. 

The output of `make install` is a binary you can invoke on your PATH named `readcommend`.

## Running/Deploying

Readcommend API is deployed via a CLI command, `readcommend serve` There are various flags you can pass this
to configure connection settings. The following described how to configure the server:

1. Config File

Use a YAML file at `$HOME/.readcommend.yaml`. The application will pick this up automatically.

```yaml
---
database:
  host: localhost
  port: 5432
  database: readcommend
  schema: public
  ssl-mode: disable
  username: postgres
  password: password123
api:
  port: 5000
  host: 0.0.0.0
```

2. Environment Variables

Environment variables will overwrite config file values. Here is a list of the possible env vars that may be used:

| Parameter         	| Default     	| Description                                                	|
|-------------------	|-------------	|------------------------------------------------------------	|
| DATABASE_HOST     	| localhost   	| The database host to connect to.                           	|
| DATABASE_PORT     	| 5432        	| The port the database is listening on.                     	|
| DATABASE_NAME     	| readcommend 	| The name of the database to connect to.                    	|
| DATABASE_SCHEMA   	| public      	| The schema to connect to in the database.                  	|
| DATABASE_SSL      	| disable     	| whether or not to use ssl-model. Should align with sql.DB. 	|
| DATABASE_USERNAME 	| postgres    	| username to connect with.                                  	|
| API_HOST          	| 0.0.0.0   	| The host at which the API should listen on.                	|
| API_PORT          	| 5000        	| The port at which the API should listen on.                	|

3. CLI Flags

Using flags, you can opt to overwrite config files and/or environment variables. Here is a list of
the possible flags you may provide:

| Parameter    	| Default     	                | Description                                                	|
|--------------	|------------------------------ |------------------------------------------------------------	|
| --db-host      	| localhost   	            | The database host to connect to.                           	|
| --db-port      	| 5432        	            | The port the database is listening on.                     	|
| --db-name      	| readcommend 	            | The name of the database to connect to.                    	|
| --db-schema    	| public      	            | The schema to connect to in the database.                  	|
| --db-ssl-model 	| disable     	            | whether or not to use ssl-model. Should align with sql.DB. 	|
| --db-username  	| postgres    	            | username to connect with.                                  	|
| -db-password  	| false       	            | If true, prompts the user to input a hidden password.      	|
| --api-host     	| 0.0.0.0   	            | The host at which the API should listen on.                	|
| --api-port     	| 5000        	            | The port at which the API should listen on.                	|
| -v            	| false       	            | Verbose logging, if toggled to true.                       	|
| -config           | $HOME/.readcommend       	| Absolute path to your config file.                       	|

#### Examples

With a default config in `$HOME/.readcommend`
> readcommend serve

Overwritting the database name
> readcommend server --db-name=foobar

or

> DATABASE_NAME=foobar readcommend server

Of course, you can always just use
> go run service/main.go serve
