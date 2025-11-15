package pkg

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
)

// ClassificationResult represents the result of image classification
type ClassificationResult struct {
	ItemID     string  `json:"itemId"`
	ItemName   string  `json:"itemName"`
	Price      float64 `json:"price"`
	Confidence float64 `json:"confidence"`
	Matched    bool    `json:"matched"`
}

// VisionClassifier handles image classification using Google Vision API
type VisionClassifier struct {
	client *vision.ImageAnnotatorClient
}

// NewVisionClassifier creates a new Vision API classifier
func NewVisionClassifier(ctx context.Context) (*VisionClassifier, error) {
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create vision client: %v", err)
	}
	return &VisionClassifier{client: client}, nil
}

// Close closes the Vision API client
func (vc *VisionClassifier) Close() error {
	return vc.client.Close()
}

// ClassifyImage classifies an image and matches it to predefined items
func (vc *VisionClassifier) ClassifyImage(ctx context.Context, imageData []byte) (*ClassificationResult, error) {
	// Create image object
	image, err := vision.NewImageFromReader(strings.NewReader(string(imageData)))
	if err != nil {
		return nil, fmt.Errorf("failed to create image: %v", err)
	}

	// Perform label detection
	labels, err := vc.client.DetectLabels(ctx, image, nil, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to detect labels: %v", err)
	}

	// Perform logo detection
	logos, err := vc.client.DetectLogos(ctx, image, nil, 10)
	if err != nil {
		// Logo detection failure is not critical, continue
		logos = nil
	}

	// Perform text detection
	texts, err := vc.client.DetectTexts(ctx, image, nil, 1)
	if err != nil {
		// Text detection failure is not critical, continue
		texts = nil
	}

	// Extract detected features
	var detectedLabels []string
	for _, label := range labels {
		detectedLabels = append(detectedLabels, strings.ToLower(label.Description))
	}

	var detectedLogos []string
	for _, logo := range logos {
		detectedLogos = append(detectedLogos, strings.ToLower(logo.Description))
	}

	var detectedText string
	if len(texts) > 0 {
		detectedText = strings.ToLower(texts[0].Description)
	}

	// Log what Vision API detected
	fmt.Printf("\n=== VISION API RESULTS ===\n")
	fmt.Printf("Labels detected (%d): %v\n", len(detectedLabels), detectedLabels)
	fmt.Printf("Logos detected (%d): %v\n", len(detectedLogos), detectedLogos)
	fmt.Printf("Text detected: %q\n", detectedText)
	fmt.Printf("========================\n\n")

	// Match to predefined items
	return matchToPredefinedItems(detectedLabels, detectedLogos, detectedText)
}

// ClassifyImageBase64 classifies a base64-encoded image
func (vc *VisionClassifier) ClassifyImageBase64(ctx context.Context, base64Image string) (*ClassificationResult, error) {
	// Decode base64 image
	imageData, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 image: %v", err)
	}

	return vc.ClassifyImage(ctx, imageData)
}

// matchToPredefinedItems matches detected features to predefined items
func matchToPredefinedItems(labels []string, logos []string, text string) (*ClassificationResult, error) {
	bestMatch := &ClassificationResult{
		Matched:    false,
		Confidence: 0.0,
	}

	fmt.Printf("=== MATCHING SCORES ===\n")
	// Iterate through predefined items and calculate match scores
	for _, item := range PredefinedItems {
		score := calculateMatchScore(item, labels, logos, text)
		fmt.Printf("%-30s: %.2f%%\n", item.Name, score*100)

		if score > bestMatch.Confidence {
			bestMatch.ItemID = item.ID
			bestMatch.ItemName = item.Name
			bestMatch.Price = item.Price
			bestMatch.Confidence = score
			bestMatch.Matched = score > 0.3 // Threshold for matching
		}
	}
	fmt.Printf("======================\n")
	fmt.Printf("Best match: %s (%.2f%% confidence)\n", bestMatch.ItemName, bestMatch.Confidence*100)
	fmt.Printf("Matched: %v (threshold: 30%%)\n\n", bestMatch.Matched)

	if !bestMatch.Matched {
		return &ClassificationResult{
			Matched:    false,
			Confidence: 0.0,
		}, nil
	}

	return bestMatch, nil
}

// calculateMatchScore calculates a match score for an item based on detected features
func calculateMatchScore(item PredefinedItem, labels []string, logos []string, text string) float64 {
	score := 0.0
	maxScore := 0.0

	// Check labels (weight: 1.0 per match)
	for _, label := range labels {
		maxScore += 1.0
		for _, keyword := range item.Keywords {
			if strings.Contains(label, strings.ToLower(keyword)) || strings.Contains(strings.ToLower(keyword), label) {
				score += 1.0
				break
			}
		}
	}

	// Check logos (weight: 3.0 per match - logos are more reliable)
	for _, logo := range logos {
		maxScore += 3.0
		for _, keyword := range item.Keywords {
			if strings.Contains(logo, strings.ToLower(keyword)) || strings.Contains(strings.ToLower(keyword), logo) {
				score += 3.0
				break
			}
		}
	}

	// Check text (weight: 2.0 if found - text is fairly reliable)
	if text != "" {
		maxScore += 2.0
		for _, keyword := range item.Keywords {
			if strings.Contains(text, strings.ToLower(keyword)) {
				score += 2.0
				break
			}
		}
	}

	// Normalize score to 0-1 range
	if maxScore > 0 {
		return score / maxScore
	}

	return 0.0
}
