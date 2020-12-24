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

var _api_schema_schema_graphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x52\xcb\x6a\x03\x31\x0c\xbc\xfb\x2b\x14\x72\x49\xa0\x5f\xe0\x5b\x1b\x5a\x08\xf4\xd0\x92\x63\xc9\x41\xb5\xc5\x46\xe0\xc8\xdb\xb5\x52\xba\x94\xfc\x7b\xb1\x9d\x6c\x36\x0f\x4a\x6e\xda\xd1\x48\x33\xa3\x35\x4b\xbb\x53\x78\x8b\x2c\xba\x2c\xe5\xaf\x01\xf8\xb1\xf0\x12\x22\xea\xc4\x00\xf4\x43\xbd\x37\x46\xfb\x96\x2a\xf9\x0e\xde\xaa\x25\xc7\x94\x16\x51\x84\x9c\x72\x94\x32\xe3\xe2\x4e\xd4\xc2\x52\xca\x14\xf9\x86\x92\x85\x8f\x2b\xee\xb3\x6f\x68\xb2\xfe\x67\x59\x26\x94\x85\x12\x3d\x9d\x36\xdc\x98\x29\x2c\xf6\x83\xa6\xe0\x96\x2c\xac\xb4\x63\x69\x8a\x87\x2d\x75\x0d\xf9\xc7\x93\x2b\xcf\xa4\x79\x65\xa5\xac\x33\xe4\x28\x84\x6c\x74\x41\x21\x8c\x35\xf2\xf7\x5d\x09\xcf\x89\x97\xf1\xae\xbb\xe3\x6c\xb7\x44\x2f\x53\x61\xe0\x6f\xb2\xf0\x14\x63\x20\x94\x8c\x7c\xc6\x4e\x46\xa1\x1c\xb6\xe8\x58\xfb\x01\xd8\xb4\x43\xd9\xc6\xc4\x59\xd8\xd6\x9f\x9b\xa1\x84\xca\x58\xb1\x03\x69\x0a\xa9\x1e\xd4\x1e\x2f\x7b\x72\xf4\xbe\xa3\xae\x2f\x96\xa6\x80\x1d\xe1\x2c\x29\x76\x6a\x47\x0f\x6b\xf2\x00\x24\xfe\x0c\x99\x8f\xa2\xd5\x0b\xcf\x8e\x81\xe6\xb6\x84\x1c\xab\xe6\x5e\x6d\x1d\xd4\x4d\xe9\xe6\xb1\x57\x4e\x6a\x2f\x6e\x78\xe6\xb8\x12\xae\xde\x50\xf1\x9f\xdc\x86\xb6\x58\xbc\x7f\xe5\x14\xb6\x86\x31\x7b\xf3\x17\x00\x00\xff\xff\xd3\x84\x0e\xe7\x1c\x03\x00\x00")

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
