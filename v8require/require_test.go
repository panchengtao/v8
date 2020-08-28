package v8require

import (
	"github.com/panchengtao/v8"
	"testing"
)

// Suppose there is a file containing the following code under /usr/local/module/indexOf named indexOf.js
const _ = `function indexOf(array, item) {
    return array.indexOf(item)
}

exports.indexOf = indexOf;`

// Suppose there is a file containing the following content under /usr/local/module/indexOf named package.json
const __ = `{
  "name": "indexOf",
  "version": "1.0.0",
  "description": "",
  "main": "indexOf.js",
  "scripts": {
    "test": "echo \"Error: no test specified\" && exit 1"
  },
  "keywords": [],
  "author": "",
  "license": "ISC"
}`

func TestRegisterRequire(t *testing.T) {
	ctx := v8.NewIsolate().NewContext()
	requireCtx, err := EnableRequiring(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	requireCtx.RegisterModulePath("/usr/local/module")
	res, err := ctx.Eval(`
		var _ = require('indexOf');
		let c = _.indexOf([1, 2, 1, 2], 2);
		c;
	`, "test.js")
	if err != nil {
		t.Fatalf("Error evaluating javascript, err: %v", err)
	}
	if num := res.Int64(); num != 1 {
		t.Errorf("Expected 1, got %v", res)
	}
}
