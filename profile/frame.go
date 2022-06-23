package profile

type BoxMsg struct {
	Id        string `json:"id"`
	BoxName   string `json:"boxName"`
	UserName  string `json:"userName" `
	Password  string `json:"password" `
	Scenarios string `json:"scenarios" `
	Remarks   string `json:"remarks" `
}

type FileItem struct {
	Id             string `json:"id"`
	FileName       string `json:"fileName"`
	ServeId        string `json:"serveId"`
	ServerFileName string `json:"serverFileName"`
}
type List struct {
	Id string `json:"id"`
}
