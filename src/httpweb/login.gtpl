
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN"
  "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd"> 
<html> 
  <head>
    <meta http-equiv="Content-type" content="text/html; charset=utf-8">
    <meta content="authenticity_token" name="csrf-param" />
<meta content="ny64GX7B2welwuwq3yC6a0fgvMt9GRM6+FjQ7227zN0=" name="csrf-token" />
    <title>用户登录_天天笔记</title> <meta name="keywords" content="记事本,网络记事本,在线记事本" /> 
    <link href="/static/bootstrap.min.css" rel="stylesheet" media="screen">
    <link rel="stylesheet" href="/static/master.css" type="text/css" media="screen" charset="utf-8"/>
    <script src="/static/jquery.min.js"></script>
  </head> 
  <body> 
    <div id='wrapper'> 
      <div id='header'>
        <div class="container"> 
          <h1><a href="/">天天笔记</a></h1> 
          <h2>在线记事本</h2> 
            <a href="/login/" class="login">登录</a>
        </div>
      </div>
      
<div id="welcome">
	<a id="toggle_register_login" href="/register/">创建新账号</a>
	<h3>登录</h3>

	<form accept-charset="UTF-8" action="/login" class="new_user" id="new_user" method="post"><div style="margin:0;padding:0;display:inline"><input name="utf8" type="hidden" value="&#x2713;" /><input name="authenticity_token" type="hidden" value="ny64GX7B2welwuwq3yC6a0fgvMt9GRM6+FjQ7227zN0=" /></div>

		<table align="center" class="table" style="" rules=none>
		<tr>
			<th>用 户：</th>
			<td><input class="w150" id="user_name" name="user[name]" size="30" type="text" /></td>
		</tr>
		<tr>
			<th>密 码：</th>
			<td><input class="w150" id="user_pass" name="user[pass]" size="30" type="password" />
			</td>
		</tr>
		<tr>
			<td colspan="2" align="center" style="padding-left: 50px;"><br/>
				<input class="btn btn-primary btn-large" name="commit" type="submit" value="安全登录" /> &nbsp;&nbsp;
				<button class="btn btn-large" name="button" onclick="reset(); return false;" type="button">重新填写</button>
			<td>
		</tr>
		</table>
</form>
</div>
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
