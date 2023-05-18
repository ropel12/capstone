package pkg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/education-hub/BE/errorr"
	"github.com/sirupsen/logrus"
)

type Quiz struct {
	Client *http.Client
	Auth   string
}
type TestResult struct {
	Email  string `json:"email"`
	Result string `json:"result"`
}

func (q *Quiz) CreateQuiz(schname string, log *logrus.Logger) (string, string, string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.flexiquiz.com/Create/New?surveyName=%s&setupType=FULL", strings.ReplaceAll(schname, " ", "%20")), nil)
	if err != nil {
		log.Errorf("[ERROR]WHEN CREATE QUIZ, Err: %v", err)
		return "", "", "", errorr.NewBad("Cannot Create Quiz")
	}
	cookie := &http.Cookie{Name: ".ASPXAUTH", Value: q.Auth}
	req.AddCookie(cookie)

	resp, err := q.Client.Do(req)
	if err != nil {
		log.Errorf("[ERROR]WHEN CREATE QUIZ, Err: %v", err)
		return "", "", "", errorr.NewBad("Cannot Create Quiz")
	}

	Location := strings.Replace(resp.Header.Get("Location"), "/Create/Edit/", "", 1)
	resp.Body.Close()
	//link
	var regex, _ = regexp.Compile(`previewSurvey\('([^?]+)`)
	var regex2, _ = regexp.Compile(`showEditSurveyPageTitle\('([^']+)`)
	req2, err := http.NewRequest("GET", "https://www.flexiquiz.com/Create/Edit/"+Location, nil)
	if err != nil {
		log.Errorf("[ERROR]WHEN CREATE QUIZ, Err: %v", err)
		return "", "", "", errorr.NewBad("Cannot Create Quiz")
	}
	req2.AddCookie(cookie)
	resp2, err := q.Client.Do(req2)
	if err != nil {
		log.Errorf("[ERROR]WHEN CREATE QUIZ, Err: %v", err)
		return "", "", "", errorr.NewBad("Cannot Create Quiz")
	}
	defer resp2.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp2.Body)
	// Quiz Link
	var quizpath = regex.FindAllString(string(bodyBytes), 1)
	// create a new page
	var prevpath = regex2.FindAllString(string(bodyBytes), 1)
	quizlink := strings.ReplaceAll(quizpath[0], "previewSurvey('", "")
	prevlink := strings.ReplaceAll(prevpath[0], "showEditSurveyPageTitle('", "")
	return quizlink, prevlink, Location, nil
}

func (q *Quiz) GetResult(prevlink string, log *logrus.Logger) ([]TestResult, error) {
	req, err := http.NewRequest("GET", "https://www.flexiquiz.com/Analyze/Index/"+prevlink, nil)
	fmt.Println("linkprev", prevlink)
	if err != nil {
		log.Errorf("[ERROR]WHEN GETTING TEST RESULT, Err: %v", err)
		return nil, errorr.NewInternal("Internal Server Error")
	}
	cookie := &http.Cookie{Name: ".ASPXAUTH", Value: q.Auth}
	req.AddCookie(cookie)

	resp, err := q.Client.Do(req)
	if err != nil {
		log.Errorf("[ERROR]WHEN GETTING TEST RESULT, Err: %v", err)
		return nil, errorr.NewInternal("Internal Server Error")
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Errorf("[ERROR]WHEN GETTING TEST RESULT, Err: %v", err)
		return nil, errorr.NewInternal("Internal Server Error")
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Errorf("[ERROR]WHEN GETTING TEST RESULT, Err: %v", err)
		return nil, errorr.NewInternal("Internal Server Error")
	}
	responses := []TestResult{}
	index := 0
	doc.Find("td span").Each(func(i int, s *goquery.Selection) {
		id, _ := s.Attr("id")
		if strings.Contains(id, "resultsPageGrade") {
			result := s.Text()
			email := doc.Find("td.columnEmail").Nodes[index].FirstChild.Data
			response := TestResult{
				Email:  email,
				Result: result,
			}
			responses = append(responses, response)
			index++
		}
	})
	return responses, nil
}
