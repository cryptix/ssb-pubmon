package assetfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"

	"github.com/cryptix/go/logging"
	"github.com/pkg/errors"
)

// Interface assetfs interface
type Interface interface {
	PrependPath(path string) error
	RegisterPath(path string) error
	Asset(name string) ([]byte, error)
	Glob(pattern string) (matches []string, err error)
	Compile() error

	NameSpace(nameSpace string) Interface
}

// AssetFS default assetfs
var assetFS Interface = &AssetFileSystem{}
var used bool

// AssetFS get AssetFS
func AssetFS() Interface {
	used = true
	return assetFS
}

// SetAssetFS set assetfs
func SetAssetFS(fs Interface) {
	if used {
		os.Mkdir("panics", os.ModePerm)
		b, tmpErr := ioutil.TempFile("panics", "assetFSwarnUsed")
		if tmpErr != nil {
			panic(errors.Wrap(tmpErr, "failed to create setAssetFS log"))
		}
		fmt.Fprintf(b, "warning! SetAssetFS called after use!")
		fmt.Fprintf(b, "Stack:\n%s", debug.Stack())

		logging.Logger("qor/SetAssetFS").Log("event", "assetFSwarnUsed", "panicLog", b.Name())

		if err := b.Close(); err != nil {
			panic(errors.Wrap(err, "failed to close setAssetFS log"))
		}
	}

	assetFS = fs
}
