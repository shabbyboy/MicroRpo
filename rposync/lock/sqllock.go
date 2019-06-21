package lock

import (
	"database/sql"
	"errors"
	"sync"
)

/*
基于关系型数据库的分布式锁
 */

const (
	driverMySQL string = "mysql"
	driverPostgres = "postgres"

	dbusername = ""
	dbpassword = ""
	dbhostip = ""
	dbname = ""
	dataSource = dbusername+":"+dbpassword+"@tcp("+dbhostip+")/"+dbname+"?charset=utf8"
)

const(
	sqlDBLockSelect = iota
	sqlDBLockInsert
	sqlDBLockUpdate
)

var sqlDB *sql.DB

//func init(){
//	sqlDB,_ = sql.Open(driverMySQL,dataSource)
//	if err := sqlDB.Ping();err != nil{
//		return
//	}
//}

var sqlStmts = []string{
	"SELECT id, tick from StoreLock FOR UPDATE",      // sqlDBLockSelect
	"INSERT INTO StoreLock (id, tick) VALUES (?, ?)", // sqlDBLockInsert
	"UPDATE StoreLock SET id=?, tick=?",              // sqlDBLockUpdate
}
type SQLLock struct {
	sync.RWMutex
}

func (sl *SQLLock) sQLLock(steal bool) (bool,error){
	sl.Lock()

	var (
		lockId uint64
		tick uint64
		lock bool = false
	)

	tx,err := sqlDB.Begin()

	if err != nil{
		return false,err
	}

	id := getGID()

	row := tx.QueryRow(sqlStmts[sqlDBLockSelect])

	err = row.Scan(&lockId,&tick)

	if err != nil && err != sql.ErrNoRows{
		return false,err
	}

	if err == sql.ErrNoRows || steal || lockId == id || lockId == 0{
		if steal {
			tick = 0
		}
		stmt := sqlStmts[sqlDBLockUpdate]
		if err == sql.ErrNoRows{
			stmt = sqlStmts[sqlDBLockInsert]
		}

		if _, err := tx.Exec(stmt,id,tick+1); err != nil{
			return false,err
		}
		lock = true
	}else {
		return false,nil
	}

	if err = tx.Commit(); err != nil{
		return false,err
	}
	tx = nil


	defer func() {
		if tx != nil{
			tx.Rollback()
		}
		sl.Unlock()
	}()
	return lock,nil
}

func (sl *SQLLock) AquireLock(){
	lock, err := sl.sQLLock(false)

	if err != nil  {
		panic(err)
	}

	if lock == true {
		return
	}else {

		for {
			lock, err = sl.sQLLock(false)
			if err != nil {
				panic(err)
			}
			if lock == true {
				return
			}
		}
	}
}


func (sl *SQLLock) sQLUnLock() error {
	sl.Lock()
	defer sl.Unlock()
	var (
		lockId uint64
		tick uint64
	)
	tx,err := sqlDB.Begin()

	if err != nil{
		return err
	}
	row := tx.QueryRow(sqlStmts[sqlDBLockSelect])
	err = row.Scan(&lockId,&tick)

	if getGID() == lockId{
		tx.Exec(sqlStmts[sqlDBLockUpdate],0,0)

		if err := tx.Commit(); err != nil {
			return err
		}
	}else{
		tx.Rollback()
		return errors.New("unlock of unlocked mutex")
	}
	return nil
}

func (sl *SQLLock) ReleaseLock(){

	if err := sl.sQLUnLock(); err != nil {
		panic(err)
	}
}
