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

var _api_schema_schema_graphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x54\xc1\x8e\xd4\x30\x0c\xbd\xf7\x2b\x32\xb7\x5d\x89\x2f\xc8\x0d\x46\x20\x46\xe2\xb0\x08\x24\x0e\x88\x83\x27\x35\xad\xa5\xd6\x2e\xa9\xbb\xcc\x08\xed\xbf\xa3\x24\x9d\x36\x69\x47\xab\xbd\x25\xee\xb3\xfd\xde\x8b\x5d\xe2\x61\x52\xf3\x24\xc4\x7a\x8a\xc7\x7f\x95\x31\x17\x6b\x3e\x75\x02\x7a\xa8\x8c\xb9\x2e\xe7\x97\xaa\xd2\xeb\x80\x09\xfc\x06\xdc\x49\xd1\x83\x92\xf0\x93\x17\xe7\x31\x1e\x63\x9a\x83\x25\x84\xd6\x7c\x10\xe9\x10\x38\xd4\xe8\xe1\x72\xac\xad\x39\xb1\xce\xb7\xcf\x48\x4d\xab\x59\x97\x9e\x38\x47\x10\xef\x10\xe3\x80\x8e\x70\xb4\xe6\x5b\x3a\x1c\x85\x19\x5d\xe8\x7d\x87\xd9\x0f\x18\x15\x23\xa7\x1e\x2e\xdf\xa5\x43\x0f\xec\xb0\xec\x77\x2f\xac\x72\x21\x47\xfa\x8a\xe6\x58\x14\x3a\x7a\xc6\x23\x76\xdd\x51\x26\xd6\x85\xb6\xdb\x45\x78\xea\xcf\xe8\x97\xeb\xb0\x1a\x66\xef\xda\x18\x40\x7f\x03\x77\xbb\xd1\xb2\x32\xd9\xc9\x4f\xd6\x17\x6d\xb1\x6e\x82\x53\x3f\x77\xd8\x8f\x75\x83\x87\x5f\xaf\x14\x0b\x80\x58\x90\xa5\xc6\xc5\xeb\x5d\x42\x84\xd0\xfa\x60\x0e\x3c\xd3\xb3\x78\x2c\xbc\x08\x14\x82\x4b\xa1\x63\xa0\xc8\xe3\xd4\x0f\xb3\xf8\x84\xaa\x09\x35\xf0\x54\x4f\xdc\x24\x18\xf6\xe8\x1b\xac\xdf\xaf\x6a\x7e\x4f\xdc\xb4\xb4\x5c\x5b\xf4\xe7\xb2\x17\x43\x1f\xb8\xa6\x22\x0b\xd5\xf4\x3e\x6f\x70\xa9\x04\x6e\x2d\xda\x7f\xcd\xfc\x89\xea\x0a\xe8\xd6\x99\x38\x2a\xc5\x32\x9c\xc5\x73\xa6\xce\xc1\x00\x69\xe4\x6e\x72\x45\xd6\xf4\x76\x58\x87\x47\x46\x4a\xe6\xc5\x45\x8d\x3b\x01\x4a\x50\x18\xba\xdd\x92\x95\xdc\xd7\x09\xfd\x35\x4d\xaf\x47\x78\x18\x15\xbc\xda\xec\x07\x71\x78\x67\x90\xeb\x22\xf2\x98\xbd\xdf\xfc\xa4\x0f\x37\x69\x8f\x49\xfc\x1c\xfe\x42\xa3\xda\x8d\x55\x31\x67\xe6\x93\xa7\xcd\xcc\xd6\x8f\x29\xf9\xce\x5a\x07\x27\x6f\x5b\x90\x2d\x44\xd4\x34\xba\x16\x7b\x88\x7a\xfe\x04\x65\x36\x09\xac\x5e\xaa\xff\x01\x00\x00\xff\xff\xeb\x75\x06\x16\xfa\x04\x00\x00")

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
