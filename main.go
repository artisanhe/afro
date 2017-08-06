package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

var IgnoreTableName = "@IgnoreTableName"

type pkg struct {
	fset  *token.FileSet
	files map[string]*file
}

type Field struct {
	Name                string
	Type                string
	DbFieldName         string
	IndexName           string
	IndexNumber         int
	IsEnable            bool
	IsSpecifyIndexOrder bool
}

type Model struct {
	Name                string
	IsDbModel           bool
	UniqueIndex         map[string][]Field
	NormalIndex         map[string][]Field
	PrimaryIndex        []Field
	HasCreateTimeField  bool
	HasUpdateTimeField  bool
	HasEnabledField     bool
	EnabledFieldType    string
	CreateTimeFieldType string
	UpdateTimeFieldType string
	DbCreateTimeField   string
	DbEnabledField      string
	DbUpdateTimeField   string
	FuncMapContent      map[string][]byte
}

type file struct {
	pkg            *pkg
	f              *ast.File
	src            []byte
	pkgPath        string
	ModelList      []Model
	TypeMapComment map[string]string
}

func ParseTagSetting(str string) map[string][]string {
	tags := strings.Split(str, ";")
	setting := map[string][]string{}
	for _, value := range tags {
		v := strings.Split(value, ":")
		k := strings.TrimSpace(strings.ToUpper(v[0]))
		if _, ok := setting[k]; !ok {
			setting[k] = make([]string, 0, 10)
		}
		if len(v) == 2 {
			setting[k] = append(setting[k], v[1])
		} else {
			setting[k] = append(setting[k], k)
		}
	}
	return setting
}

// ParseIndex indexName[0] -> return indexname, 0
func ParseIndex(index string) (string, int, bool) {
	var isSpecifyIndexOrder bool
	if len(index) == 0 {
		return "", -1, isSpecifyIndexOrder
	}
	tmpStrSlice := strings.Split(index, "[")
	if len(tmpStrSlice) != 2 {
		return tmpStrSlice[0], 0, isSpecifyIndexOrder
	}

	pos, err := strconv.ParseInt(tmpStrSlice[1][0:1], 10, 64)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	isSpecifyIndexOrder = true
	return tmpStrSlice[0], int(pos), isSpecifyIndexOrder
}

// fetchBaseInfoOfGenFuncForNormalIndex fetch part of function name, function input param,
// orm query format, orm query parameters.
func fetchBaseInfoOfGenFuncForNormalIndex(indexList []Field) *BaseInfoOfGenCode {
	var partFuncName, inputParam, ormQueryFormat, ormQueryParam string
	for _, field := range indexList {
		if len(partFuncName) == 0 {
			partFuncName = field.Name
			inputParam = fmt.Sprintf("db *gorm.DB, %s %s", convertFirstLetterToLower(field.Name), field.Type)
			ormQueryFormat = fmt.Sprintf("%s = ?", field.DbFieldName)
			ormQueryParam = fmt.Sprintf("%s", convertFirstLetterToLower(field.Name))
		} else {
			partFuncName += "And" + field.Name
			inputParam += fmt.Sprintf(", %s %s", convertFirstLetterToLower(field.Name), field.Type)
			ormQueryFormat += fmt.Sprintf(" and %s = ?", field.DbFieldName)
			ormQueryParam += fmt.Sprintf(", %s", convertFirstLetterToLower(field.Name))
		}
	}

	return &BaseInfoOfGenCode{
		PartFuncName:   partFuncName,
		FuncInputParam: inputParam,
		OrmQueryFormat: ormQueryFormat,
		OrmQueryParam:  ormQueryParam,
        BoolTrue : gTrue,
	}
}

func fetchBaseInfoOfGenFuncForUniqueIndex(model *Model, indexList []Field) *BaseInfoOfGenCode {
	var partFuncName, inputParam, ormQueryFormat, ormQueryParam string
	inputParam = fmt.Sprintf("%s", "db *gorm.DB")
	for _, field := range indexList {
		if len(partFuncName) == 0 {
			partFuncName = field.Name
			ormQueryFormat = fmt.Sprintf("%s = ?", field.DbFieldName)
			ormQueryParam = fmt.Sprintf("%s.%s", fetchUpLetter(model.Name), field.Name)
		} else {
			partFuncName += "And" + field.Name
			ormQueryFormat += fmt.Sprintf(" and %s = ?", field.DbFieldName)
			ormQueryParam += fmt.Sprintf(", %s.%s", fetchUpLetter(model.Name), field.Name)
		}
	}

	return &BaseInfoOfGenCode{
		PartFuncName:   partFuncName,
		FuncInputParam: inputParam,
		OrmQueryFormat: ormQueryFormat,
		OrmQueryParam:  ormQueryParam,
        BoolTrue : gTrue,
	}
}

func replaceUpperWithLowerAndUnderscore(src string) string {
	var dst string
	for index, letter := range src {
		if index == 0 && isUpperLetter(letter) {
			dst += fmt.Sprintf("%c", letter)
		} else if isUpperLetter(letter) {
			dst += fmt.Sprintf("_%c", letter)
		} else {
			dst += fmt.Sprintf("%c", letter)
		}
	}

	return strings.ToLower(dst)
}

type walker func(ast.Node) bool

func (w walker) Visit(node ast.Node) ast.Visitor {
	if w(node) {
		return w
	}
	return nil
}

func (f *file) walk(fn func(ast.Node) bool) {
	ast.Walk(walker(fn), f.f)
}

func (f *file) IsSpecifyIndexSequence(fields []Field) bool {
	var defaultIndexOrder, specifyIndexOrder bool
	for _, field := range fields {
		if field.IsSpecifyIndexOrder {
			specifyIndexOrder = true
		} else {
			defaultIndexOrder = true
		}
	}

	if defaultIndexOrder && specifyIndexOrder {
		fmt.Printf("Some fields are specified index order but other not in same index[%s].\n", fields[0].IndexName)
		os.Exit(1)
	}

	return specifyIndexOrder
}

func (f *file) sortFieldsByIndexNumber(fields []Field, model *Model) []Field {
	if len(fields) == 0 && len(fields) == 1 {
		return fields
	}

	if !f.IsSpecifyIndexSequence(fields) {
		return fields
	}

	var markIndexNubmer = make(map[string]string)
	var sortFieldSlice = []Field{}
	for _, field := range fields {
		if fieldName, ok := markIndexNubmer[fmt.Sprintf("%d", field.IndexNumber)]; ok {
			fmt.Println("Field[%s] and Field[%s] has same index number[%d] in Mode[%s]", fieldName, field.Name,
				field.IndexNumber, model.Name)
			os.Exit(1)
		}
		if field.IndexNumber >= 0 {
			tmpSlice := make([]Field, field.IndexNumber+1)
			tmpSlice[field.IndexNumber] = field
			if len(sortFieldSlice) > field.IndexNumber+1 {
				sortFieldSlice = append(tmpSlice, sortFieldSlice[field.IndexNumber+1:]...)
			} else if len(sortFieldSlice) < field.IndexNumber+1 {
				sortFieldSlice = append(sortFieldSlice, tmpSlice[len(sortFieldSlice):]...)
			} else {
				fmt.Printf("Field[%s] wrong index sequence, may be same index number.\n", field.Name)
				os.Exit(1)
			}
		}
	}

	var notEmptySlice = []Field{}
	for index, field := range sortFieldSlice {
		if len(field.Name) > 0 {
			notEmptySlice = append(notEmptySlice, sortFieldSlice[index])
		}
	}

	return notEmptySlice
}

func (f *file) genBatchFetchFuncBySingleIndex(model *Model, field Field) {
	if err := genBatchFetchFunc(model, field.Name, field.DbFieldName, field.Type, gTrue); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
}

func (f *file) handleGenFetchCodeBySubIndex(model *Model, fieldList []Field) {
	if len(fieldList) == 1 {
		f.genBatchFetchFuncBySingleIndex(model, fieldList[0])
	} else if len(fieldList) > 1 {
		// [x, y, z, e] Split to [x, y, z], [x, y], [x]
		for i := 1; i < len(fieldList); i++ {
			subSortFieldSlice := fieldList[:len(fieldList)-i]
			baseInfoGenCode := fetchBaseInfoOfGenFuncForNormalIndex(subSortFieldSlice)
			if err := genFetchFuncByNormalIndex(model, baseInfoGenCode); err != nil {
				fmt.Printf("%s\n", err.Error())
				os.Exit(1)
			}

			if len(subSortFieldSlice) == 1 {
				f.genBatchFetchFuncBySingleIndex(model, subSortFieldSlice[0])
			}

		}
	}
}

func (f *file) handleGenCodeForUniqueIndex(model *Model, sortFieldList []Field) {
	baseInfoGenCode := fetchBaseInfoOfGenFuncForUniqueIndex(model, sortFieldList)
	if err := genFetchFuncByUniqueIndex(model, baseInfoGenCode); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	f.handleGenFetchCodeBySubIndex(model, sortFieldList)

	if err := genFetchForUpdateFuncByUniqueIndex(model, baseInfoGenCode); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	if err := genUpdateWithStructFuncByUniqueIndex(model, baseInfoGenCode); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	if err := genUpdateWithMapFuncByUniqueIndex(model, baseInfoGenCode); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	if err := genSoftDeleteFuncByUniqueIndex(model, baseInfoGenCode, gFalse); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	if err := genPhysicsDeleteFuncByUniqueIndex(model, baseInfoGenCode); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
}

func (f *file) genCodeByUniqueIndex(model *Model) {
	for _, fieldList := range model.UniqueIndex {
		sortFieldList := f.sortFieldsByIndexNumber(fieldList, model)
		f.handleGenCodeForUniqueIndex(model, sortFieldList)
	}
}

func (f *file) genCodeByPrimaryKeyIndex(model *Model) {
	if len(model.PrimaryIndex) > 0 {
		f.handleGenCodeForUniqueIndex(model, model.PrimaryIndex)
	}
}

func (f *file) genCodeByNormalIndex(model *Model) {
	for _, fieldList := range model.NormalIndex {
		sortFieldList := f.sortFieldsByIndexNumber(fieldList, model)
		baseInfoGenCode := fetchBaseInfoOfGenFuncForNormalIndex(sortFieldList)
		if err := genFetchFuncByNormalIndex(model, baseInfoGenCode); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}

		f.handleGenFetchCodeBySubIndex(model, sortFieldList)

	}
}

func (f *file) outputFile(model *Model) {
	filename := path.Join(f.pkgPath, "generate_"+replaceUpperWithLowerAndUnderscore(model.Name)) + "_interface.go"
	fp, err := os.Create(filename)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	defer fp.Close()

	var funcNameList []string

	// first part of file
	var tableName = "tableName"
	for key := range model.FuncMapContent {
		if key == tableName {
			continue
		}
		funcNameList = append(funcNameList, key)
	}
	sort.Strings(funcNameList)

	_, err = fp.Write(model.FuncMapContent[tableName])
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	for _, funcName := range funcNameList {
		_, err := fp.Write(model.FuncMapContent[funcName])
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}
	}

	fp.Sync()
	exec.Command("gofmt", "-w", filename).CombinedOutput()
    exec.Command("goimports", "-w", filename).CombinedOutput()
}

func (f *file) output(model *Model) error {
	var ignoreCreateTableNameFunc bool
	if doc, ok := f.TypeMapComment[model.Name]; ok && IgnoreTableName == strings.TrimSpace(doc) {
		ignoreCreateTableNameFunc = true
	}

	if err := genTableNameFunc(model, f.f.Name.Name, ignoreCreateTableNameFunc); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	if err := genCreateFunc(model, gTrue); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	f.genCodeByNormalIndex(model)
	f.genCodeByUniqueIndex(model)
	f.genCodeByPrimaryKeyIndex(model)

	if err := genFetchListFunc(model, gTrue); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	f.outputFile(model)
	return nil
}

func (f *file) generateFile() error {
	for _, module := range f.ModelList {
		if !module.IsDbModel {
			continue
		}

		f.output(&module)
	}
	return nil
}

func (f *file) paraseStructField(structType *ast.StructType, structModel *Model) {
	for _, field := range structType.Fields.List {
		if field.Tag == nil {
			continue
		}
		tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
		sqlSettings := ParseTagSetting(tag.Get("sql"))
		gormSettings := ParseTagSetting(tag.Get("gorm"))
		if len(gormSettings) != 0 && len(sqlSettings) != 0 {
			structModel.IsDbModel = true
		} else {
			continue
		}

		var tmpField = Field{}
		if dbFieldName, ok := gormSettings["COLUMN"]; ok {
			tmpField.DbFieldName = dbFieldName[0]
		} else {
			tmpField.DbFieldName = tmpField.Name
		}
		tmpField.Name = field.Names[0].Name
		tmpField.Type = string(f.src[field.Type.Pos()-1 : field.Type.End()-1])

		if tmpField.Name == "Enabled" || tmpField.DbFieldName == "F_enabled" {
			// enabled field don't join into index slice
			tmpField.IsEnable = true
			structModel.HasEnabledField = true
			structModel.EnabledFieldType = tmpField.Type
			structModel.DbEnabledField = tmpField.DbFieldName
			continue
		} else if tmpField.Name == "UpdateTime" || tmpField.DbFieldName == "F_update_time" {
			structModel.HasUpdateTimeField = true
			structModel.UpdateTimeFieldType = tmpField.Type
			structModel.DbUpdateTimeField = tmpField.DbFieldName
		} else if tmpField.Name == "CreateTime" || tmpField.DbFieldName == "F_create_time" {
			structModel.HasCreateTimeField = true
			structModel.CreateTimeFieldType = tmpField.Type
			structModel.DbCreateTimeField = tmpField.DbFieldName
		}

		if _, ok := gormSettings["PRIMARY_KEY"]; ok {
			structModel.PrimaryIndex = append(structModel.PrimaryIndex, tmpField)
		}
		if indexName, ok := sqlSettings["INDEX"]; ok && len(indexName) > 0 {
			tmpField.IndexName, tmpField.IndexNumber, tmpField.IsSpecifyIndexOrder = ParseIndex(indexName[0])
			if _, ok := structModel.NormalIndex[tmpField.IndexName]; ok {
				structModel.NormalIndex[tmpField.IndexName] = append(
					structModel.NormalIndex[tmpField.IndexName],
					tmpField)
			} else {
				structModel.NormalIndex[tmpField.IndexName] = []Field{tmpField}
			}
		}
		if indexName, ok := sqlSettings["UNIQUE_INDEX"]; ok && len(indexName) > 0 {
			tmpField.IndexName, tmpField.IndexNumber, tmpField.IsSpecifyIndexOrder = ParseIndex(indexName[0])
			if _, ok := structModel.UniqueIndex[tmpField.IndexName]; ok {
				structModel.UniqueIndex[tmpField.IndexName] = append(
					structModel.UniqueIndex[tmpField.IndexName],
					tmpField)
			} else {
				structModel.UniqueIndex[tmpField.IndexName] = []Field{tmpField}
			}
		}
	}
}

func (f *file) fetchAndParseStruct() {
	walkFunc := func(node ast.Node) bool {
		v, ok := node.(*ast.GenDecl)
		if !ok {
			return true
		}

		if v.Tok == token.IMPORT || v.Tok == token.CONST || v.Tok == token.VAR {
			return true
		}

		for _, spec := range v.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				return true
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				return true
			}

			var structModel = Model{
				Name:        typeSpec.Name.Name,
				UniqueIndex: make(map[string][]Field),
				NormalIndex: make(map[string][]Field),
				//UnionUniqueSubFuncName: make(map[string]string),
				FuncMapContent: make(map[string][]byte),
			}
			f.paraseStructField(structType, &structModel)

			f.ModelList = append(f.ModelList, structModel)
			return true
		}
		return true
	}
	f.walk(walkFunc)
}

func isGoFile(filename string) bool {
	if strings.HasSuffix(filename, ".go") {
		return true
	} else {
		return false
	}
}

func isGoTestFile(filename string) bool {
	if strings.HasSuffix(filename, "_test.go") {
		return true
	} else {
		return false
	}
}

func isGeneratedFile(filename string) bool {
	if strings.HasPrefix(filename, "generate_") && strings.HasSuffix(filename, "_interface.go") {
		return true
	} else {
		return false
	}
}

func handleFile(filename, pathName string) {
	if !isGoFile(filename) {
		return
	}

	if isGoTestFile(filename) {
		return
	}

	if isGeneratedFile(filename) {
		return
	}

	if len(pathName) == 0 {
		pathName, filename = filepath.Split(filename)
	}

	src, err := ioutil.ReadFile(pathName + filename)
	if err != nil {
		fmt.Sprintf("Failed reading %s: %v\n", filename, err)
		return
	}

	pkg := &pkg{
		fset:  token.NewFileSet(),
		files: make(map[string]*file),
	}

	f, err := parser.ParseFile(pkg.fset, filename, src, parser.ParseComments)
	if err != nil {
		fmt.Printf("ParseFile:%s\n", err.Error())
		return
	}

	pkg.files[filename] = &file{
		pkg:            pkg,
		f:              f,
		src:            src,
		pkgPath:        pathName,
		TypeMapComment: parseComments(f, filename, pathName),
	}

	pkg.files[filename].fetchAndParseStruct()
	pkg.files[filename].generateFile()
}

// fetch the package path under $GOPATH/src || $GOROOT/src
//func fetchPkgPath(filename string) string {
	//absPath, err := filepath.Abs(filename)
	//if err != nil {
		//fmt.Printf("Fetch absolute path of pkg error:%s\n", err.Error())
		//return ""
	//}

	//itemSlice := strings.Split(absPath, "/")

	//// filter the last two element.
	//if len(itemSlice) < 2 {
		//return ""
	//}

	//itemSlice = itemSlice[0 : len(itemSlice)-2]
	//for index, item := range itemSlice {
		//if item == "src" && index < len(itemSlice)-1 {
			//return strings.Join(itemSlice[index+1:], "/")
		//}
	//}
	//return ""
//}

func parseComments(src *ast.File, filename, packagePath string) map[string]string {
	pkg := &ast.Package{
		Name: src.Name.Name,
		Files: map[string]*ast.File{
			src.Name.Name: src,
		},
	}

	var commentMapType = make(map[string]string)
	p := doc.New(pkg, packagePath, 1)
	for _, t := range p.Types {
		if _, ok := commentMapType[t.Name]; ok {
			fmt.Printf("Has same type[%s] in file[%s].", t.Name, filename)
		} else {
			commentMapType[t.Name] = t.Doc
		}
	}

	return commentMapType
}

func parseDir(pathName string) {
	fis, err := ioutil.ReadDir(pathName)
	if err != nil {
		fmt.Printf("ioutil.ReadDir: %s\n", err.Error())
		os.Exit(1)
	}

	for _, fi := range fis {
		handleFile(fi.Name(), pathName)
	}
}


var (
    Version string
    Build   string
    gPathName string
    gFileName string
    gConfigFile string
    gVersion bool
    gBuild bool
)

var (
    True string
    False string
    gTrue int64
    gFalse int64
)

func InitFlag() {
	flag.StringVar(&gPathName, "d", "", "-d direcotry name")
	flag.StringVar(&gFileName, "f", "", "-f file name")
	flag.BoolVar(&gVersion, "version", false, "prints current version")
	flag.BoolVar(&gBuild, "build", false, "prints build git version")
	flag.Parse()

	if gVersion {
        fmt.Println("Version: ", Version)
        os.Exit(0)
	}
	if gBuild {
        fmt.Println("Git commit hash: ", Build)
        os.Exit(0)
	}

    var err error
	if len(True) > 0 {
	    gTrue, err  = strconv.ParseInt(True, 10, 64)
	    if err != nil {
	        panic(err)
	    }
	} else {
	    gTrue = 1
	}

	if len(False) > 0 {
	    gFalse, err = strconv.ParseInt(False, 10, 64)
	    if err != nil {
	        panic(err)
	    }
	} else {
	    gFalse = 2
	}
}

func main() {
	InitFlag()
	if len(gPathName) != 0 {
		parseDir(gPathName)
	} else if len(gFileName) != 0 {
		handleFile(gFileName, "")
	} else {
		fmt.Printf("must specify generated directory or file.\n")
	}

	return
}
