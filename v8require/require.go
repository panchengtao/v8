package v8require

import (
	"encoding/json"
	"errors"
	"github.com/panchengtao/v8"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

//const modulePath = "/home/panchengtao/go/src/v8/node_modules"

var moduleCache = make(map[string]*v8.Value)

func load(in v8.CallbackArgs, packageName string, require *RequireContext) (*v8.Value, error) {
	// 计算绝对路径
	abPath, err := resolvePackageName(in, packageName, require)
	if err != nil {
		return nil, err
	}

	// 如果有缓存，取出缓存
	if moduleCache[abPath] != nil {
		return moduleCache[abPath], nil
	}

	// 手动生成 module 和 module.exports 对象，用于取代 package 中的 module, module.exports 对象
	pkgEntryCodeByte, _ := ioutil.ReadFile(abPath)
	pkgEntryCode := "function __wrapperModule(module,exports){\n" + string(pkgEntryCodeByte) + "\n};__wrapperModule;"
	moduleFunc, err := in.Context.Eval(pkgEntryCode, abPath)
	if err != nil {
		return nil, err
	}

	module, err := in.Context.Eval("var __module={exports:{}};__module;", "")
	if err != nil {
		return nil, err
	}
	exports, err := module.Get("exports")
	if err != nil {
		return nil, err
	}
	_, err = moduleFunc.Call(nil, module, exports)
	if err != nil {
		return nil, err
	}
	exports, _ = module.Get("exports")

	// 更新缓存
	moduleCache[abPath] = exports

	return exports, nil
}

func resolvePackageName(in v8.CallbackArgs, packageName string, require *RequireContext) (string, error) {
	return resolveLookupPath(in, packageName, require)
}

func resolveLookupPath(in v8.CallbackArgs, packageName string, require *RequireContext) (string, error) {
	var abPath string
	if strings.HasPrefix(packageName, "./") {
		abPath = path.Join(path.Dir(in.Caller.Filename), packageName)
		if !strings.HasSuffix(abPath, ".js") && !strings.HasSuffix(abPath, ".json") && !strings.HasSuffix(abPath, ".node") {
			file, err := os.Stat(abPath)
			if err != nil {
				abPath += ".js"
			} else {
				if file.IsDir() {
					abPath = path.Join(abPath, "index.js")
				}
			}
		}

		return abPath, nil
	} else {
		var searched bool
		for modulePath := range require.path {
			packageJsonPath := path.Join(modulePath, packageName, "package.json")
			data, _ := ioutil.ReadFile(packageJsonPath)
			var packageJson struct {
				Main string `json:"main"`
			}
			err := json.Unmarshal(data, &packageJson)
			if err != nil || packageJson.Main == "" {
				abPath = path.Join(modulePath, packageName, "index.js")
			} else {
				abPath = path.Join(modulePath, packageName, packageJson.Main)
			}

			_, err = os.Stat(abPath)
			if err != nil {
				continue
			} else {
				searched = true
				break
			}
		}

		if searched {
			return abPath, nil
		}

		return "", errors.New("unable to locate path of " + packageName)
	}
}

type RequireContext struct {
	*v8.Context
	path map[string]bool
}

func (ctx *RequireContext) RegisterModulePath(path string) error {
	if _, ok := ctx.path[path]; !ok {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return err
		}

		if ctx.path == nil {
			ctx.path = make(map[string]bool)
		}
		ctx.path[path] = true
	}

	return nil
}

func EnableRequiring(ctx *v8.Context) (*RequireContext, error) {
	require := &RequireContext{Context: ctx}
	f := ctx.Bind("require", func(in v8.CallbackArgs) (*v8.Value, error) {
		if len(in.Args) != 1 {
			return nil, errors.New("invalid arguments count for require, expected to be 1")
		}
		requirePackage := in.Arg(0)
		if !requirePackage.IsKind(v8.KindString) {
			return nil, errors.New("invalid argument type for require, expected to be string")
		}
		packageName := requirePackage.String()

		return load(in, packageName, require)
	})
	return require, ctx.Global().Set("require", f)
}
