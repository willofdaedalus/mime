# Mime Specification

## Entities

* Declared with `entity <name> ->` and closed with `end`.
* Each entity must have at least one field.
* Fields follow the format: `<name> <type> [constraint]*`.
* Entities are referenced using `@entity` syntax.

## Attributes (Fields)

* Types include: `uuid`, `int`, `text`, `bool`, `datetime`, etc.
* Constraints:
    * `required`: must be present.
    * `unique`: must be unique across records.
    * `default <value>`: default if none provided.
    * `&<enum_name>`: value must match a defined enum.

## Enums

* Declared with `enum <name> ->` and closed with `end`.
* Values are simple identifiers.
* Enums are referenced using `&enum_name` syntax in field definitions.

```mime
enum user_role ->
	admin
	user
	guest
end

entity user ->
	id uuid
	name text
	role &user_role
end
```

## Routing

* Syntax: `<VERB> <route> -> <match/expression> || <fallback>`
* Supported verbs: `GET`, `POST`, `PATCH`, `DELETE`.
* Colon-prefixed segments (e.g., `:id`) match path params.
* Query parameters are automatically bound to a `params` map.
* `@entity == params` matches any entity field that appears in `params`.
* Routes fail fast if `params` includes fields not present in the referenced entity.

## Response Handling

* `respond <status_code> <message>` ends the route.
* Supports both success and error responses.

## Query Matching

* `@entity == params` attempts to match all fields.
* If a query param does not exist in the referenced entity, request is rejected early.
* Dot notation is disallowed for now.

## DB Mapped Constraints

* Constraints like `required`, `unique`, `default`, and enum inclusion (`&X`) are enforced at the database level.
* Runtime will defer to the DB to catch these errors wherever possible.

## Runtime-only Constraints

* Cross-entity checks
* Authorization logic
* Complex validations not supported by SQL

## Example (Notes App)

```mime
enum user_role ->
	admin
	user
end

entity user ->
	id uuid
	name text
	password text
	role &user_role
end

entity note ->
	id uuid
	title text
	content text
	owner uuid @user.id
end

routes @user ->
    POST /signup -> create self || respond 400 "signup failed"
    POST /signin -> find self == params || respond 401 "invalid credentials"
    GET /notes/:id -> @note.id == :id || respond 404 "note not found"
    GET /notes -> @note == params || respond 404 "no notes found"
    POST /notes -> create @note || respond 400 "note creation failed"
    PATCH /notes/:id -> update @note.id == :id || respond 400 "update failed"
    DELETE /notes/:id -> delete @note.id == :id || respond 400 "delete failed"
end
```

## Example (LMS)

```mime
entity student ->
	id uuid
	name text
end

entity course ->
	id uuid
	title text
	instructor text
end

entity enrollment ->
	id uuid
	student_id uuid @student.id
	course_id uuid @course.id
end

routes @student ->
    POST /students -> create @student || respond 400 "create failed"
    POST /courses -> create @course || respond 400 "create failed"
    POST /enrollments -> create @enrollment || respond 400 "create failed"
    GET /students/:id -> @student.id == :id || respond 404 "student not found"
end
```
