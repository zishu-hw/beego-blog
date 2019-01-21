package controllers

import (
	"beeblog/models"
	"path"
	"strings"

	"github.com/astaxie/beego"
)

// TopicController 文章
type TopicController struct {
	beego.Controller
}

// Get 方法
func (this *TopicController) Get() {
	this.Data["IsTopic"] = true
	this.TplName = "topic.html"
	this.Data["IsLogin"] = checkAccount(this.Ctx)

	topics, err := models.GetAllTopic("", "", false)
	if err != nil {
		beego.Error(err)
		return
	}
	this.Data["Topics"] = topics
}

func (this *TopicController) Post() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}
	// 解析表单
	tid := this.Input().Get("tid")
	title := this.Input().Get("title")
	category := this.Input().Get("category")
	labels := this.Input().Get("labels")
	content := this.Input().Get("content")

	// 获取附件
	beego.Info("getFile")
	_, fh, err := this.GetFile("attachment")
	beego.Info("fdafda")
	var attachment string
	if err != nil {
		beego.Error(err)
	} else {
		// 保存附件
		attachment = fh.Filename
		beego.Info(attachment)
		err = this.SaveToFile("attachment", path.Join("attachment", attachment))
		if err != nil {
			beego.Error(err)
		}
	}

	if len(tid) == 0 {
		err = models.AddTopic(title, category, labels, content, attachment)
	} else {
		err = models.MonifyTopic(tid, title, category, labels, content, attachment)
	}
	if err != nil {
		beego.Error(err)
	}
	this.Redirect("/topic", 302)
}

func (this *TopicController) Add() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}
	this.Data["IsLogin"] = checkAccount(this.Ctx)
	this.TplName = "topic_add.html"
}

func (this *TopicController) Modify() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}
	tid := this.Input().Get("tid")
	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		this.Redirect("/", 302)
		return
	}
	this.Data["IsLogin"] = true
	this.Data["Topic"] = topic
	this.Data["Tid"] = tid

	this.TplName = "topic_modify.html"
}

func (this *TopicController) Delete() {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}

	err := models.DelTopic(this.Input().Get("tid"))
	if err != nil {
		beego.Error(err)
	}
	this.Redirect("/topic", 302)
}

func (this *TopicController) View() {
	this.TplName = "topic_view.html"

	tid := this.Ctx.Input.Params()["0"]
	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		this.Redirect("/", 302)
		return
	}
	this.Data["Topic"] = topic
	this.Data["Labels"] = strings.Split(topic.Labels, " ")

	replies, err := models.GetAllReplies(tid)
	if err != nil {
		beego.Error(err)
		return
	}
	this.Data["Replies"] = replies
	this.Data["IsLogin"] = checkAccount(this.Ctx)
}
