package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"

	"github.com/cucumber/godog"
	"github.com/docker/go-connections/nat"
	"github.com/gofiber/fiber/v2"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/hedon-go-road/bdd-demo/database"
	"github.com/hedon-go-road/bdd-demo/models"
	"github.com/hedon-go-road/bdd-demo/routes"
)

type godogsResponseCtxKey = struct{}

func init() {
	ctx := context.Background()

	const dbname = "test-db"
	const user = "postgres"
	const password = "password"

	port, _ := nat.NewPort("tcp", "5432")

	container, err := startContainer(ctx,
		WithPort(port.Port()),
		WithInitialDatabase(user, password, dbname),
		WithWaitStrategy(wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2),
		),
	)
	if err != nil {
		panic(err)
	}

	containerPort, _ := container.MappedPort(ctx, port)
	host, _ := container.Host(ctx)

	_ = os.Setenv("DB_HOST", host)
	_ = os.Setenv("DB_PORT", containerPort.Port())
	_ = os.Setenv("DB_USER", user)
	_ = os.Setenv("DB_PASS", password)
	_ = os.Setenv("DB_NAME", dbname)
}

type apiFeature struct {
	app *fiber.App
}

type response struct {
	status int
	body   any
}

func (a *apiFeature) resetResponse(*godog.Scenario) {
	a.app = fiber.New()
	routes.SetupRoutes(a.app)
	database.ConnectDB()
}

func (a *apiFeature) iSendRequestToWithPayload(ctx context.Context, method, route string,
	payloadDoc *godog.DocString) (context.Context, error) {
	var reqBody []byte

	if payloadDoc != nil {
		payloadMap := models.Book{}
		err := json.Unmarshal([]byte(payloadDoc.Content), &payloadMap)
		if err != nil {
			panic(err)
		}

		reqBody, _ = json.Marshal(payloadMap)
	}

	req := httptest.NewRequest(method, route, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := a.app.Test(req)
	defer func() { _ = resp.Body.Close() }()
	var createdBooks []models.Book
	_ = json.NewDecoder(resp.Body).Decode(&createdBooks)

	actual := response{
		status: resp.StatusCode,
		body:   createdBooks,
	}

	return context.WithValue(ctx, godogsResponseCtxKey{}, actual), nil //nolint:staticcheck
}

func (a *apiFeature) theResponseCodeShouldBe(ctx context.Context, expectedStatus int) error {
	resp, ok := ctx.Value(godogsResponseCtxKey{}).(response)
	if !ok {
		return errors.New("there are no godogs available in the context")
	}

	if expectedStatus != resp.status {
		if resp.status >= http.StatusBadRequest {
			return fmt.Errorf("expected status %d but got %d, response message: %s", expectedStatus, resp.status, resp.body)
		}
		return fmt.Errorf("expected status %d but got %d", expectedStatus, resp.status)
	}

	return nil
}

func (a *apiFeature) theResponsePayloadShouldMatchJson(ctx context.Context, expectedBody *godog.DocString) error {
	actualResp, ok := ctx.Value(godogsResponseCtxKey{}).(response)
	if !ok {
		return errors.New("there are no godogs available in the context")
	}

	books := make([]models.Book, 0)

	err := json.Unmarshal([]byte(expectedBody.Content), &books)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(actualResp.body, books) {
		return fmt.Errorf("expected JSON dose not match actual,  %v vs. %v", actualResp.body, books)
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	api := &apiFeature{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		api.resetResponse(sc)
		return ctx, nil
	})

	ctx.Step(`^I send "([^"]*)" request to "([^"]*)" with payload:$`, api.iSendRequestToWithPayload)
	ctx.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
	ctx.Step(`^the response payload should match json:$`, api.theResponsePayloadShouldMatchJson)
}
