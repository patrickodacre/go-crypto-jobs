package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chromedp/chromedp"
	"regexp"
)

type Job struct {
	Id             string
	JobTitle       string
	JobDescription string
	CompanyName    string
	CompanyAbout   string
	Remote         bool
	EmploymentType []string
	Skills         string
}

type PageData struct {
	Props struct {
		PageProps struct {
			Job Job
		}
	}
}

func scrape() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	jobs := []map[string]string{}

	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://cryptojobslist.com/blockchain-developer-jobs"),
		chromedp.AttributesAll(`a.jobTitle`, &jobs),
	); err != nil {
		panic(err)
	}

	finalJobsList := []Job{}

	for _, job := range jobs {

		content := ""
		url := job["href"]

		fmt.Println("Getting url", url)
		if err := chromedp.Run(ctx,
			chromedp.Navigate("https://cryptojobslist.com"+url),
			// actual page content is in a JSON string.
			// this page is probably rendered with React.
			chromedp.InnerHTML("#__NEXT_DATA__", &content, chromedp.ByID),
		); err != nil {
			panic(err)
		}

		var data PageData
		json.Unmarshal([]byte(content), &data)

		re := regexp.MustCompile(`(?i)javascript`)

		if re.MatchString(data.Props.PageProps.Job.JobDescription) || re.MatchString(data.Props.PageProps.Job.Skills) {
			finalJobsList = append(finalJobsList, data.Props.PageProps.Job)
		}
	}

	for _, j := range finalJobsList {
		fmt.Println(j.CompanyName)
	}

}
