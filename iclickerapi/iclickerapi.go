package iclickerapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const iClickerApiUrl = "https://api.iclicker.com"
const iClickerTrogonApiUrl = "https://api.iclicker.com/trogon"

type IClickerClient struct {
	Client *http.Client
	Token  string
	UserId string
}

func Client(token string, userId string) *IClickerClient {
	return &IClickerClient{
		Client: &http.Client{},
		Token:  token,
		UserId: userId,
	}
}

type Course struct {
	EnrollmentId string `json:"enrollmentId"`
	CourseId     string `json:"courseId"`
	Name         string `json:"name"`
}

func (client *IClickerClient) newRequest(url string, path string, method string, body *string) (*http.Request, error) {
	var readerBody io.Reader
	if body != nil {
		readerBody = strings.NewReader(*body)
	} else {
		readerBody = nil
	}
	request, err := http.NewRequest(method, url+path, readerBody)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Priority", "u=0, i")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Reef-Auth-Type", "oauth")
	request.Header.Add("Authorization", "Bearer "+client.Token)
	return request, nil
}

type CoursesGetResponse struct {
	Enrollments []Course `json:"enrollments"`
}

func (client *IClickerClient) GetCourses() ([]Course, error) {
	request, err := client.newRequest(iClickerApiUrl, "/v1/users/"+client.UserId+"/views/student-courses", "GET", nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	stringBody := string(body)
	// log.Printf(stringBody)
	var coursesGetResponse CoursesGetResponse
	err = json.Unmarshal([]byte(stringBody), &coursesGetResponse)
	if err != nil {
		return nil, err
	}
	return coursesGetResponse.Enrollments, nil
}

func (client *IClickerClient) JoinCourseAttendance(courseId string) (*string, error) {
	bodyString := fmt.Sprintf("{\"id\":\"%s\",\"geo\":{\"accuracy\":\"%s\",\"lat\":\"%s\",\"lon\":\"%s\"}}", courseId, string(rune(0)), string(rune(0)), string(rune(0)))

	request, err := client.newRequest(iClickerTrogonApiUrl, "/v2/course/attendance/join/", "POST", &bodyString)
	if err != nil {
		return nil, err
	}
	response, err := client.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	stringBody := string(body)
	// log.Printf(stringBody)

	return &stringBody, nil
}
