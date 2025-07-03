package model

type GenerateOrderImageSignedUploadURLResponse struct {
	SignedURL string `json:"signed_url"` // URL để upload
	ImageKey  string `json:"image_key"`  // S3 object key
	ImageID   int    `json:"image_id"`   // ID của hình ảnh được tạo
}

type OrderImage struct {
	ID       int    `json:"id"`
	OrderID  int    `json:"order_id"`
	ImageURL string `json:"image_url"`
	ImageKey string `json:"image_key"`
}
