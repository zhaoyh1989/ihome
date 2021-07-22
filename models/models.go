package models

import (
	"github.com/beego/beego/v2/adapter/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func init() {
	// set default database
	orm.RegisterDataBase("default", "mysql", "root:root@tcp(192.168.2.104:3306)/lovehome?charset=utf8&loc=Local")

	// register model, new(House), new(Area), new(Facility), new(HouseImage), new(OrderHouse)
	orm.RegisterModel(new(User), new(House), new(Area), new(Facility), new(HouseImage), new(OrderHouse))

	// create table
	orm.RunSyncdb("default", false, false)
}

/* 用户 table_name = user */
type User struct {
	Id            int           `json:"user_id"`                    //用户编号
	Name          string        `json:"name" orm:"size(32);unique"` //用户昵称
	Password_hash string        `json:"password" orm:"size(128)"`   //用户密码加密的
	Mobile        string        `json:"mobile" orm:"size(11)"`      //手机号
	Real_name     string        `json:"real_name" orm:"size(32)"`   //真实姓名
	Id_card       string        `json:"id_card" orm:"size(20)"	`    //身份证号
	Avatar_url    string        `json:"avatar_url" orm:"size(256)"` //用户头像路径
	Houses        []*House      `json:"houses" orm:"reverse(many)"` //用户发布的房屋信息
	Orders        []*OrderHouse `json:"orders" orm:"reverse(many)"` //用户下的订单
}

/* 户层信息 table_name = house */
type House struct {
	Id              int           `json:"house_id"`                                          //房屋编号
	User            *User         `json:"user_id" orm:"rel(fk)"`                             //房屋主人的用户编号
	Area            *Area         `json:"area_id" orm:"rel(fk)"`                             //归属地的区域编号
	Title           string        `json:"title" orm:"size(64)"`                              //房屋标题
	Price           int           `json:"price" orm:"default(0)"`                            //单价，单位：分
	Address         string        `json:"address" orm:"size(512)" orm:"default("")"`         //地址
	Room_count      int           `json:"room_count" orm:"default(1)"	`                      //房间数目
	Acreage         int           `json:"acreage" orm:"default(0)"`                          //房屋总面积
	Unit            string        `json:"unit" orm:"size(32)" orm:"default("")"`             //房屋单元，如 几室几厅
	Capacity        int           `json:"capacity"orm:"default(1)"`                          //房屋容纳的总人数
	Beds            string        `json:"beds"orm:"size(64)"	orm:"default("")"`              //房屋床铺的位置
	Deposit         int           `json:"deposit" orm:"default(0)"`                          //押金
	Min_days        int           `json:"min_days" orm:"default(1)"`                         //最少入住天数
	Max_days        int           `json:"max_days" orm:"default(0)"`                         //最多住天数 0表示不限制
	Order_count     int           `json:"order_count" orm:"default(0)"`                      //预定完成的该房屋的订单数
	Index_image_url string        `json:"index_image_url" orm:"size(256)" orm:"default("")"` //房屋主图片路径
	Facilities      []*Facility   `json:"facilities" orm:"reverse(many)"`                    //房屋设施
	Images          []*HouseImage `json:"img_urls" orm:"reverse(many)"`                      //房屋的图片
	Orders          []*OrderHouse `json:"orders" orm:"reverse(many)"`                        //房屋的订单
	Ctime           time.Time     `json:"ctime" orm:"auto_now_add;type(datetime)"`
}

//首页最高展示的房屋数量
var HOME_PAGE_MAX_HOUSES int = 5

//房屋列表页面每页显示条目数
var HOUSE_LIST_PAGE_CAPACITY int = 2

/* 区域信息 table_name = area */
type Area struct {
	Id   int    `json:"aid"`                  //区域编号
	Name string `orm:"size(32)" json:"aname"` //区域名字
}

/* 设施信息 table_name = "facility" */
type Facility struct {
	Id     int      `json:"fid"`     //设施编号
	Name   string   `orm:"size(32)"` //设施名字
	Houses []*House `orm:"rel(m2m)"` //都有哪些房屋有此设施
}

/* 房屋图片 table_name = "house_image" */
type HouseImage struct {
	Id    int    `json:"house_image_id"`         //图片id
	Url   string `orm:"size(256)" json:"url"`    //图片url
	House *House `orm:"rel(fk)" json:"house_id"` //图片所属房屋编号
}

/* 订单状态常量 */
const (
	ORDER_STATUS_WAIT_ACCEPT  = "WAIT_ACCEPT"  //待接单
	ORDER_STATUS_WAIT_PAYMENT = "WAIT_PAYMENT" //待支付
	ORDER_STATUS_PAID         = "PAID"         //已支付
	ORDER_STATUS_WAIT_COMMENT = "COMMENT"      //待评价
	ORDER_STATUS_COMPLETE     = "COMPLETE"     //已完成
	ORDER_STATUS_CANCELED     = "CANCELED"     //已取消
	ORDER_STATUS_REJECTED     = "REJECTED"     //已拒单
)

/*  订单 table_name = order */
type OrderHouse struct {
	Id          int       `json:"order_id"`               //订单编号
	User        *User     `orm:"rel(fk)"	json:"user_id"`  //下单的用户编号
	House       *House    `orm:"rel(fk)"	json:"house_id"` //预定的房间编号
	Begin_data  time.Time `orm:"type(datetime)"`          //预定的起始时间
	End_data    time.Time `orm:"type(datetime)"`          //预定的结束时间
	Days        int       //预定总天数
	House_price int       //房屋的单价
	Amount      int       //订单总金额
	Status      string    `orm:"default(WAIT_ACCEPT)"`                     //订单状态
	Comment     string    `orm:"size(512)"`                                //订单评论
	Ctime       time.Time `orm:"auto_now_add;type(datetime)"	json:"ctime"` //
}
