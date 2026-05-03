# Developer documentation

This document is intended as a short introduction for anyone unfamiliar with the project, it covers the basic of how all of Peanut's core systems fit together.

## Authentication

Users authenticate with Peanut using a username and password. Once a user has successfully authenticated they are given a securely-generated session key which is a bearer token. The session key is attached to all requests sent to Peanut. Upon receiving a request with a session key, middleware looks up the appropriate user from the database and attaches their user ID and other user-specific information to the request.

Session keys are stored in the database and expired server-side after a configurable amount of time so all sessions will terminate automatically even if the user hangs on to their id.

## Authorization

Each user belongs to some amount of groups and each group conveys a list of individual permissions. Whenever a user hits an endpoint that performs a priveleged operation, the endpoint first checks that the user has the appropriate permissions. These permissions are enforced once again at the service layer when possible.

## Data access

The only data store for this application is a single Postgres database.

Data is accessed through DAO structs that are defined in the `internal/data` directory. Each DAO is responsible for creating, reading from, and writing to a single table. If a single query needs to read data from multiple tables it should be defined in `dao_multi.go`.

DAOs are created in `main.go` and are used to construct service instances.

DVAOs are a special type of DAO that accesses data strictly through a view. These are inherently read-only and are the primary way this application looks up data spanning multiple tables.

## Javascript

Javascript is purely optional for clients, no page will require JS to work. Javascript may be used to build optional dynamic GUI elements.

## Services

Most business logic is implemented in the service layer. Like DAOs these are also created in `main.go`.

Only services are allowed to use DAOs and DVAOs. If any other part of the application needs to access stored data it must go through some service.

## Templates

Each view in this web application is composed of several HTML template fragments. Each view has a unique name which is registered along with the list of required HTML template fragments in the file `internal/template/template.go`. The endpoint layer renders based on these view names, not the file names of the templates. Whenever a new view is added, it must be registered in `template.go` to be usable.

In general a template will need to at least include files for:
* A base tmplate
* A view-specific template with unique information
* Embedded CSS
* Embedded JavaScript

Common widgets may also be implemented as a template fragment than can be included on other pages.
