# House Projects

This is a simple REST API implemented in Go. The main purpose of this API is to provide simple data
to be consumed by frontend code. I wrote this to use with React and Flutter frontends as a simple
learning tool. This API has no persistent data storage and makes use of an in-memory database which
will persist until the API session is closed.

## Quick Start

To start the API, clone this repo, cd into the repo root, and run `$ go run cmd/*.go`. The API will
start running on localhost:8080.

## API

**GET: /projects** 
Fetches all projects in the database

**GET: /projects/:title**
Fetches the project with the given title

**POST: /projects
{
	"title"         string
	"duration_days" number
	"cost"          number
	"description"   string
}**
Creates a new project. Body must be JSON encoded.

**PUT: /projects/:title
{
	"title"         string
	"duration_days" number
	"cost"          number
	"description"   string
}**
Updates an existing project with the given title. Body must be JSON encoded.

**DELETE: /projects/:title**
Deletes the project with the given title.

## License
This software is free to use under the MIT license.

## Contact
For more information, or to report issues, please email jeff.moorhead1@gmail.com.
