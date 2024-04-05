package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	dir := os.Args[1]
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/:cid", func(c echo.Context) error {
		cid := c.Param("cid")

		return c.File(filepath.Join(dir, cid))
	})
	e.POST("/:cid", func(c echo.Context) error {
		cid := c.Param("cid")
		if strings.Contains(cid, "/") {
			return fmt.Errorf("must only be a single path element name")
		}

		b, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}

		fi, err := os.CreateTemp(dir, cid)
		if err != nil {
			return err
		}

		_, err = fi.Write(b)
		if err != nil {
			return err
		}

		fi.Close()

		return os.Rename(fi.Name(), filepath.Join(dir, cid))
	})

	e.Logger.Fatal(e.Start(":1323"))
}
