{{template "../public/header.tpl"}}
<script type="text/javascript" src="/static/js/md5.js"></script>
<script type="text/javascript">
var URL="/public"
$(function() {
    $.ajax({
        type: "get",
        url: URL + "/isloggedin",
        success: function (data) {
            if (data.status == 1) {
                if (data.protocol.status) {
                    location.href = data.protocol.redirect;
                }
            } else if (data.status == 307 || data.status == 302) {
                location.href = data.protocol;
            }
        }
    });

    $("#dialog").dialog({
        closable:false,
        buttons:[{
            text:'登录',
            iconCls:'icon-save',
            handler:loginSubmit
        },
        {
            text:'重置',
            iconCls:'icon-cancel',
            handler:function(){
                $("#form").from("reset");
            }
        }]
    });
});

function loginSubmit() {
    if (!$("#form").form('validate')) {
        vac.alert("请完整填写必要信息！");
        return;
    }
    
    var pwd = hex_md5($("#pwdTextBox").val());
    var un = $("#unTextBox").val()

    var info = new Object();
    info["uid"] = un;
    info["password"] = pwd;

    $.ajax({
        type: "post",
        dataType: "json",
        data: info,
        url: URL + "/admin/login",
        success: function (data) {
            if (data.status == 1) {
                location.href = data.protocol.redirect;
            } else if (data.status == 307 || data.status == 302) {
                location.href = data.protocol;
            } else {
                vac.alert(data.info);
            }
        },
        error: function (data) {
            var tip = JSON.stringify(data);
            $.messager.alert('登陆失败', tip, 'error');
        }
    });
}
    //这个就是键盘触发的函数
var SubmitOrHidden = function(evt){
    evt = window.event || evt;
    if(evt.keyCode==13){//如果取到的键值是回车
          loginSubmit();       
     }
                
}
window.document.onkeydown=SubmitOrHidden;//当有键按下时执行函数
</script>
<body>
<div style="text-align:center;margin:0 auto;width:350px;height:250px;" id="dialog" title="后台用户登录">
<div style="padding:20px 20px 20px 40px;" >
<form id="form" method="post">
<table >
    <tr>
        <td>用户ID或手机：</td><td><input id="unTextBox" type="text" class="easyui-validatebox" required="true" name="uid" missingMessage="请输入用户ID或手机"/></td>
    </tr>
    <tr>
        <td>密码：</td><td><input id="pwdTextBox" type="password" class="easyui-validatebox" required="true" name="password" missingMessage="请输入密码"/></td>
    </tr>
</table>
</form>
</div>
</div>
</body>
</html>