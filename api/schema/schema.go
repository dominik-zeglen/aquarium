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

var _api_schema_schema_graphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x54\xcb\x6e\xdc\x30\x0c\xbc\xfb\x2b\x94\xdb\xe6\x17\x74\xdb\x2e\x9a\x76\x81\x02\x4d\xd1\x00\x3d\x04\x3d\x68\x6d\xc6\x26\xa0\x87\x4b\xd1\xe9\x2e\x8a\xfc\x7b\x21\x59\x96\xe5\x07\xf6\x46\x53\xd4\x70\x66\x48\x19\x6d\x3f\xb0\x78\x76\x68\xf9\x1c\xc3\x7f\x95\x10\x57\x29\x9e\xb4\x53\xfc\x50\x09\x71\xcb\xf1\x47\x55\x8d\xd5\x47\x02\x35\x17\x7b\x56\xc4\xb2\x80\x08\xb7\xc0\x36\xeb\x94\xaf\x95\x06\x29\xce\x96\x67\xa4\xef\xd4\x2a\x8b\xde\x3c\xa1\x66\xa0\x08\xa7\x08\x94\x9c\x5b\x84\x5a\xbe\xf5\x30\x82\xdd\x67\x17\xeb\xce\x0c\xa4\x18\x9d\x7d\x26\x57\x13\xc4\x30\x5e\xab\x55\x4e\x81\x14\x9f\x9c\xd3\xa0\x6c\xc0\x30\xea\x7a\x6a\x22\xb1\xf4\xf5\x15\xb0\xed\xb8\xe8\x62\xd0\x96\x15\x68\x37\x15\xbe\x87\x1a\xc1\x4b\xf1\xfa\x73\x8c\x1e\x7e\xef\x50\xfa\xa5\x3c\x43\x24\x63\xd4\xf5\xc5\x69\x20\x65\x6b\x58\x36\xda\x4b\xb3\xbb\x62\x8d\x7c\x47\xec\x68\x9d\xc6\x77\x38\x81\xd6\x27\x37\x58\xce\x7c\xeb\x4d\xc6\x0e\xe6\x02\x94\x3f\xfb\xd9\x29\xb9\xeb\x5f\x28\xfa\x1b\xb8\xcb\x95\x96\x99\x49\x92\x1d\x79\x60\xb3\xe8\xfd\x72\xeb\xa3\x33\xa7\x14\x07\x6b\x84\x68\x10\x38\xd8\xc5\x84\xb6\x1d\x53\x60\x80\x5a\x68\x8e\x05\x51\x65\x40\x8a\x54\x53\x09\xe1\xd2\xbe\x04\xb8\x69\x77\x4a\xa7\xa7\xdc\x9a\xc6\xc5\x91\x3d\x2e\x1d\x99\x18\x8d\xad\x7b\xe7\x71\xd4\x1f\xf7\x6c\x31\xd2\x69\xa2\xb9\xcb\x24\x64\xdd\x65\x47\xd2\xdb\x60\xdb\x0e\x73\x45\x07\x74\xc1\x77\x47\x90\x32\x25\xe2\x1a\x2d\x4e\x73\xb1\xa8\x5d\x3f\x8f\x6c\xcb\x37\x00\xc9\xcc\x6d\x33\x99\x2f\x84\xcd\x67\x0d\x06\xd2\x2b\xba\xa7\x78\x6f\x89\x7f\x0c\x40\xb7\x78\x73\x1a\xc2\x61\x22\xfb\x28\xb3\xf1\xc5\xf1\x37\xf4\x7c\x78\x8b\x0f\x5b\xae\x1e\xfa\xe3\x6a\x7c\x73\xef\x12\x33\x91\x98\x0f\x03\xe2\x92\x5c\x3e\x0a\xea\x0e\xab\x5f\x47\xc0\x78\xdd\xaa\x4f\x0d\x71\x5a\xe4\x62\xa7\xa3\x5a\x5f\x77\x60\x54\x54\xfa\x27\x68\x96\xa3\xf4\xea\xa3\xfa\x1f\x00\x00\xff\xff\xf1\x23\xaa\xb7\x2c\x05\x00\x00")

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
