package zaplogger

import (
	"encoding/json"
	"net/http"
	"strings"
	"go.uber.org/zap/zapcore"
)

// fixme: unmarshal as map, to keep fields alive
type GrayLogger struct {
	Version      string `json:"version"`
	Host         string `json:"host"`
	ShortMessage string `json:"short_message"`
	FullMessage  string `json:"-"` // do not used
	Timestamp    int64  `json:"-"` // do not used
	Level        int    `json:"level"`
	Line         int    `json:"-"` // do not used
	File         string `json:"file"`
	LoggerLevel  string `json:"_logger_level"`
	LoggerTime   string `json:"server_time"`
}

type GrayLoggerWriter struct {
	url string
}

func (glw *GrayLoggerWriter) Write(logDetail []byte) (int, error) {

	gl := make(map[string]interface{}, 0)
	err := json.Unmarshal(logDetail, &gl)
	if err != nil {
		return 0, err
	}
	gl["version"] = "1.0"

	var lv zapcore.Level
	lvStr, ok := gl["_logger_level"].(string)
	if ok {
		lv.UnmarshalText([]byte(lvStr))
		gl["level"] = int(lv)
	}
	logDetail, _ = json.Marshal(&gl)
	//fmt.Println(string(logDetail))

	res, err := http.Post(glw.url, "Content-Type: application/json", strings.NewReader(string(logDetail)))
	if err != nil {
		return 0, err
	}
	//fmt.Println(res.Status)

	if res != nil && res.StatusCode == 404 {
		glw.url = "http://localhost:12201/gelf"
	}

	res.Body.Close()

	return 0, nil
}

func (*GrayLoggerWriter) getJsonEncoderConfig() zapcore.EncoderConfig {

	return zapcore.EncoderConfig{
		TimeKey:        "server_time",
		LevelKey:       "_logger_level",
		NameKey:        "host",
		CallerKey:      "file",
		MessageKey:     "short_message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
