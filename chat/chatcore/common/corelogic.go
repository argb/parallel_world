package common

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"log"
	"parallel_world/chat/myerrors"
	"parallel_world/chat/storage/persistentstore"
)

const (
	maxRoomLen = 30
	maxUnameLen = 30
)

// CheckRoomName 检查房间名字的合法性，不要太长
func CheckRoomName(name string) (error, bool) {
	if len(name) > maxRoomLen {
		err := myerrors.NewRoomError()
		return err, false
	}

	return nil, true
}

type UserInfo struct {
	NickName string
	PID peer.ID


	FavorRooms []string // 我收藏的房间
	HistoryRooms []string // 我曾经加入过的房间
	MyRooms []string // 我创建的房间
}

// 用户信息的进程内缓存
var userMap map[string]UserInfo

func NewUser() *UserInfo {
	userInfo := new(UserInfo)

	return userInfo
}

func init() {
	userMap = make(map[string]UserInfo)
}

// SavePid
// 先写缓存，在持久化到db
func SavePid(pid peer.ID, uname string) error{

	userMap[uname] = UserInfo{
		NickName: uname,
		PID: pid,
	}
	err := persistentstore.SavePID(pid, uname)

	return err
}

/*
// GetUserNameByPID
// 先看内存，没有从db里取
func GetUserNameByPID(pid peer.ID) string {
	if userInfo, ok := userMap[pid]; ok{
		if userInfo.NickName != "" {
			return userInfo.NickName
		}
	}

	return persistentstore.GetUserNameFromLocalStore(pid.Pretty())
}

 */

// GetPIDByUsername
// 先看内存，没有从db里取
func GetPIDByUsername(uname string) string {
	if userInfo, ok := userMap[uname]; ok{
		if userInfo.NickName != "" {
			return userInfo.NickName
		}
	}

	return persistentstore.GetPIDFromLocalStorage(uname)
}

func SetupDiscoveryServices() {

}

func SaveFavorRoom(uname string, roomname string) {
	if err, ok := CheckRoomName(uname); !ok {
		log.Fatal("invalid room name:", err)
	}
	persistentstore.SaveFavorRoom(uname, roomname)
}

func GetFavorRooms() []string{
	return persistentstore.GetFavorRooms()
}