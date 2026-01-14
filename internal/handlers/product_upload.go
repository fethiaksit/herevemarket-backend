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

type MultipartProductInput struct {
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

func parseMultipartProductRequest(c *gin.Context) (MultipartProductInput, error) {
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		return MultipartProductInput{}, err
	}

	input := MultipartProductInput{}

	if value, ok := c.GetPostForm("name"); ok {
		input.Name = value
		input.NameSet = true
	}

	if value, ok := c.GetPostForm("price"); ok {
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return MultipartProductInput{}, err
		}
		input.Price = parsed
		input.PriceSet = true
	}

	categoryIDs := c.PostFormArray("category_id")
	if _, ok := c.Request.MultipartForm.Value["category_id"]; ok {
		input.CategoryIDs = categoryIDs
		input.CategoryIDSet = true
	}

	if value, ok := c.GetPostForm("description"); ok {
		input.Description = value
		input.DescriptionSet = true
	}

	if value, ok := c.GetPostForm("barcode"); ok {
		input.Barcode = value
		input.BarcodeSet = true
	}

	if value, ok := c.GetPostForm("brand"); ok {
		input.Brand = value
		input.BrandSet = true
	}

	if value, ok := c.GetPostForm("stock"); ok {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return MultipartProductInput{}, err
		}
		input.Stock = parsed
		input.StockSet = true
	}

	if value, ok := c.GetPostForm("isActive"); ok {
		parsed, err := parseBoolValue(value)
		if err != nil {
			return MultipartProductInput{}, err
		}
		input.IsActive = parsed
		input.IsActiveSet = true
	}

	if value, ok := c.GetPostForm("isCampaign"); ok {
		parsed, err := parseBoolValue(value)
		if err != nil {
			return MultipartProductInput{}, err
		}
		input.IsCampaign = parsed
		input.IsCampaignSet = true
	}

	file, err := c.FormFile("image")
	if err == nil {
		imagePath, err := saveImage(file)
		if err != nil {
			return MultipartProductInput{}, err
		}
		input.ImagePath = imagePath
		input.ImageSet = true
	} else if !errors.Is(err, http.ErrMissingFile) {
		return MultipartProductInput{}, err
	}

	return input, nil
}

func saveImage(file *multipart.FileHeader) (string, error) {
	extension := strings.ToLower(filepath.Ext(file.Filename))
	filename := primitive.NewObjectID().Hex() + extension
	dir := filepath.Join("public", "uploads", "products")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, filename)
	output, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer output.Close()

	input, err := file.Open()
	if err != nil {
		return "", err
	}
	defer input.Close()

	if _, err := io.Copy(output, input); err != nil {
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
