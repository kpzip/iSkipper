package main

import (
	"fmt"
	"iskipper/iclickerapi"
	"log"
	"os"
	"strings"
)

const logFileName = "log.txt"

func main() {

	argsWithProg := os.Args
	argsWithoutProg := argsWithProg[1:]
	//logPath := argsWithProg[0] + logFileName
	//fmt.Printf("%s\n", argsWithoutProg[1])

	displayClasses := false
	var joinCourseId *string = nil
	retry := false
	var token *string = nil

	for i, arg := range argsWithoutProg {
		if strings.HasPrefix(arg, "-") {
			if arg == "--courses" {
				displayClasses = true
			} else if arg == "--join" {
				joinCourseId = &argsWithoutProg[i+1]
			} else if arg == "--retry" {
				retry = true
			} else if arg == "--token" {
				token = &argsWithoutProg[i+1]
			} else {
				if strings.Contains(arg, "c") {
					displayClasses = true
				}
				if strings.Contains(arg, "j") {
					joinCourseId = &argsWithoutProg[i+1]
				}
				if strings.Contains(arg, "r") {
					retry = true
				}
				if strings.Contains(arg, "t") {
					token = &argsWithoutProg[i+1]
				}
			}
		}
		// ignore if this is part of another argument
	}

	// get config
	cfg, err := getConfig()
	if err != nil {
		log.Fatalln(err)
	}
	if token == nil {
		token = &cfg.Token
	}

	// setup client
	iClickerClient := iclickerapi.Client(*token, cfg.UserId)

	// Get the list of enrolled courses
	courses, err := iClickerClient.GetCourses()
	if err != nil {
		log.Fatalln(err)
	}

	//TODO error handling if token/userid is incorrect
	if displayClasses {
		// do this chicanery for nice formatting
		var names = []string{"Name"}
		var ids = []string{"Course Id"}
		var enrollmentIds = []string{"Enrollment Id"}
		var namesMaxLen = len(names[0])
		var idsMaxLen = len(ids[0])
		var enrollmentIdsMaxLen = len(enrollmentIds[0])
		for _, course := range courses {
			names = append(names, course.Name)
			ids = append(ids, course.CourseId)
			enrollmentIds = append(enrollmentIds, course.EnrollmentId)
			if len(course.Name) > namesMaxLen {
				namesMaxLen = len(course.Name)
			}
			if len(course.EnrollmentId) > enrollmentIdsMaxLen {
				enrollmentIdsMaxLen = len(course.EnrollmentId)
			}
			if len(course.CourseId) > idsMaxLen {
				idsMaxLen = len(course.CourseId)
			}
		}
		const padding = 2
		fmt.Printf("Courses:\n")
		for i := 0; i < len(names); i++ {
			fmt.Printf("| %s%s| %s%s| %s%s|\n", names[i], strings.Repeat(" ", namesMaxLen+padding-len(names[i])), ids[i], strings.Repeat(" ", idsMaxLen+padding-len(ids[i])), enrollmentIds[i], strings.Repeat(" ", enrollmentIdsMaxLen+padding-len(enrollmentIds[i])))
		}
	}
	if joinCourseId != nil {

		contains := false
		for _, course := range courses {
			if *joinCourseId == course.CourseId {
				contains = true
				break
			}
		}
		if !contains {
			fmt.Printf("Error: course %s not found.\n", *joinCourseId)
		}

		for {
			attendance, err := iClickerClient.JoinCourseAttendanceWithoutGps(*joinCourseId)
			if err != nil {
				log.Fatalln(err)
			}
			if attendance.Error != nil && attendance.Error.Code == 409 {
				if retry {
					fmt.Printf("Failed to join course because course has not started yet. Retrying...\n")
					continue
				} else {
					fmt.Printf("Failed to join course because course has not started yet.\n")
					return
				}
			} else if attendance.Result == "PRESENT" {
				fmt.Printf("Course %s has been successfully joined with attendance id %s.\n", *joinCourseId, attendance.AttendanceId)
				return
			} else {
				log.Fatalln(attendance.Error)
			}
		}

	}

}
