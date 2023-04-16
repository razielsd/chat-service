package serverdebug

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/razielsd/chat-service/internal/logger"
)

type applyLogLevelHandler struct {
	lg *zap.Logger
}

func newApplyLogLevelHandler(lg *zap.Logger) *applyLogLevelHandler {
	return &applyLogLevelHandler{
		lg: lg,
	}
}

func (l *applyLogLevelHandler) handler(eCtx echo.Context) error {
	level := eCtx.FormValue("level")
	l.lg.Info("change log level", zap.String("level", level))
	err := logger.LogLevel.UnmarshalText([]byte(level))
	if err != nil {
		return err
	}
	return eCtx.String(http.StatusOK, level)
}
