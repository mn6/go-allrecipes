package main

import (
	"bytes"

	"github.com/valyala/fasthttp"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// ParseResults goes through search result page and
// grabs every card
func ParseResults(query string) []Recipe {
	// Create request string and get search page
	var joinedQuery bytes.Buffer
	joinedQuery.WriteString("https://www.allrecipes.com/search/results/?wt=")
	joinedQuery.WriteString(query)
	root := ParseHTML(joinedQuery.String())
	// Match recipe cards
	matchCards := func(n *html.Node) bool {
		if n.DataAtom == atom.A && n.Parent != nil {
			return scrape.Attr(n.Parent, "class") == "fixed-recipe-card__h3"
		}
		return false
	}

	var recipes []Recipe
	// For each recipe card, make a request for the recipe page
	// and then collect information into Recipe
	rcps := scrape.FindAll(root, matchCards)
	hm := make(chan bool)
	for _, rcp := range rcps {
		href := scrape.Attr(rcp, "href")
		go ParsePage(&recipes, href, hm)
	}
	for i := 0; i < len(rcps); {
		select {
		case <-hm:
			i++
		}
	}
	return recipes
}

// ParsePage takes every found recipe and parses the
// page for needed information
func ParsePage(recipes *[]Recipe, url string, hm chan bool) {
	defer func() {
		hm <- true
	}()

	page := ParseHTML(url)
	matchStars := func(n *html.Node) bool {
		if n.DataAtom == atom.Div && n.Parent != nil {
			return scrape.Attr(n.Parent, "class") == "recipe-summary__stars" && scrape.Attr(n, "class") == "rating-stars"
		}
		return false
	}
	matchThumb := func(n *html.Node) bool {
		if n.DataAtom == atom.Meta && n.Parent != nil {
			return scrape.Attr(n, "property") == "og:image"
		}
		return false
	}
	matchTimes := func(n *html.Node) bool {
		if n.DataAtom == atom.Li && n.Parent != nil {
			return scrape.Attr(n, "class") == "prepTime__item"
		}
		return false
	}
	matchSteps := func(n *html.Node) bool {
		if n.DataAtom == atom.Span && n.Parent != nil && n.Parent.Parent != nil {
			return scrape.Attr(n.Parent.Parent.Parent, "class") == "directions--section__steps" && scrape.Attr(n, "class") == "recipe-directions__list--item"
		}
		return false
	}
	matchCals := func(n *html.Node) bool {
		if n.DataAtom == atom.Span && n.Parent != nil {
			return scrape.Attr(n, "class") != "calorie-count__desc" && scrape.Attr(n.Parent, "class") == "calorie-count"
		}
		return false
	}

	rtng, _ := scrape.Find(page, matchStars)
	name, _ := scrape.Find(page, scrape.ByClass("recipe-summary__h1"))
	thumb, _ := scrape.Find(page, matchThumb)
	author, _ := scrape.Find(page, scrape.ByClass("submitter__name"))
	times := scrape.FindAll(page, matchTimes)
	servings, _ := scrape.Find(page, scrape.ByClass("servings-count"))
	calories, _ := scrape.Find(page, matchCals)
	ingredients := scrape.FindAll(page, scrape.ByClass("recipe-ingred_txt"))
	steps := scrape.FindAll(page, matchSteps)

	var parsedTimes bytes.Buffer
	for _, time := range times {
		parsedTimes.WriteString(scrape.Text(time))
		parsedTimes.WriteString(" | ")
	}

	var parsedSteps []string
	for _, step := range steps {
		parsedSteps = append(parsedSteps, scrape.Text(step))
	}

	var parsedIngredients []string
	for _, ingredient := range ingredients {
		ingredientText := scrape.Text(ingredient)
		if len(ingredientText) > 1 && ingredientText != "Add all ingredients to list" {
			parsedIngredients = append(parsedIngredients, ingredientText)
		}
	}

	rrtng := RoundToUint8(scrape.Attr(rtng, "data-ratingstars"))
	rcalories := RoundToUint16(scrape.Text(calories))

	recipe := Recipe{
		Rating:      rrtng,
		Name:        scrape.Text(name),
		URL:         url,
		Thumb:       scrape.Attr(thumb, "content"),
		Author:      scrape.Text(author),
		Times:       parsedTimes.String(),
		Servings:    scrape.Text(servings),
		Calories:    rcalories,
		Ingredients: parsedIngredients,
		Steps:       parsedSteps,
	}
	*recipes = append(*recipes, recipe)
}

// ParseHTML returns a document
func ParseHTML(url string) *html.Node {
	_, body, err := fasthttp.Get(nil, url)
	if err != nil {
		panic(err)
	}
	document, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	return document
}
