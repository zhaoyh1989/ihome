package controllers

import (
	"encoding/json"
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"ihome/models"
	"strconv"
	"time"
)

type HouseController struct {
	beego.Controller
}

func (c *HouseController) retData(data interface{}) {
	c.Data["json"] = data
	c.ServeJSON()
}

// GetHouse
func (c *HouseController) GetHouse() {
	resp := make(map[string]interface{}, 6)
	defer c.retData(&resp)
	userid := c.GetSession("userid")
	if userid == nil {
		logs.Error("get session userid failed")
		resp["errno"] = models.RECODE_SESSIONERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SESSIONERR)
		return
	}
	houses := []*models.House{}
	o := orm.NewOrm()
	sql := `select id, user_id, area_id, title, price, address, room_count, acreage, unit, capacity, beds, deposit, min_days, max_days, order_count, index_image_url, ctime
from house where user_id = ?;`
	rows, err := o.Raw(sql, userid).QueryRows(&houses)
	if err != nil {
		logs.Error("select house failed err = ", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	if rows == 0 {
		logs.Info("select house nodatafound!")
		resp["errno"] = models.RECODE_NODATA
		resp["errmsg"] = models.RecodeText(models.RECODE_NODATA)
		return
	}
	user := models.User{
		Id: userid.(int),
	}
	err = o.Read(&user)
	if err != nil {
		logs.Info("select user failed! err = ", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	for _, house := range houses {
		house.User = &user
		area := models.Area{
			Id: house.Area.Id,
		}
		err := o.Read(&area)
		if err != nil {
			logs.Info("select area failed! err = ", err)
			resp["errno"] = models.RECODE_DBERR
			resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
			return
		}
		house.Area = &area
		logs.Debug("get house user is %#v", house.User)
		logs.Debug("get house area is %#v", house.Area)
	}

	logs.Debug("read house sucess!!")
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	housesMap := make(map[string]*[]*models.House)
	housesMap["houses"] = &houses
	resp["data"] = housesMap
}

// PostHouse
func (c *HouseController) PostHouse() {
	resp := make(map[string]interface{}, 6)
	defer c.retData(&resp)
	house := make(map[string]interface{}, 16)
	//house := models.House{}
	err2 := json.Unmarshal(c.Ctx.Input.RequestBody, &house)
	if err2 != nil {
		logs.Error("获取请求房屋数据失败")
		resp["errno"] = models.RECODE_SERVERERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SERVERERR)
		return
	}
	logs.Debug("请求中的房屋信息：%#v", house)
	o := orm.NewOrm()
	b, err := o.Begin()
	if err != nil {
		logs.Error("开始事务失败，err = ", err)
	}
	userid := c.GetSession("userid")
	if userid == nil {
		logs.Error("get session userid failed")
		resp["errno"] = models.RECODE_SESSIONERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SESSIONERR)
		return
	}
	sql := `insert into house (user_id, area_id, title, price, address, room_count, acreage, unit, capacity, beds, deposit, min_days, max_days, ctime)
values (?,?,?,?,?,?,?,?,?,?,?,?,?,?);`
	exec, err2 := b.Raw(sql, userid, house["area_id"], house["title"], house["price"], house["address"], house["room_count"],
		house["acreage"], house["unit"], house["capacity"], house["beds"], house["deposit"], house["min_days"], house["max_days"],
		time.Now()).Exec()
	if err2 != nil {
		logs.Error("insert house failed err = ", err2)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		b.Rollback()
		return
	}
	houseid, err2 := exec.LastInsertId()
	if err2 != nil {
		resp["errno"] = models.RECODE_UNKNOWERR
		resp["errmsg"] = models.RecodeText(models.RECODE_UNKNOWERR)
		b.Rollback()
		return
	}
	logs.Debug("插入的房屋id = ", houseid)
	sql = `insert into facility_houses (facility_id, house_id)
values (?,?);`
	for _, i2 := range house["facility"].([]interface{}) {
		_, err2 = b.Raw(sql, i2, houseid).Exec()
		if err2 != nil {
			logs.Error("insert facility_houses failed err =", err2)
			resp["errno"] = models.RECODE_DBERR
			resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
			b.Rollback()
			return
		}
	}
	b.Commit()
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	housesMap := make(map[string]int64)
	housesMap["house_id"] = houseid
	resp["data"] = housesMap
}

// PostHouseImages  上传房屋图片
func (c *HouseController) PostHouseImages() {
	resp := make(map[string]interface{}, 6)
	defer c.retData(&resp)
	userid := c.GetSession("userid")
	if userid == nil {
		logs.Error("获取session中的userid失败")
		resp["errno"] = models.RECODE_SESSIONERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SESSIONERR)
		return
	}
	houseid, err := c.GetInt("house_id")
	if err != nil {
		logs.Error("获取前端给的house_id失败，err=", err)
		resp["errno"] = models.RECODE_SERVERERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SERVERERR)
		return
	}
	logs.Debug("获取到的hoseid = ", houseid)
	file, header, err := c.GetFile("house_image")
	if err != nil {
		logs.Error("获取前端给的house_image失败，err=", err)
		resp["errno"] = models.RECODE_SERVERERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SERVERERR)
		return
	}
	b := make([]byte, header.Size)
	_, err = file.Read(b)
	if err != nil {
		logs.Error("file.Read(b)失败，err=", err)
		resp["errno"] = models.RECODE_SERVERERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SERVERERR)
		return
	}
	fileurl, err := models.UploadByBuffer(b, header.Filename)
	if err != nil {
		logs.Error("models.UploadByBuffer(b, header.Filename)失败，err=", err)
		resp["errno"] = models.RECODE_SERVERERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SERVERERR)
		return
	}
	logs.Info("fileurl = ", fileurl)
	o := orm.NewOrm()
	sql := `update house
set index_image_url = ?
where id = ?;`
	_, err = o.Raw(sql, fileurl, houseid).Exec()
	if err != nil {
		logs.Error("o.Raw(sql, fileurl, houseid).Exec()失败，err=", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	urlMap := make(map[string]string, 1)
	urlMap["url"] = "http://192.168.2.110/" + fileurl
	logs.Debug("url = ", urlMap["url"])
	resp["data"] = &urlMap
}

// GetDetailHouseData  获取房屋的详细信息
func (c *HouseController) GetDetailHouseData() {
	resp := make(map[string]interface{}, 6)
	defer c.retData(&resp)
	userid := c.GetSession("userid")
	if userid == nil {
		logs.Error("get session userid failed")
		resp["errno"] = models.RECODE_SESSIONERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SESSIONERR)
		return
	}
	houseid := c.Ctx.Input.Param(":houseid")
	house := models.House{}
	o := orm.NewOrm()
	sql := `select id, user_id, area_id, title, price, address, room_count, acreage, unit, capacity, beds, deposit, min_days, max_days, order_count, index_image_url, ctime
from house where id = ?;`
	hid, err := strconv.Atoi(houseid)
	if err != nil {
		logs.Error("strconv.Atoi(houseid) failed, err = ", err)
		resp["errno"] = models.RECODE_SERVERERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SERVERERR)
		return
	}
	err = o.Raw(sql, hid).QueryRow(&house)
	if err != nil {
		logs.Error("select house failed err = ", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	user := models.User{
		Id: userid.(int),
	}
	err = o.Read(&user)
	if err != nil {
		logs.Info("select user failed! err = ", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}

	house.User = &user
	area := models.Area{
		Id: house.Area.Id,
	}
	err = o.Read(&area)
	if err != nil {
		logs.Info("select area failed! err = ", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	house.Area = &area
	logs.Debug("get house user is %#v", house.User)
	logs.Debug("get house area is %#v", house.Area)

	logs.Debug("read house sucess!!")
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	housesMap := make(map[string]*models.House)
	housesMap["house"] = &house
	resp["data"] = housesMap
}
