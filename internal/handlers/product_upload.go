package handlers

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const maxImageSizeBytes = 5 << 20

var (
	errImageRequired     = errors.New("image file required")
	errImageTooLarge     = errors.New("image exceeds max size")
	errInvalidImageType  = errors.New("invalid image type")
	errMissingCloudinary = errors.New("missing cloudinary configuration")
)

type productFormInput struct {
	Name             string
	NameSet          bool
	Price            float64
	PriceSet         bool
	Category         []string
	CategorySet      bool
	Description      string
	DescriptionSet   bool
	Barcode          string
	BarcodeSet       bool
	Brand            string
	BrandSet         bool
	Stock            int
	StockSet         bool
	IsActive         bool
	IsActiveSet      bool
	IsCampaign       bool
	IsCampaignSet    bool
	ImageData        []byte
	ImageContentType string
	ImageFilename    string
	ImageSet         bool
}

func parseMultipartProductRequest(c *gin.Context, requireImage bool) (productFormInput, error) {
	reader, err := c.Request.MultipartReader()
	if err != nil {
		return productFormInput{}, err
	}

	input := productFormInput{}

	for {
		part, err := reader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return productFormInput{}, err
		}

		name := part.FormName()
		if name == "" {
			continue
		}

		if name == "image" {
			data, contentType, filename, err := readImagePart(part)
			if err != nil {
				return productFormInput{}, err
			}
			input.ImageData = data
			input.ImageContentType = contentType
			input.ImageFilename = filename
			input.ImageSet = true
			continue
		}

		value, err := readStringPart(part)
		if err != nil {
			return productFormInput{}, err
		}

		switch name {
		case "name":
			input.Name = value
			input.NameSet = true
		case "price":
			parsed, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return productFormInput{}, err
			}
			input.Price = parsed
			input.PriceSet = true
		case "category":
			if value != "" {
				input.Category = append(input.Category, value)
			}
			input.CategorySet = true
		case "description":
			input.Description = value
			input.DescriptionSet = true
		case "barcode":
			input.Barcode = value
			input.BarcodeSet = true
		case "brand":
			input.Brand = value
			input.BrandSet = true
		case "stock":
			parsed, err := strconv.Atoi(value)
			if err != nil {
				return productFormInput{}, err
			}
			input.Stock = parsed
			input.StockSet = true
		case "isActive":
			parsed, err := parseBoolValue(value)
			if err != nil {
				return productFormInput{}, err
			}
			input.IsActive = parsed
			input.IsActiveSet = true
		case "isCampaign":
			parsed, err := parseBoolValue(value)
			if err != nil {
				return productFormInput{}, err
			}
			input.IsCampaign = parsed
			input.IsCampaignSet = true
		}
	}

	if requireImage && !input.ImageSet {
		return productFormInput{}, errImageRequired
	}

	return input, nil
}

func readStringPart(part *multipart.Part) (string, error) {
	data, err := io.ReadAll(part)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func readImagePart(part *multipart.Part) ([]byte, string, string, error) {
	limited := io.LimitedReader{R: part, N: maxImageSizeBytes + 1}
	data, err := io.ReadAll(&limited)
	if err != nil {
		return nil, "", "", err
	}
	if int64(len(data)) > maxImageSizeBytes {
		return nil, "", "", errImageTooLarge
	}
	if len(data) == 0 {
		return nil, "", "", errImageRequired
	}

	detected := http.DetectContentType(data)
	if !strings.HasPrefix(detected, "image/") {
		return nil, "", "", errInvalidImageType
	}

	filename := part.FileName()
	if filename == "" {
		filename = "upload"
	}

	return data, detected, filename, nil
}

func parseBoolValue(value string) (bool, error) {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "on" {
		return true, nil
	}
	return strconv.ParseBool(value)
}

type cloudinaryUploadResponse struct {
	SecureURL string `json:"secure_url"`
	Error     *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func uploadToCloudinary(ctx context.Context, data []byte, filename string, contentType string) (string, error) {
	cloudName := strings.TrimSpace(os.Getenv("CLOUDINARY_CLOUD_NAME"))
	apiKey := strings.TrimSpace(os.Getenv("CLOUDINARY_API_KEY"))
	apiSecret := strings.TrimSpace(os.Getenv("CLOUDINARY_API_SECRET"))
	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return "", errMissingCloudinary
	}

	timestamp := time.Now().Unix()
	signature := signCloudinaryParams(timestamp, apiSecret)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	fileWriter, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	if _, err := fileWriter.Write(data); err != nil {
		return "", err
	}

	if err := writer.WriteField("api_key", apiKey); err != nil {
		return "", err
	}
	if err := writer.WriteField("timestamp", strconv.FormatInt(timestamp, 10)); err != nil {
		return "", err
	}
	if err := writer.WriteField("signature", signature); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	uploadURL := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", cloudName)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadURL, &body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	if contentType != "" {
		req.Header.Set("X-File-Content-Type", contentType)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var payload cloudinaryUploadResponse
	if err := json.Unmarshal(responseBody, &payload); err != nil {
		return "", err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		if payload.Error != nil && payload.Error.Message != "" {
			return "", errors.New(payload.Error.Message)
		}
		return "", fmt.Errorf("cloudinary upload failed with status %d", resp.StatusCode)
	}

	if payload.SecureURL == "" {
		return "", errors.New("missing secure_url from cloudinary")
	}

	// Cloudinary upload keeps the file off local disk; only the secure_url is stored.
	return payload.SecureURL, nil
}

func signCloudinaryParams(timestamp int64, secret string) string {
	base := fmt.Sprintf("timestamp=%d%s", timestamp, secret)
	hash := sha1.Sum([]byte(base))
	return hex.EncodeToString(hash[:])
}

func respondMultipartError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errImageRequired):
		c.JSON(http.StatusBadRequest, gin.H{"error": "image required"})
	case errors.Is(err, errImageTooLarge):
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "image too large"})
	case errors.Is(err, errInvalidImageType):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image type"})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data"})
	}
}
