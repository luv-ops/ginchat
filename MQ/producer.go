package MQ

import (
	"context"
	"log"
	"strconv"

	"github.com/segmentio/kafka-go"
)

// SendCommonMsg 目前无较大业务差距
func (k *KafkaClient) SendCommonMsg(ctx context.Context, dto *MsgDTO, topic string) error {
	data, err := dto.Marshal()
	if err != nil {
		log.Printf("私聊消息序列化失败 msgId=%s err=%v", dto.MsgID, err)
		return err
	}
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(dto.GetPartitionKey()),
		Value: data,
	}
	err = k.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("消息序列化失败 msgId=%s err=%v", dto.MsgID, err)
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
func (k *KafkaClient) ProduceGroup(ctx context.Context, dto *GroupDTO, topic string) error {
	data, err := dto.Marshal()
	if err != nil {
		log.Printf("群聊生产者序列化失败 groupName=%s,groupId=%d,err=%v", dto.GroupName, dto.GroupID, err)
		return err
	}
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(dto.GetPartitionKey(dto.Type)),
		Value: data,
	}
	err = k.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("群聊生产者失败 groupName=%s,groupId=%d,err=%v", dto.GroupName, dto.GroupID, err)
		return err
	}
	return nil
}
func (k *KafkaClient) SendDlqGroup(ctx context.Context, dto *GroupDTO, errReason string) error {
	data, _ := dto.Marshal()
	var key string
	if dto.Type == GroupCreate {
		key = dto.GroupName
	} else {
		key = strconv.Itoa(int(dto.GroupID))
	}
	msg := kafka.Message{
		Headers: []kafka.Header{
			{Key: "err_reason", Value: []byte(errReason)},
		},
		Key:   []byte(key),
		Value: data,
	}
	return k.dlqW.WriteMessages(ctx, msg)
}

func (k *KafkaClient) ProduceUser(ctx context.Context, dto *UserDTO, topic string) error {
	data, err := dto.Marshal()
	if err != nil {
		log.Printf("用户生产者序列化失败 userName=%s,err=%v", dto.Name, err)
		return err
	}
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(dto.GetPartitionKey()),
		Value: data,
	}
	err = k.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("用户生产者失败 userName=%s,err=%v", dto.Name, err)
		return err
	}
	return nil
}
func (k *KafkaClient) SendDlqUser(ctx context.Context, dto *UserDTO, errReason string) error {
	data, _ := dto.Marshal()
	msg := kafka.Message{
		Headers: []kafka.Header{
			{Key: "err_reason", Value: []byte(errReason)},
		},
		Key:   []byte(dto.Name),
		Value: data,
	}
	return k.dlqW.WriteMessages(ctx, msg)
}
