package models

import (
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"

	"github.com/Unknwon/com"
	"github.com/astaxie/beego/orm"

	_ "github.com/mattn/go-sqlite3"
)

const (
	// _DB_NAME 设置数据库路径
	_DB_NAME = "data/beeblog.db"
	// _SQLITE3_DRIVER 设置数据库驱动
	_SQLITE3_DRIVER = "sqlite3"
)

// Category 分类
type Category struct {
	ID              int64
	Title           string
	Created         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	TopicTime       time.Time `orm:"index"`
	TopicCount      int64
	TopicLastUserID int64
}

// Topic 文章
type Topic struct {
	ID              int64
	UID             int64
	Title           string
	Category        string
	Labels          string
	Content         string `orm:"size(5000)"`
	Attachment      string
	Created         time.Time `orm:"index"`
	Updated         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	Author          string
	ReplyTime       time.Time `orm:"index"`
	ReplyCount      int64
	ReplyLastUserID int64
}

// Comment 评论
type Comment struct {
	ID      int64
	Tid     int64
	Name    string
	Content string    `orm:"size(1000)"`
	Created time.Time `orm:"index"`
}

// RegisterDB 注册数据库
func RegisterDB() {
	// 检查数据库文件
	if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}
	// 注册模型
	orm.RegisterModel(new(Category), new(Topic), new(Comment))
	// 注册驱动（“sqlite3” 属于默认注册，此处代码可省略）
	orm.RegisterDriver(_SQLITE3_DRIVER, orm.DRSqlite)
	// 注册默认数据库
	orm.RegisterDataBase("default", _SQLITE3_DRIVER, _DB_NAME, 10)
}

// AddCategory 添加分类
func AddCategory(name string) error {
	o := orm.NewOrm()

	cate := &Category{
		Title:     name,
		Created:   time.Now(),
		TopicTime: time.Now(),
	}
	// 查询数据
	qs := o.QueryTable("category")
	err := qs.Filter("title", name).One(cate)
	if err == nil {
		return err
	}
	// 插入数据
	_, err = o.Insert(cate)
	if err != nil {
		return err
	}
	return nil
}

// DelCategory 删除分类
func DelCategory(id string) error {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	cate := &Category{ID: cid}
	_, err = o.Delete(cate)
	return err
}

// GetAllCategories 获取所有分类
func GetAllCategories() ([]*Category, error) {
	o := orm.NewOrm()
	cates := make([]*Category, 0)
	qs := o.QueryTable("category")
	_, err := qs.All(&cates)
	return cates, err
}

// GetTopic 获取指定文章
func GetTopic(tid string) (*Topic, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()

	topic := new(Topic)

	qs := o.QueryTable("topic")
	err = qs.Filter("id", tidNum).One(topic)
	if err != nil {
		return nil, err
	}
	topic.Views++
	_, err = o.Update(topic)

	topic.Labels = strings.Replace(strings.Replace(
		topic.Labels, "#", " ", -1), "$", "", -1)
	return topic, err
}

// GetAllTopic 获取所有文章
func GetAllTopic(category, label string, isDesc bool) (topics []*Topic, err error) {
	o := orm.NewOrm()

	topics = make([]*Topic, 0)

	qs := o.QueryTable("topic")
	if isDesc {
		if len(category) > 0 {
			qs = qs.Filter("category", category)
		}
		if len(label) > 0 {
			qs = qs.Filter("labels__contains", "$"+label+"#")
		}
		_, err = qs.OrderBy("-created").All(&topics)
	} else {
		_, err = qs.All(&topics)
	}
	return topics, err
}

// AddTopic 添加文章
func AddTopic(title, category, labels, content, attachemnt string) error {

	o := orm.NewOrm()

	// 处理标签
	mlabels := "$" + strings.Join(strings.Split(labels, " "), "#$") + "#"

	topic := &Topic{
		Title:      title,
		Category:   category,
		Labels:     mlabels,
		Content:    content,
		Attachment: attachemnt,
		Created:    time.Now(),
		Updated:    time.Now(),
		ReplyTime:  time.Now(),
	}
	// 插入数据
	_, err := o.Insert(topic)
	if err != nil {
		beego.Error(err)
		return err
	}
	// 更新分类统计
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	if err != nil {
		beego.Error(err)
		return err
	}
	cate.TopicCount++
	_, err = o.Update(cate)

	return err
}

// DelTopic 删除文章
func DelTopic(tid string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	var oldCate string
	topic := &Topic{ID: tidNum}
	if o.Read(topic) == nil {
		oldCate = topic.Category
		_, err = o.Delete(topic)
		if err != nil {
			return err
		}
	}

	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title", topic.Category).One(cate)
		if err == nil {
			cate.TopicCount--
			_, err = o.Update(cate)
		}
	}

	return err
}

// MonifyTopic 修改文章
func MonifyTopic(tid, title, category, labels, content, attachment string) error {
	// 处理标签
	mlabels := "$" + strings.Join(strings.Split(labels, " "), "#$") + "#"

	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	var oldCate string
	var oldAttach string
	o := orm.NewOrm()
	topic := &Topic{ID: tidNum}
	if err = o.Read(topic); err == nil {
		oldCate = topic.Category
		oldAttach = topic.Attachment
		topic.Title = title
		topic.Category = category
		topic.Labels = mlabels
		topic.Content = content
		topic.Attachment = attachment
		topic.Updated = time.Now()
		o.Update(topic)
	}
	// 更新分类统计
	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title", oldCate).One(cate)
		if err != nil {
			beego.Error(err)
		} else {
			cate.TopicCount--
			_, err = o.Update(cate)
		}
	}

	if len(oldAttach) > 0 {
		os.Remove(path.Join("attachment", oldAttach))
	}

	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	if err != nil {
		beego.Error(err)
	} else {
		cate.TopicCount++
		_, err = o.Update(cate)
	}
	return err
}

// AddReply 添加评论
func AddReply(tid, nickname, conent string) error {
	tidNum, err := strToInt64(tid)
	if err != nil {
		return err
	}

	reply := &Comment{
		Tid:     tidNum,
		Name:    nickname,
		Content: conent,
		Created: time.Now(),
	}
	o := orm.NewOrm()
	_, err = o.Insert(reply)
	if err != nil {
		return err
	}

	topic := &Topic{ID: tidNum}
	if o.Read(topic) == nil {
		topic.ReplyCount++
		topic.ReplyTime = time.Now()
		_, err = o.Update(topic)
	}
	return err
}

// GetAllReplies 获取所有评论
func GetAllReplies(tid string) (replies []*Comment, err error) {
	tidNum, err := strToInt64(tid)
	if err != nil {
		return nil, err
	}
	replies = make([]*Comment, 0)
	o := orm.NewOrm()
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidNum).All(&replies)
	return replies, err
}

// DelReply 删除评论
func DelReply(rid string) error {
	ridNum, err := strToInt64(rid)
	if err != nil {
		return err
	}
	o := orm.NewOrm()

	var tidNum int64
	reply := &Comment{ID: ridNum}
	if o.Read(reply) == nil {
		tidNum = reply.Tid
		_, err = o.Delete(reply)
		if err != nil {
			return err
		}
	}

	replies := make([]*Comment, 0)
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidNum).OrderBy("-created").All(&replies)
	if err != nil {
		return err
	}
	topic := &Topic{ID: tidNum}
	if o.Read(topic) == nil {
		if len(replies) > 0 {
			topic.ReplyTime = replies[0].Created
		}
		topic.ReplyCount = int64(len(replies))
		_, err = o.Update(topic)
	}

	return err
}

func strToInt64(nstr string) (int64, error) {
	n, err := strconv.ParseInt(nstr, 10, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}
