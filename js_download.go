// +build js

package hush

import (
	"syscall/js"
)

var (
	jsBlob       = js.Global().Get("Blob")
	jsDocument   = js.Global().Get("document")
	jsURL        = js.Global().Get("URL")
	jsUint8Array = js.Global().Get("Uint8Array")
)

func startDownload(contentType, fileName string, buf []byte) {
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	jsBuf := jsUint8Array.New(len(buf))
	js.CopyBytesToJS(jsBuf, buf)
	blobInstance := jsBlob.New([]interface{}{jsBuf}, map[string]interface{}{
		"type": contentType,
	})
	link := jsDocument.Call("createElement", "a")
	link.Set("href", jsURL.Call("createObjectURL", blobInstance))
	link.Set("download", fileName)
	link.Call("click")
}
