package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (h handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "success", "message": "server is healthy"}
	json.NewEncoder(w).Encode(response)
}

func (h handler) JobSubmit(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var images_payload ImagesPayload

	err := decoder.Decode(&images_payload)
	checkErr(err)

	if images_payload.Count != len(images_payload.Visits) {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]string{}
		json.NewEncoder(w).Encode(response)
		return
	}

	var job Job
	job.Status = "ongoing"

	if result := h.DB.Create(&job); result.Error != nil {
		fmt.Println(result.Error)
	}

	// goroutine
	go h.ProcessJob(images_payload, job)

	var response = map[string]int{"job_id": job.Id}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h handler) JobStatus(w http.ResponseWriter, r *http.Request) {
	jobid := r.URL.Query().Get("jobid")

	var response map[string]string

	if jobid == "" {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]string{"status": "error", "message": "you are missing jobid query parameter in the job status api"}
		json.NewEncoder(w).Encode(response)
		return
	}

	jobid_int, err := strconv.Atoi(jobid)
	checkErr(err)

	var job Job

	if result := h.DB.First(&job, jobid_int); result.Error != nil {
		fmt.Println(result.Error)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]string{}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if job.Status == "completed" || job.Status == "ongoing" {
		response = map[string]string{"status": job.Status, "job_id": jobid}
	}

	if job.Status == "failed" {
		joberrors := "nil"

		response = map[string]string{"status": job.Status, "job_id": jobid, "error": joberrors}
	}

	json.NewEncoder(w).Encode(response)
}
