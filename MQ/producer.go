package MQ

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func (k *KafkaClient) SendPrivateMsg(ctx context.Context, dto *MsgDTO) error {
	data, err := dto.Marshal()
	if err != nil {
		log.Printf("私聊消息序列化失败 msgId=%s err=%v", dto.MsgID, err)
		return err
	}
	msg := kafka.Message{
		Topic: TopicPrivateMsg,
		Key:   []byte(dto.GetPartitionKey()),
		Value: data,
	}
	err = k.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("私聊消息序列化失败 msgId=%s err=%v", dto.MsgID, err)
		return err
	}
	return nil
}
func (k *KafkaClient) SendGroupMsg(ctx context.Context, dto *MsgDTO) error {
	data, err := dto.Marshal()
	if err != nil {
		log.Printf("群聊消息序列化失败 msgId=%s err=%v", dto.MsgID, err)
		return err
	}
	msg := kafka.Message{
		Topic: TopicGroupMsg,
		Key:   []byte(dto.GetPartitionKey()),
		Value: data,
	}
	err = k.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("群聊消息序列化失败 msgId=%s err=%v", dto.MsgID, err)
		return err
	}
	return nil
}
func (k *KafkaClient) SendFriendReq(ctx context.Context, dto *MsgDTO) error {
	data, err := dto.Marshal()
	if err != nil {
		log.Printf("好友请求序列化失败 msgId=%s err=%v", dto.MsgID, err)
		return err
	}
	msg := kafka.Message{
		Topic: TopicFriendReq,
		Key:   []byte(dto.GetPartitionKey()),
		Value: data,
	}
	err = k.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("好友请求序列化失败 msgId=%s err=%v", dto.MsgID, err)
		return err
	}
	return nil
}
func (k *KafkaClient) SendDlqMsg(ctx context.Context, dto *MsgDTO, errReason string) error {
	data, _ := dto.Marshal()
	msg := kafka.Message{
		Headers: []kafka.Header{
			{Key: "err_reason", Value: []byte(errReason)},
		},
		Key:   []byte(dto.MsgID),
		Value: data,
	}
	return k.dlqW.WriteMessages(ctx, msg)
}
