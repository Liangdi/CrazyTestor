
function ChatBox(conf){
	var id = conf.el;
	var leftTplId=conf.left;
	var rightTplId=conf.right;
	var richTplId = conf.rich;
	var richsTplId = conf.richs;
	var tplFirst = conf.first;
	var tplSecond = conf.second;
	this.user = conf.user;
	this.robot =conf.robot;
	this.leftTpl = $(leftTplId).html();
	this.rightTpl = $(rightTplId).html();
	this.richTpl = $(richTplId).html();
	this.richsTpl = $(richsTplId).html();
	this.firstTpl = $(tplFirst).html();
	this.secondTpl = $(tplSecond).html();
	this.el = $(id);
	
}
ChatBox.prototype.addUserChat = function(msg){
	var info = this.user;
    console.log(this.leftTpl)
	var tpl = this.leftTpl.replace("{USER}",info.user)
		.replace("{AVATAR}",info.avatar)
		.replace("{TIME}",new Date())
		.replace("{CONTENT}",msg);
    console.log(tpl)
	$(tpl).fadeIn(200).appendTo(this.el);
	console.log(this.el);
	this.scroll();
};
ChatBox.prototype.addRobotChat = function(msg){
	var info = this.robot;
	var serverId = $("#wxServers").val();
	serverId =$("#wxServers").find("option:selected").text()
	var tpl = this.rightTpl.replace("{USER}",serverId)
		.replace("{AVATAR}",info.avatar)
		.replace("{TIME}",new Date())
		.replace("{CONTENT}",msg);
	$(tpl).fadeIn(200).appendTo(this.el);
	this.scroll();
};
ChatBox.prototype.addRichContent = function(title,desc,pic,url){
	console.log("addRichContent");
	var info = this.robot;
	var serverId = $("#wxServers").val();
	serverId =$("#wxServers").find("option:selected").text();
	var tpl = this.richTpl.replace("{USER}",serverId)
		.replace("{AVATAR}",info.avatar)
		.replace("{TIME}",new Date())
		.replace("{DESC}",desc)
		.replace(new RegExp("{URL}","gm"),url)
		.replace("{PIC}",pic)
		.replace("{TITLE}",title);

	$(tpl).fadeIn(200).appendTo(this.el);
	this.scroll();
};
ChatBox.prototype.addRichContents = function(first,others){
	console.log(others);
	var info = this.robot;
	var serverId = $("#wxServers").val();
	serverId =$("#wxServers").find("option:selected").text();
	var firstTplHTML = this.firstTpl.replace("{TITLE}",first.title)
			.replace(new RegExp("{URL}","gm"),first.url)
			.replace("{PIC}",first.pic);
	console.log(firstTplHTML);
	var secondTplHTML = "";
	for (var i = 0; i < others.length; i++) {
		var item = others[i];
		secondTplHTML+=this.secondTpl.replace("{TITLE}",item.title)
			.replace(new RegExp("{URL}","gm"),item.url)
			.replace("{PIC}",item.pic);
	}
	var tpl = this.richsTpl.replace("{USER}",serverId)
		.replace("{AVATAR}",info.avatar)
		.replace("{FIRST}",firstTplHTML)
		.replace("{SECOND}",secondTplHTML)
		.replace("{TIME}",new Date());
	$(tpl).fadeIn(200).appendTo(this.el);
	this.scroll();
};
ChatBox.prototype.scroll = function(){
	this.el.animate({scrollTop:this.el[0].scrollHeight}, "slow");
};
ChatBox.prototype.clean = function(){
	this.el.children().fadeOut( function() { $(this).remove(); });
};
var user={
	avatar:"static/images/customer.png",
	user:"debug"
};

var robot = {
	avatar:"static/images/customer-2.png",
	user:"OHTER"
};
var url = "Chat?action=get";
function getRobotRespone(){
	//console.log("ajax poll...")
	$.ajax({
		url:url,
		success:function(data){
			
			if(data.success){
				if(data.messages && data.messages.length){
					var msgs = data.messages;
					for (var i = 0; i < msgs.length; i++) {
						chatBox && chatBox.addRobotChat(msgs[i]);
					}
				}
			}
		},
		complete:function(){
			setTimeout(getRobotRespone,1000);
		}
	});
}
var sendUrl = "//localhost:8081?signature=1223&timestamp=12&nonce=1";
function sendToRobot(msg){
	var tpl = $("#message_tpl").val();
	console.log(tpl);
	var serverId = $("#wxServers").val();
	var time = new Date().getTime();
	var message =tpl.replace("{toUser}",serverId)
		.replace("{createTime}",time)
		.replace("{content}",msg);
	console.log(message);
	$.ajax({
		type: "POST",
				url: sendUrl,
				//async: false,
				contentType:"application/xml; charset=utf-8",
				data: message,
		success:function(xml){
			console.log("xml",xml);
			if(xml.success == false){
				chatBox.addRobotChat(xml.message);
				return;
			}
			var typeNode = xml.getElementsByTagName("MsgType")[0];
			var type = typeNode.textContent;
			var content = "";
			switch(type){
				case "news":
					content = "[图文消息]";
					var count =  xml.getElementsByTagName("ArticleCount")[0].textContent;
					console.log("count:",count);
					var items = xml.getElementsByTagName("item");
					console.log("items",items);
					if(count == 1){
						var item = items[0];
						var title = item.getElementsByTagName("Title")[0].textContent;
						var desc = item.getElementsByTagName("Description")[0].textContent.replace(new RegExp("\n","gm"),"<br/>");
						var pic = item.getElementsByTagName("PicUrl")[0].textContent;
						var url = item.getElementsByTagName("Url")[0].textContent;
						console.log("title",title);
						console.log("desc",desc);
						console.log("pic",pic);
						console.log("url",url);
						chatBox.addRichContent(title,desc,pic,url);
						return;
					} else {
						var item = items[0];
						var title = item.getElementsByTagName("Title")[0].textContent;
						var pic = item.getElementsByTagName("PicUrl")[0].textContent;
						var url = item.getElementsByTagName("Url")[0].textContent;
						var first = {
							title:title,
							pic:pic,
							url:url
						};
						var others = [];
						for (var i = 1; i < items.length; i++) {
							var item = items[i];
							var title = item.getElementsByTagName("Title")[0].textContent;
							var pic = item.getElementsByTagName("PicUrl")[0].textContent;
							var url = item.getElementsByTagName("Url")[0].textContent;
							var obj = {
								title:title,
								pic:pic,
								url:url
							};
							others.push(obj);
						}
						chatBox.addRichContents(first,others);
						return;
					}
					break;
				case "text":
					var contentNode = xml.getElementsByTagName("Content")[0];
					content = contentNode.textContent;
					break;
				default :
					content = "消息类型:" + type + " 不支持显示"; 
			}
			chatBox.addRobotChat(content);
			//console.log(type)
		},
		complete:function(){
			setTimeout(function(){
				$("#message").prop('disabled', false);
			},500);
		}
	});
}
$(function(){
	 var chatBox =window.chatBox= new ChatBox({
		user:user,
		robot:robot,
		el:"#chat-box",
		left:"#leftContent",
		right:"#rightContent",
		rich:"#richContent",
		richs:"#richContents",
		first:"#tplFirst",
		second:"#tplSecond"
	 });
	//getRobotRespone();
	/*$.ajax({
		url:"wxserver?action=list",
		success:function(data){
			if(data.success){
				//console.log(data);
				var list = data.list;
				for (var i = 0; i < list.length; i++) {
					var s = list[i];
					//console.log(s); 
					var title = s.title;
					var serverId = s.serverId;
					$('<option value="'+serverId+'">'+title+'</option>').appendTo("#wxServers");
				}
				var serverId = $.cookie('serverId');
				if(serverId != ""){
					$("#wxServers").val(serverId);
				}
			} else {
				msg("错误:" + data.message);
			}
		}
	}) */
	 function send(){
		var content = $("#message").val();
         console.log(content);
		if(content == ""){
			 $("#message").focus();
			return false;
		}
		chatBox.addUserChat(content);
		$("#message").prop('disabled', true);
		sendToRobot(content);
		 $("#message").val("");
	 }
	 $("#message").on("keypress",function(e){
		 if(e.which == 13) {
		   send();
	  }
		
	 });
	 $("#btn-send").on("click",function(){
		send();
	 });
	 $("#btn-send-robot").on("click",function(){
		 chatBox.clean();
		//chatBox.addRobotChat(content);
	 });
	 $("#wxServers").change(function(){
		var serverId = $(this).val();
		//$.cookie('serverId', serverId);
		//console.log($.cookie('serverId'));
	 });
});