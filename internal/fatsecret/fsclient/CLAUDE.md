# internal/fatsecret/fsclient

Typed HTTP client for the FatSecret REST API. Wraps `fsauth.FSOAuth1Client` and handles URL construction, `format=json` injection, and response unmarshalling into typed Go structs.

## Files

- **client.go** — `Client` struct, constructor, and `get`/`post`/`put`/`delete` transport helpers
- **types.go** — shared types: `MealType`, `WeightType`, `HeightType` enums and `SuccessResp`
- **profile_foods_add_favorite.go** — `FoodAddFavorite`
- **profile_foods_delete_favorite.go** — `FoodDeleteFavorite`
- **profile_foods_get_all_favorites.go** — `GetAllFavorites`
- **profile_foods_get_most_eaten.go** — `GetMostEaten`
- **profile_foods_get_recently_eaten.go** — `GetRecentlyEaten`
- **profile_recipes_add_favorite.go** — `RecipeAddFavorite`
- **profile_recipes_delete_favorite.go** — `RecipeDeleteFavorite`
- **profile_recipes_get_favorites.go** — `GetRecipeFavorites`
- **profile_food_diary_saved_meals_create.go** — `SavedMealCreate`
- **profile_food_diary_saved_meals_delete.go** — `SavedMealDelete`
- **profile_food_diary_saved_meals_edit.go** — `SavedMealEdit`
- **profile_food_diary_saved_meals_get.go** — `SavedMealsGet`
- **profile_food_diary_saved_meal_items_add.go** — `SavedMealItemAdd`
- **profile_food_diary_saved_meal_items_edit.go** — `SavedMealItemEdit`
- **profile_food_diary_saved_meal_items_delete.go** — `SavedMealItemDelete`
- **profile_food_diary_saved_meal_items_get.go** — `SavedMealItemsGet`
- **profile_exercise_diary_exercises_get.go** — `ExercisesGet`
- **profile_exercise_diary_entries_commit_day.go** — `ExerciseEntriesCommitDay`
- **profile_exercise_diary_entries_get.go** — `ExerciseEntriesGet`
- **profile_exercise_diary_entries_get_month.go** — `ExerciseEntriesGetMonth`
- **profile_exercise_diary_entries_save_template.go** — `ExerciseEntriesSaveTemplate`
- **profile_exercise_diary_entry_edit.go** — `ExerciseEntryEdit`
- **profile_weight_diary_weight_update.go** — `WeightUpdate`
- **profile_weight_diary_weights_get_month.go** — `WeightsGetMonth`

## Client

`NewClient(cfg)` wraps a `*fsauth.Config` and exposes typed methods per resource:

```go
client := fsclient.NewClient(cfg)
```

Internally, all requests go through four private helpers:

```go
get(path string, params map[string]string) (map[string]any, error)
post(path string, params map[string]string) (map[string]any, error)
put(path string, params map[string]string) (map[string]any, error)
delete(path string, params map[string]string) (map[string]any, error)
```

All four inject `format=json` automatically and delegate signing to `fsauth.FSOAuth1Client`. Resource methods call these helpers with the endpoint-specific path and params, then unmarshal the result into typed structs.

## Implemented endpoints

| File | Method | HTTP | Path |
|------|--------|------|------|
| `profile_foods_add_favorite.go` | `FoodAddFavorite` | POST | `food/favorite/v1` |
| `profile_foods_delete_favorite.go` | `FoodDeleteFavorite` | DELETE | `food/favorite/v1` |
| `profile_foods_get_all_favorites.go` | `GetAllFavorites` | GET | `food/favorites/v2` |
| `profile_foods_get_most_eaten.go` | `GetMostEaten` | GET | `food/most-eaten/v2` |
| `profile_foods_get_recently_eaten.go` | `GetRecentlyEaten` | GET | `food/recently-eaten/v2` |
| `profile_recipes_add_favorite.go` | `RecipeAddFavorite` | POST | `recipe/favorites/v1` |
| `profile_recipes_delete_favorite.go` | `RecipeDeleteFavorite` | DELETE | `recipe/favorites/v1` |
| `profile_recipes_get_favorites.go` | `GetRecipeFavorites` | GET | `recipe/favorites/v2` |
| `profile_food_diary_saved_meals_create.go` | `SavedMealCreate` | POST | `saved-meals/v1` |
| `profile_food_diary_saved_meals_delete.go` | `SavedMealDelete` | DELETE | `saved-meals/v1` |
| `profile_food_diary_saved_meals_edit.go` | `SavedMealEdit` | PUT | `saved-meals/v1` |
| `profile_food_diary_saved_meals_get.go` | `SavedMealsGet` | GET | `saved-meals/v2` |
| `profile_food_diary_saved_meal_items_add.go` | `SavedMealItemAdd` | POST | `saved-meals/item/v1` |
| `profile_food_diary_saved_meal_items_edit.go` | `SavedMealItemEdit` | PUT | `saved-meals/item/v1` |
| `profile_food_diary_saved_meal_items_delete.go` | `SavedMealItemDelete` | DELETE | `saved-meals/item/v1` |
| `profile_food_diary_saved_meal_items_get.go` | `SavedMealItemsGet` | GET | `saved-meals/item/v2` |
| `profile_exercise_diary_exercises_get.go` | `ExercisesGet` | GET | `exercises/v2` |
| `profile_exercise_diary_entries_commit_day.go` | `ExerciseEntriesCommitDay` | POST | `exercise-entries/day/v1` |
| `profile_exercise_diary_entries_get.go` | `ExerciseEntriesGet` | GET | `exercise-entries/v2` |
| `profile_exercise_diary_entries_get_month.go` | `ExerciseEntriesGetMonth` | GET | `exercise-entries/month/v2` |
| `profile_exercise_diary_entries_save_template.go` | `ExerciseEntriesSaveTemplate` | POST | `exercise-entries/template/v1` |
| `profile_exercise_diary_entry_edit.go` | `ExerciseEntryEdit` | PUT | `exercise-entries/template/v1` |
| `profile_weight_diary_weight_update.go` | `WeightUpdate` | POST | `weight/v1` |
| `profile_weight_diary_weights_get_month.go` | `WeightsGetMonth` | GET | `weight/month/v2` |

## File naming convention

Files are named `<prefix>_<operation>.go` where the prefix identifies the FatSecret API resource group:

| Prefix | Resource group |
|--------|---------------|
| `profile_foods_` | Profile food favorites, most eaten, recently eaten |
| `profile_recipes_` | Profile recipe favorites |
| `profile_food_diary_` | Saved meals and saved meal items |
| `profile_exercise_diary_` | Exercise types and exercise entries |
| `profile_weight_diary_` | Weight and height tracking |

Add a new prefix when introducing a new top-level resource group.

## Adding a new endpoint

1. Create a file named `<prefix>_<operation>.go` using the convention above (e.g. `profile_foods_search.go`).
2. Define Go structs matching the FatSecret response schema.
3. Add a method on `*Client` that calls `c.get`, `c.post`, `c.put`, or `c.delete` with the versioned path, then unmarshals the result into typed structs.

The unmarshal pattern is: `json.Marshal(resp)` → `json.Unmarshal(data, &wrapper)` → return inner slice/struct.

Example:

```go
func (c *Client) SearchFoods(query string) ([]Food, error) {
    resp, err := c.get("foods/search/v1", map[string]string{
        "search_expression": query,
    })
    if err != nil {
        return nil, err
    }
    data, _ := json.Marshal(resp)
    var wrapper struct {
        Foods struct { Food []Food `json:"food"` } `json:"foods"`
    }
    json.Unmarshal(data, &wrapper)
    return wrapper.Foods.Food, nil
}
```

## Key invariants

- Base URL is `https://platform.fatsecret.com/rest/` — paths must not include a leading `/`.
- `format=json` is always injected; do not pass it in `params`.
- OAuth signing and token injection are handled entirely by `fsauth` — this package has no signing logic.

## What this package does NOT do

- OAuth signing or token management — that lives in `internal/fatsecret/fsauth`.
- Interactive OAuth flow — that lives in `internal/fatsecret/fsfetchtoken`.