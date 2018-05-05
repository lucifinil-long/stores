{{template "../public/header.tpl"}}
<script type="text/javascript" src="/static/js/md5.js"></script>
<script type="text/javascript">
    // 关闭浏览器页面的时候，发起退出登录操作
    window.onbeforeunload = function() {
        $.ajax({
            type: "post",
            url: "/public/admin/logout"
        });
    } 
    
    var URL="/public/admin"
    $.extend($.fn.validatebox.defaults.rules, {  
        /*必须和某个字段相等*/
        equalTo: {
            validator:function(value,param){
                return $(param[0]).val() == value;
            },
            message:'字段不匹配'
        }
            
    });

    $( function() {
        //生成树
        $("#tree").tree({
            loader:function(param, success, error){
                return loadEasyUIData(URL + "/treemenu", param, success, error);
            },
            onClick:function(node){
                if(node.attributes.url == ""){
                    $(this).tree("toggle",node.target);
                    return false;
                }
                var href = node.attributes.url;
                var tabs = $("#tabs");
                if(href){
                    var content = '<iframe scrolling="auto" frameborder="0"  src="'+href+'" style="width:100%;height:100%;"></iframe>';
                }else{
                    var content = '未实现';
                }
                //已经存在tabs则选中它
                if(tabs.tabs('exists',node.text)){
                    //选中
                    tabs.tabs('select',node.text);
                    //refreshTab(node.text);
                }else{
                    //添加
                    tabs.tabs('add',{
                        title:node.text,
                        content:content,
                        closable:true,
                        cache:false,
                        fit:'true'
                    });
                }
            }
        });

        // 生成选项页
        $("#tabs").tabs({
            width: $("#tabs").parent().width(),
            height: "auto",
            fit:true,
            border:false,
            onContextMenu : function(e, title) {
                e.preventDefault();
                $("#mm").menu('show', {
                    left : e.pageX,
                    top : e.pageY
                }).data('tabTitle', title);
            }
        });

        // 生成快捷菜单
        $('#mm').menu({
            onClick : function(item) {
                var curTabTitle = $(this).data('tabTitle');
                var type = $(item.target).attr('type');

                if (type === 'refresh') {
                    refreshTab(curTabTitle);
                    return;
                }

                if (type === 'close') {
                    var t = $("#tabs").tabs('getTab', curTabTitle);
                    if (t.panel('options').closable) {
                        $("#tabs").tabs('close', curTabTitle);
                    }
                    return;
                }

                var allTabs = $("#tabs").tabs('tabs');
                var closeTabsTitle = [];

                $.each(allTabs, function() {
                    var opt = $(this).panel('options');
                    if (opt.closable && opt.title != curTabTitle && type === 'closeOther') {
                        closeTabsTitle.push(opt.title);
                    } else if (opt.closable && type === 'closeAll') {
                        closeTabsTitle.push(opt.title);
                    }
                });
                for ( var i = 0; i < closeTabsTitle.length; i++) {
                    $("#tabs").tabs('close', closeTabsTitle[i]);
                }
            }
        });
        //修改配色方案
        $("#changetheme").change(function(){
            var theme = $(this).val();
            $.cookie("theme",theme); //新建cookie
            location.reload();
        });
 
    });

    function refreshTab(title) {
        var tab = $("#tabs").tabs("getTab", title);
        $("#tabs").tabs("update", {tab: tab, options: tab.panel("options")});
    }
    function undo(){
        $('#tree').tree('expandAll');
    }
    function redo(){
        $('#tree').tree('collapseAll');
    }

    function modifypassword(){
        $("#dialog").dialog({
            modal:true,
            title:"修改密码",
            width:400,
            height:250,
            buttons:[{
                text:'保存',
                iconCls:'icon-save',
                handler:modifyPwdSubmit
            },
            {
                text:'取消',
                iconCls:'icon-cancel',
                handler:function(){
                    $("#dialog").dialog("close");
                }
            }]
        });
    }

    // 修改密码提交
    function modifyPwdSubmit() {
        if (!$("#modifyPwdForm").form('validate')) {
            vac.alert("请按要求完整填写必要信息！");
            return;
        }
        
        var oldPwd = hex_md5($("#oldTextBox").val());
        var newPwd = hex_md5($("#newTextBox").val());
        var repeatPwd = hex_md5($("#repeatTextBox").val());

        var info = new Object();
        info["old"] = oldPwd;
        info["new"] = newPwd;
        info["repeat"] = repeatPwd;

        $.ajax({
            type: "post",
            dataType: "json",
            data: info,
            url: URL + "/changepwd",
            success: function (data) {
                if(data.status){
                    $.messager.alert("提示", data.info,'info',function(){
                        Logout();
                    });
                }else {
                    vac.alert(data.info);
                }
            },
            error: function (data) {
                var tip = JSON.stringfy(data);
                $.messager.alert('修改密码失败', tip, 'error');
            }
        });
    }

    function Logout() {
        $.ajax({
            type: "post",
            url: URL + "/logout",
            success: function (data) {
                if (data.status == 307 || data.status == 302) {
                    parent.document.location = data.protocol
                } else {
                    $.messager.alert('退出失败', data.info, 'error');
                }
            },
            error: function (data) {
                var tip = JSON.stringfy(data);
                $.messager.alert('退出失败', tip, 'error');
            }
        });
    }
</script>

<style>
.ht_nav {
    float: left;
    overflow: hidden;
    padding: 0 0 0 10px;
    margin: 0;
}
.ht_nav li{
    font:700 16px/2.5 'microsoft yahei';
    float: left;
    list-style-type: none;
    margin-right: 10px;

}
.ht_nav li a{
    text-decoration: none;
    color:#333;
}
.ht_nav li a.current, .ht_nav li a:hover{
    color:#F20;

}
</style>
<body class="easyui-layout" style="text-align:left">
<div region="north" border="false" style="overflow: hidden; width: 100%; height:82px; background:#D9E5FD;">
    <div style="overflow: hidden; width:200px; padding:2px 0 0 5px;">
        <h2>仓储管理系统</h2>
    </div>
    <div id="header-inner" style="float:right; overflow:hidden; height:80px; width:300px; line-height:25px; text-align:right; padding-right:20px;margin-top:-50px; ">
        {{.userinfo.Nickname}}, 欢迎你！ <a href="javascript:void(0);" onclick="modifypassword()"> 修改密码</a>
        <a href="javascript:void(0);" onclick="Logout()"> 退 出</a>
    </div>
</div>
<div id="dialog" >
    <div style="padding:20px 20px 40px 80px;" >
        <form id="modifyPwdForm" method="post">
            <table>
                <tr>
                    <td>旧密码</td>
                    <td><input type="password"  id="oldTextBox" class="easyui-validatebox"  required="true" validType="password[5,20]" missingMessage="请填写当前使用的密码"/></td>
                </tr>
                <tr>
                    <td>新密码：</td>
                    <td><input type="password"  id="newTextBox" class="easyui-validatebox" required="true" validType="password[5,20]" missingMessage="请填写需要修改的密码"  /></td>
                </tr>
                <tr>
                    <td>重复密码：</td>
                    <td><input type="password"  id="repeatTextBox"  class="easyui-validatebox" required="true" validType="equalTo['#newTextBox']" missingMessage="请重复填写需要修改的密码" invalidMessage="两次输入密码不匹配"/></td>
                </tr>
            </table>
        </form>
    </div>
</div>
</div>
<div region="west" border="false" split="true" title="菜单"  tools="#toolbar" style="width:200px;padding:5px;">
    <ul id="tree"></ul>
</div>
<div region="center" border="false" >
    <div id="tabs" >
    </div>
</div>
<div id="toolbar">
    <a href="#" class="icon-undo" title="全部展开"  onclick="undo()"></a>
    <a href="#" class="icon-redo" title="全部关闭"  onclick="redo()"></a>
</div>
<!--右键菜单-->
<div id="mm" style="width: 120px;display:none;">
    <div iconCls='icon-reload' type="refresh">刷新</div>
    <div class="menu-sep"></div>
    <div  type="close">关闭</div>
    <div type="closeOther">关闭其他</div>
    <div type="closeAll">关闭所有</div>
</div>
</body>
</html>