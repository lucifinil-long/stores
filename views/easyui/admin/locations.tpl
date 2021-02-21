{{template "../public/header.tpl"}}
<script type="text/javascript">
var URL="/admin/depot";
var depotID = 0;

$(function(){
    
    // 库房列表
    $("#locationDg").datagrid({
        title:'库房列表',
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
            {field:'name',title:'库房名称',width:100,align:'center'},
            {field:'detail',title:'详情',width:100,align:'center'},
            {field:'shelfs',title:'货架列表',width:300,align:'center',
                formatter:function(value,row,index){
                    if (value) {
                        var arr = new Array();
                        for (var i = 0; i < value.length; i++) {
                            arr.push(value[i].name);
                        }

                        return arr.join(", ")
                    }
                    return "";
                }
            }
        ]],
        onSelect:function(index,row){
            onRowSelected();
        },
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

    // 库房货架列表
    $("#shelfsDg").datagrid({
        title:'货架列表',
        loader: function(param, success, error) {
            return loadEasyUIData(URL + "/shelf/list", param, success, error);
        },
        queryParams:{
            depot_id: depotID
        },
        fitColumns:true,
        striped:true,
        rownumbers:true,
        singleSelect:true,
        idField:'id',
        pagination:true,
        pageSize:20,
        pageList:[10,20],
        columns:[[
            {field:'name',title:'货架名称',width:100,align:'center'},
            {field:'layers',title:'层数',width:100,align:'center'},
            {field:'detail',title:'货架详情',width:300,align:'center'},
        ]],
        onDblClickRow:function(index,row){
            editShelf();
        },
        onRowContextMenu:function(e, index, row){
            e.preventDefault();
        },
        onHeaderContextMenu:function(e, field){
            e.preventDefault();
        }
    });
    $("#shelfsList").hide();

    // 创建仓库窗口
    $("#depotDlg").dialog({
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
                $("#depotDlg").dialog("close");
            }
        }],
        onOpen: onDepotDlgOpen
    });

    // 创建仓库货架窗口
    $("#shelfDlg").dialog({
        modal:true,
        resizable:true,
        top:150,
        closed:true,
        buttons:[{
            id:'btnAddShelf',
            text:'添加',
            iconCls:'icon-add',
            handler:onBtnAddShelf
        },
        {
            id:'btnUpdateShelf',
            text:'更新',
            iconCls:'icon-edit',
            handler:onBtnUpdateShelf
        },
        {
            text:'取消',
            iconCls:'icon-cancel',
            handler:function(){
                $("#shelfDlg").dialog("close");
            }
        }],
        onOpen: onShelfDlgOpen
    });
});

function onRowSelected() {
    selectedRow = $("#locationDg").datagrid("getSelected")
    if(!selectedRow){
        return;
    }

    depotID = selectedRow.id;
    $("#shelfsList").show();
    reloadShelfs();
}

function reloadShelfs() {
    var queryParams = $('#shelfsDg').datagrid('options').queryParams;
    queryParams.depot_id = depotID;
    $("#shelfsDg").datagrid("reload");
    $("#shelfsDg").datagrid("resize");
}

function editRow() {
    editDataRow = $("#locationDg").datagrid("getSelected")
    if(!editDataRow){
        vac.alert("请选择要编辑的行");
        return;
    }
    addNewDepot = false;

    var dlg = $("#depotDlg")
    dlg.dialog('open');
    dlg.panel({title:"修改仓库信息"});
}

//刷新
function reload() {
    $("#locationDg").datagrid("reload");
}

//添加弹窗
function addRow() {
    addNewDepot = true;

    var dlg = $("#depotDlg")
    dlg.dialog('open');
    dlg.panel({title:"添加新仓库"});
}

//删除
function delRow() {
    var row = $("#locationDg").datagrid("getSelected");
    if(!row){
        vac.alert("请选择要删除的库房");
        return;
    }

    var deletedRows = new Array();
    deletedRows.push(row.id);
    var modified = new Object();
    modified["deletes"] = JSON.stringify(deletedRows);
    var tip = '你确定要删除库房【' + row.name +'】? 删除成功后将不可恢复。';
    $.messager.confirm('Confirm', tip, function(r){
        if (r){
            vac.ajax(URL+'/delete', modified, 'POST', function(r){
                if(!r.status){
                    vac.alert(r.info);
                } else {
                    $("#shelfsList").hide();
                }

                reload();
            })
        }
    });
}

function onBtnAdd() {
    $.messager.progress({
        title: '提示',
        msg: '正在提交到服务器，请稍候……',
    });

    if (!$("#depotForm").form('validate')) {
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
                $("#depotDlg").dialog("close");
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
    
    if (!$("#depotForm").form('validate')) {
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
                $("#depotDlg").dialog("close");
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

function onDepotDlgOpen() {
    $("#depotForm").form('clear');

    // 根据新建还是修改设置正确的按钮显示
    if (addNewDepot) {
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

// shelf related
function editShelf() {
    editDataShelfRow = $("#shelfsDg").datagrid("getSelected")
    if(!editDataShelfRow){
        vac.alert("请选择要编辑的行");
        return;
    }
    addNewShelf = false;

    var dlg = $("#shelfDlg")
    dlg.dialog('open');
    dlg.panel({title:"修改仓库货架信息"});
}

function addShelf() {
    addNewShelf = true;

    var dlg = $("#shelfDlg")
    dlg.dialog('open');
    dlg.panel({title:"添加新仓库货架"});
}

//删除
function delShelf() {
    var row = $("#shelfsDg").datagrid("getSelected");
    if(!row){
        vac.alert("请选择要删除的货架");
        return;
    }

    var deletedRows = new Array();
    deletedRows.push(row.id);
    var modified = new Object();
    modified["deletes"] = JSON.stringify(deletedRows);
    var tip = '你确定要删除货架【' + row.name +'】? 删除成功后将不可恢复。';
    $.messager.confirm('Confirm', tip, function(r){
        if (r){
            vac.ajax(URL+'/shelf/delete', modified, 'POST', function(r){
                if(!r.status){
                    vac.alert(r.info);
                } 

                reloadShelfs();
            })
        }
    });
}

function onBtnAddShelf() {
    $.messager.progress({
        title: '提示',
        msg: '正在提交到服务器，请稍候……',
    });

    if (!$("#shelfForm").form('validate')) {
        $.messager.progress('close');
        vac.alert("请按要求填写信息！");
        return;
    }

    var info = {
        name:       $("#shelfNameTextBox").val(),
        layers:     Number($("#shelfLayersTextBox").val()),
        detail:     $("#shelfDetailTextBox").val()
    }

    var shelfs = new Array();
    shelfs.push(info);

    var req =  {
        depot_id : depotID,
        shelfs : shelfs
    }

    var data = new Object();
    data["req"] = JSON.stringify(req);

    $.ajax({
        type: "post",
        dataType: "json",
        data: data,
        url: URL + "/shelf/add",
        success: function (res) {
            $.messager.progress('close');
            if (res.status) {
                $("#shelfDlg").dialog("close");
                reload();
                reloadShelfs();
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

function onBtnUpdateShelf() {
    $.messager.progress({
        title: '提示',
        msg: '正在提交到服务器，请稍候……',
    });
    
    if (!$("#shelfForm").form('validate')) {
        $.messager.progress('close');
        vac.alert("请按要求填写信息！");
        return;
    }

    var info = {
        id:         editDataShelfRow.id,
        name:       $("#shelfNameTextBox").val(),
        layers:     Number($("#shelfLayersTextBox").val()),
        detail:     $("#shelfDetailTextBox").val()
    }

    var data = new Object();
    data["updated"] = JSON.stringify(info);

    $.ajax({
        type: "post",
        dataType: "json",
        data: data,
        url: URL + "/shelf/update",
        success: function (res) {
            $.messager.progress('close');
            if (res.status) {
                $("#shelfDlg").dialog("close");
                reload();
                reloadShelfs();
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

function onShelfDlgOpen() {
    $("#shelfForm").form('clear');

    // 根据新建还是修改设置正确的按钮显示
    if (addNewShelf) {
        $("#btnAddShelf").show();
        $("#btnUpdateShelf").hide();//隐藏更新按钮
    } else {
        $("#btnAddShelf").hide();//隐藏添加按钮
        $("#btnUpdateShelf").show();

        setEditData2ShelfDlg();
    }
}

function setEditData2ShelfDlg() {
    if (!editDataShelfRow) {
        return;
    }

    $("#shelfNameTextBox").val(editDataShelfRow.name);
    $("#shelfLayersTextBox").val(editDataShelfRow.layers);
    $("#shelfDetailTextBox").val(editDataShelfRow.detail);
}

</script>

<body>

<table id="locationDg" toolbar="#tb"></table>
<div id="tb" style="padding:5px;height:auto">
    <a href="#" icon='icon-add' plain="true" onclick="addRow()" class="easyui-linkbutton" >新增</a>
    <a href="#" icon='icon-remove' plain="true" onclick="delRow()" class="easyui-linkbutton" >删除</a>
    <a href="#" icon='icon-edit' plain="true" onclick="editRow()" class="easyui-linkbutton" >编辑</a>
    <a href="#" icon='icon-reload' plain="true" onclick="reload()" class="easyui-linkbutton" >刷新</a>
</div>

<div id="shelfsList">
<table id="shelfsDg" toolbar="#tb2"></table>
<div id="tb2" style="padding:5px;height:auto">
    <a href="#" icon='icon-add' plain="true" onclick="addShelf()" class="easyui-linkbutton" >新增</a>
    <a href="#" icon='icon-remove' plain="true" onclick="delShelf()" class="easyui-linkbutton" >删除</a>
    <a href="#" icon='icon-edit' plain="true" onclick="editShelf()" class="easyui-linkbutton" >编辑</a>
    <a href="#" icon='icon-reload' plain="true" onclick="reloadShelfs()" class="easyui-linkbutton" >刷新</a>
</div>
</div>

<div id="depotDlg" title="添加新仓库" style="width:300px;height:200px;">
    <div style="padding:30px 30px 30px 30px;" >
        <form id="depotForm">
            <table>
                <tr>
                    <td>仓库名称：</td>
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

<div id="shelfDlg" title="添加新仓库货架" style="width:300px;height:200px;">
    <div style="padding:30px 30px 30px 30px;" >
        <form id="shelfForm">
            <table>
                <tr>
                    <td>货架名称：</td>
                    <td><input id="shelfNameTextBox" class="easyui-validatebox" required="true" validType="length[0,64]"/></td>
                </tr>
                <tr>
                    <td>货架层数：</td>
                    <td><input id="shelfLayersTextBox" class="easyui-validatebox" validType="number"/></td>
                </tr>
                <tr>
                    <td>详情：</td>
                    <td><input id="shelfDetailTextBox" class="easyui-validatebox" validType="length[0,64]"/></td>
                </tr>
            </table>
        </form>
    </div>
</div>

</body>
</html>