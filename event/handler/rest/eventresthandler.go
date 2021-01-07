package eventresthandler

import (
	"database/sql"
	"net/http"
	eventsqlrepository "samase/event/repository/sql"
	eventservice "samase/event/service"

	"github.com/labstack/echo"
)

func InjectEventRESTHandler(conn *sql.DB, ee *echo.Echo) {
	getEvents := eventservice.GetEvents(eventsqlrepository.GetEvents(conn))
	getOngoingEvents := eventservice.GetOngoingEvents(eventsqlrepository.GetEvents(conn))
	ee.GET("/events", GetEvents(getEvents))
	ee.GET("/events/ongoing/", GetOngoingEvents(getOngoingEvents))

}

func GetEvents(getEvents eventservice.GetEventsFunc) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		evs, err := getEvents(ctx)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusOK, evs)
	}
}
func GetOngoingEvents(getOngoingEvents eventservice.GetOngoingEventsFunc) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		evs, err := getOngoingEvents(ctx)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusOK, evs)
	}
}
