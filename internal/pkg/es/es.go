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
	Client ClientType //连接类型
	Config Options    // 配置信息
)

type ClientType struct {
	esConn *elastic.Client
}

// Options ES配置
type Options struct {
	Host        string        `yaml:"host"`
	Port        int           `yaml:"port"`
	HealthCheck time.Duration `yaml:"healthCheck"`
	Sniff       bool          `yaml:"sniff"`
	Gzip        bool          `yaml:"gzip"`
	Timeout     string        `yaml:"timeout"`
}

//定义常用的聚合用的一些参数
type Aggregations struct {
	AvgMetric AvgMetric `json:"AVG_Metric"`
}

type AvgMetric struct {
	Buckets []Metric `json:"buckets"`
}

type Metric struct {
	AvgTime Value `json:"avg_time"`
}

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
func New(o *Options) (esConn *elastic.Client, err error) {
	url := fmt.Sprintf("%s:%d", o.Host, o.Port)
	esConn, err = elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(o.Sniff),
		elastic.SetHealthcheckInterval(o.HealthCheck*time.Second),
		elastic.SetGzip(o.Gzip),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil {
		panic(err)
	}
	info, code, err := esConn.Ping(url).Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	esVersion, err := esConn.ElasticsearchVersion(url)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esVersion)
	fmt.Println("conn es success", esConn)
	Client.esConn = esConn

	return esConn, nil
}

// Create 创建
func Create(Params map[string]string) (string, error) {
	var (
		res *elastic.IndexResponse
		err error
	)
	res, err = Client.esConn.Index().
		Index(Params["index"]).
		Id(Params["id"]).BodyJson(Params["bodyJson"]).
		Do(context.Background())

	if err != nil {
		fmt.Printf("create error %s\n", err)
		panic(err)
	}
	return res.Result, err

}

// Delete 删除
func Delete(Params map[string]string) (string, error) {
	var (
		res *elastic.DeleteResponse
		err error
	)

	res, err = Client.esConn.Delete().Index(Params["index"]).
		Id(Params["id"]).
		Do(context.Background())

	if err != nil {
		println(err.Error())
		return "", err
	}

	fmt.Printf("delete result %s\n", res.Result)
	return res.Result, nil
}

// Update 修改
func Update(Params map[string]string, Doc map[string]interface{}) string {
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
		fmt.Printf("update error %s\n", err.Error())
		return ""
	}
	fmt.Printf("update success %s\n", res.Result)
	return res.Result

}

// GetByID 通过ID查找
func GetByID(Params map[string]string) *elastic.GetResult {
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
		panic(err)
	}

	return res
}

// Query 搜索
func Query(Params map[string]string) *elastic.SearchResult {
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
		println(err.Error())
	}

	//if res.Hits.TotalHits > 0 {
	//	fmt.Printf("Found a total of %d Employee \n", res.Hits.TotalHits)
	//}
	return res
}

// List 分页列表
func List(Params map[string]string) *elastic.SearchResult {
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
	//fmt.Printf(" sort info  %s,%s\n", Params["sort"],Params["sort_type"])
	if size < 0 || page < 0 {
		fmt.Printf("param error")
		return res
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
			//SortBy(elastic.NewFieldSort("add_time").UnmappedType("long").Desc(), elastic.NewScoreSort()).
			Timeout(Config.Timeout).
			Do(context.Background())
	}

	if err != nil {
		println("func list error:" + err.Error())
	}
	return res

}

// Aggregation 聚合 平均
func Aggregation(Params map[string]string) *elastic.SearchResult {
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
		//Sort(Params["sort"],sort_type).
		Timeout(Config.Timeout).
		Do(context.Background())

	if err != nil {
		println("func Aggregation error:" + err.Error())
	}
	println("func Aggregation here 297")

	return res
}

// ProviderSet inject es settings
var ProviderSet = wire.NewSet(New, NewOptions)
