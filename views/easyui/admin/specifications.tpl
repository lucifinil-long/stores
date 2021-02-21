{{template "../public/header.tpl"}}
<script type="text/javascript">
var URL="/admin/specifications";
var specID = 0;

$(function() {
    
    // 规格列表
    $("#specsDg").datagrid({
        title:'规格列表',
        loader: function(param, success, error){
            return loadEasyUIData(URL + "/list", param, success, error);
        },
        fitColumns:true,
        striped:true,
        rownumbers:true,
        singleSelect:true,
        idField:'id',
        pagination:true,
        pageSize:20,
        pageList:[10,20,30,50,100],
        columns:[[
            {field:'name',title:'名称',width:100,align:'center'},
            {field:'parent',title:'父级规格',width:100,align:'center',
                formatter:function(value,row,index) {
                    if (value) {
                        return value.name + " - ID:" + value.id;
                    }

                    return "";
                }
            },
            {field:'sub',title:'可拆分为子级规格',width:100,align:'center',
                formatter:function(value,row,index) {
                    if (value) {
                        return value.name + " - ID:" + value.id;
                    }

                    return "";
                }
            },
            {field:'sub_amount',title:'可拆分为子级规格数量',width:100,align:'center'},
            {field:'detail',title:'详情',width:300,align:'center'}
        ]],
        onSelect:function(index,row){
            onRowSelected();
        },
        onDblClickRow:function(index,row){
        },
        onRowContextMenu:function(e, index, row){
            e.preventDefault();
        },
        onHeaderContextMenu:function(e, field){
            e.preventDefault();
        }
    });

    // 创建规格窗口
    $("#specDlg").dialog({
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
                $("#specDlg").dialog("close");
            }
        }],
        onOpen: onSpecDlgOpen
    });
});

function addRow() {
    addNewSpec = true;

    var dlg = $("#specDlg")
    dlg.dialog('open');
    dlg.panel({title:"添加新规格"});
}

function delRow() {
}

function editRow() {
    editDataRow = $("#specsDg").datagrid("getSelected")
    if(!editDataRow){
        vac.alert("请选择要编辑的行");
        return;
    }
    addNewSpec = false;

    var dlg = $("#specDlg")
    dlg.dialog('open');
    dlg.panel({title:"修改规格信息"});
}

function reload() {
    $("#specsDg").datagrid("reload");
}

function onRowSelected() {
    selectedRow = $("#specsDg").datagrid("getSelected")
    if(!selectedRow){
        return;
    }
}

function onBtnAdd() {
    $.messager.progress({
        title: '提示',
        msg: '正在提交到服务器，请稍候……',
    });

    if (!$("#specForm").form('validate')) {
        $.messager.progress('close');
        vac.alert("请按要求填写信息！");
        return;
    }

    var info = {
        name:       $("#nameTextBox").val(),
        detail:     $("#detailTextBox").val()
    }

    var data = new Object();
    data["insert"] = JSON.stringify(info);

    $.ajax({
        type: "post",
        dataType: "json",
        data: data,
        url: URL + "/add",
        success: function (res) {
            $.messager.progress('close');
            if (res.status) {
                $("#specDlg").dialog("close");
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
    
    if (!$("#specForm").form('validate')) {
        $.messager.progress('close');
        vac.alert("请按要求填写信息！");
        return;
    }

    var info = {
        id:         editDataRow.id,
        name:       $("#nameTextBox").val(),
        detail:     $("#detailTextBox").val()
    }

    var data = new Object();
    data["updated"] = JSON.stringify(info);

    $.ajax({
        type: "post",
        dataType: "json",
        data: data,
        url: URL + "/update",
        success: function (res) {
            $.messager.progress('close');
            if (res.status) {
                $("#specDlg").dialog("close");
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

function onSpecDlgOpen() {
    $("#specForm").form('clear');

    // 根据新建还是修改设置正确的按钮显示
    if (addNewSpec) {
        $("#btnAdd").show();
        $("#btnUpdate").hide();//隐藏更新按钮
    } else {
        $("#btnAdd").hide();//隐藏添加按钮
        $("#btnUpdate").show();

        setEditData2Dlg();
    }
}

function setEditData2Dlg() {
    if (!editDataRow) {
        return;
    }

    $("#nameTextBox").val(editDataRow.name);
    $("#detailTextBox").val(editDataRow.detail);
}

</script>

<body>

<table id="specsDg" toolbar="#tb"></table>
<div id="tb" style="padding:5px;height:auto">
    <a href="#" icon='icon-add' plain="true" onclick="addRow()" class="easyui-linkbutton" >新增</a>
    <a href="#" icon='icon-remove' plain="true" onclick="delRow()" class="easyui-linkbutton" >删除</a>
    <a href="#" icon='icon-edit' plain="true" onclick="editRow()" class="easyui-linkbutton" >编辑</a>
    <a href="#" icon='icon-reload' plain="true" onclick="reload()" class="easyui-linkbutton" >刷新</a>
</div>

<div id="specDlg" title="添加新规格" style="width:300px;height:200px;">
    <div style="padding:30px 30px 30px 30px;" >
        <form id="specForm">
            <table>
                <tr>
                    <td>规格名称：</td>
                    <td><input id="nameTextBox" class="easyui-validatebox" required="true" validType="length[0,64]"/></td>
                </tr>
                <tr>
                    <td>详情：</td>
                    <td><input id="detailTextBox" class="easyui-validatebox" validType="length[0,64]"/></td>
                </tr>
            </table>
        </form>
    </div>
</div>
</body>
</html>