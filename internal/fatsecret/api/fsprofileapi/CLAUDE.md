# internal/fatsecret/api/fsprofileapi

HTTP client for FatSecret profile (user-scoped) API endpoints authenticated with OAuth 1.0a 3-legged (`DelegatedAuth`). All operations act on behalf of a specific user profile.

## Files

| File | Purpose |
|---|---|
| `client.go` | `FSProfileClient` struct; `get`, `post`, `put`, `delete` transport methods; `encodeParams` helper |
| `shared_types.go` | `SuccessResponse` — common response for mutating endpoints that return no resource |
| `profile_create.go` | `ProfileCreate` — `profile/v1` (POST); also defines `ProfileAuthResponse` |
| `profile_get.go` | `ProfileGet` — `profile/v1` (GET) |
| `profile_get_auth.go` | `ProfileGetAuth` — `profile/auth/v1` (GET) |
| `food_entries_get.go` | `FoodEntriesGet` — `food-entries/v2` (GET); defines `FoodEntry`, `FoodEntriesResult` |
| `food_entry_create.go` | `FoodEntryCreate` — `food-entries/v1` (POST) |
| `food_entry_edit.go` | `FoodEntryEdit` — `food-entries/v1` (PUT) |
| `food_entry_delete.go` | `FoodEntryDelete` — `food-entries/v1` (DELETE) |
| `food_entries_get_month.go` | `FoodEntriesGetMonth` — `food-entries/month/v1` (GET) |
| `food_entries_copy.go` | `FoodEntriesCopy` — `food-entries/copy/v1` (POST) |
| `food_entries_copy_saved_meal.go` | `FoodEntriesCopySavedMeal` — `food-entries/copy/saved-meal/v1` (POST) |
| `food_create.go` | `FoodCreate` — `food/v2` (POST) |
| `food_add_favorite.go` | `FoodAddFavorite` — `food/favorite/v1` (POST) |
| `food_delete_favorite.go` | `FoodDeleteFavorite` — `food/favorite/v1` (DELETE) |
| `foods_get_favorites.go` | `FoodsGetFavorites` — `food/favorites/v2` (GET) |
| `foods_get_most_eaten.go` | `FoodsGetMostEaten` — `food/most-eaten/v2` (GET) |
| `foods_get_recently_eaten.go` | `FoodsGetRecentlyEaten` — `food/recently-eaten/v2` (GET) |
| `exercise_entries_get.go` | `ExerciseEntriesGet` — `exercise-entries/v2` (GET) |
| `exercise_entries_get_month.go` | `ExerciseEntriesGetMonth` — `exercise-entries/month/v2` (GET) |
| `exercise_entries_commit_day.go` | `ExerciseEntriesCommitDay` — `exercise-entries/day/v1` (POST) |
| `exercise_entries_save_template.go` | `ExerciseEntriesSaveTemplate` — `exercise-entries/template/v1` (POST) |
| `exercise_entry_edit.go` | `ExerciseEntryEdit` — `exercise-entries/template/v1` (PUT) |
| `exercises_get.go` | `ExercisesGet` — `exercises/v2` (GET) |
| `saved_meals_get.go` | `SavedMealsGet` — `saved-meals/v2` (GET) |
| `saved_meal_create.go` | `SavedMealCreate` — `saved-meals/v1` (POST) |
| `saved_meal_edit.go` | `SavedMealEdit` — `saved-meals/v1` (PUT) |
| `saved_meal_delete.go` | `SavedMealDelete` — `saved-meals/v1` (DELETE) |
| `saved_meal_items_get.go` | `SavedMealItemsGet` — `saved-meals/item/v2` (GET) |
| `saved_meal_item_add.go` | `SavedMealItemAdd` — `saved-meals/item/v1` (POST) |
| `saved_meal_item_edit.go` | `SavedMealItemEdit` — `saved-meals/item/v1` (PUT) |
| `saved_meal_item_delete.go` | `SavedMealItemDelete` — `saved-meals/item/v1` (DELETE) |
| `recipe_add_favorite.go` | `RecipeAddFavorite` — `recipe/favorites/v1` (POST) |
| `recipe_delete_favorite.go` | `RecipeDeleteFavorite` — `recipe/favorites/v1` (DELETE) |
| `recipes_get_favorites.go` | `RecipesGetFavorites` — `recipe/favorites/v2` (GET) |
| `weight_update.go` | `WeightUpdate` — `weight/v1` (POST) |
| `weights_get_month.go` | `WeightsGetMonth` — `weight/month/v2` (GET) |

## Client construction

```go
auth := fsauth.NewDelegatedAuth(consumerKey, consumerSecret, accessToken, tokenSecret)
client := fsprofileapi.NewProfileClient(auth)
```

`NewProfileClient` takes `*fsauth.DelegatedAuth` and defaults to `http.DefaultClient`.

## Transport methods

All four methods are private — endpoint handlers call them internally. All params are passed as `map[string]string` and always include `"format": "json"`.

| Method | HTTP | Body encoding | OAuth body signing |
|---|---|---|---|
| `get` | GET | query string | params signed |
| `post` | POST | `application/x-www-form-urlencoded` | params signed (per OAuth 1.0a spec) |
| `put` | PUT | `application/x-www-form-urlencoded` | params signed |
| `delete` | DELETE | query string | params signed |

Note: `post` and `put` have a path-building bug — they prepend `baseURL + "/" + path`, producing a double slash (`https://platform.fatsecret.com/rest//path`). `get` and `delete` use `baseURL + path` (correct). Do not "fix" these without verifying against the live API.

## Endpoint reference

| Method | REST path | HTTP |
|---|---|---|
| `ProfileCreate` | `profile/v1` | POST |
| `ProfileGet` | `profile/v1` | GET |
| `ProfileGetAuth` | `profile/auth/v1` | GET |
| `FoodEntriesGet` | `food-entries/v2` | GET |
| `FoodEntryCreate` | `food-entries/v1` | POST |
| `FoodEntryEdit` | `food-entries/v1` | PUT |
| `FoodEntryDelete` | `food-entries/v1` | DELETE |
| `FoodEntriesGetMonth` | `food-entries/month/v1` | GET |
| `FoodEntriesCopy` | `food-entries/copy/v1` | POST |
| `FoodEntriesCopySavedMeal` | `food-entries/copy/saved-meal/v1` | POST |
| `FoodCreate` | `food/v2` | POST |
| `FoodAddFavorite` | `food/favorite/v1` | POST |
| `FoodDeleteFavorite` | `food/favorite/v1` | DELETE |
| `FoodsGetFavorites` | `food/favorites/v2` | GET |
| `FoodsGetMostEaten` | `food/most-eaten/v2` | GET |
| `FoodsGetRecentlyEaten` | `food/recently-eaten/v2` | GET |
| `ExerciseEntriesGet` | `exercise-entries/v2` | GET |
| `ExerciseEntriesGetMonth` | `exercise-entries/month/v2` | GET |
| `ExerciseEntriesCommitDay` | `exercise-entries/day/v1` | POST |
| `ExerciseEntriesSaveTemplate` | `exercise-entries/template/v1` | POST |
| `ExerciseEntryEdit` | `exercise-entries/template/v1` | PUT |
| `ExercisesGet` | `exercises/v2` | GET |
| `SavedMealsGet` | `saved-meals/v2` | GET |
| `SavedMealCreate` | `saved-meals/v1` | POST |
| `SavedMealEdit` | `saved-meals/v1` | PUT |
| `SavedMealDelete` | `saved-meals/v1` | DELETE |
| `SavedMealItemsGet` | `saved-meals/item/v2` | GET |
| `SavedMealItemAdd` | `saved-meals/item/v1` | POST |
| `SavedMealItemEdit` | `saved-meals/item/v1` | PUT |
| `SavedMealItemDelete` | `saved-meals/item/v1` | DELETE |
| `RecipeAddFavorite` | `recipe/favorites/v1` | POST |
| `RecipeDeleteFavorite` | `recipe/favorites/v1` | DELETE |
| `RecipesGetFavorites` | `recipe/favorites/v2` | GET |
| `WeightUpdate` | `weight/v1` | POST |
| `WeightsGetMonth` | `weight/month/v2` | GET |

`format=json` is always injected by each endpoint handler — callers never set it.

## Date parameters

Several endpoints take a date as days since the Unix epoch (1970-01-01). Optional date params use `*int` (nil → omit from request); required dates use `int`.

## Adding a new endpoint

1. Create `<endpoint_name>.go` in this package.
2. Define a `Params` struct (pointer fields for optional params, value fields for required ones).
3. Define response structs. Use `SuccessResponse` from `shared_types.go` for mutating endpoints that return no resource.
4. Inject `"format": "json"` into the params map.
5. Call the appropriate transport method (`c.get`, `c.post`, `c.put`, or `c.delete`) and unmarshal the result.

## Numeric values

The FatSecret API serialises all numeric values (IDs, counts, nutrient amounts) as **JSON strings**. All numeric fields in response structs use `string` Go type. Callers convert as needed.

## What this package does NOT do

- OAuth token management or fetching — handled by `internal/fatsecret/fsauth` and `internal/fatsecret/fsfetchtoken`.
- OAuth 1.0a signing — delegated to `fsauth.DelegatedAuth`.
- Non-user-scoped requests — use `fsnativeapi` for public endpoints.
- Response caching or retry logic.
- Config loading — callers extract credentials from `internal/config` and pass them to `NewProfileClient`.