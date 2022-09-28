package workerchannel

import (
	"testing"
)

func TestInitChannelHandlers(t *testing.T) {
	InitChannelHandlers()

	RegisterHandler("abc", nil)
	UnregisterHandler("")
}
