package main

// json responses
type StorePayload struct {
	StoreId   string   `json:"store_id"`
	ImageUrls []string `json:"image_url"`
	VisitTime string   `json:"visit_time"`
}

type ImagesPayload struct {
	Count  int            `json:"count"`
	Visits []StorePayload `json:"visits"`
}

type JobError struct {
	StoreId string `json:"store_id"`
	Error   string `json:"error"`
}

type JobStatusResponse struct {
	JobId  string     `json:"job_id"`
	Status string     `json:"status"`
	Error  []JobError `json:"error"`
}

// db models
type Job struct {
	Id     int    `json:"id" gorm:"primaryKey"`
	Status string `json:"status"`
}

type Image struct {
	Id      int    `json:"id" gorm:"primaryKey"`
	StoreId string `json:"store_id"`
	Url     string `json:"url"`
}
