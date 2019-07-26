package actions

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/plush"
	"github.com/machinebox/graphql"
)

// EntityPageHandler - Use naming convention to load a template and
// corresponding graphql query. Execute the query and use results as
// the model for the template. File locations are:
// Template: entity_pages/{entityType}/{entityType}.html
// Query:    entity_pages/{entityType}/{entityType}.graphql
func EntityPageHandler(c buffalo.Context) error {
	entityType := c.Params().Get("type")
	entityID := c.Params().Get("id")
	
	// could do something like this:
	// theme := envy.Get("THEME", "") 
	viewTemplatePath := fmt.Sprintf("entity_pages/%s/%s.html", entityType, entityType)
	queryTemplatePath := fmt.Sprintf("templates/entity_pages/%s/%s.graphql", entityType, entityType)

	queryTemplate, err := ioutil.ReadFile(queryTemplatePath)
	if err != nil {
		return errors.Wrap(err, "finding query")
	}

	ctx := plush.NewContext()
	query, err := plush.Render(string(queryTemplate), ctx)
	if err != nil {
		return errors.Wrap(err, "processing query")
	}

	endpoint, err := envy.MustGet("GRAPHQL_ENDPOINT")
	if err != nil {
		return errors.Wrap(err, "finding endpoint")
	}
	client := graphql.NewClient(endpoint)
	req := graphql.NewRequest(query)
	req.Var("id", entityID)

	var results map[string]interface{}
	if err := client.Run(ctx, req, &results); err != nil {
		return errors.Wrap(err, "running query")

	}
	c.Set("data", results)
	return c.Render(200, r.HTML(viewTemplatePath))
}
