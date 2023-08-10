package msg

// DeviceLevel 设备等级
type DeviceLevel uint8

const (
	// DeviceLevelSlave 从设备
	DeviceLevelSlave DeviceLevel = 0
	// DeviceLevelMaster 主设备
	DeviceLevelMaster DeviceLevel = 1
)

// DeviceFlag 设备类型
type DeviceFlag uint8

const (
	// APP APP
	APP DeviceFlag = iota
	// Web Web
	Web
	// PC在线
	PC
)

type Channel struct {
	ChannelID   string
	ChannelType uint8
}
