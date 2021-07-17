package controllers

import (
	"encoding/json"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"ihome/models"
	"regexp"
)

type UserController struct {
	beego.Controller
}

func (c *UserController) retData(data interface{}) {
	c.Data["json"] = data
	c.ServeJSON()
}
func (c *UserController) Reg() {
	// 定义请求的map
	req := make(map[string]interface{}, 16)
	// 定义返回数据的map
	resp := make(map[string]interface{}, 16)
	defer c.retData(&resp)
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		logs.Error("json.Unmarshal(c.Ctx.Input.RequestBody, &req) failed err =", err)
		resp["errno"] = models.RECODE_SERVERERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SERVERERR)
		return
	}
	ok, err := regexp.MatchString("^1[3|4|5|7|8][0-9]{9}$", req["mobile"].(string))
	if !ok {
		logs.Error("手机号不合规")
		resp["errno"] = models.RECODE_MOBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_MOBERR)
		return
	}
	user := models.User{Name: req["mobile"].(string), Password_hash: req["password"].(string), Mobile: req["mobile"].(string)}
	o := orm.NewOrm()
	id, err := o.Insert(&user)
	if err != nil {
		logs.Error("insert user failed err =", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	c.SetSession("name", user.Name)
	c.SetSession("userid", id)
}

// Login
func (c *UserController) Login() {
	// 定义请求的map
	req := make(map[string]interface{}, 16)
	// 定义返回数据的map
	resp := make(map[string]interface{}, 16)
	defer c.retData(&resp)
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		logs.Error("json.Unmarshal(c.Ctx.Input.RequestBody, &req) failed err =", err)
		resp["errno"] = models.RECODE_SERVERERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SERVERERR)
		return
	}
	//regexp.Compile("^1[3|4|5|7|8][0-9]{9}$", req["mobile"].(string))
	ok, err := regexp.MatchString("^1[3|4|5|7|8][0-9]{9}$", req["mobile"].(string))
	if !ok {
		logs.Error("手机号不合规")
		resp["errno"] = models.RECODE_MOBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_MOBERR)
		return
	}
	user := models.User{Password_hash: req["password"].(string), Mobile: req["mobile"].(string)}
	o := orm.NewOrm()
	err = o.Read(&user, "mobile", "Password_hash")
	if err != nil {
		logs.Error("Read user failed err =", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	c.SetSession("name", user.Name)
	c.SetSession("userid", user.Id)
}
func (c *UserController) Postavatar() {
	resp := make(map[string]interface{}, 16)
	defer c.retData(&resp)
	id := c.GetSession("userid")
	if id == nil {
		logs.Error("未获取到session")
		resp["errno"] = models.RECODE_SESSIONERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SESSIONERR)
		return
	}
	file, header, err := c.GetFile("avatar")
	if err != nil {
		logs.Error("获取前端提交的文件avatar失败")
		resp["errno"] = models.RECODE_SERVERERR
		resp["errmsg"] = models.RecodeText(models.RECODE_SERVERERR)
		return
	}
	fb := make([]byte, header.Size)
	file.Read(fb)
	logs.Debug("header = ", header.Filename)
	//logs.Info("file = ", fb)
	fileid, err := models.UploadByBuffer(fb, header.Filename)
	o := orm.NewOrm()
	sql := `update user
	set avatar_url = ?
	where id = ?`
	_, err = o.Raw(sql, fileid, id).Exec()
	logs.Debug("提交头像更新数据库成功，更新userid=%v", id)
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)
	avatarMap := make(map[string]string, 1)
	avatarMap["avatar_url"] = "http://192.168.2.110/" + fileid
	resp["data"] = avatarMap
}
