{{template "../public/header.tpl"}}
<script type="text/javascript">
var URL="/admin/operations";
$(function(){
    // 仓储管理系统操作日志列表
    $("#datagrid").datagrid({
        title:'仓储管理系统操作日志列表',
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
            {field:'username',title:'用户名',width:80,sortable:true},
            {field:'user_id',title:'用户ID',width:80,sortable:true},
            {field:'from',title:'用户来源',width:100,formatter:formatCellTooltip},
            {field:'action',title:'动作',width:100,align:'center',formatter:formatCellTooltip},
            {field:'detail',title:'详情',width:350,align:'center',formatter:formatCellTooltip},
            {field:'created_time',title:'发生时间',width:100,align:'center',
                formatter:function(value,row,index){
                    if(value) return phpjs.date("Y-m-d H:i:s",phpjs.strtotime(value));
                    return value;
                },sortable:true
            }
        ]],
        onAfterEdit:function(index, data, changes){
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

})

//刷新
function reloadrow(){
    $("#datagrid").datagrid("reload");
}

// 格式化单元格提示信息  
function formatCellTooltip(value){  
    return "<span title='" + value + "'>" + value + "</span>";  
} 

</script>
<body>
<table id="datagrid" toolbar="#tb"></table>
<div id="tb" style="padding:5px;height:auto">
    <a>指定用户</a><input type="text" name="uname">
    <a href="#" icon='icon-search' plain="true" onclick="reloadrow()" class="easyui-linkbutton" >搜索</a>
</div>
</body>
</html>