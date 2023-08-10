package config

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
}
