## General info

You may use any HTTP methods (supported by Go HTTP server) to make API requests.

Server response will always have `Cache-Control: must-revalidate` header. Most responses return `applicaton/json` content type, but there are exceptions. All API endpoints are safe for concurrent usage.

To start HTTP server, run tiedot with CLI parameters: `-mode=httpd -dir=path_to_db_directory -port=port_number`

To enable HTTPS and disable HTTP, add additional parameters: `-tlskey=keyfile -tlscrt=crtfile`.

To enable mandatory JWT (Javascript Web Token) authorization on all API calls, add additional parameters: `-jwtprivatekey=keyfile2 -jwtpubkey=pubkeyfile`.

The "rsa-test" key-pair in tiedot source code is for testing purpose only, please refrain from using it to start HTTPS server or to enable JWT.

## General error response

Server may respond with HTTP status 400 when:

- A required parameter does not have a value (e.g. ID is required but not given).
- A parameter does not contain correct value data type (e.g. ID should be a number, but letter S is given).

When internal error occurs, server will respond with an error message (plain text) and HTTP status 500; it may also log more details in standard output and/or standard error.

## Collection management
| Function | Method | Routes                                              | Parameters| Legacy Route |  Request | Response Code | Response |
|----------|--------|-----------------------------------------------------|-----------|-----------------|----------|----------|-----------|
| Get all collection names   | GET    | /collections                      |   N/A     | /all |  N/A | 200  | `[ "c_name1", "c_name2", ... ]` |
| Create a collection   | POST   |/collection/`collection_name`            | Collection name | /create |  N/A | 200  | `{ "done" : true }` |
| Rename an existing collection   | PUT    | /collection/`collection_name` | Original collection Name    | /rename |  `{ "new_collection_name" : "new_name"}` | 200  | `{ "done" : true }` |
| Delete a collection   | DELETE |/collection/`collection_name`              |  Collection name     | /drop   |   N/A  | 200  |`{ "done" : true }`  |
| Scrub (compact and Repair) collection   | POST   | /collection/`collection_name`/scrub  |  Collection name    | /scrub |  N/A | 200  |`{ "done" : true }`  |

\* All data files are automatically synchronized every 2 seconds.

## Document management

| Function | Method | Routes                                              | Parameters| Legacy Route |  Request | Response Code | Response |
|----------|--------|-----------------------------------------------------|-----------|-----------------|----------|----------|-----------|
| Insert document | POST | /collection/`collection_name`/doc | Collection name | /insert | `{...}` | 200 | `{ "id" : 2092822085146405929 }` |
| Get a document  | GET  | /collection/`collection_name`/doc/`id` | Collection name  and document id | /get | N/A | 200 | `{...}`
| Update a document | PUT | /collection/`collection_name`/doc/`id` | Collection name and document id | /update | `{...}` | 200 |`{ "done" : true }`|
| Delete a document | DELETE | /collection/`collection_name`/doc/`id` | Collection name and document id | /delete | N/A | 200 | `{ "done" : true }` |
| Get approx. cound of documents | GET | /collection/`collection_name`/count/approx | Collection name | /approxdoccount | N/A | 200 | `{ "count" : 100 }` |
| Get a page of documents | GET | /collection/`collection_name`/page/`page`/of/`total`| Collection name, page number and total number of pages| /getpage | N/A | 200 |  `{ "<id>" : {...}, ... }`|
| Immediately synchronize all data files* | POST | /sync | N/A | /sync | N/A | 200 | `{ "done" : true }`|

\* Document ID is an automatically generated unique ID. It remains unchanged for the document until the document is deleted.

\** "getpage" divides all documents roughly equally large "pages". It is useful for doing collection scan. To calculate total number of pages, first decide how many documents you would like to see in a page, then calculate `"approxdoccount" / DOCS_PER_PAGE`. The documents in HTTP response reflect storage layout and are not ordered.

## Index management

| Function | Method | Routes                                              | Parameters| Legacy Route |  Request | Response Code | Response |
|----------|--------|-----------------------------------------------------|-----------|-----------------|----------|----------|-----------|
| Create an Index | POST |  /collection/`collection_name`/index | Collection name | /index | `{ "path" : "" }` | 200 | `{ "done" : true }`|
| Get list of all indexes in a collection | GET | /collection/`collection_name`/indexes | Collection name | /indexes | N/A | 200 | `{ "", ... }`  |
| Remove an Index | DELETE | /collection/`collection_name`/index | Collection name | /unindex | `{ "path" : "" }` | 200 | `{ "done" : true }`|

## Server management

| Function | Method | Routes                                              | Parameters| Legacy Route |  Request | Response Code | Response |
|----------|--------|-----------------------------------------------------|-----------|-----------------|----------|----------|-----------|
| Dump (backup) database | POST | /dump  | Destination | /dump | `{ "destination : "" }`| 200 | `{ "done": true }`|
| Shutdown Server | POST | /shutdown | N/A | /shutdown | N/A | 200 | N/A | 200 |  `{ "done" : true }` |
| Get go memory allocator statistics | GET | /memstats | N/A | /memstats | N/A | 200 | runtime.MemStats `{...}` |
| Version number | GET | /version | N/A | /version | N/A | 200 | `{ "version" : "6" }` |

## JWT - Javascript Web Token

Launch tiedot HTTP server with JWT will enable mandatory JWT authorization on all API endpoints. The general operation flow is following:

- tiedot HTTP server starts up with JWT enabled.
- Client calls `/getjwt?user=user_name&pass=password` and then memorize the "Authorization" header from response. This is your JW token containing your access rights. The token is encrypted by tiedot server and can only be decrypted by tiedot server. The token is valid for 72 hours.
- Client makes subsequent API calls to JWT-enabled HTTP endpoints, with an additional Authorization header (the return value of `getjwt`) in each request.

| Function | Method | Routes                                              | Parameters| Legacy Route |  Request | Response Code | Response |
|----------|--------|-----------------------------------------------------|-----------|-----------------|----------|----------|-----------|
| Aquire Token | GET | /getjwt | Username and password through Authentication Headers  | /getjwt | N/A | 200 | N/A |
| Check Tokens Validity  | GET | /checkjwt | Authorization Header  | /checkjwt | N/A  | 200 | N/A  |


### User and access rights management in JWT

An HTTP client acquires JW token by calling `/getjwt` with parameter `user` and `pass`. Upon enabling JWT for the first time on a database, the default user "admin" with empty password is created, and the following message will be shown during tiedot startup:

    JWT: successfully initialized DB for JWT features. The default user 'admin' has been created.

The database collection `jwt` stores user, password, and access rights information for all users, including "admin". "admin" is a special user not restrained by access rights assignment.

You are free to manipulate collection `jwt` to manage users, passwords, and access rights. Each document should look like:

    {
        "user": "user_name",
        "pass": "password_plain_text",
        "endpoints": [
            "create",
            "drop",
            "insert",
            "query",
            "update",
            "other_api_endpoint_names..."
        ],
        "collections: [
            "collection_name_A",
            "collection_name_B",
            "other_collection_names..."
        ]
    }

Password is in plain-text, you are free to use a randomly generated password, or a hashed password in an algorithm of your choice.

## Query

| Function | Method | Routes                                              | Parameters| Legacy Route |  Request | Response Code | Response |
|----------|--------|-----------------------------------------------------|-----------|-----------------|----------|----------|-----------|
| Execute query and return documents | POST | /collection/`collection_name`/query | Collection name and query | /query |`{ "c" : [...] }` | 200 | `{ "<id>" : {}, ... }`|
| Executed query count results       | POST | /collection/`collection_name`/count | Collection name and query | /query |`{ "c" : [...] }` | 200 | `{ "count" : 100 }`  |

### Query syntax

Query string is in JSON; it may consist of operators, query parameters, sub-queries and bare-strings. These are the supported query operations (from fastest to slowest):

- Direct document ID (no processing involved)
- Value lookup (field=value)
- Value lookup over integer range (field=1,2,3,4)
- Path existence test (field has value)
- Get all document IDs

There are also set operations - intersect, union, difference, complement; the set operations are very fast.

#### Bare strings (document IDs)

Bare strings are Document IDs that go directly into query result. For example: `["23101561275236320", "2461300515680780859"]`.

#### Basic operations

Lookup finds documents with a specific value in a path: `{"in": [ path ... ], "eq": loookup_value}`.

For example: `{"in": ["Author", "Name", "First Name"], "eq": "John"}`.

Another operation, "has", finds any document with not-null value in the path: `{"has": [ path ...] }`.

For example: `{"has": ["Author", "Name", "Pen Name"]}`.

Integer range query is also supported: `{"in": [ path ... ], "int-from": xx, "int-to": yy}`

For example: `{"in": ["Publish", "Year"], "int-from": 1993, "int-to": 2013, "limit": 10}`

All of the above queries may use an optional "limit" key (for example "limit": 10) to limit number of returned result.

Note that:

- Use "limit": 1 if you intend to get only one result document, this will significantly improve performance.
- Query paths involved in lookup and "has" queries must be indexed beforehand.
- A special operation "all" (bare-string) will return all document IDs; it is the slowest operation of all, but may prove useful in certain set operations such as complement of sets.

#### Set operations

Set operations take a list of sub-queries as parameter, the sub-queries may be arbitrarily complex.

- Intersection: `{"n": [ sub-queries ... ]}`
- Complement: `{"c": [ sub-queries ... ]}`
- Union: `[ sub-queries ...]`

Here is a complicated example: Find all books which were not written by John and published between 1993 and 2013, but include those written by John in 2000.

    [
		{
			"n": [
				{ "in": [ "Author", "Name" ], eq": "John" },
				{ "in": [ "Publish", "Year" ], "eq": 2000 }
			]
		},
		{
			"c": [
				"all",
				{ "n": [
						{ "in": [ "Author", "Name" ], "eq": "John" },
						{ "in": [ "Publish", "Year" ], "int-from": 1993, "int-to": 2013 }
					]
				}
			]
		}
	]

## Embedded usage

tiedot is designed for ease-of-use in both HTTP API and embedded usage. Embedded usage is demonstrated in `example.go`, see the source code comments for details.