package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"ihome/models"
)

type HousesIndexController struct {
	beego.Controller
}

func (c *HousesIndexController) retData(data interface{}) {
	c.Data["json"] = data
	c.ServeJSON()
}
func (c *HousesIndexController) HousesIndex() {
	// 定义返回用的map
	resp := make(map[string]interface{})
	defer c.retData(resp)
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
}
