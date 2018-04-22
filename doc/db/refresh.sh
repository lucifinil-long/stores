#!/bin/bash

set -e

#FIXME test
User=dev
Password=DevP@ssw0rd
Host=127.0.0.1
Port=3306

DBName=stores

if command -v pbcopy >/dev/null 2>&1; then
	echo $Password | pbcopy
	echo "mysql password copied to clipboard!"
fi

if [ "$1" == "db" ]; then
    mysql -u $User -h $Host -P $Port -p $DBName < gw_admin.1.0.0.sql
    # patch
fi

DataSourceName="$User:$Password@tcp($Host:$Port)/$DBName?charset=utf8"
DestDir=../../models/db
TmplDir=./template

mkdir -p $TmplDir
cat > $TmplDir/config << EOF
lang=go
genJson=1
prefix=cos_
EOF

cat > $TmplDir/struct.go.tpl <<EOF
// generated automatically when refresh database models
// don't modify the code manually, changes might be lost

package {{.Model}}

{{\$ilen := len .Imports}}
{{if gt \$ilen 0}}
import (
	{{range .Imports}}"{{.}}"{{end}}
)
{{end}}

{{range .Tables}}
// {{Mapper .Name}} is the database model struct
type {{Mapper .Name}} struct {
{{\$table := .}}
{{range .ColumnsSeq}}{{\$col := \$table.GetColumn .}}	{{Mapper \$col.Name}}	{{Type \$col}} {{Tag \$table \$col}}
{{end}}
}

{{end}}
EOF

xorm reverse mysql $DataSourceName $TmplDir $DestDir "stores_*"

rm -rf $TmplDir
