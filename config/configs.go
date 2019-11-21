package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql" // xorm will use this driver
	"github.com/go-xorm/xorm"
	"github.com/lucifinil-long/stores/utils"
	_ "github.com/mattn/go-sqlite3" // xorm will use this driver
	"github.com/mkideal/log"
)

// EnumUserAuthType is user auth enum type
type EnumUserAuthType int

const (
	// UserAuthTypeNone indicates not nee auth
	UserAuthTypeNone EnumUserAuthType = iota
	// UserAuthTypeLogin indicates need verify user information when logins
	UserAuthTypeLogin
	// UserAuthTypeRealtime indicates need verify user information realtime
	UserAuthTypeRealtime
)

const (
	cDBDriver           = "db_driver"
	cMysqlDBDriver      = "mysql"
	cSqliteDBDriver     = "sqlite3"
	cDBHost             = "db_host"
	cDBUser             = "db_user"
	cDBUserPwd          = "db_pwd"
	cDBName             = "db_name"
	cLogPath            = "log_path"
	cLogLevel           = "log_level"
	cNotAuthPackage     = "not_auth_package"
	cUserAuthType       = "user_auth_type"
	cAuthGateway        = "auth_gateway"
	cHomepage           = "homepage"
	cSuperAdminUsers    = "super_admin_users"
	cSysRefreshInterval = "sys_refresh_interval"
)

const (
	cDefaultGateway  = "/pages/admin/login"
	cDefaultHomepage = "/pages/admin/index"
	cDefaultSep      = ","
)

const (
	cCurrentDirPrefix = "./"
	cParentDirPrefix  = "../"
)

var (
	configsInstance *Configs
	once            sync.Once
)

// InitConfigs initializes configs and log environment
func InitConfigs() {
	once.Do(func() {
		configsInstance = newConfigs()
	})
}

// GetConfigs returns the singleton instance of configs
func GetConfigs() *Configs {
	// ensure that configs is always initialized
	InitConfigs()

	return configsInstance
}

func initLog(cfg *Configs) error {
	logPath := beego.AppConfig.String(cLogPath)
	var appPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	if len(logPath) == 0 {
		logPath = filepath.Join(appPath, "logs")
	} else if strings.HasPrefix(logPath, cCurrentDirPrefix) || strings.HasPrefix(logPath, cParentDirPrefix) {
		logPath = filepath.Join(appPath, logPath)
	}

	if err := log.Init("multifile/console", log.M{
		"rootdir":     logPath,
		"suffix":      ".txt",
		"date_format": "%04d-%02d-%02d",
	}); err != nil {
		return err
	}

	logLevel := beego.AppConfig.String(cLogLevel)
	log.Info("log level: %v, wanted level: %v", log.SetLevelFromString(logLevel), logLevel)
	return nil
}

func newConfigs() *Configs {
	cfg := &Configs{}

	// init log first
	err := initLog(cfg)
	if err != nil {
		log.Fatal("%v", err)
	}

	// init orm engine
	err = cfg.initOrmEngine()
	if err != nil {
		log.Fatal("%v", err)
	}

	return cfg
}

// Configs should be used as singletone model
type Configs struct {
	OrmEngine *xorm.Engine // database orm engine
}

func (cfg *Configs) initOrmEngine() error {
	var err error
	dbDriver := beego.AppConfig.String(cDBDriver)
	dbHost := beego.AppConfig.String(cDBHost)
	dbUser := beego.AppConfig.String(cDBUser)
	encryptedDbUserPwd := beego.AppConfig.String(cDBUserPwd)
	dbUserPwd := utils.RC4Base64Descrypt(encryptedDbUserPwd, utils.DefaultEncryptKey)
	dbName := beego.AppConfig.String(cDBName)

	if strings.EqualFold(dbDriver, cMysqlDBDriver) {
		dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true", dbUser, dbUserPwd, dbHost, dbName)
		if cfg.OrmEngine, err = xorm.NewEngine(dbDriver, dataSourceName); err != nil {
			return err
		}
	} else if strings.EqualFold(dbDriver, cSqliteDBDriver) {
		if cfg.OrmEngine, err = xorm.NewEngine(dbDriver, dbName); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("not supportted db driver: %v", dbDriver)
	}

	log.Info("initOrmEngine done successfully.")
	return nil
}

// GetConfigValue get the value of specified field
// @param field is the specified field
// @param refresh indicates whether reload config file before read
// @return value if found; otherwise return empty string
func (cfg Configs) GetConfigValue(field string, refresh bool) string {
	if refresh {
		var appPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
		var appConfigPath = filepath.Join(appPath, "conf", "app.conf")
		beego.LoadAppConfig("ini", appConfigPath)
	}

	return beego.AppConfig.String(field)
}

// NotAuthPackages get packages list that not need auth
// @return not auth packages list
func (cfg Configs) NotAuthPackages() []string {
	return []string{"pages", "index", "public", "static", "system"}
}

// UserAuthType get user auth type
// @return user auth type
// Note that this might be old data if config file is changed
func (cfg Configs) UserAuthType() EnumUserAuthType {
	val := cfg.GetConfigValue(cUserAuthType, false)
	ret, _ := strconv.Atoi(val)
	return EnumUserAuthType(ret)
}

// AuthGateway get auth gateway path
// @return auth gateway path
// Note that this might be old data if config file is changed
func (cfg Configs) AuthGateway() string {
	val := cfg.GetConfigValue(cAuthGateway, false)
	if len(val) == 0 {
		val = cDefaultGateway
	}

	return val
}

// AdminHomepage get homepage path
// @return homepage path
// Note that this might be old data if config file is changed
func (cfg Configs) AdminHomepage() string {
	val := cfg.GetConfigValue(cHomepage, false)
	if len(val) == 0 {
		val = cDefaultHomepage
	}

	return val
}

// SuperAdminList get super admin user name list
// @return super user name list
// Note that this might be old data if config file is changed
func (cfg Configs) SuperAdminList() []string {
	val := cfg.GetConfigValue(cSuperAdminUsers, false)
	return strings.Split(val, cDefaultSep)
}

// CollectSysInfoInterval get interval of collection
func (cfg Configs) CollectSysInfoInterval() int {
	return beego.AppConfig.DefaultInt(cSysRefreshInterval, 1)
}
