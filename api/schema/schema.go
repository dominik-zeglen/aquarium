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

var _api_schema_schema_graphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x54\xc1\x8e\xd3\x30\x10\xbd\xe7\x2b\xdc\xdb\xae\xc4\x17\xf8\x06\x15\x88\x4a\x1c\x16\x81\xc4\x01\x71\x98\x3a\x43\x32\x52\x3c\x0e\xf6\x64\x69\x85\xfa\xef\xc8\x76\x1a\xc7\x4d\x55\xed\xcd\x19\xbf\x99\x79\xef\x79\x26\xc4\xe3\x24\xea\xc5\x11\xcb\x21\x1d\xff\x35\x4a\x9d\xb4\xfa\x34\x38\x90\x5d\xa3\xd4\x79\x39\x5f\x9a\x46\xce\x23\x66\xf0\x1b\x70\x07\x41\x0f\x42\x8e\x5f\xbc\x33\x1e\xd3\x31\xa5\x19\x58\x42\xa8\xd5\x07\xe7\x06\x04\x8e\x35\x2c\x9c\xf6\xad\x56\x07\x96\xf9\xeb\x33\x52\xd7\xcb\xaa\x8b\x25\x5e\x23\x88\x37\x88\x30\xa2\x21\x0c\x5a\x7d\xcb\x87\xbd\x63\x46\x13\x7b\xdf\x61\xf6\x03\x82\x60\xe2\x64\xe1\xf4\xdd\x0d\xe8\x81\x0d\xd6\xfd\xee\x85\xc5\x9d\xc8\x90\x3c\xd0\x9c\x8a\xc2\x40\xaf\xb8\xc7\x61\xd8\xbb\x89\x65\xa1\x6d\x36\x11\x9e\xec\x11\xfd\xf2\x39\x16\xc3\xf4\x5d\x1b\x23\xe8\x6f\xe4\xae\x6f\xb4\x14\x26\x1b\xf9\xd9\xfa\xaa\x2d\xb6\x5d\x74\xea\xe7\x06\xfb\xb1\xed\x70\xf7\xeb\x41\xb1\x08\x48\x05\xd9\xb5\xb8\x78\xbd\x49\x48\x10\x2a\x0f\x66\xc0\x33\xbd\x3a\x8f\x95\x17\x41\xab\x6c\x52\x79\xa9\xc8\x94\xc3\x64\xc7\xd9\x83\x0c\x6e\x09\x25\xd2\x15\x4f\xdc\x45\x7e\x4a\xa1\x45\xdf\x61\xfb\xbe\x88\xfa\x3d\x71\xd7\xd3\xf2\xd9\xa3\x3f\xd6\x2d\x19\x6c\xa4\x9c\x8b\x2c\x8c\x6b\x06\x0f\xcc\xaa\x81\xb7\x4e\x6d\x6f\x57\x36\xc5\xcb\x1a\x7a\x6b\x50\x9a\x98\x6a\x27\x8e\xce\xf3\x4a\x9d\x81\x11\xf2\xe4\x5d\xe5\x3a\x57\xd2\xfb\xb1\xcc\x90\x0b\x94\xcd\x4b\xfb\x9a\x56\x03\x84\xa0\x32\xf4\x76\x59\x0a\xb9\xaf\x13\xfa\x73\x1e\x62\x8f\xf0\x14\x04\xbc\xe8\xd5\x7f\x62\xf7\x4e\x21\xb7\x55\xe4\x79\x36\x27\xda\x31\xbf\xec\xd3\x55\xda\x73\x16\x3f\x87\xbf\x50\x90\xed\x9b\x17\x3e\xeb\xb4\x99\x59\xb9\xcc\xc9\x77\xb6\x3b\x3a\x79\x5d\x86\xd5\x5e\x24\x4d\xc1\xf4\x68\x21\xe9\xf9\x13\x95\xe9\x2c\xb0\xb9\x34\xff\x03\x00\x00\xff\xff\x6c\x1f\x40\x3e\x01\x05\x00\x00")

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
