package event

import (
	"github.com/CanftIn/gothafoss/pkg/im/config"
	"github.com/CanftIn/gothafoss/pkg/log"
	"github.com/CanftIn/gothafoss/server/modules/file"
)

const (
	// GroupCreate 群创建
	GroupCreate string = "group.create"
	// GroupUnableAddDestroyAccount 无法添加注销账号到群聊
	GroupUnableAddDestroyAccount string = "group.unable.add.destroy.account"
	// GroupUpdate 群更新
	GroupUpdate string = "group.update"
	// GroupMemberAdd 群成员添加
	GroupMemberAdd string = "group.memberadd"
	// GroupMemberScanJoin 扫码加入群
	GroupMemberScanJoin string = "group.member.scan.join"
	// GroupMemberTransferGrouper 转让群主
	GroupMemberTransferGrouper string = "group.member.transfer.grouper"
	// GroupAvatarUpdate 群头像更新
	GroupAvatarUpdate string = "group.avatar.update"
	// GroupMemberRemove 群成员移除
	GroupMemberRemove string = "group.memberremove"
	// FriendApply 好友申请
	FriendApply string = "friend.apply"
	// GroupMemberInviteRequest 群邀请请求
	GroupMemberInviteRequest string = "group.member.invite"
	// ConversationDelete 删除最近会话
	ConversationDelete string = "conversation.delete"
	// EventTransfer 转账
	EventTransfer string = "transfer"
	// EventRedpacketReceive 领取红包
	EventRedpacketReceive string = "redpacket.receive"
	// EventUserRegister 用户注册
	EventUserRegister string = "user.register"
	// EventUserPublishMoment 用户发布动态
	EventUserPublishMoment string = "moment.publish"
	// EventUserDeleteMoment 用户删除动态
	EventUserDeleteMoment string = "moment.delete"
	// FriendSure 好友确认
	FriendSure string = "friend.sure"
	// FriendDelete 好友删除
	FriendDelete string = "friend.delete"
)

// Event 事件
type Event struct {
	db  *DB
	ctx *config.Context
	log.Log
	fileService file.IService
}
