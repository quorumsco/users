package router

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func funcEqual(a, b interface{}) bool {
	av := reflect.ValueOf(&a).Elem()
	bv := reflect.ValueOf(&b).Elem()

	return av.InterfaceData() == bv.InterfaceData()
}

func TestNodeInsert(t *testing.T) {
	var f http.HandlerFunc = func(w http.ResponseWriter, req *http.Request) {}
	node := newNode(0)
	node.insert("/test/itsme", f)
	node.insert("/test/itsmedarling", f)
	node.insert("/test/:yolo/test", f)

	g, found, p := node.find("/test/findme/test")
	assert.True(t, found, "Path not found after being inserted")
	assert.True(t, funcEqual(f, g), "Node has wrong value")
	assert.True(t, len(p) == 1, "Wrong number of params returned")
	assert.Equal(t, param{name: "yolo", value: "findme"}, p[0], "Wrong param returned")

	g, found, p = node.find("/test/itsme")
	assert.True(t, found, "Path not found after being inserted")
	assert.True(t, funcEqual(f, g), "Node has wrong value")
	assert.True(t, len(p) == 0, "Wrong number of params returned")

	g, found, p = node.find("/test/itsmedarling")
	assert.True(t, found, "Path not found after being inserted")
	assert.True(t, funcEqual(f, g), "Node has wrong value")
	assert.True(t, len(p) == 0, "Wrong number of params returned")
}
