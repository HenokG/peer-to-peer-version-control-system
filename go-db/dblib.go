package main
import (
"database/sql"
_ "github.com/go-sql-driver/mysql"
"log"
)
// global variables

// database connection related properties
var database  = "localsvn";
var user = "root";
var password = "";
var driver = "mysql";

//database pointer
var db *sql.DB
// last logged error
var err error

func connect(){
		// connnect to the database and set the pointer
	db, err = sql.Open(driver,fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s",user,password,database))	
	if err != nil {
		log.Fatal(err)
	}

}
func closeConeection(){
		// close the database connection
	db.Close();

}

func addRepository(name string, admin string) (int64, bool){
	connect()

	stmt, err := db.Prepare("INSERT INTO repository (name,admin) VALUES(?,?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(name, admin)
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
		//close the database connection
	closeConeection();

	return lastId, true;


}

func addCommit(repositoryId int, commitState string) (int64, bool){
	connect()

	stmt, err := db.Prepare("INSERT INTO commit (repid,state) VALUES(?,?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(repositoryId, commitState)
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
		//close the database connection
	closeConeection();

	return lastId, true;



}

func getRepository(admin string, name string) (int)  {
	connect()
	//table:="repository"
	rows, err := db.Query("select id from repository where admin=? and name=?",admin,name)
	if(err!=nil){
		log.Fatal(err)
		return -1
	}
	
	var (repid int)
	defer rows.Close()
	for rows.Next() {
		
		err := rows.Scan(&repid)
		if err != nil {
			log.Fatal(err)
		}
		
		
	}
	
	closeConeection()
	return repid

}
func getCommit(repid int) (int, string) {
	connect()
	rows, err := db.Query("select id,state from commit where repid = ?",repid)
	if(err!=nil){
		log.Fatal(err)
		return -1, ""
	}
	var (
			commitId int
			commitState string
		)
	defer rows.Close()
	for rows.Next() {
		
		err := rows.Scan(&commitId,&commitState)
		if err != nil {
			log.Fatal(err)
		}
		
		
	}
	
	closeConeection()
	return commitId, commitState

}