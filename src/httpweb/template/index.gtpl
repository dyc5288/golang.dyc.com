
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN"
  "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd"> 
<html> 
  <head>
    <meta http-equiv="Content-type" content="text/html; charset=utf-8">
    <meta content="authenticity_token" name="csrf-param" />
<meta content="/BoPOb532Df7asPf1ltX5AM5MXA9z1I9mPXfDeSues4=" name="csrf-token" />
    <title>【天天笔记】网络记事本_在线电子记事本软件</title> <meta name="keywords" content="记事本,网络记事本,在线记事本" /> <meta name="description" content="天天笔记是当前国内最好用的网络记事本、在线电子记事本软件，包括文档、便笺、日记本、收藏夹等服务，功能简洁、操作方便、每天备份数据，并且支持手机浏览，让您随时随地拥有云笔记服务。" />
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
      <div id='menu' style="background:url('/images/banner.jpg');"> 
  <div class="container">
    <div id="hero" style="filter:alpha(opacity=0); opacity:0; -moz-opacity:0; -khtml-opacity:0;";><img src="/images/home-hero.png" alt="天天笔记" data-pinit="registered"></div>
    <div class="slogan"><h1>天天笔记</h1><small>一款风格简洁、操作方便、体验良好的免费网络在线记事本|专门针对电信和联通进行优化，加快数据访问速度|高质量硬盘、定期数据备份，用户数据加密，支持数据导出|文档共享功能，方便您将信息分享给家人、同事、或朋友们|支持手机在线访问，随时随地浏览及编辑、速度快、节省流量</small><p>
      <a href="/register/" class="button noline">立即注册</a>
      </p></div>
  </div>
</div>
<div id='content'>
  <div class="container">
    <ul>
      <li>
        <h2>文档</h2>
        <p>您的生活、工作、学习笔记</p>
      </li>
      <li>
        <h2>便笺</h2>
        <p>记录一些简短的文字</p>
      </li>
      <li>
        <h2>备份</h2>
        <p>支持导出功能、方便备份</p>
      </li>
      <li class="last">
        <h2>手机</h2>
        <p>手机在线浏览、随时随地使用</p>
      </li>
    </ul>
  </div>
</div>

<script type="text/javascript">
$(document).ready(function(){
  var index = 0;
  var arr = $('#menu small').html().split('|');
  $('#menu small').html('');
  $(arr).each(function(index, item){
    $('#menu small').append('<span style="display:none;">' + item + '</span>');
  });

  showTip();
  setInterval(showTip, 6000);
  setTimeout(function(){
    $('#hero').fadeTo(2500, 1);
  }, 1500);
});

// 呈现个性化文字
var tipIndex = -1;
function showTip() {
  var list = $('#menu small span');
  tipIndex ++;
  if(tipIndex >= list.length)
    tipIndex = 0;
  list.hide();
  $('#menu small span:eq('+ tipIndex +')').fadeIn(800);
}
</script>
      <div id='footer'> 
        <div class="container">
        © 2012-2013 – <a href="/">天天笔记</a> - TTBIJI.com - Layout based on bootstrap.
        <br/>友情链接：<a href="http://www.hoowoo.com" target="_blank">HOOWOO</a>&nbsp;
<a href="http://www.quttang.com" target="_blank">趣糖网</a>&nbsp;
<a href="http://www.csjita.com/" target="_blank">吉他谱</a>&nbsp;
<a href="http://www.7lizhi.com/" target="_blank">青春励志签名</a>
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
