package Autowired

import (
	"GinChat/Mysql"
	"GinChat/mapper"
)

var (
	UserMapper         *mapper.UserMapper
	FriendMapper       *mapper.FriendMapper
	GroupMapper        *mapper.GroupMapper
	ConversationMapper *mapper.ConversationMapper
	MessageMapper      *mapper.MessageMapper
)

func InitMapper() {
	UserMapper = mapper.NewUserMapper(Mysql.DB)
	FriendMapper = mapper.NewFriendMapper(Mysql.DB)
	GroupMapper = mapper.NewGroupMapper(Mysql.DB)
	ConversationMapper = mapper.NewConversationMapper(Mysql.DB)
	MessageMapper = mapper.NewMessageMapper(Mysql.DB)
}
