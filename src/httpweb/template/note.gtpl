
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN"
  "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd"> 
<html> 
  <head>
    <meta http-equiv="Content-type" content="text/html; charset=utf-8">
    <meta content="authenticity_token" name="csrf-param" />
<meta content="j3j1p9zYLanR6VCqtcwCx1sJBSUOQBnMUWHrsEFg2F4=" name="csrf-token" />
    <title>文档库_天天笔记</title>  
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
          <div class="user">
    &nbsp;
                当前用户：dyc5288 &nbsp;
                <a href="/config">个人设置</a>
                &nbsp;</div>
            <div style="float:right; margin-top:-35px;"><a href="/logout/" class="login">退出</a></div>
        </div>
      </div>
      <link href="/static/pagination.css" media="all" rel="stylesheet" type="text/css" />
<style type="text/css">
#page table td a { color: #333; }
#page table td span { font-size: 12px; color: purple; }
#page table td span a { font-size: 12px; color:purple; }
.gray { background-color: #f7f7f7; } 


.grumble{position:absolute;background-image:url(/images/bubble-sprite.png);background-repeat:no-repeat;z-index:99999}.grumble-text{position:absolute;text-align:center;z-index:99999;display:table;overflow:hidden;text-transform:uppercase;font-size:14px;line-height:14px}.ie7 .grumble-text,.ie6 .grumble-text{display:block}.grumble-text .outer{display:table-cell;vertical-align:middle;color:white}.ie7 .grumble-text .outer,.ie6 .grumble-text .outer{display:block;width:85%;position:absolute;top:48%;left:0}.ie7 .inner,.ie6 .inner{position:relative;top:-50%}.grumble-text50 .outer{padding:6px}.grumble-text100 .outer{padding:8px}.grumble-text150 .outer{padding:10px}.grumble-text200 .outer{padding:12px}.grumble50{background-position:0 0}.grumble100{background-position:0 -50px}.grumble150{background-position:0 -150px}.grumble200{background-position:0 -300px}.alt-grumble50{background-position:-200px 0}.alt-grumble100{background-position:-200px -50px}.alt-grumble150{background-position:-200px -150px}.alt-grumble200{background-position:-200px -300px}.grumble-button{position:absolute;width:20px;height:12px;-moz-border-radius:3px;border-radius:3px;background:#ff5c00;color:#fff;text-align:center;font-size:.8em;line-height:.9em;-moz-box-shadow:1px 1px 1px #ccc;-webkit-box-shadow:1px 1px 1px #ccc;box-shadow:1px 1px 1px #ccc;cursor:pointer;z-index:99999}
div.inner { color: #fff; }

tr.head { background-color: #fff; border-bottom: 1px solid #ccc; }
td.w1 { width: 45px }
td.w2 { width: 120px }
</style>
<script src="/static/bootstrap.min.js"></script>
<script src="/static/jquery.cookie.js"></script>
<script src="/static/jquery.grumble.min.js"></script>
<div id="page">
	<div class="mod-note-button">
		<input type="button" id="btnAllDoc" class="btn btn-mini btn-info" value="≡ 所有文档（a）" onclick="location.href='/doc/';" /> &nbsp;

		<div class="btn-group">
            <button id="btnCatelist" class="btn btn-mini dropdown-toggle" data-toggle="dropdown">≈ 类别切换（c） <span class="caret"></span></button>
            <ul class="dropdown-menu">
			  <li><a href="/doc/list/15314/">默认类别</a></li>
			  <li><a href="/doc/list/15315/">学习资料</a></li>
			  <li><a href="/doc/list/15316/">工作文档</a></li>
			  <li><a href="/doc/list/15317/">生活资讯</a></li>
			  <li class="divider"></li>
			  <li><a href="/doc/cate">类别配置</a></li>
            </ul>
          </div>

		<button class="btn btn-mini btn-primary" type="button" onclick="__new();">♀ 新建文档（n）</button>
		<input type="button" class="btn btn-mini btn-info" value="便笺" onclick="location.href='/memo/';" />
		<input type="button" id="btnFeedback" class="btn btn-mini btn-info" value="反馈" onclick="location.href='/feedback/';" />
		<input type="button" class="btn btn-mini btn-info hide" value="捐助" onclick="location.href='/donate/';" />
	</div>

	<div class="mod-note-title">
		<h1>文档列表</h1>
	</div>
	<div class="content">

	<table class="table table-condensed">
		<tr class="head">
			<td class="w1">公开</td>
			<td class="w2">类别</td>
			<td>标题</td>
			<td style="width:100px;">字数</td>
			<td style="width:120px;">更新日期</td>
		</tr>
		<tr><td colspan="4"><br/><br/>暂无任何文档</td></tr>
	</table>

	<div class='flickr_pagination' style="text-align:center;"></div>

	</div>
</div>


<script type="text/javascript">

$().ready(function(){

	// 定义快捷键
	var keyCodes = { 'code.65': 'a', 'code.66': 'b', 'code.67': 'c', 'code.68':'d', 'code.69': 'e', 'code.78': 'n', 'code.90': 'z' };
	var hotKeys = {
		'a': function(){ $('#btnAllDoc').click(); },
		'n': __new,
		'c': function(){ $('#btnCatelist').click(); }
	};

	// 绑定快捷键
	$(document).bind('keyup', function(evt){
		// alert(evt.keyCode)
		var letter = keyCodes['code.' + evt.keyCode];
		if(letter != null) {
			var fun = hotKeys[letter];
			if(fun != null) fun();
		}
	});

	if($.cookie('feedback_tip') != '1') {
		$('#btnFeedback').grumble({
			text: '亲爱的，天天笔记新版已上线，请多提意见哦～～！', 
			angle: 85, 
			distance: 30, 
			showAfter: 500,
			hideOnClick: true,
			onHide: function(){
				$.cookie('feedback_tip', '1');
			}
		});
	}

	// 表格行的隐现
	$('#page table tr:not([class])').hover(function(){
		$(this).addClass('gray');
		$(this).find('span').removeClass('hide');
	}, function(){
		$(this).removeClass('gray');
		$(this).find('span').addClass('hide');
	})

});

function __new(){ location.href = '/doc/new/'; }

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
