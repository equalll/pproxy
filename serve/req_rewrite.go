package serve

import (
    "fmt"
    "log"
    "net/http"
    "net/url"
    "strings"
)

var rewriteJsTpl = "function pproxy_rewrite(req){\n%s\nreturn req;\n}"

func (ser *ProxyServe) parseAndSaveRewriteJs(jsStr string) error {
    rewriteJs := fmt.Sprintf(rewriteJsTpl, jsStr)
    js.Run(rewriteJs)
    jsFn, err := js.Get("pproxy_rewrite")
    if err == nil {
        ser.RewriteJs = jsStr
        ser.RewriteJsFn = jsFn
    }
    return err
}

func (ser *ProxyServe) reqRewrite(req *http.Request) {
    if ser.RewriteJs == "" {
        return
    }
    urlObj, _ := js.Object(`req={}`)
    urlObj.Set("url", req.URL.String())
    urlObj.Set("schema", req.URL.Scheme)
    
    host_info:=strings.Split(req.URL.Host,":")
    urlObj.Set("host", host_info[0])
    if(len(host_info)==2){
	    urlObj.Set("port", host_info[1])
    }else{
	    urlObj.Set("port", "")
    }
    
    urlObj.Set("path", req.URL.Path)
    urlObj.Set("rawquery", req.URL.RawQuery)
    username := ""
    psw := ""
    if req.URL.User != nil {
        username = req.URL.User.Username()
        psw, _ = req.URL.User.Password()
    }
    urlObj.Set("username", username)
    urlObj.Set("password", psw)

    js_ret, err_js := ser.RewriteJsFn.Call(ser.RewriteJsFn, urlObj)

    if err_js == nil {
        if js_ret.IsObject() {
            obj, export_err := js_ret.Export()
            if export_err == nil {
                url_obj := obj.(map[string]interface{})
                schema:=getMapValStr(url_obj,"schema")
                url_new := schema+ "://"
                username :=getMapValStr(url_obj,"username")
                if username != "" {
                    url_new += fmt.Sprintf("%s:%s@", username, getMapValStr(url_obj,"password"))
                }
                host:=getMapValStr(url_obj,"host")
                port:=getMapValStr(url_obj,"port")
                if(port!=""){
                   host+=":"+port
                }
                url_new += fmt.Sprintf("%s%s", host, getMapValStr(url_obj,"path"))
               
                rawquery:=getMapValStr(url_obj,"rawquery")
                if( rawquery!=""){
                   url_new+="?"+rawquery
                }
                if url_new == req.URL.String() {
                    return
                }
                host_addr:=getMapValStr(url_obj,"host_addr")

                var url_err error
                req.URL, url_err = url.Parse(url_new)
                if ser.Debug {
                    log.Println("DEBUG req_rewrite,url_new:", url_new, "req_new:", req.URL)
                }
                if url_err != nil {
                    log.Println("js filter err:", js_ret, url_err)
                } else {
                    req.Host = req.URL.Host
                    if(host_addr!=""){
                       req.URL.Host=host_addr
                    }
                }
            } else {
                log.Println("js filter result wrong", js_ret.String())
            }
        }
    } else {
        log.Println("js filter err:", err_js, js_ret)
    }
}
