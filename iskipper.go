package main

import (
	"iskipper/iclickerapi"
	"log"
	"os"
	"strings"
)

const logFileName = "log.txt"

func main() {

	argsWithProg := os.Args
	argsWithoutProg := argsWithProg[1:]
	logPath := argsWithoutProg[0] + logFileName

	displayClasses := false
	var joinCourseId *string = nil
	retry := false
	for i, arg := range argsWithProg {
		if strings.HasPrefix(arg, "-") {
			if arg == "--courses" {
				displayClasses = true
			} else if arg == "--join" {
				joinCourseId = &argsWithProg[i+1]
			} else if arg == "--retry" {
				retry = true
			} else if strings.Contains(arg, "c") {
				displayClasses = true
			} else if strings.Contains(arg, "j") {
				joinCourseId = &argsWithProg[i+1]
			} else if strings.Contains(arg, "r") {
				retry = true
			}
		}
		// ignore if this is part of another argument
	}

	// get config
	cfg, err := getConfig()
	if err != nil {
		log.Fatalln(err)
	}
	// setup client
	iClickerClient := iclickerapi.Client(cfg.Token, cfg.UserId)
	courses, err := iClickerClient.GetCourses()
	if err != nil {
		log.Fatalln(err)
	}
	
	course := courses[0]
	attendance, err := iClickerClient.JoinCourseAttendance(course.CourseId, iclickerapi.FromLatLon(34.41135355898463, -119.85483767851561))
	if err != nil {
		return
	}
	log.Printf(attendance.String())

}
