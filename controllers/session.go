package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"ihome/models"
)

type SessionController struct {
	beego.Controller
}

func (c *SessionController) retData(data interface{}) {
	c.Data["json"] = data
	c.ServeJSON()
}
func (c *SessionController) GetSessionData() {
	// 定义返回用的map
	resp := make(map[string]interface{})
	defer c.retData(&resp)
	sname := c.GetSession("name")
	name := map[string]interface{}{"name": sname}
	if sname == nil {
		logs.Info("get session name is nil")
		resp["errno"] = models.RECODE_SESSIONERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SESSIONERR)
		return
	}
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	resp["data"] = &name
}

// DelSessionData
func (c *SessionController) DelSessionData() {
	// 定义返回用的map
	resp := make(map[string]interface{})
	defer c.retData(&resp)
	c.DelSession("name")
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
}
