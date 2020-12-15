package load

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"loadOnlineData/utils"
	"sort"
	"strconv"
	"strings"
	"time"
)

func LoadData(serverName string,openServerTime int64) (newServerId int,NewIp string,endServername string) {
	db,err:= sql.Open("mysql","root:Qwert123!@tcp(172.25.128.114)/gm")
	utils.CheckErr(err)
	defer db.Close()
	rows,err := db.Query("select serverId,serverName,showSid,ip from server_list;")
	utils.CheckErr(err)

	serverIdList := make([]int,0)
	for rows.Next() {
		var serverId int32
		var serverName string
		var showSid int32
		var ip string
		err := rows.Scan(&serverId,&serverName,&showSid,&ip)
		utils.CheckErr(err)
		s := strconv.Itoa(int(serverId))
		if strings.HasPrefix(s,"100042") {
			serverIdList = append(serverIdList,int(serverId))
		}
		//fmt.Println(serverId,serverName,showSid,ip)
	}
	sort.Ints(serverIdList)
	newServerId = serverIdList[len(serverIdList)-1]+1

	rows2,err := db.Query("select ip from server_list where serverId = ?",serverIdList[len(serverIdList)-1])
	utils.CheckErr(err)
	var ip string
	for rows2.Next() {
		err := rows2.Scan(&ip)
		utils.CheckErr(err)
	}
	ipInt := utils.Ip2Long(ip)
	NewIp = utils.BackToIP4(int64(ipInt)+1)

	endServername = serverName + "(" + time.Now().Format("01")+"月"+time.Now().Format("02")+")"

	sqlStr := "insert into server_list(serverId,serverName,showSid,opgameId,ip,grpcPort,serverKey,logServer,firstOpenTime,isQA,maxOnline,maxLoginQ,maxUser,webAddr) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_,err = db.Exec(sqlStr,newServerId,endServername,newServerId,1000,NewIp,"10000","WDxG4oDqbVdstCjB","172.25.128.114",openServerTime,1,111111,111111,111111,"http://web.sanguo.bj/gmc/go/")
	utils.CheckErr(err)
	fmt.Println("数据写入成功!")
	fmt.Println("serverID:\t",newServerId)
	fmt.Println("ip:\t",NewIp)
	fmt.Println("serverName:\t",endServername)
	return
}

type ServerData struct {
	ServerId int32 `json:"serverId"`
	ServerName string `json:"serverName"`
	ShowSid int32 `json:"showSid"`
	OpgameId int32 `json:"opgameId"`
	Ip string `json:"ip"`
	GrpcPort string `json:"grpcPort"`
	NeedCode int32 `json:"needCode"`
	ServerKey string `json:"serverKey"`
	LogServer string `json:"logServer"`
	FirstOpenTime int32 `json:"firstOpenTime"`
	MaintainBeginTime int32 `json:"maintainBeginTime"`
	MaintainEndTime int32 `json:"maintainEndTime"`
	WhiteList string `json:"whiteList"`
	WhiteListOpen int32 `json:"whiteListOpen"`
	Status int32 `json:"status"`
	IsQA int32 `json:"isQA"`
	MaxOnline int32 `json:"maxOnline"`
	MaxLoginQ int32 `json:"maxLoginQ"`
	MaxUser int32 `json:"maxUser"`
	Available int32 `json:"available"`
	Param string `json:"param"`
	Debug int32 `json:"debug"`
	NeedActivation int32 `json:"needActivation"`
	WebAddr string `json:"webAddr"`
}

var S *ServerData

func QueryServerData(serverId int32) (value []byte) {
	db,err:= sql.Open("mysql","root:Qwert123!@tcp(172.25.128.114)/gm")
	utils.CheckErr(err)
	defer db.Close()
	rows,err := db.Query("select * from server_list where serverId = ?",serverId)
	if err != nil{
		panic(err)
	}

	S = &ServerData{}
	for rows.Next() {
		rows.Scan(&S.ServerId,&S.ServerName,&S.ShowSid,&S.OpgameId,&S.Ip,&S.GrpcPort,&S.NeedCode,&S.ServerKey,&S.LogServer,&S.FirstOpenTime,
			&S.MaintainBeginTime,&S.MaintainEndTime,&S.WhiteList,&S.WhiteListOpen,&S.Status,&S.IsQA,&S.MaxOnline,&S.MaxLoginQ,&S.MaxUser,
			&S.Available,&S.Param,&S.Debug,&S.NeedActivation,&S.WebAddr)
	}
	S.WebAddr = "http://web.sanguo.bj/gmc/go/"

	value,err = json.Marshal(S)
	utils.CheckErr(err)
	return

}




