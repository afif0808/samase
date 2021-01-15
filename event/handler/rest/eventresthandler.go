package eventresthandler

import (
	"database/sql"
	"log"
	"net/http"
	"samase/event"
	eventsqlrepository "samase/event/repository/sql"
	eventservice "samase/event/service"
	"strconv"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func InjectEventRESTHandler(conn *sql.DB, gormDB *gorm.DB, ee *echo.Echo) {
	getEvents := eventservice.GetEvents(eventsqlrepository.GetEvents(conn))
	getOngoingEvents := eventservice.GetOngoingEvents(eventsqlrepository.GetEvents(conn))
	ee.GET("/events", GetEvents(getEvents))
	ee.GET("/events/ongoing/", GetOngoingEvents(getOngoingEvents))
	createEvent := eventservice.CreateEvent(eventsqlrepository.CreateEvent(gormDB))
	ee.POST("/events", CreateEvent(createEvent))
	updateEvents := eventsqlrepository.UpdateEvents(gormDB)
	updateEventByID := eventservice.UpdateEventByID(updateEvents)
	ee.POST("/events/:id", UpdateEventByID(updateEventByID))

	deleteEvents := eventsqlrepository.DeleteEvents(gormDB)
	deleteEventByID := eventservice.DeleteEventByID(deleteEvents)
	ee.DELETE("/events/:id", DeleteEventByID(deleteEventByID))

}

func GetEvents(getEvents eventservice.GetEventsFunc) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		keyword := ectx.Request().URL.Query().Get("q")
		evs, err := getEvents(ctx, keyword)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusOK, evs)
	}
}
func GetOngoingEvents(getOngoingEvents eventservice.GetOngoingEventsFunc) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		keyword := ectx.Request().URL.Query().Get("q")
		evs, err := getOngoingEvents(ctx, keyword)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusOK, evs)
	}
}

func CreateEvent(createEvent eventservice.CreateEventFunc) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		var post struct {
			Event event.Event `json:"event"`
		}
		err := ectx.Bind(&post)
		log.Println(err, post)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		ev, err := createEvent(ctx, post.Event)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusCreated, ev)
	}
}

func UpdateEventByID(updateEvent eventservice.UpdateEventByIDFunc) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		id, err := strconv.ParseInt(ectx.Param("id"), 10, 64)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		var post struct {
			Event event.Event `json:"event"`
		}
		err = ectx.Bind(&post)
		post.Event.ID = id
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		err = updateEvent(ctx, post.Event)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusOK, nil)
	}
}

func DeleteEventByID(deleteEvent eventservice.DeleteEventByIDFunc) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		id, err := strconv.ParseInt(ectx.Param("id"), 10, 64)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, nil)
		}
		err = deleteEvent(ctx, id)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusOK, nil)
	}
}
