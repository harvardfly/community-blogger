package kafka

import (
	"community-blogger/internal/pkg/es"
	"community-blogger/internal/pkg/requests"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strconv"
	"sync"
)

var asyncProducerOnce sync.Once
var asyncProducer sarama.AsyncProducer

var syncProducerOnce sync.Once

// ClientType 定义kafka client 结构体
type ClientType struct {
	Producer sarama.SyncProducer // 默认是生产者同步模式
	Settings *Options
	logger   *zap.Logger
}

// Client  kafka连接类型
var Client ClientType

// Options Kafka配置
type Options struct {
	Brokers     string `yaml:"brokers"`
	TryTimes    int    `yaml:"tryTimes"`
	SyncEsTopic string `yaml:"syncEsTopic"`
}

// NewOptions for Kafka
func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("kafka", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal kafka option error")
	}
	logger.Info("load kafka options success", zap.Any("kafka options", o))
	Client.Settings = o
	Client.logger = logger
	return o, err
}

// New 初始化Kafka连接信息
func New(o *Options, logger *zap.Logger) (producer sarama.SyncProducer, err error) {
	config := sarama.NewConfig()
	//等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	//Hash向partition发送消息
	config.Producer.Partitioner = sarama.NewHashPartitioner
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoResponse这里才有用.
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	asyncProducerOnce.Do(func() {
		//使用配置,新建一个异步生产者
		asyncProducer, err = sarama.NewAsyncProducer([]string{o.Brokers}, config)
		if err != nil {
			logger.Error("NewAsyncProducer init error", zap.Error(err))
			return
		}
	})

	syncProducerOnce.Do(func() {
		//使用配置,新建一个同步生产者
		syncProducer, err := sarama.NewSyncProducer([]string{o.Brokers}, config)
		if err != nil {
			logger.Error("NewSyncProducer init error", zap.Error(err))
			return
		}
		Client.Producer = syncProducer
	})

	return producer, nil
}

// AsyncProduce 生产者以异步的方式发送数据到kafka
func (cli *ClientType) AsyncProduce(topic string, payload interface{}) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		cli.logger.Error("json.Marshal error", zap.Error(err))
		return
	}
	asyncProducer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(payloadJSON),
	}
}

// SyncProduce 生产者以同步方式发送数据到kafka
func (cli *ClientType) SyncProduce(topic string, payload interface{}) (result bool) {
	key := uuid.NewV4().String()
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		cli.logger.Error("json.Marshal error", zap.Error(err))
	}
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(payloadJSON),
	}
	//重试TryTimes次
	for i := 0; i < cli.Settings.TryTimes; i++ {
		_, _, err := cli.Producer.SendMessage(msg)
		if err != nil {
			result = false
			cli.logger.Error("SendMessage error", zap.Error(err))
		} else {
			result = true
			break
		}
	}
	return
}

// Subscribe 消费kafka数据到es
func (cli *ClientType) Subscribe() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V0_11_0_2

	// 初始化consumer
	consumer, err := sarama.NewConsumer([]string{cli.Settings.Brokers}, config)
	if err != nil {
		cli.logger.Error("Subscribe create consumer error", zap.Error(err))
		return
	}

	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(cli.Settings.SyncEsTopic, 0, sarama.OffsetOldest)
	if err != nil {
		cli.logger.Error("try create partitionConsumer error", zap.Error(err))
		return
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			cli.logger.Info(
				"kafka msg",
				zap.Any("msg offset", msg.Offset),
				zap.Any("partition", msg.Partition),
				zap.Any("timestamp", msg.Timestamp.String()),
				zap.Any("value", string(msg.Value)),
			)
			var esData requests.ArticleES
			err := json.Unmarshal(msg.Value, &esData)
			if err != nil {
				cli.logger.Error("json Unmarshal error", zap.Error(err))
				return
			}
			Params := make(map[string]string)
			Params["index"] = "article"
			Params["id"] = strconv.Itoa(esData.ID)
			Params["bodyJson"] = string(msg.Value)
			_, err = es.Client.Insert(Params)
			if err != nil {
				cli.logger.Error("create es error", zap.Error(err))
				return
			}
		case err := <-partitionConsumer.Errors():
			cli.logger.Error("partitionConsumer error", zap.Error(err))
			return
		}
	}
}

// ProviderSet inject kafka settings
var ProviderSet = wire.NewSet(New, NewOptions)
