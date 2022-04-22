package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

func DownloadImage(url string) {
	response, err := http.Get(url)
	checkErr(err)

	defer response.Body.Close()

	file, err := os.Create(path.Base(url))

	checkErr(err)

	defer file.Close()

	_, err = io.Copy(file, response.Body)

	checkErr(err)
}

func (h handler) ProcessJob(images_payload ImagesPayload, job Job) {
	for _, visit := range images_payload.Visits {
		for _, url := range visit.ImageUrls {
			DownloadImage(url)

			var image Image

			image.StoreId = visit.StoreId
			image.Url = url

			if result := h.DB.Create(&image); result.Error != nil {
				fmt.Println(result.Error)
			}
		}
	}

	job.Status = "completed"
	h.DB.Save(&job)
}
