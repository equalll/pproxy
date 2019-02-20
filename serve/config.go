package serve
import "github.com/equalll/mydebug"

import (
	"crypto/tls"
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/hidu/goutils"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Config  pproxy's config
type Config struct {
	Port         int
	AdminPort    int
	Title        string
	Notice       string
	AuthType     int
	DataDir      string
	FileDir      string
	ResponseSave int
	SessionView  int
	DataStoreDay float64
	ParentProxy  *url.URL

	SslOn bool

	SslCert tls.Certificate

	ModifyRequest bool
}

const (
	authTypeNO           = 0
	authTypeBasic        = 1
	authTypeBasicWithAny = 2
	authTypeBasicTry     = 3

	responseSaveAll      = 0
	responseSaveHasBroad = 1 //has show

	sessionViewALL      = 0
	sessionViewIPOrUser = 1
)

// User user struct
type User struct {
	Name         string
	Psw          string
	PswMd5       string
	IsAdmin      bool
	SkipCheckPsw bool
}

// String string format
func (u *User) String() string {mydebug.INFO()
	return fmt.Sprintf("Name:%s,Psw:%s,isAdmin:%v,SkipCheckPsw:%v", u.Name, u.Psw, u.IsAdmin, u.SkipCheckPsw)
}

// ConfigString one line in file
func (u *User) ConfigString() string {mydebug.INFO()
	return fmt.Sprintf("name:%s\tpsw:%s\tis_admin:%v\tpsw_md5:%s", u.Name, u.Psw, u.IsAdmin, u.PswMd5)
}

const (
	contentEncoding = "Content-Encoding"
)

//"0:no auth | 1:basic auth | 2:basic auth with any name"

// GetVersion get current version
func GetVersion() string {mydebug.INFO()
	return Assest.GetContent("res/version")
}

// GetDemoConf get the demo config
func GetDemoConf() string {mydebug.INFO()
	return strings.TrimSpace(Assest.GetContent("res/conf/demo.conf"))
}

func (u *User) isPswEq(psw string) bool {mydebug.INFO()
	return u.PswMd5 == utils.StrMd5(psw)
}

// LoadConfig load the pproxy's config
func LoadConfig(confPath string) (*Config, error) {mydebug.INFO()
	gconf, err := goconfig.LoadConfigFile(confPath)
	if err != nil {
		log.Println("load config", confPath, "failed,err:", err)
		return nil, err
	}
	config := new(Config)
	config.Port = gconf.MustInt(goconfig.DEFAULT_SECTION, "port", 8080)
	config.AdminPort = gconf.MustInt(goconfig.DEFAULT_SECTION, "adminPort", 0)

	if config.AdminPort == 0 {
		config.AdminPort = config.Port
	}

	config.DataStoreDay = gconf.MustFloat64(goconfig.DEFAULT_SECTION, "dataStoreDay", 0)
	if config.DataStoreDay < 0 {
		log.Println("wrong DataStoreDay,skip")
		config.DataStoreDay = 0
	}

	config.Title = gconf.MustValue(goconfig.DEFAULT_SECTION, "title")
	config.Notice = gconf.MustValue(goconfig.DEFAULT_SECTION, "notice")
	config.DataDir = gconf.MustValue(goconfig.DEFAULT_SECTION, "dataDir", "../data/")

	config.FileDir = gconf.MustValue(goconfig.DEFAULT_SECTION, "fileDir", "../file/")

	_authType := strings.ToLower(gconf.MustValue(goconfig.DEFAULT_SECTION, "authType", "none"))
	authTypes := map[string]int{"none": 0, "basic": 1, "basic_any": 2, "basic_try": 3, "try_basic": 3}

	hasError := false
	if authType, has := authTypes[_authType]; has {
		config.AuthType = authType
	} else {
		hasError = true
		log.Println("conf error,unknow value authType:", _authType)
	}

	_responseSave := strings.ToLower(gconf.MustValue(goconfig.DEFAULT_SECTION, "responseSave", "all"))
	responseSaveMap := map[string]int{"all": 0, "only_broadcast": 1}

	if responseSave, has := responseSaveMap[_responseSave]; has {
		config.ResponseSave = responseSave
	} else {
		hasError = true
		log.Println("conf error,unknow value responseSave:", _authType)
	}

	_sessionView := strings.ToLower(gconf.MustValue(goconfig.DEFAULT_SECTION, "sessionView", "all"))
	sessionViewMap := map[string]int{"all": 0, "ip_or_user": 1}

	if sessionView, has := sessionViewMap[_sessionView]; has {
		config.SessionView = sessionView
	} else {
		hasError = true
		log.Println("conf error,unknow value responseSave:", _authType)
	}

	parentProxy := gconf.MustValue(goconfig.DEFAULT_SECTION, "parentProxy", "")
	if parentProxy != "" {
		_urlObj, err := url.Parse(parentProxy)
		if err != nil || _urlObj.Scheme != "http" {
			hasError = true
			log.Println("parentProxy wrong,must http proxy")
		} else {
			config.ParentProxy = _urlObj
		}
	}
	config.SslOn = gconf.MustValue(goconfig.DEFAULT_SECTION, "ssl", "off") == "on"
	if config.SslOn {
		_sslClientCert := gconf.MustValue(goconfig.DEFAULT_SECTION, "ssl_client_cert", "")
		_sslServerKey := gconf.MustValue(goconfig.DEFAULT_SECTION, "ssl_server_key", "")
		cert, err := getSslCert(_sslClientCert, _sslServerKey)
		if err != nil {
			hasError = true
			log.Println("ssl ca config error:", err)
		} else {
			config.SslCert = cert
		}
	}

	config.ModifyRequest = gconf.MustValue(goconfig.DEFAULT_SECTION, "modifyRequest", "on") == "on"

	if hasError {
		return config, fmt.Errorf("config error")
	}

	return config, nil
}

type configHosts map[string]string

//  loadHosts 读取host配置文件
func loadHosts(confPath string) (hosts configHosts, err error) {mydebug.INFO()
	hosts = make(configHosts)
	if !utils.File_exists(confPath) {
		return
	}
	hostsByte, err := utils.File_get_contents(confPath)
	if err != nil {
		log.Println("load hosts_file failed:", confPath, err)
		return nil, err
	}
	hostsArr := utils.LoadText2Slice(string(hostsByte))
	for _, v := range hostsArr {
		if len(v) != 2 {
			log.Println("hosts file line wrong,ignore,", v)
			continue
		}
		hosts[v[0]] = v[1]
	}
	return
}

func loadUsers(confPath string) (users map[string]*User, err error) {mydebug.INFO()
	users = make(map[string]*User)
	if !utils.File_exists(confPath) {
		return
	}
	userInfoByte, err := utils.File_get_contents(confPath)
	if err != nil {
		log.Println("load user file failed:", confPath, err)
		return
	}
	lines := utils.LoadText2SliceMap(string(userInfoByte))
	for _, line := range lines {
		name, has := line["name"]
		if !has || name == "" {
			continue
		}
		if _, has := users[name]; has {
			log.Println("dup name in users:", name, line)
			continue
		}

		user := new(User)
		user.Name = name
		if val, has := line["is_admin"]; has && (val == "admin" || val == "true") {
			user.IsAdmin = true
		}
		if val, has := line["psw_md5"]; has {
			user.PswMd5 = val
		}

		if user.PswMd5 == "" {
			if val, has := line["psw"]; has {
				user.Psw = val
				user.PswMd5 = utils.StrMd5(val)
			}
		}
		users[user.Name] = user
	}
	return
}

func (config *Config) getTransport() *http.Transport {mydebug.INFO()
	if config.ParentProxy == nil {
		return nil
	}
	tr := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {

			if config.ParentProxy.User.Username() == "pass" {
				user := getAuthorInfo(req)
				urlTmp, err := url.Parse(config.ParentProxy.String())
				if err != nil {
					return nil, err
				}
				urlTmp.User = url.UserPassword(user.Name, user.Psw)
				return urlTmp, nil
			}
			return config.ParentProxy, nil
		},
	}
	return tr
}
