package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/FJ1207/cache/concurrent/mutexincludobsolesence"
)

const defaultBasePath = "/fujia/"

type HttpPool struct{//httppool结构
	self string //ip 端口
	basePath string //基础路径
}

func NewHTTPPool(self string) *HttpPool {
	return &HttpPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HttpPool) Log(fomat string,v ...interface{}) {
	//用服务器的名字log
	log.Printf("server %s,%s",p.self,fmt.Sprintf(fomat ,v ...))
}

func (p *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request ) {//p响应所有的请求
	//如果url和basepath不同，则panic
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s,%s",r.Method,r.URL.Path)
	// "wdy"[1:] = "dy"
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]
	
	group := mutexincludobsolesence.GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}