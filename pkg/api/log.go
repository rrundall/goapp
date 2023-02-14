package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"goapp/config"
	"io"
	"log"
	"os"
	"time"
)

// StartLogging is for initialize the zerolog that will log into the local log file path
// By default will only log on error level if you want debug level then debug flag need to be true
func (s *Server) StartLogging(debug *bool) {
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	f, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("%s %s", config.FailToSaveLogErrMsg, err)
	}
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	s.router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] - %s \"%s %s %s %d %s \"%s\" %s\"\n",
			param.TimeStamp.Format(time.RFC1123),
			param.ClientIP,
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage)
	}))
}
