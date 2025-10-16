package interfaces

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/ecomate/backend/internal/usecase"
)

type VoiceSearchHandler struct {
	voiceSearchUseCase usecase.VoiceSearchUseCase
}

func NewVoiceSearchHandler(voiceSearchUseCase usecase.VoiceSearchUseCase) *VoiceSearchHandler {
	return &VoiceSearchHandler{
		voiceSearchUseCase: voiceSearchUseCase,
	}
}

// SearchByText handles POST /voice-search/text
func (h *VoiceSearchHandler) SearchByText(c *gin.Context) {
	var req struct {
		Transcript string `json:"transcript" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, err := h.voiceSearchUseCase.SearchByVoice(req.Transcript)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transcript": req.Transcript,
		"products":   products,
	})
}

// SearchByAudio handles POST /voice-search/audio
func (h *VoiceSearchHandler) SearchByAudio(c *gin.Context) {
	file, err := c.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "audio file required"})
		return
	}

	// Read audio file
	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read audio"})
		return
	}
	defer openedFile.Close()

	audioData, err := io.ReadAll(openedFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read audio data"})
		return
	}

	result, err := h.voiceSearchUseCase.ProcessVoiceQuery(audioData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
