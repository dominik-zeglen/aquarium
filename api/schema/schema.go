package schema

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

var _api_schema_schema_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func api_schema_schema_go() ([]byte, error) {
	return bindata_read(
		_api_schema_schema_go,
		"api/schema/schema.go",
	)
}

var _api_schema_schema_graphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x52\xcb\x6a\xeb\x30\x10\xdd\xeb\x2b\x64\xee\x26\x81\xfb\x05\xda\xb5\xa1\x85\x40\x17\x2d\x5d\x96\x2e\x14\x6b\x6a\x0f\xc8\x23\x55\x1a\x87\x9a\x92\x7f\x2f\x92\x12\x3f\xe2\x52\xb2\x93\x66\xce\xe8\x3c\x34\x48\xbe\x67\xf9\xec\x90\x78\x9f\x8f\xdf\x42\xca\x2f\x25\x1f\xad\xd3\x5c\x09\x29\x87\xf1\x7c\x12\x82\x07\x0f\x05\x7c\x03\xee\xd5\x43\x8d\x10\x77\x8e\x08\x6a\x46\x47\x79\xa6\x76\x3d\xb1\x92\x7b\xca\x53\x60\x1a\x88\x4a\xbe\xad\xb0\x0f\xa6\x81\xea\xfd\x8f\xc7\x12\x20\x3f\x48\xce\x80\xba\xf4\x57\x03\x19\x82\x66\x24\xac\x75\x20\x3c\xba\x00\x53\x05\xac\x4d\x12\x76\x60\x6d\x62\x4c\x12\x29\xf6\x9d\x4f\x2c\x23\xca\x20\x70\xd2\xc9\x01\xa9\x29\x30\xe8\x20\x34\x60\xee\x26\x37\x1f\x3d\x35\x2d\x8e\xd7\x16\xc2\x61\xc9\x45\xba\x4b\x5a\xcb\x23\xa3\xd4\xc4\x7c\x53\x4a\x4b\xe0\x75\x44\xeb\xee\x2c\x9f\xec\x6e\x01\xbd\x4e\x46\x5b\x3c\x82\x92\xf7\xce\x59\xd0\x94\x2a\x07\x17\x68\xe6\xae\xd6\x5e\xd7\xc8\xc3\x64\xd7\xb9\x69\xbc\xf5\xe3\xd1\xbb\x88\x25\xbc\xbc\x2a\xa9\x14\x35\xa3\x5e\x04\x1a\xcb\xff\xfc\xf2\x73\x2f\x3d\x84\x21\xab\xd3\x01\xf4\x26\xb2\x0e\xac\x66\x2b\x5a\xfd\x97\x40\x66\x51\xd9\x2e\xff\x0f\xac\xdd\x5c\x9c\x6d\x8b\x77\x21\xe5\xbf\x0b\x67\xea\x95\xd6\x99\x5b\x9c\x87\x9e\x30\xb2\xba\xca\xb1\x9a\x4f\x16\xc0\x6a\x17\xb3\xf6\x58\xb7\xd0\xe9\xac\xfb\x33\x39\x50\xc5\x88\x38\x89\x9f\x00\x00\x00\xff\xff\xcd\xca\xff\xe3\x64\x03\x00\x00")

func api_schema_schema_graphql() ([]byte, error) {
	return bindata_read(
		_api_schema_schema_graphql,
		"api/schema/schema.graphql",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"api/schema/schema.go": api_schema_schema_go,
	"api/schema/schema.graphql": api_schema_schema_graphql,
}
// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"api": &_bintree_t{nil, map[string]*_bintree_t{
		"schema": &_bintree_t{nil, map[string]*_bintree_t{
			"schema.go": &_bintree_t{api_schema_schema_go, map[string]*_bintree_t{
			}},
			"schema.graphql": &_bintree_t{api_schema_schema_graphql, map[string]*_bintree_t{
			}},
		}},
	}},
}}
