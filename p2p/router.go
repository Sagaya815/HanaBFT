package p2p

import (
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	"hanaBFT/hlog"
	"hanaBFT/utils"
	"os"
	"path/filepath"
	"sync"
)

type Route struct {
	IP   string `json:"ip,omitempty"`
	Port string `json:"port,omitempty"`
}

type MapRouter struct {
	peerRoutes map[string]*Route `json:"peer_routes,omitempty"`
	peerNum    int               `json:"peer_num,omitempty"`
	lock       sync.RWMutex      `json:"lock"`
}

var MR *MapRouter

var file *os.File

func Init() {
	var err error
	MR = &MapRouter{
		peerRoutes: make(map[string]*Route),
		peerNum:    0,
		lock:       sync.RWMutex{},
	}
	file, err = os.OpenFile(filepath.Base(fmt.Sprintf("%s.%s.route", os.Args[0], utils.ShortPeerID(hlog.ID))), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 777)
	if err != nil {
		hlog.Fatalf("When we create router file, an error occurred: %s", err)
	}
}

func (mr *MapRouter) Add(ID, IP, Port string) {
	if mr.IsExists(ID) {
		return
	}
	mr.lock.Lock()
	route := &Route{
		IP:   IP,
		Port: Port,
	}
	mr.peerRoutes[ID] = route
	mr.peerNum++
	mr.lock.Unlock()
}

func (mr *MapRouter) Remove(ID string) {
	mr.lock.Lock()
	delete(mr.peerRoutes, ID)
	mr.peerNum--
	mr.lock.Unlock()
}

func (mr *MapRouter) Update(ids IDs) {
	mr.lock.Lock()
	for id, _ := range mr.peerRoutes {
		if !ids.checkIsExists(id) {
			delete(mr.peerRoutes, id)
			mr.peerNum--
		}
	}
	mr.lock.Unlock()
}

func (mr *MapRouter) IsExists(id string) bool {
	mr.lock.RLock()
	defer mr.lock.RUnlock()
	_, ok := mr.peerRoutes[id]
	if ok {
		return true
	} else {
		return false
	}
}

func (mr *MapRouter) upload() {
	mr.lock.RLock()
	defer mr.lock.RUnlock()
	routerData, err := json.Marshal(mr.peerRoutes)
	if err != nil {
		hlog.Errorf("When marshal MapRouter, an error occurred: %s", err)
		return
	}
	file.Truncate(0)
	file.Seek(0, 0)
	_, err = file.Write(routerData)
	if err != nil {
		hlog.Errorf("When upload MapRouter, an error occurred: %s", err)
		return
	}
}

type IDs []peer.ID

func (ids IDs) checkIsExists(id string) bool {
	var flag = false
	for _, i := range ids {
		if i.Pretty() == id {
			flag = true
			break
		}
	}
	return flag
}
