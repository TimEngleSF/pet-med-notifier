package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/TimEngleSF/pet-med-notifier/repository"
	"github.com/TimEngleSF/pet-med-notifier/templates"
	"github.com/TimEngleSF/pet-med-notifier/templates/pages"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
)

// indexViewHandler handles a view for the index page.
func indexViewHandler(c echo.Context) error {
	results, err := repository.GetDailyMedicines(c.Request().Context(), *MedDb)
	if err != nil {
		fmt.Printf("Error getting Daily Medicines: %v\n", err)
	}
	// Set the response content type to HTML.
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	// Define template meta tags.
	metaTags := pages.MetaTags(
		"gowebly, htmx example page, go with htmx",               // define meta keywords
		"Welcome to example! You're here because it worked out.", // define meta description
	)

	// Define template body content.
	bodyContent := pages.BodyContent(
		"Welcome to example!",                // define h1 text
		"You're here because it worked out.", // define p text
		results.GroupByTime(),
		time.Now(),
	)

	// Define template layout for index page.
	indexTemplate := templates.Layout(
		"Welcome to example!", // define title text
		metaTags,              // define meta tags
		bodyContent,           // define body content
	)

	return htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, indexTemplate)
}

// showContentAPIHandler handles an API endpoint to show content.
func showContentAPIHandler(c echo.Context) error {
	// Check, if the current request has a 'HX-Request' header.
	// For more information, see https://htmx.org/docs/#request-headers
	if !htmx.IsHTMX(c.Request()) {
		// If not, return HTTP 400 error.
		c.Response().WriteHeader(http.StatusBadRequest)
		slog.Error("request API", "method", c.Request().Method, "status", http.StatusBadRequest, "path", c.Request().URL.Path)
		return echo.NewHTTPError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	// Write HTML content.
	c.Response().Write([]byte("<p>🎉 Yes, <strong>htmx</strong> is ready to use! (<code>GET /api/hello-world</code>)</p>"))

	return htmx.NewResponse().Write(c.Response().Writer)
}

func PutMedicineTakenHandler(c echo.Context) error {
	coll := MedDb.Collection("medicines")
	idStr := c.Request().URL.Query().Get("id")

	taken := c.Request().URL.Query().Get("taken")

	objId, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		// TODO: handle this error by displaying a message and sendin
		log.Println("Error converting id query into ObjectID", err)
	}
	updateValue := taken == "true"

	update := bson.D{{"$set", bson.D{{"taken", updateValue}}}}
	_, err = coll.UpdateByID(c.Request().Context(), objId, update)
	if err != nil {
		// TODO: handle this error by displaying a message and sendin
		log.Println("Error Updating Medicine", err)
	}

	updatedMedicines, err := repository.GetDailyMedicines(c.Request().Context(), *MedDb)
	if err != nil {
		// TODO: handle this error by displaying a message and sendin
		log.Println("Error Getting Medicines", err)
	}
	groupedMeds := updatedMedicines.GroupByTime()
	medicineSectionTemplate := pages.MedicineSection(groupedMeds, groupedMeds.SortKeys())
	return htmx.NewResponse().RenderTempl(c.Request().Context(), c.Response().Writer, medicineSectionTemplate)
}
