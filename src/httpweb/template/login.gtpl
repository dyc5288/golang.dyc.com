
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN"
  "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd"> 
<html> 
  <head>
    <meta http-equiv="Content-type" content="text/html; charset=utf-8">
    <meta content="authenticity_token" name="csrf-param" />
<meta content="ny64GX7B2welwuwq3yC6a0fgvMt9GRM6+FjQ7227zN0=" name="csrf-token" />
    <title>用户登录_115笔记</title> <meta name="keywords" content="记事本,网络记事本,在线记事本" /> 
    <link href="/static/bootstrap.min.css" rel="stylesheet" media="screen">
    <link rel="stylesheet" href="/static/master.css" type="text/css" media="screen" charset="utf-8"/>
    <script src="/static/jquery.min.js"></script>
  </head> 
  <body> 
    <div id='wrapper'> 
      <div id='header'>
        <div class="container"> 
          <h1><a href="/">115笔记</a></h1> 
          <h2>在线记事本</h2> 
            <a href="/login" class="login">登录</a>
        </div>
      </div>
      
<div id="welcome">
	<a id="toggle_register_login" href="/register">创建新账号</a>
	<h3>登录</h3>

	<form accept-charset="UTF-8" action="/login" class="new_user" id="new_user" method="post"><div style="margin:0;padding:0;display:inline"></div>

		<span class="label label-warning" id="js_err" style="display:none;"></span>
		<table align="center" class="table" style="" rules=none>
		<tr>
			<th>用 户：</th>
			<td><input class="w150" id="user_name" name="username" size="30" type="text" value="dyc5288" /></td>
		</tr>
		<tr>
			<th>密 码：</th>
			<td><input class="w150" id="user_pass" name="password" size="30" type="password" value="d54321" />
			</td>
		</tr>
		<tr>
			<td colspan="2" align="center" style="padding-left: 50px;"><br/>
				<input type="hidden" name="token" value="{{.}}">
				<input class="btn btn-primary btn-large" name="commit" type="button" id="js_sub" value="安全登录" /> &nbsp;&nbsp;
				<button class="btn btn-large" name="button" onclick="reset(); return false;" type="button">重新填写</button>
			<td>
		</tr>
		</table>
</form>
<script type="text/javascript">
	$(document).ready(function(){
		var err = $("#js_err");
	  	$("#js_sub").on("click", function(){
			var user_name = $("#user_name").val();
			var user_pass = $("#user_pass").val();
			if (!user_name) {
				err.html("用户名不能为空");
				err.css("display", "block");
				return false;
			}
			if (!user_pass) {
				err.html("密码不能为空");
				err.css("display", "block");
				return false;
			}
			$.ajax({
                url:"/login",
				data:{username:user_name, password:user_pass},
                dataType:"json",
				type:"post",
                success:function(json){
                    if (json.state) {
						location.href="/";
					} else {
						err.html(json.message);
						err.show();
					}
                }
            });
			return false;
		});
		$("#user_name").on("focus",function(){
			err.css("display", "none");
			return false;
		});
		$("#user_pass").on("focus",function(){
			err.css("display", "none");
			return false;
		});
	});
</script>
</div>
      <div id='footer'> 
        <div class="container">
        © 2012-2013 – <a href="/">115笔记</a> - TTBIJI.com - Layout based on bootstrap.
        <br/>
        <ul>
          <li class="bug"><a href="/">用户反馈</a></li>
          <li><a href="/">关于本站</a></li>
          <li><a href="/">常见问题</a></li>
          <li><a href="/">隐私政策</a></li>
          <li><a href="/">站点地图</a></li>
          <li class="last"><a href="http://www.miitbeian.gov.cn/" target="_blank">豫ICP备13000858号-2</a></li>
        </div>
      </div>
    </div>
  </body> 
</html>
