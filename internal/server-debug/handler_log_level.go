package serverdebug

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/razielsd/chat-service/internal/logger"
)

func (s *Server) handlerLogLevel(eCtx echo.Context) error {
	level := eCtx.FormValue("level")
	s.lg.Info("change log level", zap.String("level", level))
	err := logger.LogLevel.UnmarshalText([]byte(level))
	if err != nil {
		return err
	}
	return eCtx.String(http.StatusOK, level)
}
