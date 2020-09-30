package es

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	// Client 连接类型
	Client ClientType
	// Config 配置信息
	Config Options
)

// ClientType es 连接信息struct
type ClientType struct {
	esConn *elastic.Client
	logger *zap.Logger
}

// Options ES配置
type Options struct {
	URL         string        `yaml:"url"`
	HealthCheck time.Duration `yaml:"healthCheck"`
	Sniff       bool          `yaml:"sniff"`
	Gzip        bool          `yaml:"gzip"`
	Timeout     string        `yaml:"timeout"`
}

// Aggregations 定义常用的聚合用的一些参数
type Aggregations struct {
	AvgMetric AvgMetric `json:"AVG_Metric"`
}

// AvgMetric 平均值 聚合
type AvgMetric struct {
	Buckets []Metric `json:"buckets"`
}

// Metric 聚合
type Metric struct {
	AvgTime Value `json:"avg_time"`
}

// Value 聚合值
type Value struct {
	Value float64 `json:"value"`
}

// NewOptions for ES
func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var (
		err error
		o   = new(Options)
	)
	if err = v.UnmarshalKey("es", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal es option error")
	}

	logger.Info("load es options success", zap.Any("es options", o))
	return o, err
}

// New 初始化ES连接信息
func New(o *Options, logger *zap.Logger) (esConn *elastic.Client, err error) {
	esConn, err = elastic.NewClient(
		elastic.SetURL(o.URL),
		elastic.SetSniff(o.Sniff),
		elastic.SetHealthcheckInterval(o.HealthCheck*time.Second),
		elastic.SetGzip(o.Gzip),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil {
		logger.Error("es NewClient create error", zap.Error(err))
		return
	}
	info, code, err := esConn.Ping(o.URL).Do(context.Background())
	if err != nil {
		logger.Error("Ping esConn error", zap.Error(err))
		return
	}

	logger.Info("ES returned with code and version",
		zap.Any("code", code),
		zap.Any("version", info.Version.Number),
	)

	esVersion, err := esConn.ElasticsearchVersion(o.URL)
	if err != nil {
		logger.Error("esConn search Version error", zap.Error(err))
		return
	}
	logger.Info("conn es success",
		zap.Any("esConn", esConn),
		zap.Any("version", esVersion),
	)
	Client.esConn = esConn

	return esConn, nil
}

// Insert 创建
func (cli *ClientType) Insert(Params map[string]string) (string, error) {
	var (
		res *elastic.IndexResponse
		err error
	)
	res, err = Client.esConn.Index().
		Index(Params["index"]).
		Id(Params["id"]).BodyJson(Params["bodyJson"]).
		Do(context.Background())

	if err != nil {
		cli.logger.Error("insert error", zap.Error(err))
		return "", err
	}
	return res.Result, nil
}

// Delete 删除
func (cli *ClientType) Delete(Params map[string]string) (string, error) {
	var (
		res *elastic.DeleteResponse
		err error
	)

	res, err = Client.esConn.Delete().Index(Params["index"]).
		Id(Params["id"]).
		Do(context.Background())

	if err != nil {
		cli.logger.Error("delete error", zap.Error(err))
		return "", err
	}

	return res.Result, nil
}

// Update 更新
func (cli *ClientType) Update(Params map[string]string, Doc map[string]interface{}) string {
	var (
		res *elastic.UpdateResponse
		err error
	)
	res, err = Client.esConn.Update().
		Index(Params["index"]).
		Id(Params["id"]).
		Doc(Doc).
		Do(context.Background())

	if err != nil {
		cli.logger.Error("update error", zap.Error(err))
		return ""
	}
	return res.Result

}

// GetByID 通过ID查找
func (cli *ClientType) GetByID(Params map[string]string) *elastic.GetResult {
	var (
		res *elastic.GetResult
		err error
	)
	if len(Params["id"]) < 0 {
		fmt.Printf("param error")
		return res
	}

	res, err = Client.esConn.Get().
		Index(Params["index"]).
		Id(Params["id"]).
		Do(context.Background())

	if err != nil {
		cli.logger.Error("GetByID error", zap.Error(err))
		return nil
	}

	return res
}

// Query 搜索
func (cli *ClientType) Query(Params map[string]string) *elastic.SearchResult {
	var (
		res *elastic.SearchResult
		err error
	)
	//取所有
	res, err = Client.esConn.Search(Params["index"]).Do(context.Background())
	if len(Params["queryString"]) > 0 {
		//字段相等
		q := elastic.NewQueryStringQuery(Params["queryString"])
		res, err = Client.esConn.Search(Params["index"]).
			Query(q).
			Do(context.Background())
	}
	if err != nil {
		cli.logger.Error("Query error", zap.Error(err))
		return nil
	}

	return res
}

// List 分页列表
func (cli *ClientType) List(Params map[string]string) *elastic.SearchResult {
	var res *elastic.SearchResult
	var err error
	size, _ := strconv.Atoi(Params["size"])
	page, _ := strconv.Atoi(Params["page"])
	q := elastic.NewQueryStringQuery(Params["queryString"])

	// 排序类型 desc asc es 中只使用 bool 值  true or false
	sortType := true
	if Params["sort_type"] == "desc" {
		sortType = false
	}

	if size < 0 || page < 0 {
		cli.logger.Error("param error", zap.Error(err))
		return nil
	}
	if len(Params["queryString"]) > 0 {
		res, err = Client.esConn.Search(Params["index"]).
			Query(q).
			Size(size).
			From((page)*size).
			Sort(Params["sort"], sortType).
			Timeout(Config.Timeout).
			Do(context.Background())

	} else {
		res, err = Client.esConn.Search(Params["index"]).
			Size(size).
			From((page)*size).
			Sort(Params["sort"], sortType).
			Timeout(Config.Timeout).
			Do(context.Background())
	}

	if err != nil {
		cli.logger.Error("func list error", zap.Error(err))
		return nil
	}
	return res

}

// Aggregation 聚合 平均
func (cli *ClientType) Aggregation(Params map[string]string) *elastic.SearchResult {
	var res *elastic.SearchResult
	var err error

	//需要聚合的指标 求平均
	avg := elastic.NewAvgAggregation().Field(Params["avg"])
	//单位时间和指定字段
	aggs := elastic.NewDateHistogramAggregation().
		Field(Params["field"]).
		SubAggregation(Params["agg_name"], avg)

	res, err = Client.esConn.Search(Params["index"]).
		Size(0).
		Aggregation(Params["aggregation_name"], aggs).
		Timeout(Config.Timeout).
		Do(context.Background())

	if err != nil {
		cli.logger.Error("func Aggregation error", zap.Error(err))
		return nil
	}

	return res
}

// ProviderSet inject es settings
var ProviderSet = wire.NewSet(New, NewOptions)
