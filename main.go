package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/thedevsaddam/gojsonq"
	"io/ioutil"
	"log"
	"net/http"
)

type Product struct {
	Type             int    `json:"type" `
	ProductNo        string `json:"product_no" `
	Amount           int    `json:"amount"`
	ElectricityHours int    `json:"electricity_hours" `
	HashRateNo       string `json:"hashrate_no" `
}

type CreateOrderParams struct {
	Products    []Product `json:"products"`
	UniqueToken string    `json:"unique_token" `
}

type Response struct {
	Code string
	Type int
}

func main() {
	//获取token
	uniqueToken := GetToken()
	//获取商品类型
	t := GetGoodsInfo()
	//这里是个大坑,通过reflect.Typeof(t)发现类型是t的实际类型是*Response,所以类型断言一定要是指针类型的！
	tt := t.(*Response)

	//构造需要传递的参数
	params := new(CreateOrderParams)

	//给构造好的结构体赋值
	params.UniqueToken = uniqueToken
	var p Product
	p.Type = tt.Type
	p.ProductNo = tt.Code
	p.Amount = 10
	p.ElectricityHours = 240
	var products []Product
	products = append(products, p)
	params.Products = products

	//将结构体序列化，转化为json格式，返回的结果是[]byte类型
	res, err := json.Marshal(&params)
	if err != nil {
		println(err.Error())
	}

	//提交订单
	SubmitOrder(res)
}

// GetGoodsInfo 获取商品信息
func GetGoodsInfo() interface{} {
	//声明要调用获取商品信息的接口
	reqGoodsUrl := "https://dev.cookiehash.org:31145/v1/products/HSBOX-10days"

	//调用接口
	req, err := http.NewRequest("GET", reqGoodsUrl, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	//打印返回结果
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	// Todo:如何将结果转换成映射到Products
	jsonX := string(body)
	t := gojsonq.New().FromString(jsonX).Find("type")
	//Todo: 要点！！！接口类型的 float64 转换成int
	//reflect.TypeOf(t) => float64
	tt := int(t.(float64))

	c := gojsonq.New().FromString(jsonX).Find("code")
	code := c.(string)
	fmt.Println("test", code)

	res := new(Response)
	res.Type = tt
	res.Code = code

	return res

}

// GetToken 获取token
func GetToken() string {
	//声明要获取unique_token的接口
	reqUniqueTokenUrl := "https://dev.cookiehash.org:31145/v1/users/token"

	//设置get请求类型,所以第三个参数为nil
	req, err := http.NewRequest("GET", reqUniqueTokenUrl, nil)
	if err != nil {
		log.Println(err)
	}

	//设置请求头header,携带token身份认证
	req.Header.Set("authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiMzdlYWMwMTMtOTM5Yy00NThjLTgxMDktYmM5MDFjMTIyY2I0IiwiZXhwIjoxNzMzODQ2NTEzfQ._p6nCjIx1nq6sbMa4B-yJ9P_vThJESsF5tLIRoMRZXA")

	//获取客户端对象，执行请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		panic("ops!,get unique token url err!")
	}
	defer resp.Body.Close()

	//打印结果
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		panic("ops!,get unique token url err!")
	}
	fmt.Println(string(body))

	//在返回的JSON对象中根据key获取value
	jsonX := string(body)
	t := gojsonq.New().FromString(jsonX).Find("unique_token")
	//由于返回的t是一个空接口类型，需要类型断言转换成string类型
	return t.(string)
}

// SubmitOrder 提交订单
func SubmitOrder(params []byte) {
	//声明下单的url
	PostOrderUrl := "https://dev.cookiehash.org:31145/v1/orders"

	//设置post请求,第三个参数传byte类型
	req, err := http.NewRequest("POST", PostOrderUrl, bytes.NewBuffer(params))
	if err != nil {
		log.Println(err)
	}

	//授权
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiMzdlYWMwMTMtOTM5Yy00NThjLTgxMDktYmM5MDFjMTIyY2I0IiwiZXhwIjoxNzMzODQ2NTEzfQ._p6nCjIx1nq6sbMa4B-yJ9P_vThJESsF5tLIRoMRZXA")

	//获取客户端对象，发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	//读取返回值
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(res))
}
