// +build js

package hush

import "context"

func ttySetup() (context.CancelFunc, error) {
	return func() {}, nil
}
