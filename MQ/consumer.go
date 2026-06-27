package MQ

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// StartCommonConsumer 目前无较大业务差距
func (k *KafkaClient) StartCommonConsumer(ctx context.Context, topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: k.brokers,
		GroupID: ConsumerGroupID,
		Topic:   topic,
		MaxWait: 100 * time.Millisecond,
	})
	defer reader.Close()
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			// 程序收到退出信号，ctx被取消，直接退出消费循环
			if ctx.Err() != nil {
				log.Println("消息消费循环结束")
				return
			}
			log.Printf("读取消息异常%v", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		var dto MsgDTO
		if err = dto.Unmarshal(msg.Value); err != nil {
			log.Printf("消息解析失败 fromId=%d targetId=%d err=%v", dto.FromID, dto.TargetID, err)
			//直接丢入死信队列
			_ = k.SendDlqMsg(ctx, &dto, err.Error())
			//提交offset，跳过坏消息
			_ = reader.CommitMessages(ctx, msg)
			continue
		}

		handleErr := Retry(ctx, k, &dto)
		if handleErr != nil {
			log.Printf("消息重试耗尽，转入死信 msgId=%s err=%v", dto.MsgID, handleErr)
		}
		//消息读取成功，提交offset
		err = reader.CommitMessages(ctx, msg)
		if err != nil {
			log.Printf("消息提交offset失败 msgId=%s err=%v", dto.MsgID, err)
		}

	}
}

func Retry(ctx context.Context, k *KafkaClient, dto *MsgDTO) error {
	retryTime := 0
	var err error
	for retryTime < MaxRetryCount {
		if dto.ChatType == ChatTypeFriendRequest {
			err = FriReqHandler.HandleFReq(dto)
		} else if dto.ChatType == ChatTypeFriendRequestAccept {
			err = FriReqHandler.HandleFReqAccept(dto)
		} else if dto.ChatType == ChatTypeFriendRequestHasRead {
			err = FriReqHandler.HandleFReqHasRead(dto)
		} else {
			err = MsgHandler.HandleMsg(dto)
		}
		if err == nil {
			return nil
		}
		retryTime++
		log.Printf("消息处理失败，准备第%d次重试 msgId=%s err=%v", retryTime, dto.MsgID, err)
		time.Sleep(200 * time.Millisecond)
	}
	//3次重试全部失败，跌入死信队列
	_ = k.SendDlqMsg(ctx, dto, "消息处理3次全部失败")
	return err
}

func (k *KafkaClient) StartGroupConsumer(ctx context.Context, topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: k.brokers,
		GroupID: ConsumerGroupID,
		Topic:   topic,
		MaxWait: 100 * time.Millisecond,
	})
	defer reader.Close()
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			// 程序收到退出信号，ctx被取消，直接退出消费循环
			if ctx.Err() != nil {
				log.Println("群消费循环结束")
				return
			}
			log.Printf("读取群异常%v", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		var dto GroupDTO
		if err = dto.Unmarshal(msg.Value); err != nil {
			log.Printf("群解析失败 groupName=%s,groupId=%d,err=%v", dto.GroupName, dto.GroupID, err)
			//直接丢入死信队列
			_ = k.SendDlqGroup(ctx, &dto, err.Error())
			//提交offset，跳过坏消息
			_ = reader.CommitMessages(ctx, msg)
			continue
		}

		handleErr := RetryGroup(ctx, k, &dto)
		if handleErr != nil {
			log.Printf("群重试耗尽，转入死信 groupName=%s,groupId=%d,err=%v", dto.GroupName, dto.GroupID, handleErr)
		}
		//消息读取成功，提交offset
		err = reader.CommitMessages(ctx, msg)
		if err != nil {
			log.Printf("群提交offset失败 groupName=%s,groupId=%d,err=%v", dto.GroupName, dto.GroupID, err)
		}

	}
}

func RetryGroup(ctx context.Context, k *KafkaClient, dto *GroupDTO) error {
	retryTime := 0
	var err error
	for retryTime < MaxRetryCount {
		if dto.Type == GroupCreate {
			err = GroHandler.HandleGroupCreate(dto)
		} else if dto.Type == GroupInvite {
			err = GroHandler.HandleGroupInvite(dto)
		} else {
			return errors.New("未知群方法类型")
		}
		if err == nil {
			return nil
		}
		retryTime++
		log.Printf("消息处理失败，准备第%d次重试 groupName=%s groupId=%d err=%v", retryTime, dto.GroupName, dto.GroupID, err)
		time.Sleep(200 * time.Millisecond)
	}
	//3次重试全部失败，跌入死信队列
	_ = k.SendDlqGroup(ctx, dto, "消息处理3次全部失败")
	return err
}

func (k *KafkaClient) StartUserConsumer(ctx context.Context, topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: k.brokers,
		GroupID: ConsumerGroupID,
		Topic:   topic,
		MaxWait: 100 * time.Millisecond,
	})
	defer reader.Close()
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			// 程序收到退出信号，ctx被取消，直接退出消费循环
			if ctx.Err() != nil {
				log.Println("用户消费循环结束")
				return
			}
			log.Printf("读取用户异常%v", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		var dto UserDTO
		if err = dto.Unmarshal(msg.Value); err != nil {
			log.Printf("用户解析失败 userName=%s,err=%v", dto.Name, err)
			//直接丢入死信队列
			_ = k.SendDlqUser(ctx, &dto, err.Error())
			//提交offset，跳过坏消息
			_ = reader.CommitMessages(ctx, msg)
			continue
		}

		handleErr := RetryUser(ctx, k, &dto)
		if handleErr != nil {
			log.Printf("用户重试耗尽，转入死信 userName=%s,err=%v", dto.Name, handleErr)
		}
		//消息读取成功，提交offset
		err = reader.CommitMessages(ctx, msg)
		if err != nil {
			log.Printf("用户提交offset失败 userName=%s,err=%v", dto.Name, err)
		}

	}
}

func RetryUser(ctx context.Context, k *KafkaClient, dto *UserDTO) error {
	retryTime := 0
	var err error
	for retryTime < MaxRetryCount {

		err = UHandler.HandleUserCreate(dto)

		if err == nil {
			return nil
		}
		retryTime++
		log.Printf("用户处理失败，准备第%d次重试 userName=%s,err=%v", dto.Name, err)
		time.Sleep(200 * time.Millisecond)
	}
	//3次重试全部失败，跌入死信队列
	_ = k.SendDlqUser(ctx, dto, "用户处理3次全部失败")
	return err
}
