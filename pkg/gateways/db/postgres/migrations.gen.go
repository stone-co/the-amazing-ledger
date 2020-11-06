// Code generated by go-bindata. DO NOT EDIT.
// sources:
// migrations/000001_initialize_schema.down.sql (110B)
// migrations/000001_initialize_schema.up.sql (1.208kB)

package postgres

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var __000001_initialize_schemaDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x72\x75\xf7\xf4\xb3\xe6\xe2\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\x48\xcd\x2b\x29\xca\x4c\x2d\xb6\x86\x8a\x45\x06\xb8\x2a\x28\xe4\x17\xa4\x16\x25\x96\x64\xe6\xe7\xc5\x97\x54\x16\xa4\x5a\x23\x2b\x4f\x4c\x4e\xce\x2f\xcd\x2b\x41\x55\x0f\x15\x84\xaa\xe6\x72\xf6\xf7\xf5\xf5\x0c\xb1\xe6\x02\x04\x00\x00\xff\xff\x00\xaf\xc8\x35\x6e\x00\x00\x00")

func _000001_initialize_schemaDownSqlBytes() ([]byte, error) {
	return bindataRead(
		__000001_initialize_schemaDownSql,
		"000001_initialize_schema.down.sql",
	)
}

func _000001_initialize_schemaDownSql() (*asset, error) {
	bytes, err := _000001_initialize_schemaDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000001_initialize_schema.down.sql", size: 110, mode: os.FileMode(0644), modTime: time.Unix(1603896993, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x11, 0x4a, 0x6d, 0xf3, 0xd4, 0x4f, 0x9f, 0xb1, 0xe, 0xa5, 0x3d, 0x3d, 0x9, 0xf9, 0x9c, 0x77, 0x57, 0x4c, 0x8f, 0xbb, 0x2c, 0x8e, 0xcc, 0x64, 0x72, 0x6f, 0xd2, 0x6a, 0x64, 0x5a, 0x1d, 0xed}}
	return a, nil
}

var __000001_initialize_schemaUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x92\xc1\x8e\x9b\x30\x18\x84\xcf\xf0\x14\xff\x2d\xbb\x52\xde\x20\x27\x42\xdc\x95\x55\x30\x11\x31\xd2\xa6\x17\xe4\x60\xb7\xb2\x0a\x86\xda\xa6\xdd\xbc\x7d\x85\x09\xc8\xb0\xb4\xdd\x72\x72\xf4\x4f\x3e\x8f\x67\xfe\x23\x7a\xc1\xe4\x10\x86\x71\x8e\x22\x8a\x80\x5e\xcf\x08\x58\x55\xb5\xbd\xb2\xa5\xbd\x77\x02\xa2\x0b\x20\x52\xa4\xf0\x14\x06\x3b\x66\x8c\xb0\xbb\x7d\x18\xec\x6a\xc9\x6e\xb2\x96\xf6\xee\x7e\x49\x55\xb5\x8d\x70\x47\xf1\xd6\x09\x65\x1e\xe7\x1f\xfd\xa0\x08\x9f\x3d\x7e\x74\x4c\xe6\x0b\xcc\x00\x95\x1c\xa6\xaf\x28\xf0\x09\xce\x39\x4e\xa3\xfc\x0a\x9f\xd1\x75\x1f\x06\xce\xc2\xf8\x45\x71\x9c\x15\x84\x96\xce\x22\xc9\x28\x90\x22\x49\xf6\x61\xd0\xfe\x52\x42\x97\x0e\x43\xd1\x2b\x7d\x37\x1a\xff\xbd\x1e\x29\xd6\x4c\xe0\xf5\xa8\x11\x96\x71\x66\xd9\xc6\xa8\xd2\x82\x59\xc1\x4b\x66\x81\xe2\x14\x5d\x68\x94\x9e\xe9\x97\x59\x01\x27\xf4\x29\x2a\x12\x0a\x71\x91\xe7\x68\xf0\x3a\x89\x16\x19\xf8\x19\x57\x35\x33\x66\x11\xf2\x32\x5a\x17\xb9\xf9\xdf\x94\x87\x1b\xda\x4e\x68\x66\x65\xab\xde\xf7\xc8\xc5\x4d\x8e\x3d\x56\x5a\x70\x69\x37\x2a\x12\xca\x6a\x29\xd6\x0d\x79\x3d\x75\x5a\x36\x4c\xdf\xe1\xbb\xb8\xef\xc3\x60\xf9\x1c\xaf\xad\x38\x89\x2e\x17\x50\xad\x05\xd5\xd7\xb5\xa7\xfc\xa6\xdb\xbe\x9b\xe3\xdf\x10\x98\xfe\x36\x6a\xfe\x24\x98\x7d\xad\x05\xf3\xcb\x27\xc3\xd9\x19\xe5\x11\xc5\x19\x19\x97\xc7\x67\x35\x03\xca\x7b\x1b\x26\x0b\xd4\x4f\xa1\x8d\x07\x02\x38\xe2\x97\x95\xc4\x6a\xa6\x0c\xab\x5c\xd2\xce\x92\x8b\xc7\x13\x78\x4b\xf3\xf0\xeb\xad\xce\xa4\x03\x2e\xbe\xb2\xbe\xb6\x7f\x5f\x1d\x4c\x4e\xe8\x15\x24\x7f\x2b\x1f\xfd\x94\xcb\xe0\x33\x32\x15\xf7\xb4\x18\x3c\x1f\xfe\x4d\x18\xc3\xde\x20\xb8\xc1\x47\x08\x73\x63\x1b\x90\x69\xf6\x11\x8e\xe4\x5b\x04\xc9\x5d\x10\x59\x9a\x62\x7a\x08\x7f\x07\x00\x00\xff\xff\x54\xd6\x5a\xee\xb8\x04\x00\x00")

func _000001_initialize_schemaUpSqlBytes() ([]byte, error) {
	return bindataRead(
		__000001_initialize_schemaUpSql,
		"000001_initialize_schema.up.sql",
	)
}

func _000001_initialize_schemaUpSql() (*asset, error) {
	bytes, err := _000001_initialize_schemaUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "000001_initialize_schema.up.sql", size: 1208, mode: os.FileMode(0644), modTime: time.Unix(1604966207, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x55, 0x42, 0x73, 0xe4, 0xd4, 0x9f, 0x7d, 0x37, 0x45, 0xcb, 0xb9, 0x81, 0xa, 0xab, 0x32, 0xc0, 0xe4, 0xdb, 0xf, 0xac, 0xd2, 0xf0, 0x1d, 0x46, 0xd8, 0x44, 0x20, 0xf5, 0xa1, 0xc6, 0xd0, 0xe}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
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
var _bindata = map[string]func() (*asset, error){
	"000001_initialize_schema.down.sql": _000001_initialize_schemaDownSql,
	"000001_initialize_schema.up.sql":   _000001_initialize_schemaUpSql,
}

// AssetDebug is true if the assets were built with the debug flag enabled.
const AssetDebug = false

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
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
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"000001_initialize_schema.down.sql": &bintree{_000001_initialize_schemaDownSql, map[string]*bintree{}},
	"000001_initialize_schema.up.sql":   &bintree{_000001_initialize_schemaUpSql, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory.
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
