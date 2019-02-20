
var socket = io.connect();
var connectNum=0;
var pproxy_colors=["#FFFFFF","#CCFFFF","#FFCCCC","#99CCCC","996699","#CC9999","#0099CC","#FFFF66","#336633","#99CC00"]

var ip_colors={} 
var pproxy_session_max_len=0;
var pproxy_table_max_row=2000;


var pproxy_req_list=[];
var flag = true;

if(window.localStorage){
	pproxy_req_list=$.parseJSON(window.localStorage["reqs"]||"[]")
	pproxy_session_max_len=1000;
}

function pproxy_show_reqs_from_local(){
	for(var i=0;i<pproxy_req_list.length;i++){
		pproxy_show_req(pproxy_req_list[i]);
	}
}

window.onbeforeunload=function(){
	if(pproxy_session_max_len>0){
		window.localStorage["reqs"]=JSON.stringify(pproxy_req_list);
	}
}

function pproxy_save_req_local(dataStr64){
	if(pproxy_session_max_len<1){
		return;
	}
	var n=pproxy_req_list.push(dataStr64)
	if(n>pproxy_session_max_len){
		pproxy_req_list.shift();
	}
}

function pporxy_session_stop(){
    flag = !flag;
    if (flag) {
        $("#btn_stop").text("stop");
    } else {
        $("#btn_stop").html("&nbsp;go&nbsp;");
    }

}

function pporxy_session_clean(){
    pproxy_req_list=[];
    $('#tb_network tbody').empty();
}

function pproxy_log(msg){
	$("#log_div").append("<div>"+(new Date().toString())+":"+msg+"</div>");
}


function pproxy_net_local_filter(host,path){
	if(!_pproxy_filter(host,$("#net_local_host").val())){
		return false;
	}
	if(!_pproxy_filter(path,$("#net_local_path").val())){
		return false;
	}
	return true;
}

function _pproxy_filter(kw,kwstr){
	if(kwstr==""){
		return true;
	}
	var kws=(kwstr+"").split("|");
	kw+="";
	var tmp=[]	
	for(var i in kws){
		var _kw=$.trim(kws[i]+"");
		if(_kw!=""){
			tmp.push(_kw);
		}
	}
	if(tmp.length==0){
		return true;
	}
	for (var i in tmp){
		if(kw.indexOf(tmp[i])!=-1){
			return true;
		}
	}
	return false;
}

socket.on('connect', function() {
	connectNum++
	if(connectNum>1){
		pproxy_log("ws error.connectNum="+connectNum);
		socket.emit("disconnect");
		return;
	}
    $("#connect_status").html("online")
    $("#network_filter_form").change();
});

socket.on("disconnect", function() {
	connectNum--;
    $("#connect_status").html("<font color=pink>offline</font>");
});

socket.on("hello",function(data){
	console && console.log(new Date().toString(),data);
});

function pproxy_getColor(addr){
	var info=addr.split(":");
	if(info.length!=2){
		return "#FFFFFF";
	}
	var ip=info[0];
	var color=ip_colors[ip];
	if(!color){
		var l=0;
		for(var _t in ip_colors){
			l++;
		}
		color=pproxy_colors[l% pproxy_colors.length];
		ip_colors[ip]=color;
	}
	return color
}


function pproxy_show_req(data) {
	console && console.log("pproxy_show_req:",data)
//    var dataStr=Base64.decode(dataStr64+"");
//	var data={};
//	try{
//		data=$.parseJSON(dataStr)
//	}catch(e){
//		console && console.log("parseJSON_error-pproxy_show_req",e,"datastr:",dataStr)
//		return
//	}
	
//    console && console.log("req", data)
    var html="<tr onclick=\"get_response(this,'" + data['docid'] + "')\" ";
	var cls=[]
    if(data["replay"]){
    	cls.push("replay");
    }
	if(!pproxy_net_local_filter(data["host"],data["path"])){
		cls.push("hide");
	}
	
	if(cls.length>0){
    	html+="class='"+cls.join(" ")+"' ";
	}
    html+=" data-ip='"+data["client_ip"]+"' data-host='"+h(data["host"])+"' data-path='"+h(data["path"])+"'>" 
    + "<td bgcolor='"+pproxy_getColor(data["client_ip"])+"'>" + data["sid"] + "</td>"
    + "<td><div class='oneline' title='"+h(data["host"])+"'>" + data["host"] + "</div></td>" +
    "<td><div class='oneline' title='"+h(data["url"])+"'>" +data["method"]+"&nbsp;"+ h(data["path"])+ "</div></td>" + 
    "</tr>";
    var tb=$("#tb_network tbody");
    tb.prepend(html);
    
    if(Math.random()*100>95 && tb.find("tr").size()>pproxy_table_max_row*1.5){
    	tb.find("tr:gt("+pproxy_table_max_row+")").each(function(){
    		$(this).remove();
    	});
    }
}


socket.on("req",function(data){
	console && console.log("on.req","data:",data);
	// pproxy_save_req_local(data);
    if(flag){
        pproxy_show_req(data);
    }
});


socket.on("user_num", function(data) {
	$("#user_num_online").html(data);
});

socket.on("res",
        function(data) {
			console && console.log("on.res","data:",data)
//	        var dataStr=Base64.decode(dataStr64+"");
//			try{
//				var data=$.parseJSON(dataStr)
//			}catch(e){
//				console && console.log("parseJSON_error on.res:",e)
//				return
//			}
//			console && console.log(data)
            var req = data["req"];
            var res = data["res"];
            var html="";
            if(req){
	            var re_do_str=req["schema"]=="http"?("&nbsp;<a target='_blank' href='/replay?id="+req["id"]+"'>replay</a>"):"";
	            
	            html += "<div><table class='tb_1'><caption>Request"+re_do_str+pproxy_timeformat(req["now"])+"</caption>";
	            html += "<tr><th width='80px'>url</th><td>" + h(req["url"]) + "&nbsp;&nbsp;<a href='"+h(req["url"])+"' target='_blank'>view</a></td></tr>"
	            if (req["url_origin"]!=req["url"]) {
	                html += "<tr><th>origin</th><td><span style='color:blue'>" + h(req["url_origin"]) + "</span></td></tr>";
	            }
	            if (req["msg"]) {
	            	html += "<tr><th>msg</th><td><span style='color:red'>" + h(req["msg"])+"</span></td></tr>";
	            }
	            html += "<tr><th>proxy_urer</th>" +
	            		"<td><b>remote_addr : </b>&nbsp;" +req["client_ip"] + "&nbsp;&nbsp;<b> docid : </b>&nbsp;"+ req["id"] + 
	            		"</td></tr>";
	            html += pproxy_tr_sub_table(req["form_get"], "get_params");
	            html += pproxy_tr_sub_table(req["form_post"], "post_params");
	            if (req["dump"]) {
	                html += "<tr><th>req_dump</th><td>" + h(Base64.decode(req["dump"])).replace(/\n/g, "<br/>")
	                        + "</td></tr>";
	            }
	            html += "</table></div>";
            }else{
            	html="<br/><br/><br/><br/><center>request not exists!</center>";
            }
            var res_link = "";
            
            var hideBigBody=false;
            
            if (res) {
                res_link = "<a href='/response?id=" + res["id"] + "' target='_blank'>view</a>";
                html += "<div><table class='tb_1'><caption>Response&nbsp;" + res_link +pproxy_timeformat(res["now"])+ "</caption>"
                if (res["msg"]) {
	            	html += "<tr><th  width='80px'>msg</th><td><span style='color:red'>" + h(res["msg"])+"</span></td></tr>";
	            }
                
                if (res["dump"]) {
                    html += "<tr><th width='80px'>res_dump</th><td>" + h(Base64.decode(res["dump"])).replace(/\n/g, "<br/>")
                    + "</td></tr>";
                }
                var body_str = Base64.decode(res["body"])
                var isImg = res["header"]["Content-Type"] != undefined
                        && res["header"]["Content-Type"][0].indexOf("image") > -1;

                var isStatusOk = res["status"] == 200;

                var bd_json = pproxy_parseAsjson(body_str);

                if (bd_json) {
                	hideBigBody=true;
                    html += "<tr><th width='80px'>body_json</th><td>" + bd_json + "</td></tr>";
                }
                if (isImg) {
                	hideBigBody=true;
                    html += "<tr><th>body_img</th><td><img src='data:" + res["header"]["Content-Type"][0] + ";base64,"
                            + res["body"] + "'/></td></tr>";
                }
                if (!isImg || res["body"].length < 1000 || !isStatusOk) {
                	html += "<tr><th width='80px'>body";
                	if(res["body"].length>400){
                		html+="<div><a href='#' onclick='return pproxy_res_td_body_toggle()'>toggle</a></div>";
                	}else{
                		hideBigBody=false;
                	}
                	html+= "</th>" +
                			"<td>" +
                			"<div id='res_td_body' "+(hideBigBody?"class='res_td_body' ":"")+">" + h(body_str).replace(/\n/g, "<br/>") + 
                			"</div></td></tr>";
                }
            }

            html += "</table></div>";
            $("#right_content").empty().html(html).hide().slideDown("fast");
        })

function pproxy_res_td_body_toggle(){
	$("#res_td_body").toggleClass("res_td_body");
	return false;
}

function pproxy_timeformat(sec){
	if(!sec ||sec<1000){
		return "";
	}
	var numFill=function(num,len){
		num=num+"";
		var l=len-num.length;
		for(;l>0;l--){
			num="0"+num;
		}
		return num;
	}
	var d=new Date()
	d.setTime(sec*1000)
	return "&nbsp;<font size=-1>"+d.getFullYear()+"-"+numFill(d.getMonth()+1,2)+"-"+numFill(d.getDate(),2)+" "
		   +numFill(d.getHours(),2)+":"+numFill(d.getMinutes(),2)+":"+numFill(d.getSeconds(),2)+"</font>";
}

        
function pproxy_parseAsjson(str) {
    try {
    	str=str+"";
    	if(str[0]!="{" && str[0]!="["){
    		return false;
    	}
        var jsonObj = JSON.parse(str);
        if (jsonObj) {
            var json_str = JSON.stringify(jsonObj, null, 4);
            return "<pre>" + json_str + "</pre>";
        }
    } catch (e) {
    	console.log("pproxy_parseAsjson_error",e);
    }
    return false;
}

function pproxy_tr_sub_table(obj, name) {
    if (!obj) {
        return "";
    }
    var html = "<tr><th>" + name + "</th><td class='td_has_sub'><table class='tb_1'>";
    var i = 0;
    var max_key_len=0;
    for ( var k in obj) {
    	max_key_len=Math.max(max_key_len,(k+"").length);
    }
    for ( var k in obj) {
        html += "<tr><th  "+(max_key_len<40?"width='120px' nowrap":"width='140px'")+">" + k + "</th><td><ul class='td_ul'>";
        for ( var i in obj[k]) {
            html += "<li>";
            var json_str = pproxy_parseAsjson(obj[k][i]);
            if (json_str) {
                html += json_str;
            } else {
                html += h(obj[k][i]);
            }
            html += "</li>";
        }
        html += "</ul></td></tr>";
        i++;
    }
    if (i < 1) {
        return "";
    }
    html += "</table></td></tr>"
    return html;
}

function pproxy_show_response(docid){
	console && console.log("get_response start,docid=", docid);
    var loading_msg="loading...docid=" + docid;
    var isValidId=(docid+"").length>2;
    if(!isValidId){
    	loading_msg="https request:no data";
    }else{
    	loading_msg+="&nbsp;<a href='javascript:;' onclick=\"pproxy_show_response('"+docid+"')\">reload</a>";
    }
    $("#right_content").empty().html("<center style='margin:200px 0 auto'>"+loading_msg+"</center>");
    if(!isValidId){
    	return;
    }
	socket.emit("get_response", docid);
	console.log && console.log("emit get_response,docid=",docid);
}

function get_response(tr, docid) {
    pproxy_show_response(docid);
    $(tr).parent("tbody").find("tr").removeClass("selected");
    $(tr).addClass("selected");
    location.hash="req_"+docid;
}

function bytesToString(bytes) {
    var result = "";
    for (var i = 0; i < bytes.length; i++) {
        result += String.fromCharCode(parseInt(bytes[i], 2));
    }
    return result;
}

function h(html) {
	if(html==""){
		return "&nbsp;";
	}
	html = (html+"").replace(/&/g, '&amp;')
				.replace(/</g, '&lt;')
				.replace(/>/g, '&gt;')
			    .replace(/'/g, '&acute;')
			    .replace(/"/g, '&quot;')
	            .replace(/\|/g, '&brvbar;');
    return html;
}

$().ready(function() {
	$("#network_filter_form input:text").each(function(){
		pproxy_local_save(this,$(this).attr("name"));
	});
	var filter_form=$("#network_filter_form");
	filter_form.change(function() {
        var form_data = $(this).serialize();
        socket.emit("client_filter", form_data);
    });
    
    setTimeout(function(){filter_form.change();},600);
    setTimeout(function(){filter_form.change();},3000);
    
	filter_form.find("input:text").keyup(function(){
		filter_form.change();
    });
    
    if(location.hash.match(/req_\d+/)){
        var docid=location.hash.substr(5);
        setTimeout((function(id){
        	return function(){
        		pproxy_show_response(id);
        	}
        })(docid),500);
    }
    
    setTimeout(pproxy_show_reqs_from_local,0);
    
    
    $("#net_local_host,#net_local_path").bind("keyup change",function(){
    	$('#tb_network tbody tr').each(function(){
    		if(pproxy_net_local_filter($(this).data("host"),$(this).data("path"))){
    			$(this).removeClass("hide");
    		}else{
    			$(this).addClass("hide");
    		}
    	});
    });
});

function pproxy_local_save(target,id){
	if(!window.localStorage){
		return;
	}
	$(target).val(window.localStorage[id]||"").change(function(){
		window.localStorage[id]=$(this).val();
	});
}