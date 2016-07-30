package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	// "github.com/jinzhu/gorm"

	"github.com/hobo-go/echo-mw/session"

	"github.com/hobo-go/echo-web/models"
	"github.com/hobo-go/echo-web/modules/cache"
	"github.com/hobo-go/echo-web/modules/log"
)

func ApiHandler(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)

	u := &models.User{}
	if err != nil {
		log.DebugPrint("Render Error: %v", err)
	} else {
		var User models.User
		u = User.GetUserById(id)
	}

	// 缓存测试
	value := -1
	if err == nil {
		cacheStore := cache.Default(c)
		if id == 1 {
			value = 0
			cacheStore.Set("userId", 1, time.Minute)
		} else {
			cacheStore.Get("userId", &value)
		}
	}

	// Flash测试
	s := session.Default(c)
	s.AddFlash("0")
	s.AddFlash("1")
	s.AddFlash("10", "key1")
	s.AddFlash("20", "key2")
	s.AddFlash("21", "key2")

	request := c.Request()
	c.JSON(http.StatusOK, map[string]interface{}{
		"title":        "Api Index",
		"User":         u,
		"CacheValue":   value,
		"Scheme":       request.Scheme(),
		"Host":         request.Host(),
		"UserAgent":    request.UserAgent(),
		"Method":       request.Method(),
		"URI":          request.URI(),
		"RemoteAddr":   request.RemoteAddress(),
		"Path":         request.URL().Path(),
		"QueryString":  request.URL().QueryString(),
		"QueryParams":  request.URL().QueryParams(),
		"HeaderKeys":   request.Header().Keys(),
		"FlashDefault": s.Flashes(),
		"Flash1":       s.Flashes("key1"),
		"Flash2":       s.Flashes("key2"),
	})

	return nil
}

func JETTesterHandler(c echo.Context) error {
	t, err := getJETToken()
	if err != nil {
		return err
	}
	c.Set("tmpl", "api/jwt_tester")
	c.Set("data", map[string]interface{}{
		"title": "JWT 接口测试",
		"token": t,
	})

	return nil
}
