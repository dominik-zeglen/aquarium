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

var _api_schema_schema_graphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x92\xc1\x4a\x03\x31\x10\x86\xef\x79\x8a\x14\x2f\x2d\xf8\x04\xb9\x69\x51\x28\x78\x50\x3c\x4a\x0f\x63\x32\x6c\x07\xd2\xc9\x9a\x4c\xc5\x45\xfa\xee\x92\xa4\x5d\xb3\x5d\x91\xde\xb2\xff\x7c\x93\xf9\xff\xd9\x10\xf7\x07\xd1\xcf\x81\x58\x36\xe5\xf8\xad\xb4\xfe\x32\xfa\xd1\x07\x90\x85\xd2\x7a\x18\xcf\x47\xa5\x64\xe8\xb1\xc2\x57\x70\xaf\x3d\x5a\xc2\xb4\x0e\xcc\x68\x85\x02\x97\x1e\x1b\x0e\x2c\x46\x6f\xb8\x74\xa1\xeb\x30\x19\xfd\x36\x63\x1f\x5c\x87\x8b\xed\x3f\x97\x65\xa0\x5c\xc8\xc1\xa1\x39\xd7\x67\x0d\x05\x21\x37\x0e\x64\xd8\x67\x5a\x22\x71\x57\x0c\xec\x31\x76\xe8\xee\x7e\x2d\x39\x42\xc9\x8e\x2a\xb2\xcd\x92\x45\xef\xb3\xcb\x35\x7a\xdf\x9a\xca\xdf\x57\xc5\x9b\x82\x97\xd9\xe6\xd5\x26\x58\x99\x39\x41\x2f\x23\x81\xa7\x4f\x34\xfa\x3e\x04\x8f\xc0\x59\x79\x0f\x91\x9b\x44\x16\x7a\xb0\x24\xc3\x28\xec\xfa\xf1\xd8\x87\x44\x79\xaa\xa9\xbf\x35\x4b\x09\x84\xa0\x6a\x27\xe8\x46\xa7\xba\xcd\x3f\xf6\xfc\x72\xc0\x38\x14\x4b\x10\x11\x96\x49\x20\x8a\x69\x1e\xd4\xe2\x56\x23\xbb\x89\xb2\x6a\x56\x59\x97\xbb\x3c\xc7\x59\xd5\xc0\xed\xcc\x5c\xab\xa5\xd3\x6c\x75\x6a\x7a\xa2\x24\xe6\x62\x79\x13\xb7\x15\x98\xbd\x9c\xe2\x3d\xd9\x1d\xee\xa1\xf8\xfe\xc8\x09\x4c\x0d\xa2\x8e\xea\x27\x00\x00\xff\xff\x14\xfc\x03\x9c\x12\x03\x00\x00")

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
