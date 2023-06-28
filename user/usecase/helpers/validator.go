package helpers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

type CustomValidator struct {
	Validators *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.Validators.Struct(i)

	if err != nil {
		var sb strings.Builder
		sb.WriteString("Validation error:\n")

		for _, err := range err.(validator.ValidationErrors) {
			sb.WriteString(fmt.Sprintf("- %s\n", err))
		}

		return echo.NewHTTPError(http.StatusBadRequest, sb.String())
	}

	return nil
}
