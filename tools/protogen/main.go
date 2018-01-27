package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/mkideal/cli"
	"github.com/mkideal/log"
)

type argT struct {
	cli.Helper
	Inputs    string `cli:"i" usage:"input files, seperate by comma ','" dft:"proto.go"`
	LogLevel  string `cli:"v" usage:"log level: t,d,i,w,e,f" dft:"i"`
	ReadmeOut string `cli:"readme-file" usage:"readme out put file" dft:"./README.md"`
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		log.SetLevelFromString(argv.LogLevel)

		// parse proto first
		err := parseProto(ctx, argv)
		if err != nil {
			return err
		}

		return nil
	})
}

func parseProto(ctx *cli.Context, argv *argT) error {
	files := strings.Split(argv.Inputs, ",")
	if len(files) == 0 {
		log.Fatal("input files is empty")
	}
	structs := make([]Struct, 0)
	enums := make([]Enum, 0)
	for _, inputFile := range files {
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, inputFile, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				log.Fatal("invalid decl: %T", decl)
				continue
			}
			enum := tryParseEnumFlag(genDecl)
			for _, ispec := range genDecl.Specs {
				switch spec := ispec.(type) {
				case *ast.TypeSpec:
					structType, ok := spec.Type.(*ast.StructType)
					if !ok { //不处理即可 直接忽略 以免强类型枚举出错
						_, ok := spec.Type.(*ast.ArrayType)
						if !ok {
							log.Warn("unsupported expr of TypeSpec: %T,ispec:%v,type:%v", spec.Type, ispec, spec.Type)
							continue
						}

						typ := parseType("", spec.Type)
						obj := Struct{
							Type: &typ,
							Name: spec.Name.Obj.Name,
						}
						structs = append(structs, obj)
						continue
					}
					obj := Struct{
						Name:     spec.Name.Obj.Name,
						Fields:   make([]*Field, 0, len(structType.Fields.List)),
						Comments: parseCommentGroup(spec.Doc),
					}
					if len(genDecl.Specs) == 1 && len(obj.Comments) == 0 {
						obj.Comments = append(obj.Comments, parseCommentGroup(genDecl.Doc)...)
					}
					for _, field := range structType.Fields.List {
						var (
							name = ""
							typ  Type
						)
						if len(field.Names) == 0 {
							typ = parseType("", field.Type)
							obj.Parent = typ.Name
							continue
						} else if len(field.Names) == 1 {
							name = field.Names[0].Name
							typ = parseType(name, field.Type)
						} else {
							log.Fatal("unsupported field.Names.length=%d", len(field.Names))
						}
						if field.Tag != nil {
							tag := field.Tag.Value
							if len(tag) > 0 && tag[0] == '`' {
								tag = tag[1:]
							}
							if len(tag) > 0 && tag[len(tag)-1] == '`' {
								tag = tag[:len(tag)-1]
							}
							jsonTag := reflect.StructTag(tag).Get("json")
							if len(jsonTag) > 0 {
								tmpStrs := strings.SplitN(jsonTag, ",", 2)
								if len(tmpStrs) > 0 {
									jsonTag = tmpStrs[0]
								}
							}
							obj.Fields = append(obj.Fields, &Field{
								Name:     name,
								Type:     typ,
								JSONKey:  jsonTag,
								Comments: parseCommentGroup(field.Comment),
							})
						}
					}
					structs = append(structs, obj)
				case *ast.ValueSpec:
					if enum != nil {
						if len(spec.Names) != 1 {
							log.Fatal("spec.Names.length == %d, want 1", len(spec.Names))
						}
						if len(spec.Values) != 1 {
							log.Fatal("spec.Values.length == %d, want 1", len(spec.Values))
						}
						enumField := EnumField{Name: spec.Names[0].Name}
						if ident, ok := spec.Values[0].(*ast.Ident); ok {
							enumField.Value = ident.Name
						} else if c, ok := spec.Values[0].(*ast.BasicLit); ok {
							enumField.Value = c.Value
						} else {
							log.Fatal("unsupported enum field value type: %T", spec.Values[0])
						}
						enumField.Comment = strings.Join(parseCommentGroup(spec.Comment), "\n")
						enum.Fields = append(enum.Fields, enumField)
					}
				default:
					log.Fatal("unsupported spec: %T", ispec)
				}
			}
			if enum != nil {
				enums = append(enums, *enum)
			}
		}
	}
	if log.GetLevel().MoreVerboseThan(log.LvINFO) {
		ctx.JSONIndentln(structs, "", "    ")
	}

	readmeGen(argv, structs, enums)

	return nil
}

func tryParseEnumFlag(genDecl *ast.GenDecl) *Enum {
	if genDecl.Doc == nil {
		return nil
	}
	comments := parseCommentGroup(genDecl.Doc)
	if len(comments) == 0 {
		return nil
	}
	line0 := strings.TrimPrefix(strings.TrimPrefix(comments[0], "//"), " ")
	if strings.HasPrefix(line0, "enum:") {
		name := strings.TrimSpace(strings.TrimPrefix(line0, "enum:"))
		return &Enum{Name: name, Fields: make([]EnumField, 0), Comments: comments}
	}
	return nil
}

type Type struct {
	Kind      string
	Name      string
	IndexType *Type
	ElemType  *Type
}

type EnumField struct {
	Name    string
	Value   string
	Comment string
}

type Enum struct {
	Name     string
	Fields   []EnumField
	Comments []string
}

type Field struct {
	Name     string
	Type     Type
	JSONKey  string
	Comments []string
}

type Struct struct {
	Type     *Type
	Name     string
	Fields   []*Field
	Comments []string
	Parent   string
}

type Protocol struct {
	Struct
}

func parseCommentGroup(cg *ast.CommentGroup) []string {
	if cg == nil || len(cg.List) == 0 {
		return []string{}
	}
	comments := make([]string, 0, len(cg.List))
	for _, comment := range cg.List {
		comments = append(comments, comment.Text)
	}
	return comments
}

func parseType(name string, typ ast.Expr) Type {
	t := Type{}
	switch fieldType := typ.(type) {
	case *ast.Ident:
		obj := fieldType.Obj
		if obj != nil {
			t.Kind = "struct"
			t.Name = obj.Name
		} else {
			t.Kind = fieldType.Name
		}
	case *ast.ArrayType:
		elem := fieldType.Elt
		t.Kind = "slice"
		elemType := parseType(name, elem)
		t.ElemType = &elemType
	case *ast.InterfaceType:
		t.Kind = "any"
	case *ast.StarExpr:
		t.Kind = "empty"
	default:
		log.Fatal("unsupported field type of %s: %T", name, typ)
	}
	return t
}

func formatEnumFields(fields []EnumField) string {
	buf := new(bytes.Buffer)
	nameMaxLen, valuMaxLen := 0, 0
	for _, field := range fields {
		if len(field.Name) > nameMaxLen {
			nameMaxLen = len(field.Name)
		}
		if len(field.Value) > valuMaxLen {
			valuMaxLen = len(field.Value)
		}
	}
	format := fmt.Sprintf("\t%%-%ds = %%%ds, %%s\n", nameMaxLen, valuMaxLen)
	for _, field := range fields {
		fmt.Fprintf(buf, format, field.Name, field.Value, field.Comment)
	}
	return buf.String()
}

func readmeGen(argv *argT, structs []Struct, enums []Enum) {
	enumbuf := new(bytes.Buffer)
	reqbuf := new(bytes.Buffer)
	pagebuf := new(bytes.Buffer)
	for _, enum := range enums {
		comments := ""
		if len(enum.Comments) > 1 {
			comments = strings.Join(enum.Comments[1:], "\n") + "\n"
		}
		fmt.Fprintf(enumbuf, "```c\n%senum %s {\n%s};\n```\n\n", comments, enum.Name, formatEnumFields(enum.Fields))
	}
	for _, obj := range structs {
		if strings.HasSuffix(obj.Name, "Page") && len(obj.Comments) > 0 {
			title, comment := formatStructComment(&obj, true)
			fmt.Fprintf(pagebuf, "### %s\n %s\n\n```js\n访问指定路径可以获取指定html页面\n```\n\n", title, comment)
			continue
		}

		if !strings.HasSuffix(obj.Name, "Req") || len(obj.Comments) == 0 {
			continue
		}

		title, comment := formatStructComment(&obj, false)
		fmt.Fprintf(reqbuf, "### %s\n %s\n\n", title, comment)
		cmdName := strings.TrimSuffix(obj.Name, "Req")
		res := cmdName + "Res"
		var resObj *Struct
		for _, o := range structs {
			if o.Name == res {
				resObj = &o
				break
			}
		}
		if resObj != nil {
			fmt.Fprintf(reqbuf, "```js\n// 请求表单参数示例\n%s\n\n// 返回JSON示例(仅Response.Protocol部分)\n%s\n```\n\n", obj.FormRequest(structs, "\t"), resObj.JSON(structs, "\t"))
		} else {
			fmt.Fprintf(reqbuf, "```js\n// 请求表单参数示例\n%s\n\n// 返回JSON示例(仅Response.Protocol部分)\n{\n}\n```\n\n", obj.FormRequest(structs, "\t"))
		}
	}

	commonResponse := ""
	for _, obj := range structs {
		if obj.Name == "Response" {
			commonResponse = "```js\n" + obj.JSON(structs, "\t") + "\n```"
			break
		}
	}

	content := fmt.Sprintf(readmeTpl, enumbuf.String(), commonResponse, reqbuf.String(), pagebuf.String())
	dir, _ := filepath.Split(argv.ReadmeOut)
	os.MkdirAll(dir, 0755)
	ioutil.WriteFile(argv.ReadmeOut, []byte(content), 0666)
}

func formatStructComment(obj *Struct, isPage bool) (string, string) {
	path := ""
	title := ""
	method := ""
	for _, desc := range obj.Comments {
		desc = strings.TrimSpace(desc)
		desc = strings.TrimPrefix(desc, "//")
		desc = strings.TrimSpace(desc)
		if strings.HasPrefix(desc, obj.Name) {
			title = strings.TrimPrefix(desc, obj.Name)
			title = strings.TrimSpace(title)
		} else if strings.HasPrefix(desc, "path:") {
			path = strings.TrimPrefix(desc, "path:")
			path = strings.TrimSpace(path)
		} else if strings.HasPrefix(desc, "method:") {
			method = strings.TrimPrefix(desc, "method:")
			method = strings.TrimSpace(method)
		}
	}

	if isPage {
		title = fmt.Sprintf("访问%v可以%v", path, title)
	} else {
		title = fmt.Sprintf("API接口%v %v", path, title)
	}

	comment := "访问方法: " + method
	return title, comment
}

func findStructByName(structs []Struct, name string) (Struct, bool) {
	for _, obj := range structs {
		if obj.Name == name {
			return obj, true
		}
	}
	return Struct{}, false
}

func (obj Struct) JSON(structs []Struct, prefix string) string {
	if obj.Type != nil {
		return obj.Type.JSON(structs, prefix)
	}

	res := "{\n"
	var parents []Struct
	parent := obj
	for parent.Parent != "" {
		p, ok := findStructByName(structs, parent.Parent)
		if !ok {
			log.Fatal("parent %s not found", parent.Parent)
		}
		parent = p
		parents = append([]Struct{parent}, parents...)
	}

	hasField := false

	for _, parent := range parents {
		if !hasField && len(parent.Fields) > 0 {
			hasField = true
		}
		for _, field := range parent.Fields {
			res += fmt.Sprintf(`%s"%s": %s`, prefix, field.JSONKey, field.Type.JSON(structs, prefix+"\t"))
			res += ","
			if field.Type.Kind != "struct" {
				if len(field.Comments) > 0 {
					res += " " + field.Comments[0]
				}
			}
			res += "\n"
		}
	}
	for i, field := range obj.Fields {
		res += fmt.Sprintf(`%s"%s": %s`, prefix, field.JSONKey, field.Type.JSON(structs, prefix+"\t"))
		if i+1 != len(obj.Fields) {
			res += ","
		}
		if field.Type.Kind != "struct" {
			if len(field.Comments) > 0 {
				res += " " + field.Comments[0]
			}
		}
		res += "\n"
	}
	if strings.HasSuffix(prefix, "\t") {
		prefix = strings.TrimSuffix(prefix, "\t")
	}
	res += prefix + "}"

	if !hasField && len(obj.Fields) == 0 {
		return "本接口不返回json格式数据"
	}

	return res
}

func (obj Struct) FormRequest(structs []Struct, prefix string) string {
	params := []string{}
	var parents []Struct
	parent := obj
	for parent.Parent != "" {
		p, ok := findStructByName(structs, parent.Parent)
		if !ok {
			log.Fatal("parent %s not found", parent.Parent)
		}
		parent = p
		parents = append([]Struct{parent}, parents...)
	}
	for _, parent := range parents {
		for _, field := range parent.Fields {
			params = append(params, fmt.Sprintf(`%s=xxx`, field.JSONKey))
		}
	}
	for _, field := range obj.Fields {
		params = append(params, fmt.Sprintf(`%s=xxx`, field.JSONKey))
	}

	res := strings.Join(params, "&")
	if len(res) > 0 {
		res += "\n"
		res += "// 请求表单参数内容格式说明(参数值为json对象的，值对应为json对象的字符串内容)\n"
		res += obj.JSON(structs, prefix)
	} else {
		res = "无需参数"
	}

	return res
}

func (typ Type) JSON(structs []Struct, prefix string) string {
	switch typ.Kind {
	case "string":
		return `"value"`
	case "bool":
		return "false"
	case "slice":
		if typ.ElemType.Kind != "slice" && typ.ElemType.Kind != "struct" {
			ret := "[" + typ.ElemType.JSON(structs, prefix+"\t") + "]"
			return ret
		}
		return "[\n\t" + typ.ElemType.JSON(structs, prefix+"\t") + "\n" + strings.TrimSuffix(prefix, "\t") + "]"
	case "struct":
		for _, obj := range structs {
			if obj.Name == typ.Name {
				return obj.JSON(structs, prefix)
			}
		}
		return "{}"
	case "any":
		return "{}"
	case "empty":
		return ""
	default:
		return "0"
	}
}

const readmeTpl = `
# 仓储管理服务 接口文档

## 所有用到的枚举值

%s

## HTTP API接口 请求/响应参数
访问方式为请求路径加请求表单方式传递参数
每个请求参数结构的第一层每个字段对应一个表单参数
请求或响应参数为空的，表示无需请求数据或没有响应数据

API请求接口使用统一的返回JSON结构：

%s

%s

## 页面请求

%s
`
