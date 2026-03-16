# internal/fatsecret/api/fsnativeapi

HTTP client for FatSecret REST API endpoints authenticated with OAuth 1.0a (`SignedAuth`).
Handles request signing, parameter encoding, and JSON response decoding.

## Files

| File | Purpose |
|---|---|
| `client.go` | `FSNativeAPIClient` struct, `get` and `post` transport methods |
| `shared_types.go` | Types reused across multiple endpoints (`Serving`, `FoodAttributes`, `FoodImage`, `EatenFood`, NLP/image recognition response types) |
| `foods_search.go` | `FoodsSearch` — `v1/foods.search` |
| `foods_search_v5.go` | `FoodsSearchV5` — `foods/search/v5` |
| `foods_autocomplete.go` | `FoodsAutocomplete` — `food/autocomplete/v2` |
| `food_get.go` | `FoodGet` — `food/v5` |
| `food_find_id_for_barcode.go` | `FoodFindIDForBarcode` — `food/barcode/find-by-id/v2` |
| `food_brands_get.go` | `FoodBrandsGet` — `brands/v2` |
| `food_categories_get.go` | `FoodCategoriesGet` — `food-categories/v2` |
| `food_sub_categories_get.go` | `FoodSubCategoriesGet` — `food-sub-categories/v2` |
| `recipe_get.go` | `RecipeGet` — `recipe/v2` |
| `recipes_search.go` | `RecipesSearch` — `recipes/search/v3` |
| `recipe_types_get.go` | `RecipeTypesGet` — `recipe-types/v2` |
| `natural_language_processing.go` | `NaturalLanguageProcessing` — `natural-language-processing/v1` |
| `image_recognition.go` | `ImageRecognition` — `image-recognition/v2` |

## Client construction

```go
auth := fsauth.NewSignedAuth(consumerKey, consumerSecret)
client := fsnativeapi.NewNativeAPIClient(*auth)
```

`NewNativeAPIClient` takes `fsauth.SignedAuth` by value and defaults to `http.DefaultClient`.

## Transport methods

Both methods are private — endpoint handlers call them internally.

### `get(path, params)`

Used by all GET endpoints. Signs the request over `baseURL/path` + `params`, then appends `params` as query string.

### `post(path, body)`

Used by NLP and image recognition. Marshals `body` to JSON, sets `Content-Type: application/json`. The JSON body is **not** included in the OAuth signature base string (correct per OAuth 1.0a spec — only form-encoded bodies would be).

## Endpoint reference

| Method | REST path | HTTP | Auth scope | Notes |
|---|---|---|---|---|
| `FoodsSearch` | `v1/foods.search` | GET | premier | Lightweight search, no serving detail |
| `FoodsSearchV5` | `foods/search/v5` | GET | premier | Full serving detail, `food_type` filter |
| `FoodsAutocomplete` | `food/autocomplete/v2` | GET | premier | Suggestions only, max 10 results |
| `FoodGet` | `food/v5` | GET | — | Single food with full serving detail |
| `FoodFindIDForBarcode` | `food/barcode/find-by-id/v2` | GET | barcode | GTIN-13 / UPC-A / EAN input |
| `FoodBrandsGet` | `brands/v2` | GET | premier | Filter by letter or `*` for numeric |
| `FoodCategoriesGet` | `food-categories/v2` | GET | premier | Returns all categories |
| `FoodSubCategoriesGet` | `food-sub-categories/v2` | GET | premier | Requires `food_category_id` |
| `RecipeGet` | `recipe/v2` | GET | — | Full recipe with ingredients and directions |
| `RecipesSearch` | `recipes/search/v3` | GET | — | Paginated, calorie/macro range filters |
| `RecipeTypesGet` | `recipe-types/v2` | GET | — | Returns all type names |
| `NaturalLanguageProcessing` | `natural-language-processing/v1` | POST | nlp | Free-text food description → structured data |
| `ImageRecognition` | `image-recognition/v2` | POST | image-recognition | Base64 image → detected foods |

`format=json` is always injected by the transport layer — callers never set it.

## Adding a new endpoint

1. Create `<endpoint_name>.go` in this package.
2. Define a `Params` struct (pointer fields for optional params, value fields for required ones).
3. Define response structs. Use types from `shared_types.go` where applicable (`Serving`, `FoodAttributes`, etc.).
4. Call `c.get(path, p)` or `c.post(path, body)` and unmarshal the result.
5. Apply the single-vs-array pattern (see below) for any collection that may return one item.

## Single-vs-array quirk

FatSecret returns a JSON **object** instead of a single-element **array** for collection fields when there is exactly one result. This applies to `food`, `serving`, `suggestion`, `food_brand`, `food_sub_category`, `recipe`, `recipe_type`, etc.

Handle it by capturing the field as `json.RawMessage` in a private raw struct, then peeking at the first byte:

```go
switch raw.Items[0] {
case '[':
    json.Unmarshal(raw.Items, &result.Items)
case '{':   // object — wrap in slice
    var single Item
    json.Unmarshal(raw.Items, &single)
    result.Items = []Item{single}
case '"':   // string — wrap in slice (for string arrays like brand names)
    var single string
    json.Unmarshal(raw.Items, &single)
    result.Items = []string{single}
}
```

Newer API versions (v2+, v3+) claim consistent array formatting, but the pattern is applied defensively for top-level collections.

## Numeric values

The FatSecret REST API serialises all numeric values (IDs, counts, nutrient amounts) as **JSON strings**. All numeric fields in response structs use `string` Go type. Callers convert as needed.

## What this package does NOT do

- OAuth token management — handled by `internal/fatsecret/fsauth`.
- OAuth 2.0 / Bearer token requests — this client is OAuth 1.0a only.
- Response caching or retry logic.
- Config loading — callers extract credentials from `internal/config` and pass them to `NewNativeAPIClient`.
