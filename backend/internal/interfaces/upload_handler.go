package interfaces

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct {
	uploadDir string
}

func NewUploadHandler() *UploadHandler {
	uploadDir := "./uploads"
	if dir := os.Getenv("UPLOAD_DIR"); dir != "" {
		uploadDir = dir
	}

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(err)
	}

	return &UploadHandler{
		uploadDir: uploadDir,
	}
}

type UploadResponse struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
}

// UploadImage handles image file uploads
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// Parse multipart form (max 10MB per file)
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only image files are allowed"})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := uuid.New().String() + ext
	filePath := filepath.Join(h.uploadDir, filename)

	// Create file
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer dst.Close()

	// Copy uploaded file to destination
	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Return URL (using public endpoint)
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	c.JSON(http.StatusOK, UploadResponse{
		URL:      baseURL + "/uploads/" + filename,
		Filename: filename,
	})
}

// UploadMultipleImages handles multiple image uploads
func (h *UploadHandler) UploadMultipleImages(c *gin.Context) {
	log.Println("=== [UPLOAD] UploadMultipleImages called ===")
	log.Printf("[UPLOAD] Content-Type: %s", c.ContentType())
	log.Printf("[UPLOAD] Upload directory: %s", h.uploadDir)

	// Parse multipart form (max 50MB total)
	if err := c.Request.ParseMultipartForm(50 << 20); err != nil {
		log.Printf("[UPLOAD ERROR] Failed to parse multipart form: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Files too large"})
		return
	}
	log.Println("[UPLOAD] Multipart form parsed successfully")

	form := c.Request.MultipartForm
	files := form.File["images"]
	log.Printf("[UPLOAD] Found %d files in 'images' field", len(files))

	if len(files) == 0 {
		log.Println("[UPLOAD ERROR] No files found in request")
		log.Printf("[UPLOAD] Available form fields: %v", form.File)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
		return
	}

	if len(files) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum 10 images allowed"})
		return
	}

	var responses []UploadResponse

	for i, fileHeader := range files {
		log.Printf("[UPLOAD] Processing file %d: %s (size: %d bytes)", i+1, fileHeader.Filename, fileHeader.Size)

		file, err := fileHeader.Open()
		if err != nil {
			log.Printf("[UPLOAD ERROR] Failed to open file %s: %v", fileHeader.Filename, err)
			continue
		}
		defer file.Close()

		// Validate file type
		contentType := fileHeader.Header.Get("Content-Type")
		log.Printf("[UPLOAD] File %s Content-Type: %s", fileHeader.Filename, contentType)
		if !strings.HasPrefix(contentType, "image/") {
			log.Printf("[UPLOAD ERROR] File %s rejected: not an image (Content-Type: %s)", fileHeader.Filename, contentType)
			continue
		}

		// Generate unique filename
		ext := filepath.Ext(fileHeader.Filename)
		filename := uuid.New().String() + ext
		filePath := filepath.Join(h.uploadDir, filename)
		log.Printf("[UPLOAD] Saving to: %s", filePath)

		// Create file
		dst, err := os.Create(filePath)
		if err != nil {
			log.Printf("[UPLOAD ERROR] Failed to create file %s: %v", filePath, err)
			continue
		}
		defer dst.Close()

		// Copy uploaded file to destination
		written, err := io.Copy(dst, file)
		if err != nil {
			log.Printf("[UPLOAD ERROR] Failed to copy file %s: %v", filename, err)
			continue
		}
		log.Printf("[UPLOAD] Successfully wrote %d bytes to %s", written, filename)

		// Add to responses
		baseURL := os.Getenv("BASE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}

		response := UploadResponse{
			URL:      baseURL + "/uploads/" + filename,
			Filename: filename,
		}
		responses = append(responses, response)
		log.Printf("[UPLOAD] Added to responses: %s", response.URL)
	}

	log.Printf("[UPLOAD] Successfully uploaded %d/%d files", len(responses), len(files))

	if len(responses) == 0 {
		log.Println("[UPLOAD ERROR] No files were successfully uploaded")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload any files"})
		return
	}

	log.Printf("[UPLOAD SUCCESS] Returning %d file URLs", len(responses))
	c.JSON(http.StatusOK, gin.H{"images": responses})
}

// ServeUploadedFile serves uploaded files
func (h *UploadHandler) ServeUploadedFile(c *gin.Context) {
	filename := c.Param("filename")

	// Security: prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filename"})
		return
	}

	filePath := filepath.Join(h.uploadDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Set cache headers (1 year)
	c.Header("Cache-Control", "public, max-age=31536000")
	c.Header("Expires", time.Now().Add(365*24*time.Hour).Format(http.TimeFormat))

	c.File(filePath)
}
