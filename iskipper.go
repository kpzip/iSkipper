package main

import (
	"iskipper/iclickerapi"
	"log"
)

func main() {
	// get config
	cfg, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}
	// setup client
	iClickerClient := iclickerapi.Client(cfg.Token, cfg.UserId)
	courses, err := iClickerClient.GetCourses()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v\n", courses)
	course := courses[0]
	attendance, err := iClickerClient.JoinCourseAttendance(course.CourseId, "", "", "")
	if err != nil {
		return
	}
	log.Printf(*attendance)

}
