package controllers

import (
	"axis/ecommerce-backend/configs"
	"axis/ecommerce-backend/pkg/entities"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func next(c *gin.Context, qp entities.QueryPathParam) *entities.RequestPathData {
	queries := c.Request.URL.Query()
	path := c.FullPath()
	rData := new(entities.RequestPathData)
	nextQueries := make(map[string][]string)
	previousQueries := make(map[string][]string)
	for i, s := range queries {
		if i == "limit" {
			convLimit := strconv.Itoa(qp.Limit)
			nextQueries[i] = []string{convLimit}
			previousQueries[i] = []string{convLimit}
			if qp.Limit == -1 {
				nextQueries["limit"] = []string{"-1"}
				nextQueries["offset"] = []string{"-1"}
				previousQueries["offset"] = []string{"-1"}
				previousQueries["limit"] = []string{"-1"}
			}
			continue
		}

		if i == "offset" {
			if qp.Limit == -1 {
				continue
			}
			nextOffset := qp.Offset + qp.Limit
			previousOffset := qp.Offset - qp.Limit
			if previousOffset < 0 {
				previousOffset = 0
			}

			nf := strconv.Itoa(nextOffset)
			pf := strconv.Itoa(previousOffset)
			nextQueries["offset"] = []string{nf}
			previousQueries["offset"] = []string{pf}
			continue
		}

		nextQueries[i] = s
		previousQueries[i] = s
	}

	if len(nextQueries) > 0 {
		nextFullUrl := configs.AppConfig.ApiUrl + path + "?"
		previousFullUrl := configs.AppConfig.ApiUrl + path + "?"
		tt := len(nextQueries)
		c := 1
		for s, ss := range nextQueries {
			sep := "&"
			if c == tt {
				sep = ""
			}
			nextFullUrl = nextFullUrl + s + "=" + strings.Join(ss, ",") + sep
			previousFullUrl = previousFullUrl + s + "=" + strings.Join(previousQueries[s], ",") + sep
			c += 1
		}

		rData.Next = nextFullUrl
		rData.Previous = previousFullUrl
	}

	return rData
}

func queryParams(c *gin.Context, defValues entities.QueryPathParam) (int, int, error) {
	limit := defValues.Limit
	offset := defValues.Offset
	l, ok := c.GetQuery("limit")
	if ok {
		lt, err := strconv.Atoi(l)
		if err != nil {
			return limit, offset, err
		}
		limit = lt
	}

	o, ok := c.GetQuery("offset")
	if ok {
		of, err := strconv.Atoi(o)
		if err != nil {
			return limit, offset, err
		}
		offset = of
	}
	if limit == -1 {
		return -1, -1, nil
	}

	return limit, offset, nil
}
