# Shopify API
## Create Item
Creates a new inventory item with user-specified data.

|                  |                           |
| :---:            | :----:                    |
| URL              | /api/items                |
| Method           | `POST`                       |
| Body Fields      | Required: `sku`, `name` <br /> Optional: `description`, `price_CAD`, `quantity`   |
| Success Response | Code: `201 Created`|
| Error Responses  | Code: `400 Bad Request` <br /> OR <br /> Code: `409 Conflict` |

### Sample Request Body
```json
{
    "sku": "AB-123_abcd09",
    "name": "Thing 3",
    "description": "the third item",
    "price_CAD": 15.00,
    "quantity": 5
}
```

### Notes:
* A `sku` is 4-12 characters in length and may only contain alphanumeric digits, hyphens, or underscores. (`400 Bad Request`)
* A `sku` must be unique within the system and not currently in use. (`409 Conflict`)
* A `name` may not be the empty string or whitespace. (`400 Bad Request`).
* A `price` may only be a non-negative number. (`400 Bad Request`)
* A `quantity` may only be a non-negative integer. (`400 Bad Request`)
* The default value for a `quantity` is `0`.
* Any extra body fields (i.e. not specified above) will be ignored.
* The Header of a successful request will contain the relative path of the newly created item (`Location` field).

## Get Items
Returns json data about all inventory items.

|                  |                           |
| :---:            | :----:                    |
| URL              | /api/items                |
| Method           | `GET`                       |
| Success Response | Code: `200 OK` |
| Error Responses  | N/A |

### Sample Response Body
```json
[
    {
        "id": "abcdefghijklmnopqrst",
        "sku": "AAAAAAAA",
        "name": "Thing 1",
        "description": "first thing's first",
        "price_CAD": 15.00,
        "quantity": 5
    },
    {
        "id": "01234567890123456789",
        "sku": "BBBBBBBB",
        "name": "Thing 2",
        "quantity": 0
    },
]
```

### Notes:
* `description` and `price_CAD` are optional fields. They are omitted in the response objects if they are present.
* `quantity` is also optional but is given a default value of `0`, so it always appears in response objects.

## Get Item
Returns json data about a single inventory item.

|                  |                           |
| :---:            | :----:                    |
| URL              | /api/items/id             |
| Method           | `GET`                      |
| Success Response | Code: `200 OK` |
| Error Responses  | Code: `404 Not Found` |

### Sample Response Body

endpoint: `/api/items/01234567890123456789`

```json
{
    "id": "01234567890123456789",
    "sku": "BBBBBBBB",
    "name": "Thing 2",
    "quantity": 0
}
```
endpoint: `/api/items/not-a-real-ID`
```json
"there is no item with ID not-a-real-ID"
```

### Notes:
* `description` and `price_CAD` are optional fields. They are omitted in the response object if they are present.
* `quantity` is also optional but is given a default value of `0`, so it always appears in the response object.

## Update Item
Updates an existing inventory item's data with user-provided data. Overwrites all fields; does not perform partial updates.

|                  |                           |
| :---:            | :----:                    |
| URL              | /api/items/id             |
| Method           | `PUT`                      |
| Body Fields      | Required: `sku`, `name` <br /> Optional: `description`, `price_CAD`, `quantity`   |
| Success Response | Code: `204 No Content` |
| Error Responses  | Code: `400 Bad Request` <br /> OR <br /> Code: `404 Not Found` <br /> OR <br /> Code: `409 Conflict` |

### Sample Requests and Responses

#### Example 1
Request

endpoint: `/api/items/01234567890123456789`

```json
{
    "sku": "BBBBBBBB",
    "name": "Thing 2",
    "description": "Update this description please",
    "quantity": 0
}
```
Response

`204 No Content`

#### Example 2

Request

endpoint: `/api/items/not-a-real-ID`

Response

```json
"there is no item with ID not-a-real-ID"
```

`404 Not Found`

#### Example 3

Request

endpoint: `/api/items/01234567890123456789`

```json
{
    "sku": "AAAAAAAA",
    "name": "Thing 2",
    "description": "I'm updating the SKU as well!",
    "quantity": 0
}
```

Response

```json
"there is already an item with SKU AAAAAAAA"
```

`409 Conflict`

### Notes:
* A wholesale replacement is performed. Any optional fields omitted in the request will be overwritten to default values.
* A `sku` is 4-12 characters in length and may only contain alphanumeric digits, hyphens, or underscores. (`400 Bad Request`)
* A `sku` must not be currently in use by a different item. (`409 Conflict`)
* A `name` may not be the empty string or whitespace. (`400 Bad Request`)
* A `price` may only be a non-negative number. (`400 Bad Request`)
* A `quantity` may only be a non-negative integer. (`400 Bad Request`)
* The default value for a `quantity` is `0`.
* Any extra body fields (i.e. not specified above) will be ignored.

## Delete Item
Permanently deletes an item from inventory.

|                  |                           |
| :---:            | :----:                    |
| URL              | /api/items/id             |
| Method           | `DELETE`                 |
| Success Response | Code: `204 No Content` |
| Error Responses  | Code: `404 Not Found` |