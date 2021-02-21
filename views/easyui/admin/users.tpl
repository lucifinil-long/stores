{{template "../public/header.tpl"}}
<script type="text/javascript" src="/static/js/md5.js"></script>
<script type="text/javascript">
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
            {field:'nickname',title:'用户名',width:100,align:'center'},
            {field:'mobile',title:'手机',width:100,align:'center'},
            {field:'role',title:'角色组',width:100,align:'center'},
            {field:'remark',title:'备注',width:150,align:'center'},
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
        password:   hex_md5($("#pwdTextBox").val()),
        nickname:   $("#nickTextBox").val(),
        mobile:     Number($("#mobileTextBox").val()),
        remark:     $("#remarkTextArea").val(),
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
        password:   pwd,
        nickname:   $("#nickTextBox").val(),
        mobile:     $("#mobileTextBox").val(),
        remark:     $("#remarkTextArea").val()
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

    // 根据新建还是修改设置正确的按钮显示
    if (addNewUser) {
        $("#btnAdd").show();
        $("#btnUpdate").hide();//隐藏更新按钮
        $("#btnModifyPwd").hide();
        $("#pwdTextBox").prop('readonly', false);
        $("#repeatTextBox").prop('readonly', false);
    } else {
        $("#btnAdd").hide();//隐藏添加按钮
        $("#btnUpdate").show();
        $("#btnModifyPwd").show();
        $("#pwdTextBox").prop('readonly', true);
        $("#repeatTextBox").prop('readonly', true);

        setEditData2Dlg();
    }
}

function onBtnModifyPwd() {
    $("#pwdTextBox").val("");
    $("#repeatTextBox").val("");
    $("#pwdTextBox").prop('readonly', false);
    $("#repeatTextBox").prop('readonly', false);
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

    $("#nickTextBox").val(editDataRow.nickname);
    $("#mobileTextBox").val(editDataRow.mobile);
    $("#remarkTextArea").val(editDataRow.remark);

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
                    <td>手机：</td>
                    <td><input id="mobileTextBox" class="easyui-validatebox" validType="mobile" required="true" /></td>
                </tr>
                <tr>
                    <td>用户名：</td>
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