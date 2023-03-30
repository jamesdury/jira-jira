package routes

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"github.com/jxxxxxxxy/jira-jira/internal/jira"
)

// https://blog.logrocket.com/express-style-api-go-fiber/

func SetupRoutes(app *fiber.App) {

	web := app.Group("/")

	j := jira.Config{
		Host:  os.Getenv("HOST"),
		Email: os.Getenv("MAIL"),
		Token: os.Getenv("TOKEN"),
	}

	a, err := j.Create()

	if err != nil {
		panic(err)
	}

	// TODO change to media
	web.Get("/image", func(c *fiber.Ctx) error {

		url := c.Query("url")
		mimetype := c.Query("mimetype")

		if url == "" || mimetype == "" {
			c.Status(fiber.StatusNoContent)
			return c.JSON(url)
		}

		body, err := j.GetImageBytes(url)
		if err != nil {
			c.Status(fiber.StatusNoContent)
			panic(err)
		}

		c.Status(fiber.StatusOK)
		c.Set("Content-Type", mimetype)
		c.Write(body)

		return nil
	})

	web.Get("/", func(c *fiber.Ctx) error {
		myself, err := j.GetMyself(a)

		if err != nil {
			panic(err)
		}
		return c.Render("index", fiber.Map{
			"Myself": myself,
		}, "layouts/myself")
	})

	web.Get("/issues", func(c *fiber.Ctx) error {

		jql := c.Query("jql")

		issues, err := j.QueryJql(a, jql)

		if err != nil {
			return c.Render("index", fiber.Map{}, "layouts/issues")
		}

		//return c.JSON( issues)

		c.Status(fiber.StatusOK)

		return c.Render("index", fiber.Map{
			"Jql":    jql,
			"Issues": issues.Issues,
		}, "layouts/issues")
	})

	web.Get("/issues/:id.json", func(c *fiber.Ctx) error {

		id := c.Params("id")

		issue, err := j.GetIssue(a, id)
		if err != nil {
			panic(err)
		}
		/*
		        , err := j.GetIssueLinks(a, id)
				if err != nil {
					panic(err)
				}
		*/

		c.Status(fiber.StatusOK)
		return c.JSON(issue)
	})

	web.Get("/issues/:id", func(c *fiber.Ctx) error {

		id := c.Params("id")

		issue, err := j.GetIssue(a, id)
		if err != nil {
			panic(err)
		}

		issueLinks, err := j.GetIssueLinks(a, id)
		if err != nil {
			panic(err)
		}

		c.Status(fiber.StatusOK)

		return c.Render("index", fiber.Map{
			"Key":        id,
			"Project":    issue.Fields.Summary,
			"Issue":      issue.Fields,
			"IssueLinks": issueLinks,
		}, "layouts/issue")
	})
}
