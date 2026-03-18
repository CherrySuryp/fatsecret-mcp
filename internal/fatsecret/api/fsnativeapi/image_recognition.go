package fsnativeapi

import (
	"encoding/json"
	"fmt"
)

const imageRecognitionPath = "image-recognition/v2"

// ImageRecognitionParams holds the parameters for the image-recognition/v2 endpoint.
type ImageRecognitionParams struct {
	// ImageB64 is the Base64-encoded image (jpg, png, or webp). Required.
	// Maximum ~1.09 MB (~999,982 Base64 characters).
	// Recommended resolution: 256×256 or 512×512 pixels.
	ImageB64 string `json:"image_b64"`
	// IncludeFoodData includes full food data in each response item when true.
	IncludeFoodData *bool `json:"include_food_data,omitempty"`
	// EatenFoods provides previously consumed foods to improve matching accuracy.
	EatenFoods []EatenFood `json:"eaten_foods,omitempty"`
	// Region filters results by region (e.g. "US", "FR"). Defaults to "US".
	Region *string `json:"region,omitempty"`
	// Language specifies the response language. Requires Region to be set.
	Language *string `json:"language,omitempty"`
}

// ImageRecognitionResponse is the top-level response envelope for image-recognition/v2.
type ImageRecognitionResponse struct {
	FoodResponse []FoodResponseItem `json:"food_response"`
}

// ImageRecognition calls the image-recognition/v2 endpoint, detecting foods in the
// provided Base64-encoded image and returning structured food and serving data.
// Docs: https://platform.fatsecret.com/docs/v2/image.recognition
func (c *FSNativeAPIClient) ImageRecognition(params ImageRecognitionParams) (*ImageRecognitionResponse, error) {
	body, err := c.post(imageRecognitionPath, params)
	if err != nil {
		return nil, err
	}

	var result ImageRecognitionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("fsnativeapi: unmarshal image.recognition response: %w", err)
	}

	return &result, nil
}