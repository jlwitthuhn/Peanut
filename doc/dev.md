# Developer documentation

## Data access

The only data store for this application is a single Postgres database.

Data is accessed through DAO structs that are defined in the `internal/data` directory. Each DAO is responsible for creating, reading from, and writing to a single table. If a single query needs to read data from multiple tables it should be defined in `dao_multi.go`.

## Templates

Each view in this web application is composed of several HTML template fragments. Each view has a unique name which is registered along with the list of required HTML template fragments in the file `internal/template/template.go`. Whenever a new template is added to the project it must be registered here with a unique name.

In general a template will need to at least include files for:
* A base tmplate
* A view-specific template with unique information
* Embedded CSS
* Embedded JavaScript

Common widgets may also be implemented as a template fragment than can be included on other pages.
