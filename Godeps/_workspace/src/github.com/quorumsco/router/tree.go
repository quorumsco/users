package router

import (
	"net/http"
	"reflect"
)

type node struct {
	segment  segment
	handler  http.Handler
	children []*node
}

func newNode(key uint8) *node {
	return &node{
		children: make([]*node, 0),
	}
}

func findRec(n *node, path string, params []Param) (http.Handler, bool, []Param) {
	switch {
	case len(path) == 0:
		return n.handler, n.handler != nil, params
	default:
		for _, child := range n.children {
			remainingPath, newParam, matched := child.segment.Match(path)
			if matched {
				var (
					handler        http.Handler
					found          bool
					returnedParams []Param
				)
				if newParam.Name != "" {
					handler, found, returnedParams = findRec(child, remainingPath, append(params, newParam))
				} else {
					handler, found, returnedParams = findRec(child, remainingPath, params)
				}
				if found {
					return handler, found, returnedParams
				}
			}
		}
	}
	return nil, false, nil
}

func (n *node) find(path string) (http.Handler, bool, []Param) {
	return findRec(n, path, make([]Param, 0, 10))
}

func insertRec(n *node, segments []segment, handler http.Handler, params uint8) uint8 {
	switch {
	case len(segments) == 0:
		n.handler = handler
		return params
	default:
		for _, child := range n.children {
			if reflect.TypeOf(child.segment) == reflect.TypeOf(segments[0]) && child.segment.Is(segments[0]) {
				return insertRec(child, segments[1:], handler, params+child.segment.Params())
			}
		}
		newNode := newNode(0)
		newNode.segment = segments[0]
		n.children = append(n.children, newNode)
		return insertRec(newNode, segments[1:], handler, params+newNode.segment.Params())
	}
}

func (n *node) insert(path string, handler http.Handler) {
	segments := parseSegments(path)
	params := insertRec(n, segments, handler, 0)
	if params > 10 {
		//n.maxParams = params
	}
}
