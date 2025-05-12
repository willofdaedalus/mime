# Mime Specification

## Entities

* Declared with `entity <name> ->` and closed with `end`.
* Each entity must have at least one field.
* Fields follow the format: `<name> <type> [constraint]*`.
* Entities are referenced using `@entity` syntax.
* Types include: `uuid`, `float`, `int`, `text`, `bool`, `timestamp`

## Attributes (Fields)
Attributes are additional rules applied to fields to elicit certain behaviour. 
They are currently split in two; runtime and database constraints. List below

| Attribute         | Runtime Check?  | DB Constraint?  | Notes                                                                                                                    |
|-------------------|-----------------|-----------------|--------------------------------------------------------------------------------------------------------------------------|
| `check:<expr>`    | ❌ No           | ✅ Yes          | enforced by SQLite using `CHECK` constraints. Great for value bounds, logic rules, etc.                                  |
| `default:<val>`   | ❌ No           | ✅ Yes          | let SQLite handle defaults. You *can* prefill at runtime if you want more control.                                       |
| `foreign:<ref>`   | ❌ No           | ✅ Yes          | references another table and enforces referential integrity. You might validate foreign existence at runtime optionally. |
| `hash`            | ✅ Yes          | ❌ No           | needs runtime hashing using bcrypt. Should only apply to string fields.                                                  |
| `hidden`          | ✅ Yes          | ❌ No           | hides the field from output by default. Controlled by your runtime tooling.                                              |
| `increment`       | ❌ No           | ✅ Yes          | applied as `AUTOINCREMENT` in SQLite. Should never be done in runtime.                                                   |
| `length:min,max`  | ✅ Yes          | ❌ No           | useful for enforcing string length constraints. Could technically be duplicated in DB with `CHECK` if needed.            |
| `override`        | ✅ Yes          | ❌ No           | used to explicitly expose fields marked as `hidden`. Has no DB meaning.                                                  |
| `pattern:<regex>` | ✅ Yes          | ❌ No           | validates a value matches a regex. Only viable in runtime — SQLite regex is limited or requires extensions.              |
| `primary`         | ❌ No           | ✅ Yes          | you can let SQLite enforce it. You’ll still want to ensure only one field is marked as primary at parse time.            |
| `readonly`        | ✅ Yes          | ❌ No           | value is returned in queries but should be ignored in mutations. Logic-only attribute.                                   |
| `required`        | ✅ Yes          | ✅ Yes          | enforced in both runtime (e.g. on insert) and in the DB via `NOT NULL`.                                                  |
| `unique`          | ❌ No           | ✅ Yes          | should be left to SQLite. Runtime enforcement requires costly queries and is race-prone.                                 |

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
