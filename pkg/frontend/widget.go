package frontend

import (
	"embed"
	"fmt"
	"html/template"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/vmkteam/embedlog"
)

type WidgetManager struct {
	embedlog.Logger

	indexTmpl *template.Template
}

func NewWidgetManager(logger embedlog.Logger) *WidgetManager {
	return &WidgetManager{Logger: logger}
}

//go:embed index.html
var f embed.FS

var funcMap = template.FuncMap{}

func (wm *WidgetManager) Init() error {
	// index.html
	srvBytes, err := f.ReadFile("index.html")
	if err != nil {
		return fmt.Errorf("read index.html err=%w", err)
	}

	indexTmpl, err := template.New("index").Funcs(funcMap).Parse(string(srvBytes))
	if err != nil {
		return fmt.Errorf("parse index.html err=%w", err)
	}

	wm.indexTmpl = indexTmpl

	return nil
}

type Data struct {
	Time time.Time
}

func (wm *WidgetManager) MainHandler(c echo.Context) error {
	data := Data{Time: time.Now()}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	err := wm.indexTmpl.Execute(c.Response().Writer, &data)
	if err != nil {
		wm.Error(c.Request().Context(), "render widget failed", "err", err)
		return err
	}

	return nil
}
