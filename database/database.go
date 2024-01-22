package db

import (
	"fmt"
	misakaLogger "github.com/misakacoder/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"haibara/config"
	"haibara/util"
	"strings"
	"time"
)

const (
	urlFormat              = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	defaultMaxIdleConn     = 16
	defaultMaxOpenConn     = 32
	defaultConnMaxLifetime = time.Hour
	defaultSlowSqlTime     = 100 * time.Millisecond
)

var GORM = initGORM()

type sqlLogWriter struct{}

func (writer sqlLogWriter) Printf(format string, v ...any) {
	level := misakaLogger.DEBUG
	caller := v[0].(string)
	v = v[1:]
	index := strings.Index(format, "%s")
	format = strings.ReplaceAll(format, "\n", " ")[index+2:]
	format = strings.TrimPrefix(format, " ")
	switch tp := v[0].(type) {
	case error:
		level = misakaLogger.ERROR
	case string:
		if strings.Contains(tp, "SLOW SQL") {
			level = misakaLogger.WARN
		}
	}
	misakaLogger.GetLogger().(*misakaLogger.SimpleLogger).Push(level, caller, format, v...)
}

func (writer sqlLogWriter) Write(p []byte) (n int, err error) {
	message := string(p)
	messages := strings.Split(message, "\n")
	if len(messages) >= 2 {
		level := misakaLogger.DEBUG
		caller := messages[0]
		sql := messages[1]
		if strings.Contains(caller, "SLOW SQL") {
			level = misakaLogger.WARN
			caller, sql = getCallerAndSql(caller, sql)
		} else if strings.Contains(caller, "Error") {
			level = misakaLogger.ERROR
			caller, sql = getCallerAndSql(caller, sql)
		}
		misakaLogger.GetLogger().(*misakaLogger.SimpleLogger).Push(level, caller, sql)
	}
	return len(p), err
}

func getCallerAndSql(caller, sql string) (string, string) {
	index := strings.Index(caller, " ")
	sql = caller[index+1:] + " " + sql
	caller = caller[:index]
	return caller, sql
}

func initGORM() *gorm.DB {
	database := config.Configuration.Database
	host := database.Host
	port := database.Port
	username := database.Username
	password := database.Password
	name := database.Name
	maxIdleConn := database.MaxIdleConn
	maxOpenConn := database.MaxOpenConn
	if util.AllNotNull(host, port, username, password, name) {
		url := fmt.Sprintf(urlFormat, username, password, host, port, name)
		gormConfig := &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		}
		if database.PrintSql {
			slowSqlTime, err := time.ParseDuration(database.SlowSqlTime)
			if err != nil {
				slowSqlTime = defaultSlowSqlTime
			}
			writer := sqlLogWriter{}
			//writer := log.New(sqlLogWriter{}, "", log.Lmsgprefix)
			sqlLogger := logger.New(writer, logger.Config{
				SlowThreshold: slowSqlTime,
				Colorful:      true,
				LogLevel:      logger.Info,
			})
			gormConfig.Logger = sqlLogger
		}
		gormDB, err := gorm.Open(mysql.Open(url), gormConfig)
		if err != nil {
			misakaLogger.Panic(err.Error())
		}
		db, _ := gormDB.DB()
		db.SetMaxIdleConns(util.ConditionalExpression(maxIdleConn <= 0, defaultMaxIdleConn, maxIdleConn))
		db.SetMaxOpenConns(util.ConditionalExpression(maxOpenConn <= 0, defaultMaxOpenConn, maxOpenConn))
		connMaxLifeTime, err := time.ParseDuration(database.ConnMaxLifeTime)
		if err != nil {
			connMaxLifeTime = defaultConnMaxLifetime
		}
		db.SetConnMaxLifetime(connMaxLifeTime)
		return gormDB
	}
	misakaLogger.Panic("数据库连接失败，请检查配置文件")
	return nil
}
