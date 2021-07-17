package controllers

import (
	"context"
	"encoding/json"
	"github.com/beego/beego/v2/client/cache"
	_ "github.com/beego/beego/v2/client/cache/redis"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"ihome/models"
	"time"
)

type AreaController struct {
	beego.Controller
}

func (c *AreaController) retData(data interface{}) {
	c.Data["json"] = data
	c.ServeJSON()
}
func (c *AreaController) GetArea() {
	// 定义返回用的map
	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	areas := make([]models.Area, 0, 20)
	defer c.retData(&resp)
	rconn, err := cache.NewCache("redis", `{"key":"ihome","conn":"192.168.2.110:6379","dbNum":"0"}`)
	if err != nil {
		logs.Error("get redis conn failed err =", err)
		resp["errno"] = models.RECODE_REDISERR
		resp["errmsg"] = models.RecodeText(models.RECODE_REDISERR)
		return
	}
	area, err := rconn.Get(context.TODO(), "area")
	if err != nil {
		logs.Error("get redis failed err =", err)
		resp["errno"] = models.RECODE_REDISERR
		resp["errmsg"] = models.RecodeText(models.RECODE_REDISERR)
		return
	}
	if area != nil {
		logs.Info("get areaData from redis.")
		err := json.Unmarshal(area.([]byte), &areas)
		if err != nil {
			logs.Error("unmarshal redis data failed err =", err)
		} else {
			resp["data"] = &areas
			return
		}
	}
	o := orm.NewOrm()
	_, err = o.QueryTable("area").All(&areas)
	if err != nil {
		logs.Error("select area failed err =", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	//logs.Info("areas = ", areas)
	putarea, err := json.Marshal(&areas)
	if err != nil {
		logs.Error("json.Marshal(&areas) failed err =", err)
		resp["errno"] = models.RECODE_SERVERERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SERVERERR)
		return
	}
	err = rconn.Put(context.TODO(), "area", putarea, time.Hour*2400)
	if err != nil {
		logs.Error("put redis failed err =", err)
		resp["errno"] = models.RECODE_REDISERR
		resp["errmsg"] = models.RecodeText(models.RECODE_REDISERR)
		return
	}
	logs.Info("get areaData from mysql.")
	resp["data"] = &areas
}
