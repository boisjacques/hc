package rtp

import (
	"fmt"
	"github.com/boisjacques/hc/characteristic"
	"github.com/boisjacques/hc/tlv8"
	"testing"
)

func TestSelectedStreamConfiguration(t *testing.T) {
	c := characteristic.NewSelectedStreamConfiguration()
	c.Value = "ARUCAQABEHW8tiJ9E0F4tLlvOURdFCc="

	b := c.GetValue()

	var cfg StreamConfiguration
	err := tlv8.Unmarshal(b, &cfg)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v", cfg)
}
