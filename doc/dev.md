# Developer documentation

## Data access

The only data store for this application is a single Postgres database.

Data is accessed through DAO structs that are defined in the `internal/data` directory. Each DAO is responsible for creating, reading from, and writing to a single table. If a single query needs to read data from multiple tables it should be defined in `dao_multi.go`.

DAOs are created in `main.go` and are used to construct service instances.

## Services

Most business logic is implemented in the service layer. Like DAOs these are also created in `main.go`.

Only services are allowed to use DAOs. If any other part of the application needs to access stored data it must go through some service.

## Templates

Each view in this web application is composed of several HTML template fragments. Each view has a unique name which is registered along with the list of required HTML template fragments in the file `internal/template/template.go`. Whenever a new template is added to the project it must be registered here with a unique name.

In general a template will need to at least include files for:
* A base tmplate
* A view-specific template with unique information
* Embedded CSS
* Embedded JavaScript

Common widgets may also be implemented as a template fragment than can be included on other pages.
