package iclickerapi

import (
	"encoding/json"
	"io"
	"log"
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

type ApiError struct {
  Code int64 `json:"code"`
  Message string `json:"desc"`
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
	request.Header.Add("Client-Tag", "ICLICKER/STUDENT-WEB/2024-11-18T22:28:38.159Z/Win/NT 10.0/Chrome/Web-Browser/131.0.0.0")
	return request, nil
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

	type CoursesGetResponse struct {
		Enrollments []Course `json:"enrollments"`
	}

	var coursesGetResponse CoursesGetResponse
	err = json.Unmarshal([]byte(stringBody), &coursesGetResponse)
	if err != nil {
		return nil, err
	}
	return coursesGetResponse.Enrollments, nil
}

type AttendanceResponse struct {
  AttendanceId string `json:"attendanceId"`
  Result string `json:"result"`
  Method string `json:"method"`
  ProfessorLocation GeoData `json:"profLocation"`
  Error *ApiError `json:"error"`
}

type GeoData struct {
	Accuracy  float64 `json:"accuracy"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

func (client *IClickerClient) JoinCourseAttendance(courseId string, latitude float64, longitude float64, accuracy float64) (*AttendanceResponse, error) {

	type JoinBodyData struct {
		Id  string  `json:"id"`
		Geo GeoData `json:"geo"`
	}

	requestBodyData := JoinBodyData{
		Id: courseId,
		Geo: GeoData{
			Latitude:  latitude,
			Longitude: longitude,
			Accuracy:  accuracy,
		},
	}

	requestBody, _ := json.Marshal(requestBodyData)
	requestBodyString := string(requestBody)
	log.Printf(requestBodyString)

	request, err := client.newRequest(iClickerTrogonApiUrl, "/v2/course/attendance/join/"+courseId, "POST", &requestBodyString)
	if err != nil {
		return nil, err
	}
	response, err := client.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseBodyData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	responseBodyString := string(responseBodyData)
	// log.Printf(responseBodyString)
  var deserializedResponse AttendanceResponse

  err = json.Unmarshal([]byte(responseBodyString), &deserializedResponse)

	return &deserializedResponse, nil
}
