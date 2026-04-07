//go:build js && wasm

package main

import (
	"syscall/js"

	passgen "codeberg.org/jhot/passgen"
)

func main() {
	js.Global().Set("passgenGenerate", js.FuncOf(generate))
	select {} // keep the Go runtime alive
}

func generate(_ js.Value, args []js.Value) any {
	if len(args) == 0 {
		return errObj("no arguments provided")
	}
	opts := args[0]

	passphrase := jsStr(opts, "passphrase")
	b := passgen.New(passphrase)

	if s := jsStr(opts, "salt"); s != "" {
		b = b.WithSalt(s)
	}
	if n := jsInt(opts, "length"); n > 0 {
		b = b.WithLength(n)
	}
	if v := jsInt(opts, "version"); v > 0 {
		b = b.WithVersion(v)
	}
	if jsBool(opts, "noUppers") {
		b = b.NoUppers()
	}
	if jsBool(opts, "noNumbers") {
		b = b.NoNumbers()
	}
	if jsBool(opts, "noSpecials") {
		b = b.NoSpecials()
	}
	if cs := jsStr(opts, "customSpecials"); cs != "" {
		b = b.WithSpecials(cs)
	}
	if ca := jsStr(opts, "customAmbiguous"); ca != "" {
		b = b.WithAmbiguousChars(ca)
	} else if jsBool(opts, "noAmbiguous") {
		b = b.NoAmbiguous()
	}

	pass, err := b.Generate()
	if err != nil {
		return errObj(err.Error())
	}

	result := js.Global().Get("Object").New()
	result.Set("password", pass)
	return result
}

func errObj(msg string) js.Value {
	obj := js.Global().Get("Object").New()
	obj.Set("error", msg)
	return obj
}

func jsStr(v js.Value, key string) string {
	f := v.Get(key)
	if f.IsNull() || f.IsUndefined() {
		return ""
	}
	return f.String()
}

func jsInt(v js.Value, key string) int {
	f := v.Get(key)
	if f.IsNull() || f.IsUndefined() {
		return 0
	}
	return f.Int()
}

func jsBool(v js.Value, key string) bool {
	f := v.Get(key)
	if f.IsNull() || f.IsUndefined() {
		return false
	}
	return f.Bool()
}
