// Package middleware provides echo request and response output log
package middleware

import (
	"github.com/sirupsen/logrus"
	"time"

	"github.com/labstack/echo/v4"
)

// Logger returns a middleware that logs HTTP requests.
func Logger() echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {

			req := c.Request()
			res := c.Response()
			start := time.Now()
			startMicro := start.UnixMicro()

			var err error
			if err = next(c); err != nil {

				c.Error(err)
			}

			stop := time.Now()
			stopMicro := stop.UnixMicro()

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {

				id = res.Header().Get(echo.HeaderXRequestID)
			}

			reqSize := req.Header.Get(echo.HeaderContentLength)
			if reqSize == "" {

				reqSize = "0"
			}

			traceID := req.Header.Get("trace-id")
			if traceID == "" {

				traceID = "0"
			}

			logrus.WithFields(logrus.Fields{
				"service":   "api",
				"id":   id,
				"real_ip": c.RealIP(),
				"time": stop.Format(time.RFC3339),
				"host":req.Host,
				"method": req.Method,
				"request_uri": req.RequestURI,
				"status": res.Status,
				"request_size": reqSize,
				"referer": req.Referer(),
				"user_agent": req.UserAgent(),
				"response_time_millisecond": stopMicro - startMicro,
				"trace-id": traceID,
			}).Infof("API Response")

			return err
		}
	}
}
