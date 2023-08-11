package config

import (
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/CanftIn/gothafoss/pkg/util"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

type Mode string

const (
	DebugMode   Mode = "debug"
	ReleaseMode Mode = "release"
	BenchMode   Mode = "bench"
)

type FileService string

const (
	FileServiceAliyunOSS  FileService = "aliyunOSS"
	FileServiceTencentCOS FileService = "tencentCOS"
	FileServiceMinio      FileService = "minio"
)

func (u FileService) String() string {
	return string(u)
}

type TablePartitionConfig struct {
	MessageTableCount         int // 消息表数量
	MessageUserEditTableCount int // 用户消息编辑表
	ChannelOffsetTableCount   int // 频道偏移表
}

func newTablePartitionConfig() TablePartitionConfig {
	return TablePartitionConfig{
		MessageTableCount:         5,
		MessageUserEditTableCount: 3,
		ChannelOffsetTableCount:   3,
	}
}

type Config struct {
	vp *viper.Viper // 内部配置对象

	// ---------- 基础配置 ----------
	Mode                        Mode   // 模式 debug 测试 release 正式 bench 压力测试
	AppID                       string // APP ID
	AppName                     string // APP名称
	Version                     string // 版本
	RootDir                     string // 数据根目录
	Addr                        string // 服务监听地址 x.x.x.x:8080
	GRPCAddr                    string // grpc的通信地址 （建议内网通信）
	SSLAddr                     string // ssl 监听地址
	MessageSaveAcrossDevice     bool   // 消息是否跨设备保存（换设备登录消息是否还能同步到老消息）
	WelcomeMessage              string // 登录注册欢迎语
	PhoneSearchOff              bool   // 是否关闭手机号搜索
	OnlineStatusOn              bool   // 是否开启在线状态显示
	GroupUpgradeWhenMemberCount int    // 当成员数量大于此配置时 自动升级为超级群 默认为 1000
	EventPoolSize               int64  // 事件任务池大小

	// ---------- 外网配置 ----------
	External struct {
		IP          string // 外网IP
		BaseURL     string // 本服务的对外的基础地址
		H5BaseURL   string // h5页面的基地址 如果没有配置默认未 BaseURL + /web
		APIBaseURL  string // api的基地址 如果没有配置默认未 BaseURL + /v1
		WebLoginURL string // web登录地址
	}

	// ---------- 日志配置 ----------
	Logger struct {
		Dir     string // 日志存储目录
		Level   zapcore.Level
		LineNum bool // 是否显示代码行数
	}

	// ---------- db相关配置 ----------
	DB struct {
		MySQLAddr          string // mysql的连接信息
		Migration          bool   // 是否合并数据库
		RedisAddr          string // redis地址
		RedisPass          string // redis密码
		AsynctaskRedisAddr string // 异步任务的redis地址 不写默认为RedisAddr的地址
	}

	// ---------- 分布式配置 ----------
	Cluster struct {
		NodeID int //  节点ID 节点ID需要小于1024
	}

	// ---------- 缓存配置 ----------
	Cache struct {
		TokenCachePrefix            string        // token缓存前缀
		LoginDeviceCachePrefix      string        // 登录设备缓存前缀
		LoginDeviceCacheExpire      time.Duration // 登录设备缓存过期时间
		UIDTokenCachePrefix         string        // uidtoken缓存前缀
		FriendApplyTokenCachePrefix string        // 申请好友的token的前缀
		FriendApplyExpire           time.Duration // 好友申请过期时间
		TokenExpire                 time.Duration // token失效时间
		NameCacheExpire             time.Duration // 名字缓存过期时间
	}

	// ---------- 系统账户设置 ----------
	Account struct {
		SystemUID       string // 系统账号uid
		FileHelperUID   string // 文件助手uid
		SystemGroupID   string // 系统群ID 需求在app_config表里设置new_user_join_system_group为1才有效
		SystemGroupName string // 系统群的名字
		AdminUID        string // 系统管理员账号
	}

	// ---------- 文件服务 ----------
	FileService FileService   // 文件服务
	OSS         OSSConfig     // 阿里云oss配置
	Minio       MinioConfig   // minio配置
	Seaweed     SeaweedConfig // seaweedfs配置

	// ---------- 短信运营商 ----------
	SMSCode                string // 模拟的短信验证码
	SMSProvider            SMSProvider
	UniSMS                 UnismsConfig                 // unisms https://unisms.apistd.com/
	AliyunSMS              AliyunSMSConfig              // aliyun sms
	AliyunInternationalSMS AliyunInternationalSMSConfig // 阿里云国际短信

	// ---------- GothIM ----------
	GothIM struct {
		APIURL       string // im基地址
		ManagerToken string // im的管理者配置了就需要填写，没配置就不需要
	}

	// ---------- 头像 ----------
	Avatar struct {
		Default        string // 默认头像
		DefaultCount   int    // 默认头像数量
		Partition      int    // 头像分区数量
		DefaultBaseURL string // 默认头像的基地址
	}

	// ---------- 短编号 ----------
	ShortNo struct {
		NumOn   bool // 是否开启数字短编号
		NumLen  int  // 数字短编号长度
		EditOff bool // 是否关闭短编号编辑
	}

	// ---------- robot ----------
	Robot struct {
		MessageExpire      time.Duration // 消息过期时间
		InlineQueryTimeout time.Duration // inlineQuery事件过期时间
		EventPoolSize      int64         // 机器人事件池大小
	}

	// ---------- github ----------
	Github struct {
		OAuthURL     string // github oauth url
		ClientID     string // github client id
		ClientSecret string // github client secret
	}

	// ---------- owt ----------
	OWT struct {
		URL          string // owt api地址 例如： https://xx.xx.xx.xx:3000/v1
		ServiceID    string // owt的服务ID
		ServiceKey   string // owt的服务key （用户访问后台的api）
		RoomMaxCount int    // 房间最大参与人数
	}
	Register struct {
		Off           bool // 是否关闭注册
		OnlyChina     bool // 是否仅仅中国手机号可以注册
		StickerAddOff bool // 是否关闭注册添加表情
	}

	// ---------- push ----------
	Push struct {
		ContentDetailOn bool     // 推送是否显示正文详情(如果为false，则只显示“您有一条新的消息” 默认为true)
		PushPoolSize    int64    // 推送任务池大小
		APNS            APNSPush // 苹果推送
		MI              MIPush   // 小米推送
		HMS             HMSPush  // 华为推送
		VIVO            VIVOPush // vivo推送
		OPPO            OPPOPush // oppo推送
	}

	// ---------- wechat ----------
	Wechat struct {
		AppID     string // 微信appid 在开放平台内
		AppSecret string
	}

	// ---------- tracing ----------
	Tracing struct {
		On   bool   // 是否开启tracing
		Addr string // tracer的地址
	}

	// ---------- support ----------
	Support struct {
		Email     string // 技术支持的邮箱地址
		EmailSmtp string // 技术支持的邮箱的smtp
		EmailPwd  string // 邮箱密码
	}

	// ---------- 其他 ----------
	Test bool // 是否是测试模式

	QRCodeInfoURL    string   // 获取二维码信息的URL
	VisitorUIDPrefix string   // 访客uid的前缀
	TimingWheelTick  duration // The time-round training interval must be 1ms or more
	TimingWheelSize  int64    // Time wheel size

	// ---------- 系统配置  由系统生成,无需用户配置 ----------
	AppRSAPrivateKey string
	AppRSAPubKey     string
}

func New() *Config {
	cfg := &Config{
		// ---------- 基础配置 ----------
		Mode:                        ReleaseMode,
		AppID:                       "gothafoss",
		AppName:                     "众神瀑布",
		Addr:                        ":8090",
		GRPCAddr:                    "0.0.0.0:6979",
		PhoneSearchOff:              false,
		OnlineStatusOn:              true,
		GroupUpgradeWhenMemberCount: 1000,
		MessageSaveAcrossDevice:     true,
		EventPoolSize:               100,
		WelcomeMessage:              "欢迎使用{{appName}}",
		RootDir:                     "gothdata",

		// ---------- 外网配置 ----------
		External: struct {
			IP          string
			BaseURL     string
			H5BaseURL   string
			APIBaseURL  string
			WebLoginURL string
		}{
			BaseURL:     "",
			WebLoginURL: "",
		},

		// ---------- db配置 ----------
		DB: struct {
			MySQLAddr          string
			Migration          bool
			RedisAddr          string
			RedisPass          string
			AsynctaskRedisAddr string
		}{
			MySQLAddr: "root:demo@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=true",
			Migration: true,
			RedisAddr: "127.0.0.1:6379",
		},

		// ---------- 缓存配置 ----------
		Cache: struct {
			TokenCachePrefix            string
			LoginDeviceCachePrefix      string
			LoginDeviceCacheExpire      time.Duration
			UIDTokenCachePrefix         string
			FriendApplyTokenCachePrefix string
			FriendApplyExpire           time.Duration
			TokenExpire                 time.Duration
			NameCacheExpire             time.Duration
		}{
			TokenCachePrefix:            "token:",
			TokenExpire:                 time.Hour * 24 * 30,
			LoginDeviceCachePrefix:      "login_device:",
			LoginDeviceCacheExpire:      time.Minute * 5,
			UIDTokenCachePrefix:         "uidtoken:",
			FriendApplyTokenCachePrefix: "friend_token:",
			FriendApplyExpire:           time.Hour * 24 * 15,
			NameCacheExpire:             time.Hour * 24 * 7,
		},

		// ---------- 系统账户设置 ----------
		Account: struct {
			SystemUID       string
			FileHelperUID   string
			SystemGroupID   string
			SystemGroupName string
			AdminUID        string
		}{
			SystemUID:       "u_10000",
			SystemGroupID:   "g_10000",
			SystemGroupName: "意见反馈群",
			FileHelperUID:   "fileHelper",
			AdminUID:        "admin",
		},
		// ---------- 文件服务 ----------
		FileService: FileServiceMinio,

		// ---------- push  ----------
		Push: struct {
			ContentDetailOn bool
			PushPoolSize    int64
			APNS            APNSPush
			MI              MIPush
			HMS             HMSPush
			VIVO            VIVOPush
			OPPO            OPPOPush
		}{
			ContentDetailOn: true,
			PushPoolSize:    100,
			APNS: APNSPush{
				Dev:      true,
				Topic:    "com.xinbida.tangsengdaodao",
				Password: "123456",
			},
		},

		TimingWheelTick: duration{
			Duration: time.Millisecond * 10,
		},
		TimingWheelSize: 100,
	}
	return cfg
}

func (c *Config) ConfigFileUsed() string {
	return c.vp.ConfigFileUsed()
}

func (c *Config) configureLog() {
	logLevel := c.vp.GetInt("logger.level")
	// level
	if logLevel == 0 { // 没有设置
		if c.Mode == DebugMode {
			logLevel = int(zapcore.DebugLevel)
		} else {
			logLevel = int(zapcore.InfoLevel)
		}
	} else {
		logLevel = logLevel - 2
	}
	c.Logger.Level = zapcore.Level(logLevel)
	c.Logger.Dir = c.vp.GetString("logger.dir")
	if strings.TrimSpace(c.Logger.Dir) == "" {
		c.Logger.Dir = "logs"
	}
	if !strings.HasPrefix(strings.TrimSpace(c.Logger.Dir), "/") {
		c.Logger.Dir = filepath.Join(c.RootDir, c.Logger.Dir)
	}
	c.Logger.LineNum = c.vp.GetBool("logger.lineNum")
}

func (c *Config) getString(key string, defaultValue string) string {
	v := c.vp.GetString(key)
	if v == "" {
		return defaultValue
	}
	return v
}

func (c *Config) getBool(key string, defaultValue bool) bool {
	objV := c.vp.Get(key)
	if objV == nil {
		return defaultValue
	}
	return cast.ToBool(objV)
}
func (c *Config) getInt(key string, defaultValue int) int {
	v := c.vp.GetInt(key)
	if v == 0 {
		return defaultValue
	}
	return v
}

func (c *Config) getInt64(key string, defaultValue int64) int64 {
	v := c.vp.GetInt64(key)
	if v == 0 {
		return defaultValue
	}
	return v
}

func (c *Config) getDuration(key string, defaultValue time.Duration) time.Duration {
	v := c.vp.GetDuration(key)
	if v == 0 {
		return defaultValue
	}
	return v
}

// GetAvatarPath 获取用户头像path
func (c *Config) GetAvatarPath(uid string) string {
	return fmt.Sprintf("users/%s/avatar", uid)
}

// GetGroupAvatarFilePath 获取群头像上传路径
func (c *Config) GetGroupAvatarFilePath(groupNo string) string {
	avatarID := crc32.ChecksumIEEE([]byte(groupNo)) % uint32(c.Avatar.Partition)
	return fmt.Sprintf("group/%d/%s.png", avatarID, groupNo)
}

// GetCommunityAvatarFilePath 获取社区头像上传路径
func (c *Config) GetCommunityAvatarFilePath(communityNo string) string {
	avatarID := crc32.ChecksumIEEE([]byte(communityNo)) % uint32(c.Avatar.Partition)
	return fmt.Sprintf("community/%d/%s.png", avatarID, communityNo)
}

// GetCommunityCoverFilePath 获取社区封面上传路径
func (c *Config) GetCommunityCoverFilePath(communityNo string) string {
	avatarID := crc32.ChecksumIEEE([]byte(communityNo)) % uint32(c.Avatar.Partition)
	return fmt.Sprintf("community/%d/%s_cover.png", avatarID, communityNo)
}

// IsVisitorChannel 是访客频道
func (c *Config) IsVisitorChannel(uid string) bool {

	return strings.HasSuffix(uid, "@ht")
}

// 获取客服频道真实ID
func (c *Config) GetCustomerServiceChannelID(channelID string) (string, bool) {
	if !strings.Contains(channelID, "|") {
		return "", false
	}
	channelIDs := strings.Split(channelID, "|")
	return channelIDs[1], true
}

// 获取客服频道的访客id
func (c *Config) GetCustomerServiceVisitorUID(channelID string) (string, bool) {
	if !strings.Contains(channelID, "|") {
		return "", false
	}
	channelIDs := strings.Split(channelID, "|")
	return channelIDs[0], true
}

// 组合客服ID
func (c *Config) ComposeCustomerServiceChannelID(vid string, channelID string) string {
	return fmt.Sprintf("%s|%s", vid, channelID)
}

// IsVisitor 是访客uid
func (c *Config) IsVisitor(uid string) bool {

	return strings.HasPrefix(uid, c.VisitorUIDPrefix)
}

// GetEnv 成环境变量里获取
func GetEnv(key string, defaultValue string) string {
	v := os.Getenv(key)
	if strings.TrimSpace(v) == "" {
		return defaultValue
	}
	return v
}

// GetEnvBool 成环境变量里获取
func GetEnvBool(key string, defaultValue bool) bool {
	v := os.Getenv(key)
	if strings.TrimSpace(v) == "" {
		return defaultValue
	}
	if v == "true" {
		return true
	}
	return false
}

// GetEnvInt64 环境变量获取
func GetEnvInt64(key string, defaultValue int64) int64 {
	v := os.Getenv(key)
	if strings.TrimSpace(v) == "" {
		return defaultValue
	}
	i, _ := strconv.ParseInt(v, 10, 64)
	return i
}

// GetEnvInt 环境变量获取
func GetEnvInt(key string, defaultValue int) int {
	v := os.Getenv(key)
	if strings.TrimSpace(v) == "" {
		return defaultValue
	}
	i, _ := strconv.ParseInt(v, 10, 64)
	return int(i)
}

// GetEnvFloat64 环境变量获取
func GetEnvFloat64(key string, defaultValue float64) float64 {
	v := os.Getenv(key)
	if strings.TrimSpace(v) == "" {
		return defaultValue
	}
	i, _ := strconv.ParseFloat(v, 64)
	return i
}

// StringEnv StringEnv
func StringEnv(v *string, key string) {
	vv := os.Getenv(key)
	if vv != "" {
		*v = vv
	}
}

// BoolEnv 环境bool值
func BoolEnv(b *bool, key string) {
	value := os.Getenv(key)
	if strings.TrimSpace(value) != "" {
		if value == "true" {
			*b = true
		} else {
			*b = false
		}
	}
}

// 获取内网地址
func getIntranetIP() string {
	intranetIPs, err := util.GetIntranetIP()
	if err != nil {
		panic(err)
	}
	if len(intranetIPs) > 0 {
		return intranetIPs[0]
	}
	return ""
}

// SMSProvider 短信供应者
type SMSProvider string

const (
	// SMSProviderAliyun aliyun
	SMSProviderAliyun SMSProvider = "aliyun"
	SMSProviderUnisms SMSProvider = "unisms" // 联合短信(https://unisms.apistd.com/docs/api/send/)
)

// AliyunSMSConfig 阿里云短信
type AliyunSMSConfig struct {
	AccessKeyID  string // aliyun的AccessKeyID
	AccessSecret string // aliyun的AccessSecret
	TemplateCode string // aliyun的短信模版
	SignName     string // 签名
}

// aliyun oss
type OSSConfig struct {
	Endpoint        string
	BucketName      string // Bucket名称
	BucketURL       string // 文件下载地址域名 对应aliyun的Bucket域名
	AccessKeyID     string
	AccessKeySecret string
}

type MinioConfig struct {
	URL             string // 文件下载上传基地址 例如： http://127.0.0.1:9000
	AccessKeyID     string // minio accessKeyID
	SecretAccessKey string // minio secretAccessKey
}

type SeaweedConfig struct {
	URL string // 文件下载上传基地址
}

// UnismsConfig unisms短信
type UnismsConfig struct {
	Signature   string
	AccessKeyID string
}

// AliyunInternationalSMSConfig 阿里云短信
type AliyunInternationalSMSConfig struct {
	AccessKeyID  string // aliyun的AccessKeyID
	AccessSecret string // aliyun的AccessSecret
	SignName     string // 签名
}

// 苹果推送
type APNSPush struct {
	Dev      bool
	Topic    string
	Password string
	Cert     string
}

// 华为推送
type HMSPush struct {
	PackageName string
	AppID       string
	AppSecret   string
}

// 小米推送
type MIPush struct {
	PackageName string
	AppID       string
	AppSecret   string
	ChannelID   string
}

// oppo推送
type OPPOPush struct {
	PackageName  string
	AppID        string
	AppKey       string
	AppSecret    string
	MasterSecret string
}

type VIVOPush struct {
	PackageName string
	AppID       string
	AppKey      string
	AppSecret   string
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
