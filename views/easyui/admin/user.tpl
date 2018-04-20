{{template "../public/header.tpl"}}
<script type="text/javascript" src="/static/js/md5.js"></script>
<script type="text/javascript">
var statuslist = [
    {statusid:'0',name:'禁用'},
    {statusid:'1',name:'启用'}
];
var URL="/admin/user";

$.extend($.fn.validatebox.defaults.rules, {  
    /*必须和某个字段相等*/
    equalTo: {
        validator:function(value,param){
            return $(param[0]).val() == value;
        },
        message:'字段不匹配'
    }
});

$(function(){
    //用户列表
    $("#userdg").datagrid({
        title:'用户列表',
        loader: function(param, success, error){
            return loadEasyUIData(URL + "/list", param, success, error);
        },
        method:'POST',
        pagination:true,
        fitColumns:true,
        striped:true,
        rownumbers:true,
        singleSelect:true,
        idField:'id',
        pagination:true,
        pageSize:20,
        pageList:[10,20,30,50,100],
        columns:[[
            {field:'id',title:'ID',width:50,sortable:true},
            {field:'username',title:'用户名',width:100,sortable:true},
            {field:'nickname',title:'昵称',width:100,align:'center',editor:'text'},
            {field:'mobile',title:'手机',width:100,align:'center',editor:'text'},
            {field:'remark',title:'备注',width:150,align:'center',editor:'text'},
            {field:'status',title:'状态',width:50,align:'center',
                formatter:function(value){
                    for(var i=0; i<statuslist.length; i++){
                        if (statuslist[i].statusid == value) return statuslist[i].name;
                    }
                    return value;
                },
                editor:{
                    type:'combobox',
                    options:{
                        valueField:'statusid',
                        textField:'name',
                        data:statuslist,
                        required:true
                    }
                }
            },
            {field:'last_login_time',title:'上次登录时间',width:100,align:'center',
                formatter:function(value,row,index){
                    if(value) return phpjs.date("Y-m-d H:i:s",phpjs.strtotime(value));
                    return value;
                }
            },
            {field:'created_time',title:'添加时间',width:100,align:'center',
                formatter:function(value,row,index){
                    if(value) return phpjs.date("Y-m-d H:i:s",phpjs.strtotime(value));
                    return value;
                }
            }
        ]],
        onDblClickRow:function(index,row){
            editRow();
        },
        onRowContextMenu:function(e, index, row){
            e.preventDefault();
        },
        onHeaderContextMenu:function(e, field){
            e.preventDefault();
        }
    });
    //创建用户窗口
    $("#userDlg").dialog({
        modal:true,
        resizable:true,
        top:150,
        closed:true,
        buttons:[{
            id:'btnAdd',
            text:'添加',
            iconCls:'icon-add',
            handler:onBtnAdd
        },
        {
            id:'btnUpdate',
            text:'更新',
            iconCls:'icon-edit',
            handler:onBtnUpdate
        },
        {
            text:'取消',
            iconCls:'icon-cancel',
            handler:function(){
                $("#userDlg").dialog("close");
            }
        }],
        onOpen: onUserDlgOpen
    });

    userAccessList = new Array();
    $('#authtree').treegrid({
        idField:'id',
        treeField:'title',
        fitColumns:true,
        singleSelect:false,
        columns:[[
            {title:'标题',field:'title', width:300},
            {title:'是否需要授权访问',field:'auth', width:150,
                formatter:function(value) {
                if (value == 0) { return "无需授权";}
                if (value == 1) { return "授权访问";}
                return value;
            }
        }
        ]],
        loader:function(param,success,error){
            $.ajax({
                type: "get",
                data: param,
                url: "/public/accesslist",
                success: function (data) {
                    if(data.status == 1){
                        success(data.protocol);
                    } else if (data.status == 307 || data.status == 302) {
                        parent.document.location = data.protocol
                    } else {
                        error(data.info);
                    }
                },
                error: function (data) {
                    var tip = JSON.stringfy(data);
                    error(tip);
                }
            });
            return true
        },
        onSelect:function(row) {
            $(this).treegrid('expandAll',row.id);
            if(row._parentId != undefined && row._parentId != 0){
                $(this).treegrid('select',row._parentId);
            }
        },
        onUnselect:function(row) {
            if(row.children != undefined){
                for(var i=0;i<row.children.length;i++){
                    $(this).treegrid('unselect',row.children[i].id);
                }
            }
        }
    });

    //创建设置权限窗口
    $("#authDlg").dialog({
        modal:true,
        resizable:true,
        top:150,
        closed:true,
        buttons:[{
            text:'设置',
            iconCls:'icon-save',
            handler:function(){
                var tdata = $("#authtree").treegrid('getSelections');
                userAccessList = new Array(tdata.length);
                for(var i=0;i<tdata.length;i++){
                    userAccessList[i] = tdata[i].id;
                }
                $("#authDlg").dialog("close");
            }
        },
        {
            text:'关闭',
            iconCls:'icon-cancel',
            handler:function(){
                $("#authDlg").dialog("close");
            }
        }],
        onOpen: onAuthDlgOpen
    });
})

function onBtnAdd() {
    $.messager.progress({
        title: '提示',
        msg: '正在提交到服务器，请稍候……',
    });

    if (!$("#userForm").form('validate')) {
        $.messager.progress('close');
        vac.alert("请按要求先完整填写所有的必填项！");
        return;
    }

    var status = 0;
    if ($("#statusCheckBox").prop("checked")) {
        status = 1;
    }

    var info = {
        username:   $("#unTextBox").val(),
        password:   hex_md5($("#pwdTextBox").val()),
        nickname:   $("#nickTextBox").val(),
        mobile:     $("#mobileTextBox").val(),
        remark:     $("#remarkTextArea").val(),
        status:     status,
        accesses:   userAccessList
    }

    var data = new Object();
    data["inserts"] = JSON.stringify(info);

    $.ajax({
        type: "post",
        dataType: "json",
        data: data,
        url: URL + "/add",
        success: function (res) {
            $.messager.progress('close');
            if (res.status) {
                $("#userDlg").dialog("close");
                reload();
            } else {
                $.messager.alert('保存失败', res.info, 'error');
            }
        },
        error: function (res) {
            $.messager.progress('close');
            var tip = JSON.stringify(res);
            $.messager.alert('失败', tip, 'error');
        }
    });
}

function onBtnUpdate() {
    $.messager.progress({
        title: '提示',
        msg: '正在提交到服务器，请稍候……',
    });
    
    if (!$("#userForm").form('validate')) {
        $.messager.progress('close');
        vac.alert("请按要求先完整填写所有的必填项！");
        return;
    }

    var status = 0;
    if ($("#statusCheckBox").prop("checked")) {
        status = 1;
    }

    var pwd = "";
    if (!$("#pwdTextBox").prop('readonly')) {
        pwd = hex_md5($("#pwdTextBox").val());
    }

    var info = {
        id:         editDataRow.id,
        username:   $("#unTextBox").val(),
        password:   pwd,
        nickname:   $("#nickTextBox").val(),
        mobile:     $("#mobileTextBox").val(),
        remark:     $("#remarkTextArea").val(),
        status:     status,
        accesses:   userAccessList
    }

    var data = new Object();
    data["updates"] = JSON.stringify(info);

    $.ajax({
        type: "post",
        dataType: "json",
        data: data,
        url: URL + "/update",
        success: function (res) {
            $.messager.progress('close');
            if (res.status) {
                $("#userDlg").dialog("close");
                reload();
            } else {
                $.messager.alert('保存失败', res.info, 'error');
            }
        },
        error: function (res) {
            $.messager.progress('close');
            var tip = JSON.stringify(res);
            $.messager.alert('失败', tip, 'error');
        }
    });
}

function onUserDlgOpen() {
    $("#userForm").form('clear');

    userAccessList = new Array();

    // 根据新建还是修改设置正确的按钮显示
    if (addNewUser) {
        $("#btnAdd").show();
        $("#btnUpdate").hide();//隐藏更新按钮
        $("#btnModifyPwd").hide();
        $("#unTextBox").prop('readonly', false);
        $("#statusCheckBox").prop("checked", true);
        $("#pwdTextBox").prop('readonly', false);
        $("#repeatTextBox").prop('readonly', false);

        $.ajax({
            type: "get",
            dataType: "json",
            data: {uid: 0}, 
            url: URL + '/accesses',
            success: function (rsp) {
                if(rsp.status == 1){
                    userAccessList = rsp.protocol.list;
                } else if (rsp.status == 307 || rsp.status == 302) {
                    parent.document.location = rsp.protocol
                } else {
                    vac.alert(rsp.info);
                }
            }
        });

    } else {
        $("#btnAdd").hide();//隐藏添加按钮
        $("#btnUpdate").show();
        $("#btnModifyPwd").show();
        $("#unTextBox").prop('readonly', true);
        $("#pwdTextBox").prop('readonly', true);
        $("#repeatTextBox").prop('readonly', true);

        setEditData2Dlg();

        $.ajax({
            type: "get",
            dataType: "json",
            data: {uid: editDataRow.id}, 
            url: URL + '/accesses',
            success: function (rsp) {
                if(rsp.status == 1){
                    userAccessList = rsp.protocol.list;
                } else if (rsp.status == 307 || rsp.status == 302) {
                    parent.document.location = rsp.protocol
                } else {
                    vac.alert(rsp.info);
                }
            }
        });
    }
}

function onBtnModifyPwd() {
    $("#pwdTextBox").val("");
    $("#repeatTextBox").val("");
    $("#pwdTextBox").prop('readonly', false);
    $("#repeatTextBox").prop('readonly', false);
}

function onAuthDlgOpen() {
    var tg = $('#authtree')
    tg.treegrid('unselectAll');
    selectTreeGridItems(tg, userAccessList)
}

function selectTreeGridItems(tg, data) {
     //选中已存在的对应关系
    for(var i=0;i<data.length;i++){
        tg.treegrid('select', data[i]);
    }
}

function onUnselectRow(row) {
}

function onBtnAuth() {
    var dlg = $("#authDlg")
    dlg.dialog('open');
}

function editRow(){
    editDataRow = $("#userdg").datagrid("getSelected")
    if(!editDataRow){
        vac.alert("请选择要编辑的行");
        return;
    }

    addNewUser = false;

    var dlg = $("#userDlg")
    dlg.dialog('open');
    dlg.panel({title:"修改用户信息"});
}

function setEditData2Dlg() {
    if (!editDataRow) {
        return;
    }

    $("#unTextBox").val(editDataRow.username);
    $("#nickTextBox").val(editDataRow.nickname);
    $("#mobileTextBox").val(editDataRow.mobile);
    $("#remarkTextArea").val(editDataRow.remark);
    if (editDataRow.status == 1) {
        $("#statusCheckBox").prop("checked", true);
    }

    $("#pwdTextBox").val("12345");
    $("#repeatTextBox").val("12345");
    $("#userForm").form('validate');
}

//刷新
function reload(){
    $("#userdg").datagrid("reload");
}

//添加用户弹窗
function addRow(){
    addNewUser = true;

    var dlg = $("#userDlg")
    dlg.dialog('open');
    dlg.panel({title:"添加新用户"});
}

//删除
function delRow(){
    var row = $("#userdg").datagrid("getSelected");
    if(!row){
        vac.alert("请选择要删除的用户");
        return;
    }

    var tip = '你确定要删除用户【' + row. username +'】? 删除成功后将不可恢复。'
    $.messager.confirm('Confirm', tip, function(r){
        if (r){
            vac.ajax(URL+'/delete', {uid:row.id}, 'POST', function(r){
                if(!r.status){
                    vac.alert(r.info);
                }

                $("#userdg").datagrid('reload');
            })
        }
    });
}

</script>

<body>
<table id="userdg" toolbar="#tb"></table>
<div id="tb" style="padding:5px;height:auto">
    <a href="#" icon='icon-add' plain="true" onclick="addRow()" class="easyui-linkbutton" >新增</a>
    <a href="#" icon='icon-remove' plain="true" onclick="delRow()" class="easyui-linkbutton" >删除</a>
    <a href="#" icon='icon-edit' plain="true" onclick="editRow()" class="easyui-linkbutton" >编辑</a>
    <a href="#" icon='icon-reload' plain="true" onclick="reload()" class="easyui-linkbutton" >刷新</a>
</div>

<div id="userDlg" title="添加用户" style="width:400px;height:400px;">
    <div style="padding:20px 20px 20px 20px;" >
        <form id="userForm">
            <table>
                <tr>
                    <td>用户名：</td>
                    <td><input id="unTextBox" class="easyui-validatebox" required="true"/></td>
                </tr>
                <tr>
                    <td>昵称：</td>
                    <td><input id="nickTextBox" class="easyui-validatebox" required="true"  /></td>
                </tr>
                <tr>
                    <td>密码：</td>
                    <td><input id="pwdTextBox" name="Password" type="password" class="easyui-validatebox" required="true" validType="password[5,20]" missingMessage="请填写密码"/></td>
                    <td><a id="btnModifyPwd" href="#" icon='icon-edit' plain="true" onclick="onBtnModifyPwd()" class="easyui-linkbutton">修改密码</a></td>
                </tr>
                <tr>
                    <td>重复密码：</td>
                    <td><input id="repeatTextBox" type="password" class="easyui-validatebox" required="true" validType="equalTo['#pwdTextBox']" missingMessage="请重复填写密码" invalidMessage="两次输入密码不匹配" /></td>
                </tr>
                <tr>
                    <td>手机：</td>
                    <td><input id="mobileTextBox" class="easyui-validatebox" validType="mobile" required="true" /></td>
                </tr>
                <tr>
                    <td>状态：</td>
                    <td>启用<input id="statusCheckBox" type="checkbox" value=1/></td>
                </tr>
                <tr>
                    <td>权限列表：</td>
                    <td><a id="btnAuth" href="#" icon='icon-edit' plain="true" onclick="onBtnAuth()" class="easyui-linkbutton">设置</a></td>
                </tr>
                <tr>
                    <td>备注：</td>
                    <td><textarea id="remarkTextArea" class="easyui-validatebox" validType="length[0,512]"></textarea></td>
                </tr>
            </table>
        </form>
    </div>
</div>

<div id="authDlg" title="修改用户授权" style="width:600px;height:400px;">
    <form id="authForm">
        <table id="authtree" style="width:600px;height:400px"></table>
    </form>
</div>
</body>
</html>