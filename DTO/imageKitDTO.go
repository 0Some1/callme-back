package DTO

type ImageKit struct {
	FileID       string      `json:"fileId"`
	Name         string      `json:"name"`
	Size         int         `json:"size"`
	FilePath     string      `json:"filePath"`
	URL          string      `json:"url"`
	FileType     string      `json:"fileType"`
	Height       int         `json:"height"`
	Width        int         `json:"width"`
	ThumbnailURL string      `json:"thumbnailUrl"`
	AITags       interface{} `json:"AITags"`
}
