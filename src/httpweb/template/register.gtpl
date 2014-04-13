
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN"
  "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd"> 
<html> 
  <head>
    <meta http-equiv="Content-type" content="text/html; charset=utf-8">
    <meta content="authenticity_token" name="csrf-param" />
<meta content="j3j1p9zYLanR6VCqtcwCx1sJBSUOQBnMUWHrsEFg2F4=" name="csrf-token" />
    <title>用户注册_天天笔记</title>  
    <link href="/static/bootstrap.min.css" rel="stylesheet" media="screen">
    <link rel="stylesheet" href="/static/master.css" type="text/css" media="screen" charset="utf-8"/>
    <script src="/static/jquery.min.js"></script>
  </head> 
  <body> 
    <div id='wrapper'> 
      <div id='header'>
        <div class="container"> 
          <h1><a href="/">天天笔记</a></h1> 
          <h2>一款最好用的免费网络在线记事本</h2> 
            <a href="/login/" class="login">登录</a>
        </div>
      </div>
      <div id="page">
	<div class="mod-note-title">
		<h3>新用户注册</h3>
	</div>
	<div class="content">

		<form accept-charset="UTF-8" action="/register" class="new_user" id="new_user" method="post"><div style="margin:0;padding:0;display:inline"><input name="utf8" type="hidden" value="&#x2713;" /><input name="authenticity_token" type="hidden" value="j3j1p9zYLanR6VCqtcwCx1sJBSUOQBnMUWHrsEFg2F4=" /></div>
		<table class="table">
		<tr>
			<td>用户名：</td>
			<td><input autofocus="autofocus" class="w120" id="user_name" name="user[name]" required="required" size="30" type="text" /></td>
			<td class="tip">该名称用于登陆，不能少于4个字符</td>
		</tr>
		<tr>
			<td>登录密码：</td>
			<td><input class="w120" id="user_pass" name="user[pass]" required="required" size="30" type="password" /></td>
			<td class="tip">不能少于6个字符，请小心输入并牢记</td>
		</tr>
		<tr>
			<td>电子邮箱</td>
			<td><input class="w200" id="user_email" name="user[email]" required="required" size="30" type="email" /></td>
			<td class="tip">用于接收注册信息，以及找回密码</td>
		</tr>
		<tr>
			<td>手机号码：</td>
			<td><input class="w120" id="user_mobile" maxlength="11" name="user[mobile]" pattern="[0-9]{11}" size="11" type="text" /></td>
			<td class="tip">如果您希望通过手机找回密码，请正确输入</td>
		</tr>
		<tr>
			<td>验证码：</td>
			<td><style type="text/CSS">
#captcha{
    width: 60px !important;
    font-size: 16px;
    background-color: #fff;
  }
</style>

<img alt="captcha" src="/captcha/?code={{.}}&amp;time=1396452114" />
<input autocomplete="off" id="captcha" name="captcha" placeholder="" required="required" type="text" /><input id="captcha_key" name="captcha_key" type="hidden" value="{{.}}" />
</td>
			<td class="tip">为了确保您不是一个机器人，请正确输入</td>
		</tr>
		<tr>
			<td colspan="3" style="text-align:center;">
				<div id="regtip" style="display:none; text-align:left;">
					<i class="icon-remove"></i>
					<span class="label label-important"></span>
				</div>
				<input class="btn btn-primary btn-large" id="btnReg" name="commit" type="button" value="注册用户" /> &nbsp;
				<button class="btn btn-large" name="button" onclick="__reset(); return false;" type="submit">重新填写</button>
			</td>
		</tr>
		</table>
</form>
	</div>
</div>


<script type="text/javascript">
// 显示错误提示语
function showErr (msg) {
	$('#regtip span.label').html(msg);
	$('#regtip').show();
}

// 重置表单
function __reset(){
	$('form').get(0).reset();
	$('#regtip').hide();
}

$().ready(function(){
	// 错误提示语
	
	function get_data() {
		var res = {};
		$("input").each(function(){
			var name = $(this).attr("name");
			var value = $(this).val();
			res[name] = value;
		});
		return res;
	}

	// 绑定注册按钮事件
	$('#btnReg').click(function(){
		var result = true;
		var arrID = [
			['user_name', '用户名'],
			['user_pass', '登录密码'],
			['user_email', '电子邮箱'],
			['captcha', '验证码']
		];
		$(arrID).each(function(){
			var id = $(this)[0];
			var name = $(this)[1];
			var obj = $('#' + id);
			if(obj.val().length == 0) {
				showErr(name + '不能为空');
				result = false;			
				return false;
			}
		});
		
		var data = get_data();
		console.log(data);
		
		if (result) {
			$.ajax({
               	url:"/register/",
				data:data,
	               	dataType:"json",
					type:"post",
	               	success:function(json){
	                   if (json.state) {
						location.href="/login/?type=reg_success";
					} else {
						err.html(json.message);
						err.show();
					}
               	}
           	});
		}
		
		return false;
	});
});
</script>

      <div id='footer'> 
        <div class="container">
        © 2012-2013 – <a href="/">天天笔记</a> - TTBIJI.com - Layout based on bootstrap.
        <br/>
        <ul>
          <li class="bug"><a href="/feedback/">用户反馈</a></li>
          <li><a href="/about/">关于本站</a></li>
          <li><a href="/faq/">常见问题</a></li>
          <li><a href="/privacy/">隐私政策</a></li>
          <li><a href="/sitemap.xml">站点地图</a></li>
          <li class="last"><a href="http://www.miitbeian.gov.cn/" target="_blank">豫ICP备13000858号-2</a></li>
        </div>
      </div>
    </div>
  </body> 
</html>
