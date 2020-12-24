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

var _api_schema_schema_graphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x92\xc1\x6a\xfb\x30\x0c\xc6\xef\x7e\x0a\x95\x5e\x5a\xf8\x3f\x81\x6f\xff\x95\x0d\x0a\x3b\x6c\xec\x38\x7a\xd0\x6c\x91\x0a\x5c\x39\x8b\xd5\xb1\x30\xfa\xee\xc3\x76\x9b\x25\xcd\x18\xbd\x39\x9f\x7e\xb2\xbe\x4f\x31\x4b\x7b\x54\x78\x8a\x2c\xba\x2d\xc7\x2f\x03\xf0\x69\xe1\x21\x44\xd4\x85\x01\xe8\x87\xf3\xc9\x18\xed\x5b\xaa\xf0\x0d\xdc\x4b\x4b\x8e\x29\x6d\xa2\x08\x39\xe5\x28\xa5\xc7\xc5\xa3\xa8\x85\xad\x94\x2e\xf2\x0d\x25\x0b\xaf\x33\xf6\xde\x37\xb4\xd8\xfd\x71\x59\x06\xca\x85\x12\x3d\xd9\x4b\x7d\xd6\x50\x10\xf6\xc3\x40\xc1\x43\xa6\xb5\x63\x69\x8a\x81\x03\x75\x0d\xf9\xff\x3f\x96\x3c\x93\x66\x47\x15\xd9\x65\xc9\x51\x08\xd9\xe5\x86\x42\x18\x9b\xca\xdf\x37\xc5\x9b\x82\xd7\xd9\xe6\xd5\x51\xb0\x32\x73\x82\x5e\x47\xc2\xc0\x1f\x64\xe1\x2e\xc6\x40\x28\x59\x79\x8b\x9d\x8c\x12\x39\x6c\xd1\xb1\xf6\x83\xb0\x6f\x87\x63\x1b\x13\xe7\xa9\xb6\xfe\xd6\x2c\x25\x54\xc6\xaa\x9d\xa1\x25\xa4\xba\xcd\x5f\xf6\xfc\x7c\xa4\xae\x2f\x96\x96\x80\x1d\xe1\x2a\x29\x76\x6a\x47\x4f\x6a\xf1\x0f\x48\xfc\x44\x59\x8f\x96\x59\xd7\xbb\xba\x04\x5a\xd7\xc8\xe3\xa9\xb9\x56\x4b\xe7\xe9\xe6\xdc\xf4\xc8\x49\xed\xd5\xfa\x26\x7e\x2b\x30\x7b\x3b\xc5\x7d\x72\x7b\x3a\x60\x71\xfe\x9e\x33\xd8\x1a\xc5\x9c\xcc\x77\x00\x00\x00\xff\xff\x35\xfe\x11\x0d\x14\x03\x00\x00")

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
