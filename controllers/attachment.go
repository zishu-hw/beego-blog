package controllers

import (
	"io"
	"net/url"
	"os"

	"github.com/astaxie/beego"
)

// AttachController 附件控制器
type AttachController struct {
	beego.Controller
}

// Get 方法
func (c *AttachController) Get() {
	filePath, err := url.QueryUnescape(c.Ctx.Request.RequestURI[1:])
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}
	f, err := os.Open(filePath)
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}
	_, err = io.Copy(c.Ctx.ResponseWriter, f)
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}
}
