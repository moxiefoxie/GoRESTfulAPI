package main

import (
	"fmt"
	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"net/http"
	"time"
)

type WriteSyncer struct {
	io.Writer
}

func (ws WriteSyncer) Sync() error {
	return nil
}
func main() {
	var sqlManager *GoSqlManager
	sqlManager = new(GoSqlManager)
	sqlManager.OpenConnection("brown14", "savvytest")
	currentTime := time.Now()
	m := make(map[string]string)
	m["message"] = ("Current time: " + currentTime.String())
	m["user"] = "Savvy"
	rec1 := Record{IP: "192.00.00.1", Status: "200"}
	rec2 := Record{IP: "192.00.00.2", Status: "404"}
	var recs []Record
	recs = append(recs, rec1, rec2)
	var records Records
	records = recs
	sqlManager.BulkInsert(records)
	//sqlManager.InsertInto("savvytest.dbo.HaProxy_Test", m)
	testing, err := sqlManager.ExecuteQuery("select top(10) * from savvytest.dbo.HaProxy_Test;")
	for testing.Next() {
		var message string
		var user string
		err := testing.Scan(&message, &user)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(message, user)
	}

	e := echo.New()
	var logName = "http_access.log"
	cfg := zap.NewProductionConfig()
	cfg.Encoding = "json"
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.OutputPaths = []string{logName}
	zapLogger, err := cfg.Build()
	if err == nil {
		fmt.Print("ERROR:: ")
	}
	e.Use(echozap.ZapLogger(zapLogger))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))

}
