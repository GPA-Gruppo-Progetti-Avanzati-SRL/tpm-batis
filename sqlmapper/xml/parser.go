package xml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"strings"
)

const DepthIndicator = ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"

func WithTreePrintWalkOption() WalkOption {
	return func(walkOptions *WalkOptions) {

		walkOptions.Ante = func(depth int, me MapperNode) (bool, error) {
			switch v := me.(type) {

			case *ChardataNode:
			case *MapperRootNode:
			case *SQLNode:
			case *ForEachNode:
			case *IfNode:
			case *TrimNode:
			case *WhenNode:
			case *WhereNode:
			case *ChooseNode:
			case *SelectNode:
			case *DeleteNode:
			case *InsertNode:
			case *OtherwiseNode:
			case *UpdateNode:
			case *SetNode:
			case *IncludeNode:

			default:
				log.Trace().Msg(fmt.Sprint(v))
			}

			log.Trace().Msg(fmt.Sprint(DepthIndicator[0:depth*2], me))
			return true, nil
		}
	}
}

func ParseXML(XMLdata string) (*MapperRootNode, error) {

	var mapperNodes []MapperNode

	decoder := xml.NewDecoder(strings.NewReader(string(XMLdata)))
	var skipTreeDepth int
	for {

		// err is ignore here. IF you are reading from a XML file
		// do not ignore err and also check for io.EOF
		token, derr := decoder.Token()
		if derr != nil && derr != io.EOF {
			return nil, derr
		}

		if token == nil {
			return mapperNodes[0].(*MapperRootNode), nil
		}

		switch Element := token.(type) {
		case xml.StartElement:

			var n MapperNode = nil
			switch Element.Name.Local {
			case "choose":
				n = &ChooseNode{}
			case "delete":
				n = &DeleteNode{}
			case "foreach":
				n = &ForEachNode{}
			case "if":
				n = &IfNode{}
			case "include":
				n2 := &IncludeNode{}
				n2.SetName(Element.Name)
				err := decoder.DecodeElement(n2, &Element)
				if err != nil {
					panic(err)
				}
				mapperNodes[len(mapperNodes)-1].SetNext(n2)
			case "insert":
				n = &InsertNode{}
			case "set":
				n = &SetNode{}
			case "sqlmapper":
				n = &MapperRootNode{}
			case "otherwise":
				n = &OtherwiseNode{}
			case "trim":
				n = &TrimNode{}
			case "resultMap":
				n1 := &ResultMapNode{}
				n1.SetName(Element.Name)
				err := decoder.DecodeElement(n1, &Element)
				if err != nil {
					panic(err)
				}
				mapperNodes[len(mapperNodes)-1].SetNext(n1)
			case "select":
				n = &SelectNode{}
			case "sql":
				n = &SQLNode{}
			case "update":
				n = &UpdateNode{}
			case "when":
				n = &WhenNode{}
			case "where":
				n = &WhereNode{}

			default:
				skipTreeDepth++
			}

			if n != nil {
				n.SetName(Element.Name)
				for _, a := range Element.Attr {
					n.SetAttribute(a)
				}
				mapperNodes = append(mapperNodes, n)
			}

		case xml.EndElement:

			if skipTreeDepth > 0 {
				skipTreeDepth--
			} else {
				if len(mapperNodes) > 1 && mapperNodes[len(mapperNodes)-1].GetName().Local == Element.Name.Local {
					mapperNodes[len(mapperNodes)-2].SetNext(mapperNodes[len(mapperNodes)-1])
					mapperNodes = mapperNodes[0 : len(mapperNodes)-1]
				}
			}

		case xml.CharData:

			if skipTreeDepth == 0 {
				str := strings.TrimSpace(string([]byte(Element)))
				if str != "" {
					cdata := ChardataNode{}
					cdata.XMLName = xml.Name{"", "__chardata__"}
					cdata.Chardata = str
					if len(mapperNodes) > 0 {
						mapperNodes[len(mapperNodes)-1].SetNext(&cdata)
					}
				}
			}
		}
	}
}

func ParseMapper(data2 string) (*MapperRootNode, error) {

	mapper, perr := ParseXML(string(data2))
	if perr != nil {
		return nil, perr
	}

	mapper.Walk(0, WithTreePrintWalkOption())

	var includingElems []MapperNode
	for _, me := range mapper.GetNodes() {
		switch typednode := me.(type) {
		case *SelectNode:
			includingElems = append(includingElems, me)
			if typednode.ResultMapId != "" {
				if rmn, ok := mapper.FindById(typednode.ResultMapId); !ok {
					return nil, errors.New(fmt.Sprintf("Could Not Find Declared ResultMap %s in %s", typednode.ResultMapId, me.GetId()))
				} else {
					typednode.ResultMapNode = rmn
				}
			}
		case *InsertNode:
			includingElems = append(includingElems, me)
		case *UpdateNode:
			includingElems = append(includingElems, me)
		case *DeleteNode:
			includingElems = append(includingElems, me)
		}
	}

	log.Trace().Interface("els", includingElems).Send()
	// Con questo loop risolvo
	for _, ie := range includingElems {

		ie.Walk(0, WithProcessIncludeVisitFunction(mapper))
	}

	mapper.Walk(0, WithTreePrintWalkOption())
	log.Trace().Msg("-------------------------------------")

	for _, ie := range includingElems {
		ie.Walk(0, WithExtractVariablesWalkOption())
		ie.Walk(0, WithFlattenIncludeVisitFunction())
	}

	log.Trace().Msg("Include Flattened ----------------")
	mapper.Walk(0, WithTreePrintWalkOption())

	return mapper, nil
}

func WithProcessIncludeVisitFunction(mapper *MapperRootNode) WalkOption {
	return func(walkOptions *WalkOptions) {
		walkOptions.Ante = func(depth int, me MapperNode) (bool, error) {
			if includeElement, ok := me.(*IncludeNode); ok {
				id := includeElement.RefId
				if refel, ok := mapper.FindById(id); ok {
					if _, ok := refel.(*SQLNode); !ok {
						return false, errors.New(fmt.Sprintf("Only SQL ELement Can Be Included: %s\n", refel.GetName().Local))
					}

					cloned := refel.Clone()
					includeElement.Nodes = append(includeElement.Nodes, cloned)
					if len(includeElement.Properties) > 0 {
						for _, e1 := range cloned.GetNodes() {
							if cdata, ok := e1.(*ChardataNode); ok {
								for _, p := range includeElement.Properties {
									cdata.Chardata = strings.ReplaceAll(cdata.Chardata, "${"+p.Name+"}", p.Value)
								}
							}
						}
					}
				}
			}

			return true, nil
		}
	}
}

func WithExtractVariablesWalkOption() WalkOption {
	return func(walkOptions *WalkOptions) {
		walkOptions.Ante = func(depth int, me MapperNode) (bool, error) {
			if charData, ok := me.(*ChardataNode); ok {
				if p, e := NewStatementParamParser(); e != nil {
					return false, e
				} else {
					if charData.DollarVariables, e = p.GetDollarVariables(charData.Chardata); e != nil {
						return false, e
					}

					if charData.DashVariables, e = p.GetDashVariables(charData.Chardata); e != nil {
						return false, e
					}
				}
			}
			return true, nil
		}
	}
}

func WithFlattenIncludeVisitFunction() WalkOption {
	return func(walkOptions *WalkOptions) {
		walkOptions.Ante = func(depth int, mapper MapperNode) (bool, error) {

			resolved := true
			for resolved == true {
				newNodes := make([]MapperNode, 0, len(mapper.GetNodes()))
				resolved = false
				for _, n := range mapper.GetNodes() {
					if include, ok := n.(*IncludeNode); ok {
						resolved = true
						sqlNodes := include.Nodes[0].GetNodes()
						newNodes = append(newNodes, sqlNodes...)
					} else {
						newNodes = append(newNodes, n)
					}
				}

				if resolved {
					mapper.SetNodes(newNodes)
				}
			}

			return true, nil
		}
	}
}
