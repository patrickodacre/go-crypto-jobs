package cryptojobslist

import (
	"context"
	"encoding/json"

	"github.com/chromedp/chromedp"
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

func Scrape() []Job {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	scrapedJobs := []map[string]string{}
	jobs := []Job{}

	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://cryptojobslist.com/blockchain-developer-jobs"),
		chromedp.AttributesAll(`a.jobTitle`, &scrapedJobs),
	); err != nil {
		panic(err)
	}

	for _, job := range scrapedJobs {

		content := ""

		if err := chromedp.Run(ctx,
			chromedp.Navigate("https://cryptojobslist.com"+job["href"]),
			// actual page content is in a JSON string.
			// this page is probably rendered with React.
			chromedp.InnerHTML("#__NEXT_DATA__", &content, chromedp.ByID),
		); err != nil {
			panic(err)
		}

		var data PageData
		json.Unmarshal([]byte(content), &data)

		jobs = append(jobs, data.Props.PageProps.Job)
	}

	return jobs
}
