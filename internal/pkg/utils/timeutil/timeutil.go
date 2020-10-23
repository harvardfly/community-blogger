package timeutil

import (
	"bytes"
	"errors"
	"math"
	"strconv"
	"time"
)

// GeneratorYesterdayPeriod 生成昨日时间区间
// 2006-01-02 00:00:00 ----- 2006-01-03 00:00:00
func GeneratorYesterdayPeriod() (time.Time, time.Time) {
	now := time.Now().In(time.UTC)
	begin := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return begin, end
}

// TimeInterval 时间间隔
type TimeInterval struct {
	Begin time.Time
	End   time.Time
}

// TimeProcessInterval 时间处理成时间区间
func TimeProcessInterval(t ...time.Time) []TimeInterval {
	tis := make([]TimeInterval, 0)
	for _, v := range t {
		tm1 := time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, v.Location())
		tm2 := tm1.AddDate(0, 0, 1)
		ti := TimeInterval{
			Begin: tm1,
			End:   tm2,
		}
		tis = append(tis, ti)
	}
	return tis
}

// TimeProcessReportParams 时间方式验证
func TimeProcessReportParams(style string, ti []string) ([]time.Time, error) {
	if style != "day" && style != "week" && style != "month" {
		return nil, errors.New("unsupported_metric_calculation")
	}
	if len(ti) != 2 {
		return nil, errors.New("date_format_error")
	}
	beginStr := ti[0]
	beginTime, err := time.Parse("2006-01-02", beginStr)
	if err != nil {
		return nil, errors.New("date_format_error")
	}
	endStr := ti[1]
	endTime, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		return nil, errors.New("date_format_error")
	}
	if endTime.Before(beginTime) {
		return nil, errors.New("date_format_error")
	}
	// 按天选择不能超过90天
	if style == "day" && beginTime.AddDate(0, 0, 90).Before(endTime) {
		return nil, errors.New("date_format_error")
	}
	// 按周选择不能超过50周
	if style == "week" && beginTime.AddDate(0, 0, 50*7).Before(endTime) {
		return nil, errors.New("date_format_error")
	}
	// 按月选择不能超过10年
	if style == "month" && beginTime.AddDate(10, 0, 0).Before(endTime) {
		return nil, errors.New("date_format_error")
	}
	// 按周选择 第一天必须是周一 最后一天必须是周日
	if style == "week" && (beginTime.Weekday() != time.Weekday(1) || endTime.Weekday() != time.Weekday(0)) {
		return nil, errors.New("date_format_error")
	}
	// 按月选择 第一天必须是一号 最后一天必须是月末
	if style == "month" && (beginTime.Day() != 1 || endTime.AddDate(0, 0, 1).Day() != 1) {
		return nil, errors.New("date_format_error")
	}
	return []time.Time{beginTime, endTime.Add(time.Hour * 24)}, nil
}

// StrTime 格式化时间
func StrTime(times time.Time) string {
	var res string
	datetime := "2006-01-02 15:04:05" //待转化为时间戳的字符串
	time1, _ := time.ParseInLocation(datetime, times.Format(datetime), time.Local)
	atime := time1.Unix()
	if atime < 0 {
		return res
	}
	var byTime = []int64{365 * 24 * 60 * 60, 24 * 60 * 60, 60 * 60, 60, 1}
	var unit = []string{"年前", "天前", "小时前", "分钟前", "秒钟前"}
	now := time.Now().Unix()
	ct := now - atime
	if ct < 0 {
		return "刚刚"
	}

	if ct > 30*24*60*60 {
		return times.Format(datetime)
	}
	for i := 0; i < len(byTime); i++ {
		if ct < byTime[i] {
			continue
		}
		var temp = math.Floor(float64(ct / byTime[i]))
		ct = ct % byTime[i]
		if temp > 0 {
			var tempStr string
			tempStr = strconv.FormatFloat(temp, 'f', -1, 64)
			res = MergeString(tempStr, unit[i]) //此处调用了一个我自己封装的字符串拼接的函数（你也可以自己实现）
		}
		break
	}
	return res
}

// StrDate 格式化日期
func StrDate(times time.Time) string {
	var res string
	datetime := "2006-01-02 15:04:05" //待转化为时间戳的字符串
	time1, _ := time.ParseInLocation(datetime, times.Format(datetime), time.Local)
	atime := time1.Unix()
	if atime < 0 {
		return res
	}
	return times.Format(datetime)
}

// MergeString @des 拼接字符串  args ...string 要被拼接的字符串序列 return string
func MergeString(args ...string) string {
	buffer := bytes.Buffer{}
	for i := 0; i < len(args); i++ {
		buffer.WriteString(args[i])
	}
	return buffer.String()
}

// GetNowMonthDays 获取当前月已过的天数(包含当日)
func GetNowMonthDays() int {
	now := time.Now()
	day := now.Day()
	return day
}
