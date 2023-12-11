package messaging

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/why-xn/alap-backend/pkg/core/log"
	"github.com/why-xn/alap-backend/pkg/dto"
	"sync"
)

var wsLock sync.RWMutex
var wsMap = map[string]*websocket.Conn{}
var userWsMap = sync.Map{}

func AddToWsMap(wsId string, ws *websocket.Conn) {
	wsLock.Lock()
	defer wsLock.Unlock()
	wsMap[wsId] = ws
}

func RemoveFromWsMap(wsId string) {
	wsLock.Lock()
	defer wsLock.Unlock()
	delete(wsMap, wsId)
}

func GetWsFromWsMap(wsId string) *websocket.Conn {
	wsLock.RLock()
	defer wsLock.RUnlock()
	if val, ok := wsMap[wsId]; ok {
		return val
	}
	return nil
}

func AddWsToUser(userId string, wsId string) {
	wsIdList, found := userWsMap.Load(userId)
	if found {
		wsIdList = append(wsIdList.([]string), wsId)
	} else {
		wsIdList = []string{}
		wsIdList = append(wsIdList.([]string), wsId)
	}
	userWsMap.Store(userId, wsIdList)
}

func RemoveWsFromUser(userId string, wsId string) {
	wsIdList, found := userWsMap.Load(userId)
	if found {
		newWsIdList := []string{}
		for i := 0; i < len(wsIdList.([]string)); i++ {
			if wsIdList.([]string)[i] == wsId {
				continue
			}
			newWsIdList = append(newWsIdList, wsIdList.([]string)[i])
		}
		userWsMap.Store(userId, newWsIdList)
	}
}

func GetWsIdListOfUser(userId string) []string {
	wsIdList, found := userWsMap.Load(userId)
	if found {
		return wsIdList.([]string)
	}
	return nil
}

func ProcessIncomingMsg(in dto.WsInMessageDTO) {
	//write ws data

	wsIdList := GetWsIdListOfUser(in.To)
	var ws *websocket.Conn
	var err error

	var out = dto.WsOutMessageDTO{
		Sender:     in.To,
		ChatWindow: in.ChatWindow,
		Msg:        in.Msg,
	}

	outByte, _ := json.Marshal(out)

	for _, wsId := range wsIdList {
		ws = GetWsFromWsMap(wsId)
		if ws != nil {
			err = ws.WriteMessage(in.MessageType, outByte)
			if err != nil {
				log.Logger.Info("write:", err)
				continue
			}
		}
	}
}
