package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
	"upement/cryp"
	"upement/profile"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Static("/static", "./static")
	r.Static("/storage", "./storage")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", home)
	r.GET("/get-box-list", get_box_list)
	r.GET("/dfile", dfile)

	r.POST("/create-box", create_box)
	r.POST("/create-box-file", create_box_file)
	r.POST("/preview", preview)
	r.POST("/get-det", get_det)

	r.DELETE("/box", delete_box)
	r.PUT("/box", put_box)
	println(`
			欢迎使用
	首页地址: http://localhost:2623/
	免责声明:
		本软件免费开源，无木马后门。
		本程序仅用于学习与研究，不保证其准确性。
		本程序的作者不对任何问题及事件负责，请使用者自行承担风险。
		本程序不可用于商业用途。
		如不能接受请自行删除本程序。
	更多信息请访问: https://www.opty.fun/
		
	`)
	// 执行cmd命令
	cmd := exec.Command("start","http://localhost:2623/")
	cmd.Run()
	r.Run(":2623")
}

func home(c *gin.Context) {

	c.HTML(http.StatusOK, "index.html", nil)
}

func preview(c *gin.Context) {

	type Pre struct {
		Id       string `json:"id"`
		FileName string `json:"fileName"`
	}
	var pre Pre
	c.BindJSON(&pre)

	if strings.HasSuffix(pre.FileName, ".jpg") || strings.HasSuffix(pre.FileName, ".png") || strings.HasSuffix(pre.FileName, ".jpeg") {
		buty := profile.ReadFile("storage/" + pre.FileName)
		img := cryp.Decrypt(buty, "key")
		c.JSON(http.StatusOK, gin.H{
			"code": 1024,
			"type": "image",
			"data": img,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 1024,
			"type": "file",
			"url":  "/dfile?fn=" + pre.FileName,
		})
	}
}
func dfile(c *gin.Context) {

	fn := c.Query("fn")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+fn)
	data := profile.ReadFile("storage/" + fn)
	file := cryp.Decrypt(data, "key")
	c.Writer.Write(file)

}

func put_box(c *gin.Context) {

	recv := profile.BoxMsg{}
	c.BindJSON(&recv)

	if recv.BoxName == "" {
		code_to_cli(*c, 1025, "box name is empty")
		return
	}
	profile.UpdateBox(recv)

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}

func get_det(c *gin.Context) {

	id := c.PostForm("id")
	boxitem := profile.QueryBox(id)

	fitem := profile.QueryFiles(id)

	c.JSON(http.StatusOK, gin.H{
		"code":    1024,
		"message": "success",
		"data": gin.H{
			"box":  boxitem,
			"file": fitem,
		},
	})
}

func delete_box(c *gin.Context) {

	delList := c.PostForm("delList")

	delList = strings.Replace(delList, "[", "", -1)
	delList = strings.Replace(delList, "]", "", -1)
	delList = strings.Replace(delList, "'", "", -1)
	delList = strings.Replace(delList, `"`, "", -1)
	list := strings.Split(delList, ",")

	for _, clId := range list {
		for _, v := range profile.QueryFiles(clId) {
			os.Remove("storage/" + v.ServerFileName)
		}
		profile.DeleteFiles(clId)
		profile.DeleteBox(clId)
	}
	code_to_cli(*c, 1024, "success")
}

func get_box_list(c *gin.Context) {

	result := profile.QueryBoxList()
	c.JSON(http.StatusOK, gin.H{
		"code":    1024,
		"data":    result,
		"message": "success",
	})
}

func create_box(c *gin.Context) {

	recv := profile.BoxMsg{}
	c.BindJSON(&recv)

	if recv.BoxName == "" {
		code_to_cli(*c, 1025, "box name is empty")
		return
	}

	box_item := profile.QueryBox(recv.Id)

	if box_item.BoxName != "" {
		code_to_cli(*c, 1026, "box already exists")
		return
	}

	profile.InsertBox(recv)
	code_to_cli(*c, 1024, "success")
}
func create_box_file(c *gin.Context) {

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		code_to_cli(*c, 1025, "error")
		return
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		code_to_cli(*c, 1025, "error")
		return
	}

	fileId := c.PostForm("id")
	fileName := c.PostForm("fileName")

	box_itm := profile.QueryBox(fileId)
	if box_itm.BoxName == "" {
		code_to_cli(*c, 1030, "Box does not exist. Please create a box first")
		return
	}

	ext := path.Ext(fileName)
	var file_name string
	file_name = rd_file_name() + ext
	imgFiles, _ := ioutil.ReadDir("storage")
	for _, v := range imgFiles {
		if v.Name() == file_name {
			file_name = rd_file_name() + ext
		}
	}

	cryp.Encrpyt(buf.Bytes(), "key", "storage/"+file_name)

	profile.InsertFiles(profile.FileItem{
		Id:             rd_str16(),
		FileName:       fileName,
		ServeId:        fileId,
		ServerFileName: file_name,
	})
	code_to_cli(*c, 1024, "success")
}

func code_to_cli(r gin.Context, code uint16, msg string) {
	r.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": msg,
	})
}
func rd_str16() string {

	var str = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var bytes = []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 16; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
func rd_file_name() string {

	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var bytes = []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 20; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
