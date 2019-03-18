package media

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type apiResponse struct {
	JobID  string `json:"job_id" xml:"job_id" form:"job_id" query:"job_id"`
	TaskID string `json:"task_id" xml:"task_id" form:"task_id" query:"task_id"`
}


// AddMedia performs a POST request with credentials and a media payload and receives a JobId and a TaskId in string format
func AddMedia(url string, version string, host string, username string, password string, jobID string, apiToken string) (string, string) {

	addMediaURL := fmt.Sprintf("https://%s/api/job/add_media/?v=%d&job_id=%s&api_token=%s", url, version, jobID, apiToken)

	req, _ := http.NewRequest("POST", addMediaURL, nil)
	// payload is a buffered file
	//      req, _ := http.NewRequest("POST", addMediaUrl, payload)

	req.Header.Add("cookie", "csrftoken=6sV6LZbtELMGtof1KRs5Bcih84nydl4J")
	req.Header.Add("authorization", "Basic Y2llbG9kZXY6Zm9v")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-auth-password", password)
	req.Header.Add("x-auth-user", username)
	req.Header.Add("host", host)

	res, requestErr := http.DefaultClient.Do(req)
	if requestErr != nil {
		fmt.Printf("The request to the API server was not successful. Description: %v", requestErr)
		panic(requestErr)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var lr apiResponse

	err := json.Unmarshal(body, &lr)
	if err != nil {
		fmt.Printf("An error occured while adding media: %v", err)
		panic(err)
	}

	return lr.JobID, lr.TaskID

}
