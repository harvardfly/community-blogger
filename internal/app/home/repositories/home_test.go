package repositories

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"os"
	"testing"
	"community-blogger/internal/pkg/models"
	"community-blogger/internal/pkg/requests"
)

var (
	home HomeRepository
	db   *gorm.DB
)

//InitHome 实例化home存储
func InitHome(db *gorm.DB) (HomeRepository, error) {
	return &MysqlHomeRepository{
		db: db,
	}, nil
}

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test_blog?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = db.DB().Ping()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	db.LogMode(true)
	db.AutoMigrate(&models.Home{})
	home, _ = InitHome(db)
}

func TestHome(t *testing.T) {
	req := requests.Home{
		URL:         "www.baidu.com",
		Img:         "测试article插入",
		Title:       "GOLANG 连接Mysql的时区问题",
		Description: "测试...........",
	}
	err := home.Home(&req)
	if err != nil {
		t.Errorf("insert home failed, err:%v\n", err)
		return
	}

	t.Logf("insert home succ, Title:%s\n", req.Title)
}

func TestHomeList(t *testing.T) {
	req := requests.HomeList{
		Title:       "GOLANG 连接Mysql的时区问题",
		Description: "测试...........",
	}
	count, res, err := home.HomeList(&req)
	t.Logf("get count:%d\n", count)
	if err != nil {
		t.Errorf("get home failed")
		return
	}

	t.Logf("get home succ, homelist:%v\n", res)
}
