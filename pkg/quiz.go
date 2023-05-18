package pkg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/education-hub/BE/errorr"
	"github.com/sirupsen/logrus"
)

type Quiz struct {
	Client *http.Client
	Auth   string
}

func (q *Quiz) CreateQuiz(schname string, log *logrus.Logger) (string, string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.flexiquiz.com/Create/New?surveyName=%s&setupType=FULL", schname), nil)
	if err != nil {
		log.Errorf("[ERROR]WHEN CREATE QUIZ, Err: %v", err)
		return "", "", errorr.NewBad("Cannot Create Quiz")
	}
	cookie := &http.Cookie{Name: ".ASPXAUTH", Value: q.Auth}
	req.AddCookie(cookie)

	resp, err := q.Client.Do(req)
	if err != nil {
		log.Errorf("[ERROR]WHEN CREATE QUIZ, Err: %v", err)
		return "", "", errorr.NewBad("Cannot Create Quiz")
	}

	resp.Body.Close()
	Location := strings.Replace(resp.Header.Get("Location"), "/Create/Edit/", "", 1)
	//link
	var regex, _ = regexp.Compile(`previewSurvey\('([^?]+)`)
	var regex2, _ = regexp.Compile(`showEditSurveyPageTitle\('([^']+)`)
	req2, err := http.NewRequest("GET", "https://www.flexiquiz.com/Create/Edit/"+Location, nil)
	if err != nil {
		log.Errorf("[ERROR]WHEN CREATE QUIZ, Err: %v", err)
		return "", "", errorr.NewBad("Cannot Create Quiz")
	}
	cookie = &http.Cookie{Name: ".ASPXAUTH", Value: q.Auth}
	req2.AddCookie(cookie)
	resp2, err := q.Client.Do(req2)
	if err != nil {
		log.Errorf("[ERROR]WHEN CREATE QUIZ, Err: %v", err)
		return "", "", errorr.NewBad("Cannot Create Quiz")
	}
	bodyBytes, err := ioutil.ReadAll(resp2.Body)
	// Quiz Link
	var quizpath = regex.FindAllString(string(bodyBytes), 1)
	// create a new page
	var prevpath = regex2.FindAllString(string(bodyBytes), 1)
	quizlink := strings.ReplaceAll(quizpath[0], "previewSurvey('", "")
	prevlink := strings.ReplaceAll(prevpath[0], "showEditSurveyPageTitle('", "")
	resp2.Body.Close()
	return quizlink, prevlink, nil
}

