package userpointsqlrepo

import (
	"context"
	"database/sql"
	"log"
	"samase/userpoint"
	userpointrepo "samase/userpoint/repository"
)

const (
	userPointTable       = "user_point"
	createUserPointQuery = "INSERT " + userPointTable + " SET user_point.user_id = ? AND user_point.value = ?"
)

func CreateUserPoint(conn *sql.DB) userpointrepo.CreateUserPointFunc {
	return func(ctx context.Context, uspo userpoint.UserPoint) (userpoint.UserPoint, error) {
		_, err := conn.ExecContext(ctx, createUserPointQuery, uspo.UserID, uspo.Value)
		if err != nil {
			log.Println(err)
		}
		return uspo, err
	}
}
