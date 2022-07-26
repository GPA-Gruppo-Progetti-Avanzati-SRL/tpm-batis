package sqlmapper

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/sqlmapper/xml"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/system/util"
	"github.com/rs/zerolog/log"

	"reflect"
	"strconv"
	"strings"
	"text/template"
)

type ResultMap interface {
	MapColumnNames2TypeProperties(columns []string) ([]string, error)
}

type resultmap struct {
}

type Statement interface {
	ExecuteTemplate(params map[string]interface{}) (string, []PreparedStatementParam, error)
	GetSource() string
	GetResultMap() ResultMap
}

type statement struct {
	id        string
	source    string
	tmpl      *template.Template
	resultMap ResultMap
}

type PreparedStatementParam struct {
	Name        string
	Type        reflect.Type
	Value       interface{}
	PlaceHolder string
}

func newStatement(xmlMapper *xml.MapperRootNode, elementId string, mo mapperOptions) (Statement, error) {

	if mnode, ok := xmlMapper.FindById(elementId); !ok {
		return nil, errors.New(fmt.Sprintf("Statement %s not found in maper", elementId))
	} else {

		if src, t, err := newTemplate(xmlMapper, elementId, mo); err == nil {
			s := statement{id: elementId, source: src, tmpl: t}
			if selNode, ok := mnode.(*xml.SelectNode); ok {
				if selNode.ResultMapNode != nil {
					if s.resultMap, err = newResultMap(selNode.ResultMapNode.(*xml.ResultMapNode)); err != nil {
						return nil, err
					}
				}
			}
			return &s, nil
		} else {
			return nil, err
		}
	}
}

func (s *statement) GetSource() string {
	return s.source
}

func (s *statement) ExecuteTemplate(params map[string]interface{}) (string, []PreparedStatementParam, error) {

	var mapperCtx = mapperContext{}
	var sqlStatement bytes.Buffer

	if params == nil {
		params = make(map[string]interface{})
	}
	params["mctx"] = &mapperCtx
	if err := s.tmpl.ExecuteTemplate(&sqlStatement, s.id, params); err != nil {
		log.Error().Err(err).Send()
		return "", nil, err
	}

	return util.StripDuplicateWhiteSpaces(sqlStatement.String()), mapperCtx.ParamsList, nil
}

func (s *statement) GetResultMap() ResultMap {
	return s.resultMap
}

/*
 * Template Preparation Section.
 */
type TemplateContext interface {
	addPreparedStatementParam(aCondition string, aType reflect.Type, aValue interface{}, aPlaceHolderStyle ParamBindStyle) PreparedStatementParam
}

type mapperContext struct {
	ParamsList []PreparedStatementParam
}

func (mctx *mapperContext) addPreparedStatementParam(aCondition string, aType reflect.Type, aValue interface{}, aPlaceHolderStyle ParamBindStyle) PreparedStatementParam {

	ph := "?"
	switch aPlaceHolderStyle {
	case BINDSTYLE_DOLLAR:
		ph = "$" + strconv.Itoa(len(mctx.ParamsList)+1)
	}

	pp := PreparedStatementParam{Name: aCondition, Type: aType, Value: aValue, PlaceHolder: ph}
	mctx.ParamsList = append(mctx.ParamsList, pp)
	return pp
}

type templateBuilderStack struct {
	tmpBaseId          string
	tmplSequence       int
	accumulator        strings.Builder
	stringBuilderStack []*templateBuilder
}

type templateBuilder struct {
	tmplId string
	sb     strings.Builder
}

func buildTemplateText(xmlMapper *xml.MapperRootNode, elementId string) (string, error) {
	if me, b := xmlMapper.FindById(elementId); b {

		ts := templateBuilderStack{tmpBaseId: elementId}
		me.Walk(0, WithBuildTemplateTextFunctionAnte(&ts), WithBuildTemplateTextFunctionPost(&ts))

		return ts.accumulator.String(), nil
	} else {
		return "", errors.New(fmt.Sprintf("Element Not Found: %s", elementId))
	}
}

func WithBuildTemplateTextFunctionAnte(ts *templateBuilderStack) xml.WalkOption {
	return func(walkOptions *xml.WalkOptions) {
		walkOptions.Ante = func(depth int, mapper xml.MapperNode) (bool, error) {

			topStack := len(ts.stringBuilderStack) - 1

			switch e := mapper.(type) {
			case *xml.SelectNode, *xml.InsertNode, *xml.UpdateNode, *xml.DeleteNode:
				tb := templateBuilder{}
				tb.tmplId = ts.tmpBaseId
				tb.sb.WriteString(fmt.Sprintf("{{define \"%s\"}}\n", tb.tmplId))
				ts.stringBuilderStack = append(ts.stringBuilderStack, &tb)

			case *xml.WhereNode:
				ts.tmplSequence++
				tmplId := fmt.Sprintf("%s_%d", ts.tmpBaseId, ts.tmplSequence)

				ts.stringBuilderStack[topStack].sb.WriteString("{{ $ctx := setCtx $.mctx \"data\" . }}\n")
				ts.stringBuilderStack[topStack].sb.WriteString(fmt.Sprintf("{{ $var := tmpl \"%s\" $ctx }}\n", tmplId))
				ts.stringBuilderStack[topStack].sb.WriteString("{{ where $var | print }}\n")

				tb := templateBuilder{}
				tb.tmplId = tmplId
				tb.sb.WriteString(fmt.Sprintf("{{define \"%s\"}}\n", tmplId))
				ts.stringBuilderStack = append(ts.stringBuilderStack, &tb)

			case *xml.TrimNode:
				ts.tmplSequence++
				tmplId := fmt.Sprintf("%s_%d", ts.tmpBaseId, ts.tmplSequence)

				ts.stringBuilderStack[topStack].sb.WriteString("{{ $ctx := setCtx $.mctx \"data\" . }}\n")
				ts.stringBuilderStack[topStack].sb.WriteString(fmt.Sprintf("{{ $var := tmpl \"%s\" $ctx }}\n", tmplId))
				ts.stringBuilderStack[topStack].sb.WriteString(fmt.Sprintf("{{ trim $var %q %q %q | print }}\n", e.Prefix, e.Suffix, e.PrefixOverrides))

				tb := templateBuilder{}
				tb.tmplId = tmplId
				tb.sb.WriteString(fmt.Sprintf("{{define \"%s\"}}\n", tmplId))
				ts.stringBuilderStack = append(ts.stringBuilderStack, &tb)

			case *xml.SetNode:
				ts.tmplSequence++
				tmplId := fmt.Sprintf("%s_%d", ts.tmpBaseId, ts.tmplSequence)

				ts.stringBuilderStack[topStack].sb.WriteString("{{ $ctx := setCtx $.mctx \"data\" . }}\n")
				ts.stringBuilderStack[topStack].sb.WriteString(fmt.Sprintf("{{ $var := tmpl \"%s\" $ctx }}\n", tmplId))
				ts.stringBuilderStack[topStack].sb.WriteString(fmt.Sprintf("{{ set $var | print }}\n"))

				tb := templateBuilder{}
				tb.tmplId = tmplId
				tb.sb.WriteString(fmt.Sprintf("{{define \"%s\"}}\n", tmplId))
				ts.stringBuilderStack = append(ts.stringBuilderStack, &tb)

			case *xml.ChardataNode:
				s := e.Chardata
				for _, dollar := range e.DollarVariables {
					s = strings.ReplaceAll(s, dollar.String(), fmt.Sprintf("{{ %s }}", dollar.Expr))
				}

				for _, dash := range e.DashVariables {
					s = strings.ReplaceAll(s, dash.String(), fmt.Sprintf("{{ stmtParam $.mctx  \"%s\" %s | print}}", dash.Expr, dash.Expr))
				}

				ts.stringBuilderStack[topStack].sb.WriteString(s)
				ts.stringBuilderStack[topStack].sb.WriteRune('\n')

			case *xml.IfNode:
				s := fmt.Sprintf("{{- if %s}}", e.Test)
				ts.stringBuilderStack[topStack].sb.WriteString(s)
				ts.stringBuilderStack[topStack].sb.WriteString("\n")

			case *xml.ChooseNode:
				ts.stringBuilderStack[topStack].sb.WriteString("{{- if false}}\n")

			case *xml.WhenNode:
				s := fmt.Sprintf("{{- else if %s}}", e.Test)
				ts.stringBuilderStack[topStack].sb.WriteString(s)
				ts.stringBuilderStack[topStack].sb.WriteString("\n")

			case *xml.OtherwiseNode:
				ts.stringBuilderStack[topStack].sb.WriteString("{{- else }}\n")

			case *xml.ForEachNode:
				ts.tmplSequence++
				tmplId := fmt.Sprintf("%s_%d", ts.tmpBaseId, ts.tmplSequence)

				if e.Open != "" {
					ts.stringBuilderStack[topStack].sb.WriteString(e.Open)
				}
				ts.stringBuilderStack[topStack].sb.WriteString(fmt.Sprintf("{{- range $i, $e := %s}}\n", e.Collection))
				ts.stringBuilderStack[topStack].sb.WriteString(fmt.Sprintf("{{$ctx := setCtx $.mctx %q . }}\n", e.Item))
				ts.stringBuilderStack[topStack].sb.WriteString(fmt.Sprintf("{{ $var := tmpl \"%s\" $ctx }}\n", tmplId))
				if e.Separator != "" {
					ts.stringBuilderStack[topStack].sb.WriteString(fmt.Sprintf("{{ if $i }}%s{{end}}", e.Separator))
				}
				ts.stringBuilderStack[topStack].sb.WriteString(fmt.Sprintf("{{ $var }}\n"))
				ts.stringBuilderStack[topStack].sb.WriteString(fmt.Sprintf("{{ end }}"))
				if e.Close != "" {
					ts.stringBuilderStack[topStack].sb.WriteString(e.Close)
				}
				ts.stringBuilderStack[topStack].sb.WriteRune('\n')

				tb := templateBuilder{}
				tb.tmplId = tmplId
				tb.sb.WriteString(fmt.Sprintf("{{define \"%s\"}}\n", tmplId))
				ts.stringBuilderStack = append(ts.stringBuilderStack, &tb)
			}

			return true, nil
		}
	}
}

func WithBuildTemplateTextFunctionPost(ts *templateBuilderStack) xml.WalkOption {
	return func(walkOptions *xml.WalkOptions) {
		walkOptions.After = func(depth int, mapper xml.MapperNode) (bool, error) {

			topStack := len(ts.stringBuilderStack) - 1
			switch mapper.(type) {
			case *xml.SelectNode, *xml.InsertNode, *xml.UpdateNode, *xml.DeleteNode:
				ts.stringBuilderStack[topStack].sb.WriteString(fmt.Sprintf("{{end}}\n"))
				ts.accumulator.WriteString(ts.stringBuilderStack[topStack].sb.String())
				ts.accumulator.WriteRune('\n')
				ts.stringBuilderStack = ts.stringBuilderStack[:topStack]

			case *xml.WhereNode:
				ts.stringBuilderStack[topStack].sb.WriteString(fmt.Sprintf("{{end}}\n"))
				ts.accumulator.WriteString(ts.stringBuilderStack[topStack].sb.String())
				ts.accumulator.WriteRune('\n')
				ts.stringBuilderStack = ts.stringBuilderStack[:topStack]

			case *xml.TrimNode:
				ts.stringBuilderStack[topStack].sb.WriteString(fmt.Sprintf("{{end}}\n"))
				ts.accumulator.WriteString(ts.stringBuilderStack[topStack].sb.String())
				ts.accumulator.WriteRune('\n')
				ts.stringBuilderStack = ts.stringBuilderStack[:topStack]

			case *xml.SetNode:
				ts.stringBuilderStack[topStack].sb.WriteString(fmt.Sprintf("{{end}}\n"))
				ts.accumulator.WriteString(ts.stringBuilderStack[topStack].sb.String())
				ts.accumulator.WriteRune('\n')
				ts.stringBuilderStack = ts.stringBuilderStack[:topStack]

			case *xml.IfNode:
				ts.stringBuilderStack[topStack].sb.WriteString("{{end}}\n")

			case *xml.ChooseNode:
				ts.stringBuilderStack[topStack].sb.WriteString("{{end}}\n")

			case *xml.WhenNode:
				// ts.stringBuilderStack[topStack].sb.WriteString("{{end}}\n")

			case *xml.ForEachNode:
				ts.stringBuilderStack[topStack].sb.WriteString(fmt.Sprintf("{{end}}\n"))
				ts.accumulator.WriteString(ts.stringBuilderStack[topStack].sb.String())
				ts.accumulator.WriteRune('\n')
				ts.stringBuilderStack = ts.stringBuilderStack[:topStack]
			}

			return true, nil
		}
	}
}

func newTemplate(xmlMapper *xml.MapperRootNode, elementId string, mo mapperOptions) (string, *template.Template, error) {
	if txt, e := buildTemplateText(xmlMapper, elementId); e == nil {

		var t *template.Template
		t, err := template.New(elementId).Funcs(template.FuncMap{
			"tmpl": func(name string, data interface{}) (string, error) {
				buf := &bytes.Buffer{}
				err := t.ExecuteTemplate(buf, name, data)
				return buf.String(), err
			},
			"setCtx": func(mapperCtx TemplateContext, dataKey string, data interface{}) interface{} {

				data1 := data
				pval := reflect.ValueOf(data)
				if pval.Kind() != reflect.Map {
					data1 = map[string]interface{}{"mctx": mapperCtx, dataKey: data}
					return data1
				}

				return data
			},
			"trim": func(aString string, prefix string, suffix string, prefixOverrides string) string {
				t := strings.TrimSpace(aString)
				if prefixOverrides != "" && strings.HasPrefix(strings.ToUpper(t), strings.ToUpper(prefixOverrides)) {
					t = t[len(prefixOverrides):]
				}

				if t != "" {
					return prefix + strings.TrimSpace(t) + suffix
				}

				return ""
			},
			"set": func(aString string) string {
				t := strings.TrimSpace(aString)
				t = strings.TrimSuffix(t, ",")

				if t != "" {
					return "set " + t
				}

				return t
			},
			"where": func(aString string) string {
				t := strings.TrimSpace(aString)

				if t != "" {
					return "where " + t
				}

				return ""
			},
			"stmtParam": func(mctx TemplateContext, aName string, aValue interface{}) string {
				// fmt.Printf("%T\n", aValue)
				pval := reflect.ValueOf(aValue)
				if pval.Kind() == reflect.Slice || pval.Kind() == reflect.Array {
					var sb strings.Builder
					for i := 0; i < pval.Len(); i++ {

						iValue := pval.Index(i).Interface()
						pp := mctx.addPreparedStatementParam(aName, reflect.TypeOf(iValue), iValue, mo.bindStyle)
						if i > 0 {
							sb.WriteString(", ")
						}

						sb.WriteString(pp.PlaceHolder)
					}
					return "(" + sb.String() + ")"
				}

				pp := mctx.addPreparedStatementParam(aName, reflect.TypeOf(aValue), aValue, mo.bindStyle)
				return pp.PlaceHolder
			},
		}).Parse(txt)

		return txt, t, err

	} else {
		return "", nil, e
	}

}

/*
 *
 */
func newResultMap(xmlMapper *xml.ResultMapNode) (ResultMap, error) {
	rmap := &resultmap{}
	return rmap, nil
}

func (rsm *resultmap) MapColumnNames2TypeProperties(columns []string) ([]string, error) {
	return columns, nil
}
