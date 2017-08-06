package main

import (
	"bytes"
	"fmt"
	//"os"
	"strings"
	"text/template"
)

var golangKeyword = map[string]bool{
	"break":       true,
	"default":     true,
	"func":        true,
	"interface":   true,
	"select":      true,
	"case ":       true,
	"defer":       true,
	"go":          true,
	"map":         true,
	"struct":      true,
	"chan":        true,
	"else":        true,
	"goto":        true,
	"package":     true,
	"switch":      true,
	"const":       true,
	"fallthrough": true,
	"if":          true,
	"range":       true,
	"type":        true,
	"continue":    true,
	"for":         true,
	"import":      true,
	"return":      true,
	"var":         true,
	"nil":         true,
}

var fileReserveWord = map[string]bool{
	"time":      true,
	"golib":     true,
	"error":     true,
	"logging":   true,
	"mysql":     true,
	"gorm":      true,
	"timelib":   true,
	"httplib":   true,
	"db":        true,
	"id":        true,
	"err":       true,
	"int32":     true,
	"count":     true,
	"updateMap": true,
	"ok":        true,
}

var golangSpecialShortUpperWord = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}

func isKeyword(word string) bool {
	return golangKeyword[word] || fileReserveWord[word]
}

func isShortUpperWord(word string) bool {
	return golangSpecialShortUpperWord[word]
}

func isUpperLetterString(dstString string) bool {
	for _, letter := range dstString {
		if !isUpperLetter(letter) {
			return false
		}
	}

	return true
}

func isUpperLetter(letter rune) bool {
	if letter >= 'A' && letter <= 'Z' {
		return true
	} else {
		return false
	}
}

func fetchUpLetter(src string) string {
	var dst string
	for _, character := range src {
		if isUpperLetter(character) {
			dst = fmt.Sprintf("%s%c", dst, character)
		}
	}

	shortLetter := strings.ToLower(dst)
	if isKeyword(shortLetter) {
		shortLetter += "Allen"
	}

	return shortLetter
}

func convertFirstLetterToLower(dst string) string {
	if len(dst) > 0 {
		return strings.ToLower(fmt.Sprintf("%c", dst[0])) + dst[1:]
	} else {
		return dst
	}
}

// HelloGirl -> hello_girl
// niceCar -> nice_car
// TestIDCard->test_id-card
func convertUpLetterToUnderscorePlusLowLetter(src string) string {
	var wordList []string
	var currentWord, tmpWord string
	for i := 0; i < len(src); i++ {
		if isUpperLetter(rune(src[i])) {
			tmpWord = currentWord + fmt.Sprintf("%c", src[i])
			if isUpperLetterString(tmpWord) {
				if isShortUpperWord(tmpWord) {
					wordList = append(wordList, strings.ToLower(tmpWord))
					currentWord = currentWord[len(currentWord):]
				} else {
					currentWord += fmt.Sprintf("%c", src[i])
				}
			} else {
				wordList = append(wordList, strings.ToLower(currentWord))
				currentWord = fmt.Sprintf("%c", src[i])
			}
		} else {
			currentWord += fmt.Sprintf("%c", src[i])
		}
	}

	if len(currentWord) > 0 {
		wordList = append(wordList, strings.ToLower(currentWord))
	}

	return strings.Join(wordList, "_")
}

func genDataAndStore(model *Model, data interface{}, text, tmplName string) error {
	tmpl := template.Must(template.New(tmplName).Parse(text))

	var tmpBuf []byte
	buf := bytes.NewBuffer(tmpBuf)
	if err := tmpl.Execute(buf, data); err != nil {
		return err
	} else {
		if _, ok := model.FuncMapContent[tmplName]; ok {
			// do nothing
		} else {
			model.FuncMapContent[tmplName] = buf.Bytes()
		}
		return nil
	}
}

type BaseInfoOfGenCode struct {
	PartFuncName   string
	FuncInputParam string
	OrmQueryFormat string
	OrmQueryParam  string
	BoolTrue int64
}

type TableNameTemplateParam struct {
	BaseInfoOfGenCode
	StructName              string
	ReceiverName            string
	PackageName             string
	HasCreateTimeField      bool
	HasUpdateTimeField      bool
	HasEnabledField         bool
	EnabledFieldType        string
	CreateTimeFieldType     string
	UpdateTimeFieldType     string
	DbEnabledField          string
	DbCreateTimeField       string
	DbUpdateTimeField       string
	NeedCreateTableNameFunc bool
	TableNameStr            string
}

func genTableNameFunc(model *Model, pkgName string, ignoreCreateTableNameFunc bool) error {
	data := TableNameTemplateParam{
		StructName:              model.Name,
		ReceiverName:            fetchUpLetter(model.Name),
		PackageName:             pkgName,
		HasEnabledField:         model.HasEnabledField,
		HasCreateTimeField:      model.HasCreateTimeField,
		HasUpdateTimeField:      model.HasUpdateTimeField,
		CreateTimeFieldType:     model.CreateTimeFieldType,
		UpdateTimeFieldType:     model.UpdateTimeFieldType,
		NeedCreateTableNameFunc: ignoreCreateTableNameFunc == false,
		TableNameStr:            "t_" + convertUpLetterToUnderscorePlusLowLetter(model.Name),
	}

	return genDataAndStore(model, data, TableNameTemplate, "tableName")
}

var TableNameTemplate = `package {{.PackageName}} 
import (
    "time"

    "github.com/jinzhu/gorm"
) 

type {{.StructName}}List []{{.StructName}}

{{if .NeedCreateTableNameFunc}} 
func ({{.ReceiverName}} {{.StructName}}) TableName() string {
    return "{{.TableNameStr}}"
}{{end}}

func ({{.ReceiverName}} {{.StructName}}) PrintDuration(printParam map[string]interface{}) func() {
    now := time.Now()
    return func() {
        duration := time.Now().Sub(now)
        var printStr string
        for key, value := range printParam {
            printStr += fmt.Sprintf("%s:%s, ", key, value)
        }
        fmt.Printf("%sCost:%sms\n", printStr, fmt.Sprintf("%0.3f", float64(duration/time.Millisecond)))
    }
}
`

var CreateTemplate = `
func ({{.ReceiverName}} *{{.StructName}}) Create(db *gorm.DB) error {
    defer {{.ReceiverName}}.PrintDuration(map[string]interface{}{
            "request" :  "[DB]{{.StructName}}.Create",
        })()

    {{if .HasCreateTimeField}}
    {{if eq .CreateTimeFieldType "int32" "int64" "uint32" "uint64"}}if {{.ReceiverName}}.CreateTime == 0 { 
        {{.ReceiverName}}.CreateTime = time.Now().Unix() 
    }{{else if eq .CreateTimeFieldType "time.Time"}}if {{.ReceiverName}}.CreateTime.IsZero() {
        {{.ReceiverName}}.CreateTime = time.Now()
    }{{else if eq .CreateTimeFieldType "timelib.MySQLTimestamp"}} if time.Time({{.ReceiverName}}.CreateTime).IsZero() {
        {{.ReceiverName}}.CreateTime = timelib.MySQLTimestamp(time.Now()) 
    }{{end}}
    {{end}}
    {{if .HasUpdateTimeField}} 
    {{if eq .UpdateTimeFieldType "int32" "int64" "uint32" "uint64"}}if {{.ReceiverName}}.UpdateTime == 0 { 
        {{.ReceiverName}}.UpdateTime = time.Now().Unix() 
    }{{else if eq .UpdateTimeFieldType "time.Time"}}if {{.ReceiverName}}.UpdateTime.IsZero() {
        {{.ReceiverName}}.UpdateTime = time.Now()
    }{{else if eq .UpdateTimeFieldType "timelib.MySQLTimestamp"}} if time.Time({{.ReceiverName}}.UpdateTime).IsZero() {
        {{.ReceiverName}}.UpdateTime = timelib.MySQLTimestamp(time.Now()) 
    }{{end}}
    {{end}}
    {{if .HasEnabledField}} 
    {{.ReceiverName}}.Enabled = {{.EnabledFieldType}}({{.BoolTrue}}){{end}}
    return db.Table({{.ReceiverName}}.TableName()).Create({{.ReceiverName}}).Error
}
`

func genCreateFunc(model *Model, boolTrue int64) error {
	data := TableNameTemplateParam{
		StructName:          model.Name,
		ReceiverName:        fetchUpLetter(model.Name),
		HasEnabledField:     model.HasEnabledField,
		HasUpdateTimeField:  model.HasUpdateTimeField,
		HasCreateTimeField:  model.HasCreateTimeField,
		EnabledFieldType:    model.EnabledFieldType,
		UpdateTimeFieldType: model.UpdateTimeFieldType,
		CreateTimeFieldType: model.CreateTimeFieldType,
	}
	data.BoolTrue = boolTrue

	return genDataAndStore(model, data, CreateTemplate, "Create")
}

type FetchTemplateParam struct {
	TableNameTemplateParam
	Field            string
	DbField          string
	ReceiverListName string
	StructListName   string
}

func genFetchFuncByNormalIndex(model *Model, baseInfoGenCode *BaseInfoOfGenCode) error {
	data := new(FetchTemplateParam)
	data.StructName = model.Name
	data.ReceiverName = fetchUpLetter(model.Name)
	data.HasEnabledField = model.HasEnabledField
	data.DbEnabledField = model.DbEnabledField
	data.ReceiverListName = fetchUpLetter(model.Name) + "l"
	data.StructListName = model.Name + "List"
	data.PartFuncName = baseInfoGenCode.PartFuncName
	data.FuncInputParam = baseInfoGenCode.FuncInputParam
	data.OrmQueryFormat = baseInfoGenCode.OrmQueryFormat
	data.OrmQueryParam = baseInfoGenCode.OrmQueryParam
	data.BoolTrue = baseInfoGenCode.BoolTrue

	return genDataAndStore(model, data, FetchTemplate, "FetchBy"+data.PartFuncName)
}

var FetchTemplate = `
func ({{.ReceiverListName}} *{{.StructListName}}) FetchBy{{.PartFuncName}}({{.FuncInputParam}}) error {
    defer {{.StructName}}{}.PrintDuration(map[string]interface{}{
            "request" :  "[DB]{{.StructName}}.FetchBy{{.PartFuncName}}",
        })()

    {{if .HasEnabledField}}err := db.Table({{.StructName}}{}.TableName()).Where("{{.OrmQueryFormat}} and {{.DbEnabledField}} = ?", {{.OrmQueryParam}}, {{.BoolTrue}}).Find({{.ReceiverListName}}).Error{{else}} 
    err := db.Table({{.StructName}}{}.TableName()).Where("{{.OrmQueryFormat}}", {{.OrmQueryParam}}).Find({{.ReceiverListName}}).Error{{end}}
    return err
}
`

type BatchFetchTemplateParam struct {
	FetchTemplateParam
	FieldType        string
	ParamField       string
	ReceiverListName string
	StructListName   string
}

func genBatchFetchFunc(model *Model, field, dbField, fieldType string, boolTrue int64) error {
	data := new(BatchFetchTemplateParam)
	data.StructName = model.Name
	data.ReceiverName = fetchUpLetter(model.Name)
	data.Field = field
	data.ParamField = convertFirstLetterToLower(field)
	data.ReceiverListName = fetchUpLetter(model.Name) + "l"
	data.StructListName = model.Name + "List"
	data.DbField = dbField
	data.FieldType = fieldType
	data.HasEnabledField = model.HasEnabledField
	data.DbEnabledField = model.DbEnabledField
	data.BoolTrue = boolTrue

	return genDataAndStore(model, data, BatchFetchTemplate, "BatchFetchBy"+data.Field)
}

var BatchFetchTemplate = `
func ({{.ReceiverListName}} *{{.StructListName}}) BatchFetchBy{{.Field}}List(db *gorm.DB, {{.ParamField}}List []{{.FieldType}}) error {
    defer {{.StructName}}{}.PrintDuration(map[string]interface{}{
            "request" :  "[DB]{{.StructName}}.BatchFetchBy{{.Field}}List",
        })()

    if len({{.ParamField}}List) == 0 {
        return nil
    }

    {{if .HasEnabledField}}err := db.Table({{.StructName}}{}.TableName()).Where("{{.DbField}} in (?) and {{.DbEnabledField}} = ?", {{.ParamField}}List, {{.BoolTrue}}).Find({{.ReceiverListName}}).Error{{else}}
    err := db.Table({{.StructName}}{}.TableName()).Where("{{.DbField}} in (?)", {{.ParamField}}List).Find({{.ReceiverListName}}).Error{{end}}
    return err
}
`

type FetchListTemplateParam struct {
	ReceiverListName   string
	StructListName     string
	StructName         string
	HasEnabledField    bool
	HasCreateTimeField bool
	DbCreateTimeField  string
	DbEnabledField     string
	BoolTrue int64
}

func genFetchListFunc(model *Model, boolTrue int64) error {
	data := new(FetchListTemplateParam)
	data.StructName = model.Name
	data.ReceiverListName = fetchUpLetter(model.Name) + "l"
	data.StructListName = model.Name + "List"
	data.HasEnabledField = model.HasEnabledField
	data.HasCreateTimeField = model.HasCreateTimeField
	data.DbCreateTimeField = model.DbCreateTimeField
	data.DbEnabledField = model.DbEnabledField
	data.BoolTrue = boolTrue

	return genDataAndStore(model, data, FetchListTemplate, "FetchList")
}

var FetchListTemplate = `
func ({{.ReceiverListName}} *{{.StructListName}}) FetchList(db *gorm.DB, size, offset int32, query ...map[string]interface{}) (int32, error) {
    defer {{.StructName}}{}.PrintDuration(map[string]interface{}{
            "request" :  "[DB]{{.StructName}}.FetchList",
        })()

    var count int32{{if .HasEnabledField}} 
    if len(query) == 0 {
        query = append(query, map[string]interface{}{"{{.DbEnabledField}}": {{.BoolTrue}}})
    } else {
        if _, ok := query[0]["{{.DbEnabledField}}"]; !ok { 
            query[0]["{{.DbEnabledField}}"] = {{.BoolTrue}}
        }
    }

    if size <= 0 {
        size = -1
        offset = -1
    }
    var err error

    {{if .HasCreateTimeField}}err = db.Table({{.StructName}}{}.TableName()).Where(query[0]).Count(&count).Limit(size).Offset(offset).Order("{{.DbCreateTimeField}} desc").Find({{.ReceiverListName}}).Error{{else}}err = db.Table({{.StructName}}{}.TableName()).Where(query[0]).Count(&count).Limit(size).Offset(offset).Find({{.ReceiverListName}}).Error{{end}}
    {{else}}
    if size <= 0 {
        size = -1
        offset = -1
    }

    var err error
    if len(query) == 0 {
        {{if .HasCreateTimeField}}err = db.Table({{.StructName}}{}.TableName()).Count(&count).Limit(size).Offset(offset).Order("{{.DbCreateTimeField}} desc").Find({{.ReceiverListName}}).Error{{else}}err = db.Table({{.StructName}}{}.TableName()).Count(&count).Limit(size).Offset(offset).Find({{.ReceiverListName}}).Error{{end}}
    } else {
        {{if .HasCreateTimeField}}err = db.Table({{.StructName}}{}.TableName()).Where(query[0]).Count(&count).Limit(size).Offset(offset).Order("{{.DbCreateTimeField}} desc").Find({{.ReceiverListName}}).Error{{else}}err = db.Table({{.StructName}}{}.TableName()).Where(query[0]).Count(&count).Limit(size).Offset(offset).Find({{.ReceiverListName}}).Error{{end}}
    }{{end}} 
    return int32(count), err
}
`

type UniqueFetchTemplateParam struct {
	TableNameTemplateParam
	UnionField string
}

func genFetchForUpdateFuncByUniqueIndex(model *Model, baseInfoGenCode *BaseInfoOfGenCode) error {
	data := new(UniqueFetchTemplateParam)
	data.StructName = model.Name
	data.ReceiverName = fetchUpLetter(model.Name)
	//data.UnionField = unionField
	data.DbEnabledField = model.DbEnabledField
	data.HasEnabledField = model.HasEnabledField
	data.PartFuncName = baseInfoGenCode.PartFuncName
	data.FuncInputParam = baseInfoGenCode.FuncInputParam
	data.OrmQueryFormat = baseInfoGenCode.OrmQueryFormat
	data.OrmQueryParam = baseInfoGenCode.OrmQueryParam
    data.BoolTrue = baseInfoGenCode.BoolTrue

	return genDataAndStore(model, data, genFetchForUpdateFuncByUniqueIndeTemplate,
		"FetchBy"+data.PartFuncName+"ForUpdate")
}

var genFetchForUpdateFuncByUniqueIndeTemplate = `
func ({{.ReceiverName}} *{{.StructName}}) FetchBy{{.PartFuncName}}ForUpdate({{.FuncInputParam}}) error {
    defer {{.ReceiverName}}.PrintDuration(map[string]interface{}{
            "request" :  "[DB]{{.StructName}}.FetchBy{{.PartFuncName}}ForUpdate",
        })()

    {{if .HasEnabledField}}err := db.Table({{.ReceiverName}}.TableName()).Where("{{.OrmQueryFormat}} and {{.DbEnabledField}} = ?", {{.OrmQueryParam}}, {{.BoolTrue}}).Set("gorm:query_option", "FOR UPDATE").Find({{.ReceiverName}}).Error{{else}} 
    err := db.Table({{.ReceiverName}}.TableName()).Where("{{.OrmQueryFormat}}", {{.OrmQueryParam}}).Set("gorm:query_option", "FOR UPDATE").Find({{.ReceiverName}}).Error{{end}}
    return err
}
`

func genFetchFuncByUniqueIndex(model *Model, baseInfoGenCode *BaseInfoOfGenCode) error {
	data := new(UniqueFetchTemplateParam)
	data.StructName = model.Name
	data.ReceiverName = fetchUpLetter(model.Name)
	data.HasEnabledField = model.HasEnabledField
	//data.UnionField = unionField
	data.DbEnabledField = model.DbEnabledField
	data.PartFuncName = baseInfoGenCode.PartFuncName
	data.FuncInputParam = baseInfoGenCode.FuncInputParam
	data.OrmQueryFormat = baseInfoGenCode.OrmQueryFormat
	data.OrmQueryParam = baseInfoGenCode.OrmQueryParam
	data.BoolTrue = baseInfoGenCode.BoolTrue

	return genDataAndStore(model, data, FetchByUniqueInexTemplate, "FetchBy"+data.PartFuncName)
}

var FetchByUniqueInexTemplate = `
func ({{.ReceiverName}} *{{.StructName}}) FetchBy{{.PartFuncName}}({{.FuncInputParam}}) error {
    defer {{.ReceiverName}}.PrintDuration(map[string]interface{}{
            "request" :  "[DB]{{.StructName}}.FetchBy{{.PartFuncName}}",
        })()

    {{if .HasEnabledField}}err := db.Table({{.ReceiverName}}.TableName()).Where("{{.OrmQueryFormat}} and {{.DbEnabledField}} = ?", {{.OrmQueryParam}}, {{.BoolTrue}}).Find({{.ReceiverName}}).Error{{else}} 
    err := db.Table({{.ReceiverName}}.TableName()).Where("{{.OrmQueryFormat}}", {{.OrmQueryParam}}).Find({{.ReceiverName}}).Error{{end}}
    return err
}
`

type UniqueUpdateWithStructTemplateParam struct {
	TableNameTemplateParam
	UnionField string
}

func genUpdateWithStructFuncByUniqueIndex(model *Model, baseInfoGenCode *BaseInfoOfGenCode) error {
	data := new(UniqueUpdateWithStructTemplateParam)
	data.StructName = model.Name
	data.ReceiverName = fetchUpLetter(model.Name)
	//data.UnionField = uniqueField
	data.HasUpdateTimeField = model.HasUpdateTimeField
	data.UpdateTimeFieldType = model.UpdateTimeFieldType
	data.HasEnabledField = model.HasEnabledField
	data.DbEnabledField = model.DbEnabledField
	data.PartFuncName = baseInfoGenCode.PartFuncName
	data.FuncInputParam = baseInfoGenCode.FuncInputParam
	data.OrmQueryFormat = baseInfoGenCode.OrmQueryFormat
	data.OrmQueryParam = baseInfoGenCode.OrmQueryParam
	data.BoolTrue = baseInfoGenCode.BoolTrue

	return genDataAndStore(model, data, genUpdateWithStructFuncByUniqueIndexTemplate,
		"UpdateBy"+data.PartFuncName+"WithStruct")
}

var genUpdateWithStructFuncByUniqueIndexTemplate = `
func ({{.ReceiverName}} *{{.StructName}}) UpdateBy{{.PartFuncName}}WithStruct({{.FuncInputParam}}) error {
    defer {{.ReceiverName}}.PrintDuration(map[string]interface{}{
            "request" :  "[DB]{{.StructName}}.UpdateBy{{.PartFuncName}}WithStruct",
        })()

    {{if .HasUpdateTimeField}}
    {{if eq .UpdateTimeFieldType "int32" "int64" "uint32" "uint64"}}if {{.ReceiverName}}.UpdateTime == 0 { 
        {{.ReceiverName}}.UpdateTime = time.Now().Unix() 
    }{{else if eq .UpdateTimeFieldType "time.Time"}}if {{.ReceiverName}}.UpdateTime.IsZero() {
        {{.ReceiverName}}.UpdateTime = time.Now()
    }{{else if eq .UpdateTimeFieldType "timelib.MySQLTimestamp"}} if time.Time({{.ReceiverName}}.UpdateTime).IsZero() {
        {{.ReceiverName}}.UpdateTime = timelib.MySQLTimestamp(time.Now()) 
    }{{end}}
    {{end}}
    {{if .HasEnabledField}}dbRet := db.Table({{.ReceiverName}}.TableName()).Where("{{.OrmQueryFormat}} and {{.DbEnabledField}} = ?", {{.OrmQueryParam}}, {{.BoolTrue}}).Updates({{.ReceiverName}}){{else}} 
    dbRet := db.Table({{.ReceiverName}}.TableName()).Where("{{.OrmQueryFormat}}", {{.OrmQueryParam}}).Updates({{.ReceiverName}}){{end}}
    err := dbRet.Error
    if err != nil {
        return err
    } else {
        if dbRet.RowsAffected == 0 {
            return db.Table({{.ReceiverName}}.TableName()).Where("{{.OrmQueryFormat}} and {{.DbEnabledField}} = ?", {{.OrmQueryParam}}, {{.BoolTrue}}).Find(&{{.StructName}}{}).Error
	    } else {
		    return nil
	    }
    }
}
`

type UniqueUpdateWithMapTemplateParam struct {
	TableNameTemplateParam
	UnionField string
}

func genUpdateWithMapFuncByUniqueIndex(model *Model, baseInfoGenCode *BaseInfoOfGenCode) error {
	data := new(UniqueUpdateWithMapTemplateParam)
	data.StructName = model.Name
	data.ReceiverName = fetchUpLetter(model.Name)
	data.HasUpdateTimeField = model.HasUpdateTimeField
	data.UpdateTimeFieldType = model.UpdateTimeFieldType
	data.DbUpdateTimeField = model.DbUpdateTimeField
	data.HasEnabledField = model.HasEnabledField
	data.DbEnabledField = model.DbEnabledField
	data.PartFuncName = baseInfoGenCode.PartFuncName
	data.FuncInputParam = baseInfoGenCode.FuncInputParam
	data.OrmQueryFormat = baseInfoGenCode.OrmQueryFormat
	data.OrmQueryParam = baseInfoGenCode.OrmQueryParam
    data.BoolTrue = baseInfoGenCode.BoolTrue

	return genDataAndStore(model, data, genUpdateWithMapFuncByUniqueIndexTemplate,
		"UpdateBy"+data.PartFuncName+"WithMap")
}

var genUpdateWithMapFuncByUniqueIndexTemplate = `
func ({{.ReceiverName}} *{{.StructName}}) UpdateBy{{.PartFuncName}}WithMap({{.FuncInputParam}}, updateMap map[string]interface{}) error {
    defer {{.ReceiverName}}.PrintDuration(map[string]interface{}{
            "request" :  "[DB]{{.StructName}}.UpdateBy{{.PartFuncName}}WithMap",
        })()

    {{if .HasUpdateTimeField}}if _, ok := updateMap["{{.DbUpdateTimeField}}"]; !ok { 
        {{if eq .UpdateTimeFieldType "int32" "int64" "uint32" "uint64"}}{{.ReceiverName}}.UpdateTime = time.Now().Unix() 
        {{else if eq .UpdateTimeFieldType "time.Time"}}{{.ReceiverName}}.UpdateTime = time.Now()
        {{else if eq .UpdateTimeFieldType "timelib.MySQLTimestamp"}}{{.ReceiverName}}.UpdateTime = timelib.MySQLTimestamp(time.Now()) 
        {{end}}
    }{{end}}
    {{if .HasEnabledField}}dbRet := db.Table({{.ReceiverName}}.TableName()).Where("{{.OrmQueryFormat}} and {{.DbEnabledField}} = ?", {{.OrmQueryParam}}, {{.BoolTrue}}).Updates(updateMap){{else}} 
    dbRet := db.Table({{.ReceiverName}}.TableName()).Where("{{.OrmQueryFormat}}", {{.OrmQueryParam}}).Updates(updateMap){{end}}
    err := dbRet.Error
    if err != nil {
        return err
    } else {
        if dbRet.RowsAffected == 0 {
            return db.Table({{.ReceiverName}}.TableName()).Where("{{.OrmQueryFormat}} and {{.DbEnabledField}} = ?", {{.OrmQueryParam}}, {{.BoolTrue}}).Find(&{{.StructName}}{}).Error
	    } else {
		    return nil
	    }
    }
}
`

type UniqueDeleteTemplateParam struct {
	TableNameTemplateParam
	UnionField string
	BoolFalse int64
}

func genSoftDeleteFuncByUniqueIndex(model *Model, baseInfoGenCode *BaseInfoOfGenCode, boolFalse int64) error {
	data := new(UniqueDeleteTemplateParam)
	data.StructName = model.Name
	data.ReceiverName = fetchUpLetter(model.Name)
	//data.UnionField = unionField
	data.HasUpdateTimeField = model.HasUpdateTimeField
	data.HasEnabledField = model.HasEnabledField
	data.UpdateTimeFieldType = model.UpdateTimeFieldType
	data.DbEnabledField = model.DbEnabledField
	data.PartFuncName = baseInfoGenCode.PartFuncName
	data.FuncInputParam = baseInfoGenCode.FuncInputParam
	data.OrmQueryFormat = baseInfoGenCode.OrmQueryFormat
	data.OrmQueryParam = baseInfoGenCode.OrmQueryParam
	data.BoolTrue = baseInfoGenCode.BoolTrue
	data.BoolFalse = boolFalse

	return genDataAndStore(model, data, genSoftDeleteFuncByUniqueIndexTemplate,
		"SoftDeleteBy"+data.PartFuncName)
}

var genSoftDeleteFuncByUniqueIndexTemplate = `
func ({{.ReceiverName}} *{{.StructName}}) SoftDeleteBy{{.PartFuncName}}({{.FuncInputParam}}) error {
    defer {{.ReceiverName}}.PrintDuration(map[string]interface{}{
            "request" :  "[DB]{{.StructName}}.SoftDeleteBy{{.PartFuncName}}",
        })()

    {{if .HasEnabledField}}var updateMap = map[string]interface{}{}
    updateMap["{{.DbEnabledField}}"] = {{.BoolFalse}}
    {{if .HasUpdateTimeField}}
    {{if eq .UpdateTimeFieldType "int32" "int64" "uint32" "uint64"}}if {{.ReceiverName}}.UpdateTime == 0 { 
        {{.ReceiverName}}.UpdateTime = time.Now().Unix() 
    }{{else if eq .UpdateTimeFieldType "time.Time"}}if {{.ReceiverName}}.UpdateTime.IsZero() {
        {{.ReceiverName}}.UpdateTime = time.Now()
    }{{else if eq .UpdateTimeFieldType "timelib.MySQLTimestamp"}} if time.Time({{.ReceiverName}}.UpdateTime).IsZero() {
        {{.ReceiverName}}.UpdateTime = timelib.MySQLTimestamp(time.Now()) 
    }{{end}}
    {{end}}
    err := db.Table({{.ReceiverName}}.TableName()).Where("{{.OrmQueryFormat}} and {{.DbEnabledField}} = ?", {{.OrmQueryParam}}, {{.BoolTrue}}).Updates(updateMap).Error
    return err{{else}}
    return nil{{end}}
}
`

type UniquePhysicsDeleteTemplateParam struct {
	TableNameTemplateParam
	UnionField string
}

func genPhysicsDeleteFuncByUniqueIndex(model *Model, baseInfoGenCode *BaseInfoOfGenCode) error {
	data := new(UniqueDeleteTemplateParam)
	data.StructName = model.Name
	data.ReceiverName = fetchUpLetter(model.Name)
	data.DbEnabledField = model.DbEnabledField
	data.HasEnabledField = model.HasEnabledField
	data.PartFuncName = baseInfoGenCode.PartFuncName
	data.FuncInputParam = baseInfoGenCode.FuncInputParam
	data.OrmQueryFormat = baseInfoGenCode.OrmQueryFormat
	data.OrmQueryParam = baseInfoGenCode.OrmQueryParam
	data.BoolTrue = baseInfoGenCode.BoolTrue

	return genDataAndStore(model, data, genPhysicsDeleteFuncByUniqueIndexTemplate,
		"DeleteBy"+data.PartFuncName)
}

var genPhysicsDeleteFuncByUniqueIndexTemplate = `
func ({{.ReceiverName}} *{{.StructName}}) DeleteBy{{.PartFuncName}}({{.FuncInputParam}}) error {
    defer {{.ReceiverName}}.PrintDuration(map[string]interface{}{
            "request" :  "[DB]{{.StructName}}.DeleteBy{{.PartFuncName}}",
        })()

    {{if .HasEnabledField}}err := db.Table({{.ReceiverName}}.TableName()).Where("{{.OrmQueryFormat}} and {{.DbEnabledField}} = ?", {{.OrmQueryParam}}, {{.BoolTrue}}).Delete({{.ReceiverName}}).Error{{else}}
    err := db.Table({{.ReceiverName}}.TableName()).Where("{{.OrmQueryFormat}}", {{.OrmQueryParam}}).Delete({{.ReceiverName}}).Error{{end}}
    return err
}
`
