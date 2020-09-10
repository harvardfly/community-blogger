package repositories

import (
	"community-blogger/internal/pkg/models"
	"community-blogger/internal/pkg/requests"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"os"
	"testing"
)

var (
	art ArticleRepository
	db  *gorm.DB
)

//InitArticle 实例化article存储
func InitArticle(db *gorm.DB) (ArticleRepository, error) {
	return &MysqlArticleRepository{
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
	db.AutoMigrate(&models.Article{}, &models.Category{})
	art, _ = InitArticle(db)
}

func TestArticle(t *testing.T) {
	article := requests.Article{
		CategoryID: 1,
		Summary:    "测试article插入",
		Title:      "GOLANG 连接Mysql的时区问题",
	}
	res, err := art.Article(&article)
	if err != nil {
		t.Errorf("insert article failed, err:%v\n", err)
		return
	}

	t.Logf("insert article succ, articleId:%d\n", res.ID)
}

func TestGetArticle(t *testing.T) {
	id := 1
	res := art.GetArticle(id)
	if res.ID != id {
		t.Errorf("get article failed")
		return
	}

	t.Logf("get article succ, article:%v\n", res)
}

func TestArticleReadCount(t *testing.T) {
	id := 1
	err := art.ArticleReadCount(id)
	if err != nil {
		t.Errorf("get article failed, err:%d\n", err)
		return
	}

	t.Logf("get article succ, articleId:%#v\n", id)
}

func TestGetTOPNArticles(t *testing.T) {
	ids := []int{1, 2, 3, 5}
	articleList, err := art.GetTOPNArticles(ids)
	if err != nil {
		t.Errorf("get relative article failed, err:%v\n", err)
		return
	}

	for _, v := range articleList {
		t.Logf("id:%d title:%s\n", v.ID, v.Title)
	}
}
