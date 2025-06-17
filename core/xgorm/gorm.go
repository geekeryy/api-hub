package xgorm

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/SpectatorNan/gorm-zero/gormc/config"
	"github.com/SpectatorNan/gorm-zero/gormc/config/mysql"
	"github.com/SpectatorNan/gorm-zero/gormc/config/pg"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	gormiologger "gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

// Colors
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

type logger struct {
	gormiologger.Interface
	gormiologger.Config
	gormiologger.Writer
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

const (
	TraceStr     = "%s  [%.3fms] [rows:%v] %s"
	TraceWarnStr = "%s %s [%.3fms] [rows:%v] %s"
	TraceErrStr  = "%s %s [%.3fms] [rows:%v] %s"
)

func ConnectMysql(m mysql.Mysql) (*gorm.DB, error) {
	db, err := mysql.ConnectWithConfig(m, &gorm.Config{
		Logger: NewGormLogger(&m),
	})
	if err != nil {
		return nil, err
	}
	if err := db.Use(tracing.NewPlugin()); err != nil {
		return nil, err
	}
	return db, nil
}

func ConnectPg(m pg.PgSql) (*gorm.DB, error) {
	db, err := pg.ConnectWithConfig(m, &gorm.Config{
		Logger: NewGormLogger(&m),
	})
	if err != nil {
		return nil, err
	}
	if err := db.Use(tracing.NewPlugin(tracing.WithoutServerAddress())); err != nil {
		return nil, err
	}
	return db, nil
}

func NewGormLogger(cfg config.GormLogConfigI) gormiologger.Interface {
	if cfg.GetGormLogMode() == gormiologger.Error {
		return gormiologger.Default
	}
	loggerConfig := gormiologger.Config{
		SlowThreshold:             cfg.GetSlowThreshold(), // 慢 SQL 阈值
		LogLevel:                  cfg.GetGormLogMode(),   // 日志级别
		IgnoreRecordNotFoundError: true,                   // 忽略ErrRecordNotFound（记录未找到）错误
		Colorful:                  cfg.GetColorful(),      // 禁用彩色打印
	}

	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if loggerConfig.Colorful {
		infoStr = Green + "%s\n" + Reset + Green + "[info] " + Reset
		warnStr = BlueBold + "%s\n" + Reset + Magenta + "[warn] " + Reset
		errStr = Magenta + "%s\n" + Reset + Red + "[error] " + Reset
		traceStr = Green + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
		traceWarnStr = Green + "%s " + Yellow + "%s\n" + Reset + RedBold + "[%.3fms] " + Yellow + "[rows:%v]" + Magenta + " %s" + Reset
		traceErrStr = RedBold + "%s " + MagentaBold + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
	}
	writer := log.New(os.Stderr, "\r\n", log.LstdFlags)
	newLogger := gormiologger.New(
		writer, // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		loggerConfig,
	)

	return logger{
		Interface:    newLogger,
		Config:       loggerConfig,
		Writer:       writer,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

// Info print info
func (l logger) Info(ctx context.Context, msg string, data ...any) {
	l.Interface.Info(ctx, msg, data...)
}

// Warn print warn messages
func (l logger) Warn(ctx context.Context, msg string, data ...any) {
	l.Interface.Warn(ctx, msg, data...)
}

// Error print error messages
func (l logger) Error(ctx context.Context, msg string, data ...any) {
	l.Interface.Error(ctx, msg, data...)
}

// Trace print sql message
func (l logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormiologger.Silent {
		return
	}
	lines := fmt.Sprint(FileWithLineNum()...)
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormiologger.Error && (!errors.Is(err, gormiologger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.Printf(l.traceErrStr, lines, err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			logx.WithContext(ctx).Errorf(TraceErrStr, lines, err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceErrStr, lines, err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			logx.WithContext(ctx).Errorf(TraceErrStr, lines, err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormiologger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.Printf(l.traceWarnStr, lines, slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			logx.WithContext(ctx).Infof(TraceWarnStr, lines, slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceWarnStr, lines, slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			logx.WithContext(ctx).Infof(TraceWarnStr, lines, slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == gormiologger.Info:
		sql, rows := fc()
		if rows == -1 {
			l.Printf(l.traceStr, lines, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			logx.WithContext(ctx).Infof(TraceStr, lines, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceStr, lines, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			logx.WithContext(ctx).Infof(TraceStr, lines, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

// Trace print sql message
func (l logger) ParamsFilter(ctx context.Context, sql string, params ...any) (string, []any) {
	if l.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}

var gormSourceDir string

func init() {
	_, file, _, _ := runtime.Caller(0)
	// compatible solution to get gorm source directory with various operating systems
	gormSourceDir = sourceDir(file)
}

func sourceDir(file string) string {
	dir := filepath.Dir(file)
	dir = filepath.Dir(dir)

	s := filepath.Dir(dir)
	if filepath.Base(s) != "gorm.io" {
		s = dir
	}
	return filepath.ToSlash(s) + "/"
}

// FileWithLineNum return the file name and line number of the current file
func FileWithLineNum() []any {
	result := make([]any, 0)
	p := filepath.ToSlash(filepath.Dir(filepath.Dir(gormSourceDir))) + "/"
	// the second caller usually from gorm internal, so set i start from 2
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!strings.HasPrefix(file, gormSourceDir) || strings.HasSuffix(file, "_test.go")) && strings.HasPrefix(file, p) && len(result) < 2 {
			result = append(result, file+":"+strconv.FormatInt(int64(line), 10)+" ")
		}
	}
	if len(result) == 2 {
		result[0], result[1] = result[1], result[0]
	}
	return result
}
