package rblogger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func New() *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	config.FunctionKey = "func"
	config.EncodeTime = zapcore.ISO8601TimeEncoder   // Format de la date ISO8601
	config.EncodeLevel = zapcore.CapitalLevelEncoder // Niveau de log en majuscules (INFO, ERROR, etc.)

	// Créer un encoder qui écrit en format texte
	encoder := zapcore.NewConsoleEncoder(config)

	// Créer un Core avec ce format et écrire dans la sortie standard (stdout)
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel)

	return zap.New(core)
}
