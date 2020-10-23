package constutil

const (
	// DefaultPageSize 默认每页的条数
	DefaultPageSize int = 20
	// MaxPageSize 每页最大的的条数
	MaxPageSize int = 100
)

const (
	// Letters 数字字母组合
	Letters = "0123456789aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"
	// NumberLetters 数字
	NumberLetters = "0123456789"
)

const (
	// TablePrefix 数据库表的前缀
	TablePrefix = "blog_"
)

const (
	// Schema ETCD的Schema
	Schema string = "blog"
)

const (
	// CreateArticle kafka topic
	CreateArticle string = "first_topic"
)

const (
	WeekRank  string = "week"
	MonthRank string = "month"
	TotalRank string = "all"
)

const (
	WeekDays  int = 7
	MonthDays int = 31
)
