package profile

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func Create() {
	db, err := sql.Open("sqlite3", "./base.db")
	CheckErr(err)
	defer db.Close()

	sqlStmt := `
	create table box (
			id text not null primary key,
		 	boxName text,
		  	userName text,
		  	password text,
		  	scenarios text,
		  	remarks text
		);
	create table files (
			id text not null primary key,
		 	fileName text,
		  	serveId text,
		  	serverFileName text
		);
	`
	_, err = db.Exec(sqlStmt)
	CheckErr(err)
}

func InsertBox(box BoxMsg) {
	db, err := sql.Open("sqlite3", "./base.db")
	CheckErr(err)
	defer db.Close()

	stmt, err := db.Prepare("insert into box(id, boxName, userName, password, scenarios, remarks) values(?, ?, ?, ?, ?, ?)")
	CheckErr(err)
	_, err = stmt.Exec(box.Id, box.BoxName, box.UserName, box.Password, box.Scenarios, box.Remarks)
	CheckErr(err)

}

func UpdateBox(box BoxMsg) {
	db, err := sql.Open("sqlite3", "./base.db")
	CheckErr(err)
	defer db.Close()

	stmt, err := db.Prepare("update box set boxName=?, userName=?, password=?, scenarios=?, remarks=? where id=?")
	CheckErr(err)

	_, err = stmt.Exec(box.BoxName, box.UserName, box.Password, box.Scenarios, box.Remarks, box.Id)
	CheckErr(err)
}

func QueryBox(id string) BoxMsg {
	db, err := sql.Open("sqlite3", "./base.db")
	CheckErr(err)
	defer db.Close()

	rows, err := db.Query("select * from box where id=?", id)
	CheckErr(err)

	var boxName, userName, password, scenarios, remarks string
	for rows.Next() {
		err = rows.Scan(&id, &boxName, &userName, &password, &scenarios, &remarks)
		CheckErr(err)
	}

	return BoxMsg{
		Id: id, BoxName: boxName, UserName: userName, Password: password, Scenarios: scenarios, Remarks: remarks,
	}
}
func QueryBoxList() []BoxMsg {
	db, err := sql.Open("sqlite3", "./base.db")
	CheckErr(err)
	defer db.Close()

	rows, err := db.Query("select * from box")
	CheckErr(err)

	var id, boxName, userName, password, scenarios, remarks string
	var boxList []BoxMsg
	for rows.Next() {
		err = rows.Scan(&id, &boxName, &userName, &password, &scenarios, &remarks)
		CheckErr(err)
		boxList = append(boxList, BoxMsg{
			Id: id, BoxName: boxName, UserName: userName, Password: password, Scenarios: scenarios, Remarks: remarks,
		})
	}
	return boxList
}

func DeleteBox(id string) {
	db, err := sql.Open("sqlite3", "./base.db")
	CheckErr(err)
	defer db.Close()
	stmt, err := db.Prepare("delete from box where id=?")
	CheckErr(err)
	_, err = stmt.Exec(id)
	CheckErr(err)
}

func InsertFiles(file FileItem) {
	db, err := sql.Open("sqlite3", "./base.db")
	CheckErr(err)
	defer db.Close()

	stmt, err := db.Prepare("insert into files(id, fileName, serveId, serverFileName) values(?, ?, ?, ?)")
	CheckErr(err)
	_, err = stmt.Exec(file.Id, file.FileName, file.ServeId, file.ServerFileName)
	CheckErr(err)
}

func QueryFiles(serveerid string) []FileItem {
	db, err := sql.Open("sqlite3", "./base.db")
	CheckErr(err)
	defer db.Close()

	rows, err := db.Query("select * from files where serveId=?", serveerid)
	CheckErr(err)

	var id, fileName, serveId, serverFileName string
	var fileList []FileItem
	for rows.Next() {
		err = rows.Scan(&id, &fileName, &serveId, &serverFileName)
		CheckErr(err)
		fileList = append(fileList, FileItem{
			Id: id, FileName: fileName, ServeId: serveId, ServerFileName: serverFileName,
		})
	}
	return fileList
}

func UpdateFiles(file FileItem) {
	db, err := sql.Open("sqlite3", "./base.db")
	CheckErr(err)
	defer db.Close()

	stmt, err := db.Prepare("update files set fileName=?, serveId=?, serverFileName=? where id=?")
	CheckErr(err)

	_, err = stmt.Exec(file.FileName, file.ServeId, file.ServerFileName, file.Id)
	CheckErr(err)
}

func DeleteFiles(serverId string) {
	db, err := sql.Open("sqlite3", "./base.db")
	CheckErr(err)
	defer db.Close()
	stmt, err := db.Prepare("delete from files where serveId=?")
	CheckErr(err)
	_, err = stmt.Exec(serverId)
	CheckErr(err)
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
