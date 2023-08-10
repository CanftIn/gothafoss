package config

import (
	"time"

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
	WelcomeMessage              string //登录注册欢迎语
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
		SystemUID       string //系统账号uid
		FileHelperUID   string // 文件助手uid
		SystemGroupID   string //系统群ID 需求在app_config表里设置new_user_join_system_group为1才有效
		SystemGroupName string // 系统群的名字
		AdminUID        string //系统管理员账号
	}

	// ---------- push ----------
	Push struct {
		ContentDetailOn bool     //  推送是否显示正文详情(如果为false，则只显示“您有一条新的消息” 默认为true)
		PushPoolSize    int64    // 推送任务池大小
		APNS            APNSPush // 苹果推送
		MI              MIPush   // 小米推送
		HMS             HMSPush  // 华为推送
		VIVO            VIVOPush // vivo推送
		OPPO            OPPOPush // oppo推送
	}

	// ---------- 文件服务 ----------

	FileService FileService   // 文件服务
	OSS         OSSConfig     // 阿里云oss配置
	Minio       MinioConfig   // minio配置
	Seaweed     SeaweedConfig // seaweedfs配置

	TimingWheelTick duration // The time-round training interval must be 1ms or more
	TimingWheelSize int64    // Time wheel size
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
	BucketName      string // Bucket名称 比如 tangsengdaodao
	BucketURL       string // 文件下载地址域名 对应aliyun的Bucket域名
	AccessKeyID     string
	AccessKeySecret string
}

type MinioConfig struct {
	URL             string // 文件下载上传基地址 例如： http://127.0.0.1:9000
	AccessKeyID     string //minio accessKeyID
	SecretAccessKey string //minio secretAccessKey
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
