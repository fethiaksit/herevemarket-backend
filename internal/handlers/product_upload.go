package handlers

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type productFormInput struct {
	Name           string
	NameSet        bool
	Price          float64
	PriceSet       bool
	CategoryIDs    []string
	CategoryIDSet  bool
	Description    string
	DescriptionSet bool
	Barcode        string
	BarcodeSet     bool
	Brand          string
	BrandSet       bool
	ImagePath      string
	ImageSet       bool
	Stock          int
	StockSet       bool
	IsActive       bool
	IsActiveSet    bool
	IsCampaign     bool
	IsCampaignSet  bool
}

func parseMultipartProductRequest(c *gin.Context) (productFormInput, error) {
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

		if part.FileName() != "" {
			if name != "image" {
				if _, err := io.Copy(io.Discard, part); err != nil {
					return productFormInput{}, err
				}
				continue
			}

			imagePath, err := saveImagePart(part)
			if err != nil {
				return productFormInput{}, err
			}
			input.ImagePath = imagePath
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
		case "category_id":
			if value != "" {
				input.CategoryIDs = append(input.CategoryIDs, value)
			}
			input.CategoryIDSet = true
		case "image_url", "imageUrl":
			continue
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

	return input, nil
}

func readStringPart(part *multipart.Part) (string, error) {
	data, err := io.ReadAll(part)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func saveImagePart(part *multipart.Part) (string, error) {
	extension := strings.ToLower(filepath.Ext(part.FileName()))
	filename := primitive.NewObjectID().Hex() + extension
	dir := filepath.Join("public", "uploads", "products")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, filename)
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(file, part); err != nil {
		return "", err
	}

	return filepath.ToSlash(filepath.Join("uploads", "products", filename)), nil
}

func parseBoolValue(value string) (bool, error) {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "on" {
		return true, nil
	}
	return strconv.ParseBool(value)
}

func respondMultipartError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data"})
}
