package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseIdParam(ctx *gin.Context) (int, error) {
	param := ctx.Param("id")

	id, err := strconv.Atoi(param)

	if err != nil {
		return 0, err
	}

	return id, nil

}
