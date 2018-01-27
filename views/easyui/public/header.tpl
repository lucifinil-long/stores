<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<title>仓储管理系统</title>
    <link rel="stylesheet" type="text/css" href="/static/easyui/jquery-easyui/themes/default/easyui.css" />
    <link rel="stylesheet" type="text/css" href="/static/easyui/jquery-easyui/themes/icon.css" />
    <script type="text/javascript" src="/static/easyui/jquery-easyui/jquery.min.js"></script>
    <script type="text/javascript" src="/static/easyui/jquery-easyui/jquery.easyui.min.js"></script>
    <script type="text/javascript" src="/static/easyui/jquery-easyui/common.js"></script>
    <script type="text/javascript" src="/static/easyui/jquery-easyui/easyui_expand.js"></script>
    <script type="text/javascript" src="/static/easyui/jquery-easyui/phpjs-min.js"></script>
    <script type="text/javascript" src="/static/easyui/jquery-easyui/locale/easyui-lang-zh_CN.js"></script>
</head>
<script type="text/javascript">
function loadEasyUIData(URL, param, success, error) {
    $.ajax({
        type: "get",
        data: param,
        url: URL,
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
}
</script>