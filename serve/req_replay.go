package serve

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	REPLAY_FLAG       = "Proxy-pproxy_replay"
	REPLAY_REMOTEADDR = "Proxy-pproxy_remoteaddr"
	REPLAY_USER_NAME  = "Proxy-pproxy_user"
)

func (ser *ProxyServe) req_replay(w http.ResponseWriter, req *http.Request, values map[string]interface{}) {
	if req.Method == "POST" {
		ser.req_replayPost(w, req, values)
		return
	}
	docid_str := strings.TrimSpace(req.FormValue("id"))
	if docid_str == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("empty id param"))
		return
	}
	docid, err_int := parseDocId(docid_str)
	if err_int != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("param id[%s] error:\n%s", docid_str, err_int)))
		return
	}
	req_doc := ser.GetRequestByDocid(docid)
	if req_doc == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("request doc not found!"))
		return
	}
	_url := fmt.Sprintf("%s", req_doc["url"])
	u, err := url.Parse(_url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("parse url[%s] error\n%s", _url, err)))
		return
	}
	u.RawQuery = ""
	values["req"] = req_doc
	values["action_url"] = u.String()
	values["subTitle"] = "replay|" + u.String() + "|"
	html := render_html("replay.html", values, true)
	w.Write([]byte(html))
}

var replay_skip_headers = map[string]int{"Content-Length": 1}

func (ser *ProxyServe) req_replayPost(w http.ResponseWriter, req *http.Request, values map[string]interface{}) {
	replay := req.FormValue("replay")
	basic := make(map[string]string)
	basic["action_url"] = strings.TrimSpace(req.FormValue("basic_action_url"))
	method := strings.TrimSpace(strings.ToUpper(req.FormValue("basic_method")))
	basic["method"] = method

	host := strings.TrimSpace(req.FormValue("basic_host"))

	basic_remoteAddr := req.FormValue("basic_RemoteAddr")
	basic_user := req.FormValue("basic_user")

	header := GetFormValuesWithPrefix(req.Form, "header_")
	get := GetFormValuesWithPrefix(req.Form, "get_")
	post := GetFormValuesWithPrefix(req.Form, "post_")

	formData := make(map[string]interface{})
	formData["basic"] = basic

	formData["header"] = header
	formData["get"] = get
	formData["post"] = post

	values["form"] = formData
	if replay == "direct" {
		html := render_html("replay_direct.html", values, true)
		w.Write([]byte(html))
		return
	} else {
		req_bd := ""

		_url := basic["action_url"]

		if len(get) > 0 {
			form_values := make(url.Values)
			for k, v := range get {
				for _, _v := range v {
					form_values.Add(k, _v)
				}
			}
			if strings.Contains(_url, "?") {
				_url += "&"
			} else {
				_url += "?"
			}
			_url += form_values.Encode()
		}

		if len(post) > 0 {
			form_values := make(url.Values)
			for k, v := range post {
				for _, _v := range v {
					form_values.Add(k, _v)
				}
			}
			req_bd = form_values.Encode()
		}

		replay_req, err := http.NewRequest(method, _url, strings.NewReader(req_bd))
		if err != nil {
			w.Write([]byte("build request failed\n" + err.Error()))
			return
		}
		if host != "" {
			replay_req.Host = host
		}
		replay_req.Header.Set(REPLAY_FLAG, "replay")

		replay_req.Header.Set(REPLAY_REMOTEADDR, basic_remoteAddr)
		replay_req.Header.Set(REPLAY_USER_NAME, basic_user)

		for k, v := range header {
			if _, has := replay_skip_headers[k]; has {
				continue
			}
			replay_req.Header.Set(k, strings.Join(v, ";"))
		}
		ser.httpProxy.ServeHTTP(w, replay_req)
	}
}