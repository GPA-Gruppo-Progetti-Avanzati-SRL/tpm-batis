package xml

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type WalkOption func(walkOptions *WalkOptions)
type WalkVisitFunc func(int, MapperNode) (bool, error)

type WalkOptions struct {
	Ante  WalkVisitFunc
	After WalkVisitFunc
}

type MapperNode interface {
	SetAttribute(xml.Attr)
	SetNext(node MapperNode)
	GetName() xml.Name
	SetName(aName xml.Name)
	GetId() string
	Clone() MapperNode
	GetNodes() []MapperNode
	SetNodes(nodes []MapperNode)
	Walk(depth int, visitFuncs ...WalkOption) (bool, error)
	String() string
}

type MapperNodeBase struct {
	XMLName xml.Name
	Id      string `xml:"id,attr"`
	Nodes   []MapperNode
}

func (e *MapperNodeBase) GetName() xml.Name {
	return e.XMLName
}

func (e *MapperNodeBase) GetId() string {
	return e.Id
}

func (e *MapperNodeBase) GetNodes() []MapperNode {
	return e.Nodes
}

func (e *MapperNodeBase) SetNodes(nodes []MapperNode) {
	e.Nodes = nodes
}

/*
func (e *MapperNodeBase) String() string {
	return e.XMLName.Local
}
*/

func (e *MapperNodeBase) SetName(aName xml.Name) {
	e.XMLName = aName
}

func (e *MapperNodeBase) SetAttribute(_ xml.Attr) {

}

func (e *MapperNodeBase) SetNext(mn MapperNode) {
	e.Nodes = append(e.Nodes, mn)
}

/*
func (e *MapperNodeBase)Clone()  MapperNode {
	ch2 := *e
	ch2.Nodes = make([]MapperNode, 0, len(e.Nodes))
	for _, n := range e.Nodes {
		ch2.Nodes = append(ch2.Nodes, n.Clone())
	}
	return &ch2
}*/

/*
func (e *MapperNodeBase)Walk(depth int, f func(int, MapperNode) bool) bool {

	f(depth, e)
	for _, n := range e.Nodes {
		n.Walk(depth + 1, f)
	}

	return true
}
*/

//
// CharData
//
type ChardataNode struct {
	MapperNodeBase
	Chardata string

	DollarVariables []StatementParam
	DashVariables   []StatementParam
}

func (e *ChardataNode) Clone() MapperNode {
	ch2 := *e
	ch2.Nodes = make([]MapperNode, 0, len(e.Nodes))
	for _, n := range e.Nodes {
		ch2.Nodes = append(ch2.Nodes, n.Clone())
	}
	return &ch2
}

func (e *ChardataNode) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s - %s ", e.XMLName.Local, e.Chardata))
	sb.WriteString(" Vars: ")
	for _, v := range e.DollarVariables {
		sb.WriteString(v.String())
	}
	for _, v := range e.DashVariables {
		sb.WriteString(v.String())
	}
	return sb.String()
}

func (e *ChardataNode) Walk(depth int, visitFuncs ...WalkOption) (bool, error) {

	vfs := WalkOptions{Ante: func(depth int, me MapperNode) (bool, error) { return true, nil }, After: func(depth int, me MapperNode) (bool, error) { return true, nil }}
	for _, wv := range visitFuncs {
		wv(&vfs)
	}

	if b, e := vfs.Ante(depth, e); !b || e != nil {
		return b, e
	}

	for _, n := range e.Nodes {
		if b, e := n.Walk(depth+1, visitFuncs...); !b || e != nil {
			return b, e
		}
	}

	if b, e := vfs.After(depth, e); !b || e != nil {
		return b, e
	}

	return true, nil
}

//
// ChooseNode
//
type ChooseNode struct {
	MapperNodeBase
}

func (e *ChooseNode) Clone() MapperNode {
	ch2 := *e
	ch2.Nodes = make([]MapperNode, 0, len(e.Nodes))
	for _, n := range e.Nodes {
		ch2.Nodes = append(ch2.Nodes, n.Clone())
	}
	return &ch2
}

func (e *ChooseNode) String() string {
	return fmt.Sprintf("%s", e.XMLName.Local)
}

func (e *ChooseNode) Walk(depth int, visitFuncs ...WalkOption) (bool, error) {

	vfs := WalkOptions{Ante: func(depth int, me MapperNode) (bool, error) { return true, nil }, After: func(depth int, me MapperNode) (bool, error) { return true, nil }}
	for _, wv := range visitFuncs {
		wv(&vfs)
	}

	if b, e := vfs.Ante(depth, e); !b || e != nil {
		return b, e
	}

	for _, n := range e.Nodes {
		if b, e := n.Walk(depth+1, visitFuncs...); !b || e != nil {
			return b, e
		}
	}

	if b, e := vfs.After(depth, e); !b || e != nil {
		return b, e
	}

	return true, nil
}

//
// DeleteNode
//
type DeleteNode struct {
	MapperNodeBase
	ParameterType string
}

func (e *DeleteNode) Clone() MapperNode {
	ch2 := *e
	ch2.Nodes = make([]MapperNode, 0, len(e.Nodes))
	for _, n := range e.Nodes {
		ch2.Nodes = append(ch2.Nodes, n.Clone())
	}
	return &ch2
}

func (e *DeleteNode) String() string {
	return fmt.Sprintf("%s - Id: %s - ParameterType: %s", e.XMLName.Local, e.Id, e.ParameterType)
}

func (e *DeleteNode) SetAttribute(attr xml.Attr) {
	switch strings.ToLower(attr.Name.Local) {
	case "id":
		e.Id = attr.Value
	case "parametertype":
		e.ParameterType = attr.Value
	}
}

func (e *DeleteNode) Walk(depth int, visitFuncs ...WalkOption) (bool, error) {

	vfs := WalkOptions{Ante: func(depth int, me MapperNode) (bool, error) { return true, nil }, After: func(depth int, me MapperNode) (bool, error) { return true, nil }}
	for _, wv := range visitFuncs {
		wv(&vfs)
	}

	if b, e := vfs.Ante(depth, e); !b || e != nil {
		return b, e
	}

	for _, n := range e.Nodes {
		if b, e := n.Walk(depth+1, visitFuncs...); !b || e != nil {
			return b, e
		}
	}

	if b, e := vfs.After(depth, e); !b || e != nil {
		return b, e
	}

	return true, nil
}

//
// ForEachNode
//
type ForEachNode struct {
	MapperNodeBase
	Collection string
	Item       string
	Separator  string
	Open       string
	Close      string
}

func (e *ForEachNode) SetAttribute(attr xml.Attr) {
	switch strings.ToLower(attr.Name.Local) {
	case "collection":
		e.Collection = attr.Value
	case "item":
		e.Item = attr.Value
	case "separator":
		e.Separator = attr.Value
	case "open":
		e.Open = attr.Value
	case "close":
		e.Close = attr.Value
	}
}

func (e *ForEachNode) String() string {
	return fmt.Sprintf("%s - Collection: %s - Item: %s - Sep: %s", e.XMLName.Local, e.Collection, e.Item, e.Separator)
}

func (e *ForEachNode) Clone() MapperNode {
	ch2 := *e
	ch2.Nodes = make([]MapperNode, 0, len(e.Nodes))
	for _, n := range e.Nodes {
		ch2.Nodes = append(ch2.Nodes, n.Clone())
	}
	return &ch2
}

func (e *ForEachNode) Walk(depth int, visitFuncs ...WalkOption) (bool, error) {

	vfs := WalkOptions{Ante: func(depth int, me MapperNode) (bool, error) { return true, nil }, After: func(depth int, me MapperNode) (bool, error) { return true, nil }}
	for _, wv := range visitFuncs {
		wv(&vfs)
	}

	if b, e := vfs.Ante(depth, e); !b || e != nil {
		return b, e
	}

	for _, n := range e.Nodes {
		if b, e := n.Walk(depth+1, visitFuncs...); !b || e != nil {
			return b, e
		}
	}

	if b, e := vfs.After(depth, e); !b || e != nil {
		return b, e
	}

	return true, nil
}

//
// IfNode
//
type IfNode struct {
	MapperNodeBase
	Test string
}

func (e *IfNode) String() string {
	return fmt.Sprintf("%s - Test: %s", e.XMLName.Local, e.Test)
}

func (e *IfNode) Clone() MapperNode {
	ch2 := *e
	ch2.Nodes = make([]MapperNode, 0, len(e.Nodes))
	for _, n := range e.Nodes {
		ch2.Nodes = append(ch2.Nodes, n.Clone())
	}
	return &ch2
}

func (e *IfNode) SetAttribute(attr xml.Attr) {
	switch strings.ToLower(attr.Name.Local) {
	case "test":
		e.Test = attr.Value
	}
}

func (e *IfNode) Walk(depth int, visitFuncs ...WalkOption) (bool, error) {

	vfs := WalkOptions{Ante: func(depth int, me MapperNode) (bool, error) { return true, nil }, After: func(depth int, me MapperNode) (bool, error) { return true, nil }}
	for _, wv := range visitFuncs {
		wv(&vfs)
	}

	if b, e := vfs.Ante(depth, e); !b || e != nil {
		return b, e
	}

	for _, n := range e.Nodes {
		if b, e := n.Walk(depth+1, visitFuncs...); !b || e != nil {
			return b, e
		}
	}

	if b, e := vfs.After(depth, e); !b || e != nil {
		return b, e
	}

	return true, nil
}

//
// IncludeNode
//
type IncludeNode struct {
	MapperNodeBase
	RefId      string            `xml:"refid,attr"`
	Properties []IncludeProperty `xml:"property"`
}

type IncludeProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

func (e *IncludeNode) String() string {

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s - Refid %s Properties: [", e.XMLName.Local, e.RefId))
	for _, p := range e.Properties {
		sb.WriteString(fmt.Sprintf("{Name: %s - Options: %s}", p.Name, p.Value))
	}
	sb.WriteString("]")
	return sb.String()
}

func (e *IncludeNode) Clone() MapperNode {
	ch2 := *e
	ch2.Nodes = make([]MapperNode, 0, len(e.Nodes))
	for _, n := range e.Nodes {
		ch2.Nodes = append(ch2.Nodes, n.Clone())
	}
	return &ch2
}

func (e *IncludeNode) SetAttribute(attr xml.Attr) {
	switch strings.ToLower(attr.Name.Local) {
	case "refid":
		e.RefId = attr.Value
	}
}

func (e *IncludeNode) Walk(depth int, visitFuncs ...WalkOption) (bool, error) {

	vfs := WalkOptions{Ante: func(depth int, me MapperNode) (bool, error) { return true, nil }, After: func(depth int, me MapperNode) (bool, error) { return true, nil }}
	for _, wv := range visitFuncs {
		wv(&vfs)
	}

	if b, e := vfs.Ante(depth, e); !b || e != nil {
		return b, e
	}

	for _, n := range e.Nodes {
		if b, e := n.Walk(depth+1, visitFuncs...); !b || e != nil {
			return b, e
		}
	}

	if b, e := vfs.After(depth, e); !b || e != nil {
		return b, e
	}

	return true, nil
}

//
// InsertNode
//
type InsertNode struct {
	MapperNodeBase

	ResultMapId   string
	ParameterType string
	ResultType    string
}

func (e *InsertNode) String() string {
	return fmt.Sprintf("%s - Id: %s", e.XMLName.Local, e.Id)
}

func (e *InsertNode) Clone() MapperNode {
	ch2 := *e
	ch2.Nodes = make([]MapperNode, 0, len(e.Nodes))
	for _, n := range e.Nodes {
		ch2.Nodes = append(ch2.Nodes, n.Clone())
	}
	return &ch2
}

func (e *InsertNode) SetAttribute(attr xml.Attr) {
	switch strings.ToLower(attr.Name.Local) {
	case "id":
		e.Id = attr.Value
	case "parametertype":
		e.ParameterType = attr.Value
	}
}

func (e *InsertNode) Walk(depth int, visitFuncs ...WalkOption) (bool, error) {

	vfs := WalkOptions{Ante: func(depth int, me MapperNode) (bool, error) { return true, nil }, After: func(depth int, me MapperNode) (bool, error) { return true, nil }}
	for _, wv := range visitFuncs {
		wv(&vfs)
	}

	if b, e := vfs.Ante(depth, e); !b || e != nil {
		return b, e
	}

	for _, n := range e.Nodes {
		if b, e := n.Walk(depth+1, visitFuncs...); !b || e != nil {
			return b, e
		}
	}

	if b, e := vfs.After(depth, e); !b || e != nil {
		return b, e
	}

	return true, nil
}

//
// MapperRootNode
//
type MapperRootNode struct {
	MapperNodeBase
	Namespace string
}

func (e *MapperRootNode) String() string {
	return fmt.Sprintf("%s - Namespace: %s", e.XMLName.Local, e.Namespace)
}

func (e *MapperRootNode) SetAttribute(attr xml.Attr) {
	switch strings.ToLower(attr.Name.Local) {
	case "namespace":
		e.Namespace = attr.Value
	}
}

func (e *MapperRootNode) Clone() MapperNode {
	ch2 := *e
	ch2.Nodes = make([]MapperNode, 0, len(e.Nodes))
	for _, n := range e.Nodes {
		ch2.Nodes = append(ch2.Nodes, n.Clone())
	}
	return &ch2
}

func (e *MapperRootNode) Walk(depth int, visitFuncs ...WalkOption) (bool, error) {

	vfs := WalkOptions{Ante: func(depth int, me MapperNode) (bool, error) { return true, nil }, After: func(depth int, me MapperNode) (bool, error) { return true, nil }}
	for _, wv := range visitFuncs {
		wv(&vfs)
	}

	if b, e := vfs.Ante(depth, e); !b || e != nil {
		return b, e
	}

	for _, n := range e.Nodes {
		if b, e := n.Walk(depth+1, visitFuncs...); !b || e != nil {
			return b, e
		}
	}

	if b, e := vfs.After(depth, e); !b || e != nil {
		return b, e
	}

	return true, nil
}

func (e *MapperRootNode) FindById(id string) (MapperNode, bool) {
	for _, n := range e.Nodes {
		if n.GetId() == id {
			return n, true
		}
	}

	return nil, false
}

func (e *MapperRootNode) GetStatementsIds() []string {
	var nodes []string
	for _, me := range e.GetNodes() {
		switch me.(type) {
		case *SelectNode:
			nodes = append(nodes, me.GetId())
		case *InsertNode:
			nodes = append(nodes, me.GetId())
		case *UpdateNode:
			nodes = append(nodes, me.GetId())
		case *DeleteNode:
			nodes = append(nodes, me.GetId())
		}
	}

	return nodes
}

func (e *MapperRootNode) GetStatements() []MapperNode {
	var nodes []MapperNode
	for _, me := range e.GetNodes() {
		switch me.(type) {
		case *SelectNode:
			nodes = append(nodes, me)
		case *InsertNode:
			nodes = append(nodes, me)
		case *UpdateNode:
			nodes = append(nodes, me)
		case *DeleteNode:
			nodes = append(nodes, me)
		}
	}

	return nodes
}

//
// OtherwiseNode
//
type OtherwiseNode struct {
	MapperNodeBase
}

func (e *OtherwiseNode) Clone() MapperNode {
	ch2 := *e
	ch2.Nodes = make([]MapperNode, 0, len(e.Nodes))
	for _, n := range e.Nodes {
		ch2.Nodes = append(ch2.Nodes, n.Clone())
	}
	return &ch2
}

func (e *OtherwiseNode) String() string {
	return fmt.Sprintf("%s", e.XMLName.Local)
}

func (e *OtherwiseNode) Walk(depth int, visitFuncs ...WalkOption) (bool, error) {

	vfs := WalkOptions{Ante: func(depth int, me MapperNode) (bool, error) { return true, nil }, After: func(depth int, me MapperNode) (bool, error) { return true, nil }}
	for _, wv := range visitFuncs {
		wv(&vfs)
	}

	if b, e := vfs.Ante(depth, e); !b || e != nil {
		return b, e
	}

	for _, n := range e.Nodes {
		if b, e := n.Walk(depth+1, visitFuncs...); !b || e != nil {
			return b, e
		}
	}

	if b, e := vfs.After(depth, e); !b || e != nil {
		return b, e
	}

	return true, nil
}

//
//
//
type ResultMapNode struct {
	MapperNodeBase
	Type string `xml:"type,attr"`

	ResultColumns []ResultMapResultColumn `xml:"result"`
	IdColumns     []ResultMapResultColumn `xml:"id"`
	Collections   []ResultMapCollection   `xml:"collection"`
}

type ResultMapResultColumn struct {
	XMLName  xml.Name
	Column   string `xml:"column,attr"`
	Property string `xml:"property,attr"`
	JdbcType string `xml:"jdbcType,attr"`
}

type ResultMapCollection struct {
	XMLName  xml.Name
	Column   string `xml:"column,attr"`
	Property string `xml:"property,attr"`
	JavaType string `xml:"javaType,attr"`

	ResultColumns []ResultMapResultColumn `xml:"result"`
	IdColumns     []ResultMapResultColumn `xml:"id"`
	Collections   []ResultMapCollection   `xml:"collection"`
}

func (e *ResultMapNode) SetAttribute(attr xml.Attr) {
	switch strings.ToLower(attr.Name.Local) {
	case "id":
		e.Id = attr.Value
	case "type":
		e.Type = attr.Value
	}
}

func (e *ResultMapNode) String() string {
	return fmt.Sprintf("%s - Id: %s - Type: %s", e.XMLName.Local, e.Id, e.Type)
}

func (e *ResultMapNode) Clone() MapperNode {
	return e
}

func (e *ResultMapNode) Walk(depth int, visitFuncs ...WalkOption) (bool, error) {

	vfs := WalkOptions{Ante: func(depth int, me MapperNode) (bool, error) { return true, nil }, After: func(depth int, me MapperNode) (bool, error) { return true, nil }}
	for _, wv := range visitFuncs {
		wv(&vfs)
	}

	if b, e := vfs.Ante(depth, e); !b || e != nil {
		return b, e
	}

	if b, e := vfs.After(depth, e); !b || e != nil {
		return b, e
	}

	return true, nil
}

//
// SelectNode
//
type SelectNode struct {
	MapperNodeBase

	ResultMapId   string
	ResultMapNode MapperNode

	ParameterType string
	ResultType    string
}

func (e *SelectNode) SetAttribute(attr xml.Attr) {
	switch strings.ToLower(attr.Name.Local) {
	case "id":
		e.Id = attr.Value
	case "resultmap":
		e.ResultMapId = attr.Value
	case "parametertype":
		e.ParameterType = attr.Value
	case "resulttype":
		e.ResultType = attr.Value
	}
}

func (e *SelectNode) String() string {
	return fmt.Sprintf("%s - Id: %s - ResultMap: %s - ParameterType: %s - ResultType: %s", e.XMLName.Local, e.Id, e.ResultMapId, e.ParameterType, e.ResultType)
}

func (e *SelectNode) Clone() MapperNode {
	ch2 := *e
	ch2.Nodes = make([]MapperNode, 0, len(e.Nodes))
	for _, n := range e.Nodes {
		ch2.Nodes = append(ch2.Nodes, n.Clone())
	}
	return &ch2
}

func (e *SelectNode) Walk(depth int, visitFuncs ...WalkOption) (bool, error) {

	vfs := WalkOptions{Ante: func(depth int, me MapperNode) (bool, error) { return true, nil }, After: func(depth int, me MapperNode) (bool, error) { return true, nil }}
	for _, wv := range visitFuncs {
		wv(&vfs)
	}

	if b, e := vfs.Ante(depth, e); !b || e != nil {
		return b, e
	}

	for _, n := range e.Nodes {
		if b, e := n.Walk(depth+1, visitFuncs...); !b || e != nil {
			return b, e
		}
	}

	if b, e := vfs.After(depth, e); !b || e != nil {
		return b, e
	}

	return true, nil
}

//
// SQLNode
//
type SQLNode struct {
	MapperNodeBase
}

func (e *SQLNode) SetAttribute(attr xml.Attr) {
	switch attr.Name.Local {
	case "id":
		e.Id = attr.Value
	}
}

func (e *SQLNode) String() string {
	return fmt.Sprintf("%s - Id: %s", e.XMLName.Local, e.Id)
}

func (e *SQLNode) Clone() MapperNode {
	ch2 := *e
	ch2.Nodes = make([]MapperNode, 0, len(e.Nodes))
	for _, n := range e.Nodes {
		ch2.Nodes = append(ch2.Nodes, n.Clone())
	}
	return &ch2
}

func (e *SQLNode) Walk(depth int, visitFuncs ...WalkOption) (bool, error) {

	vfs := WalkOptions{Ante: func(depth int, me MapperNode) (bool, error) { return true, nil }, After: func(depth int, me MapperNode) (bool, error) { return true, nil }}
	for _, wv := range visitFuncs {
		wv(&vfs)
	}

	if b, e := vfs.Ante(depth, e); !b || e != nil {
		return b, e
	}

	for _, n := range e.Nodes {
		if b, e := n.Walk(depth+1, visitFuncs...); !b || e != nil {
			return b, e
		}
	}

	if b, e := vfs.After(depth, e); !b || e != nil {
		return b, e
	}

	return true, nil
}

//
// TrimNode
//
type TrimNode struct {
	MapperNodeBase
	Prefix          string
	Suffix          string
	PrefixOverrides string
}

func (e *TrimNode) Clone() MapperNode {
	ch2 := *e
	ch2.Nodes = make([]MapperNode, 0, len(e.Nodes))
	for _, n := range e.Nodes {
		ch2.Nodes = append(ch2.Nodes, n.Clone())
	}
	return &ch2
}

func (e *TrimNode) String() string {
	return fmt.Sprintf("%s - Prefix: %s - Suffix: %s PrefixOverrides: %s", e.XMLName.Local, e.Prefix, e.Suffix, e.PrefixOverrides)
}

func (e *TrimNode) SetAttribute(attr xml.Attr) {
	switch strings.ToLower(attr.Name.Local) {
	case "prefix":
		e.Prefix = attr.Value
	case "suffix":
		e.Suffix = attr.Value
	case "prefixoverrides":
		e.PrefixOverrides = attr.Value
	}
}

func (e *TrimNode) Walk(depth int, visitFuncs ...WalkOption) (bool, error) {

	vfs := WalkOptions{Ante: func(depth int, me MapperNode) (bool, error) { return true, nil }, After: func(depth int, me MapperNode) (bool, error) { return true, nil }}
	for _, wv := range visitFuncs {
		wv(&vfs)
	}

	if b, e := vfs.Ante(depth, e); !b || e != nil {
		return b, e
	}

	for _, n := range e.Nodes {
		if b, e := n.Walk(depth+1, visitFuncs...); !b || e != nil {
			return b, e
		}
	}

	if b, e := vfs.After(depth, e); !b || e != nil {
		return b, e
	}

	return true, nil
}

//
// UpdateNode
//
type UpdateNode struct {
	MapperNodeBase
	ParameterType string
}

func (e *UpdateNode) String() string {
	return fmt.Sprintf("%s - Id: %s - ParameterType: %s", e.XMLName.Local, e.Id, e.ParameterType)
}

func (e *UpdateNode) Clone() MapperNode {
	ch2 := *e
	ch2.Nodes = make([]MapperNode, 0, len(e.Nodes))
	for _, n := range e.Nodes {
		ch2.Nodes = append(ch2.Nodes, n.Clone())
	}
	return &ch2
}

func (e *UpdateNode) SetAttribute(attr xml.Attr) {
	switch strings.ToLower(attr.Name.Local) {
	case "id":
		e.Id = attr.Value
	case "parametertype":
		e.ParameterType = attr.Value
	}
}

func (e *UpdateNode) Walk(depth int, visitFuncs ...WalkOption) (bool, error) {

	vfs := WalkOptions{Ante: func(depth int, me MapperNode) (bool, error) { return true, nil }, After: func(depth int, me MapperNode) (bool, error) { return true, nil }}
	for _, wv := range visitFuncs {
		wv(&vfs)
	}

	if b, e := vfs.Ante(depth, e); !b || e != nil {
		return b, e
	}

	for _, n := range e.Nodes {
		if b, e := n.Walk(depth+1, visitFuncs...); !b || e != nil {
			return b, e
		}
	}

	if b, e := vfs.After(depth, e); !b || e != nil {
		return b, e
	}

	return true, nil
}

//
// WhenNode
//
type WhenNode struct {
	MapperNodeBase
	Test string
}

func (e *WhenNode) String() string {
	return fmt.Sprintf("%s - Test: %s", e.XMLName.Local, e.Test)
}

func (e *WhenNode) Clone() MapperNode {
	ch2 := *e
	ch2.Nodes = make([]MapperNode, 0, len(e.Nodes))
	for _, n := range e.Nodes {
		ch2.Nodes = append(ch2.Nodes, n.Clone())
	}
	return &ch2
}

func (e *WhenNode) SetAttribute(attr xml.Attr) {
	switch strings.ToLower(attr.Name.Local) {
	case "test":
		e.Test = attr.Value
	}
}

func (e *WhenNode) Walk(depth int, visitFuncs ...WalkOption) (bool, error) {

	vfs := WalkOptions{Ante: func(depth int, me MapperNode) (bool, error) { return true, nil }, After: func(depth int, me MapperNode) (bool, error) { return true, nil }}
	for _, wv := range visitFuncs {
		wv(&vfs)
	}

	if b, e := vfs.Ante(depth, e); !b || e != nil {
		return b, e
	}

	for _, n := range e.Nodes {
		if b, e := n.Walk(depth+1, visitFuncs...); !b || e != nil {
			return b, e
		}
	}

	if b, e := vfs.After(depth, e); !b || e != nil {
		return b, e
	}

	return true, nil
}

//
// WhereNode
//
type WhereNode struct {
	MapperNodeBase
}

func (e *WhereNode) String() string {
	return fmt.Sprintf("%s", e.XMLName.Local)
}

func (e *WhereNode) Clone() MapperNode {
	ch2 := *e
	ch2.Nodes = make([]MapperNode, 0, len(e.Nodes))
	for _, n := range e.Nodes {
		ch2.Nodes = append(ch2.Nodes, n.Clone())
	}
	return &ch2
}

func (e *WhereNode) Walk(depth int, visitFuncs ...WalkOption) (bool, error) {

	vfs := WalkOptions{Ante: func(depth int, me MapperNode) (bool, error) { return true, nil }, After: func(depth int, me MapperNode) (bool, error) { return true, nil }}
	for _, wv := range visitFuncs {
		wv(&vfs)
	}

	if b, e := vfs.Ante(depth, e); !b || e != nil {
		return b, e
	}

	for _, n := range e.Nodes {
		if b, e := n.Walk(depth+1, visitFuncs...); !b || e != nil {
			return b, e
		}
	}

	if b, e := vfs.After(depth, e); !b || e != nil {
		return b, e
	}

	return true, nil
}
