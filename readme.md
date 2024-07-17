# DB Taxes

This application solves managing taxes applied in different municipalities. The application uses a sqlite database, and implements a RESTful API.

The following assumptions were made:
* Given a request for a municipality with overlapping result, the most specific result is returned (eg. Daily, Weekly,  Monthly, Yearly).
* Responses are the entire record in scope for the request. Ie. a full JSON object is returned, instead of just the tax rate as a number.
* Users of the API are satisfied with non-technical error descriptions.

The project furthermore includes a [Bruno](https://www.usebruno.com/) collection for easy API testing.

### Usage

Build and start the server with
`go run .`

Or by using Docker 
```
docker build -t theilgaard/db_taxes .
docker run theilgaard/db_taxes -p 8080/8080
``` 

The server binds to http://localhost:8080

Endpoints are availble at 
* http://localhost:8080/records
* http://localhost:8080/records/Copenhagen
* http://localhost:8080/records/Copenhagen/2024-01-01

Be aware that the server drops, creates and repopulates the database on restart.

### Testing

Test the server via integration testing with
`go test`

To demonstrate multi stage build with docker, go is only included in the build step.