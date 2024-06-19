package rle_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/mikepepping/didactic-enigma/rle"
)

func TestEncode(t *testing.T) {
	// Define some test cases with expected input and output
	tests := []struct {
		name   string
		input  []byte
		output [][]byte
	}{
		{name: "Single byte", input: []byte{'A'}, output: [][]byte{{'A', 1}}},
		{name: "Repeated byte", input: []byte{'A', 'A', 'A'}, output: [][]byte{{'A', 3}}},
		{name: "Multiple bytes", input: []byte{'A', 'A', 'B', 'B', 'C'}, output: [][]byte{{'A', 2}, {'B', 2}, {'C', 1}}},
		{name: "Max run length", input: bytes.Repeat([]byte{'A'}, 256), output: [][]byte{{'A', 255}, {'A', 1}}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			inCh := make(chan byte)
			outCh := make(chan []byte)

			// Launch goroutine for encoding
			go func() {
				err := rle.Encode(inCh, outCh)
				if err != nil {
					t.Errorf("Unexpected error during encoding: %v", err)
					return
				}
			}()

			// Send test data on the input channel
			go func() {
				for _, b := range tc.input {
					inCh <- b
				}
				close(inCh)
			}()

			// Collect encoded data
			var encoded [][]byte
			for data := range outCh {
				encoded = append(encoded, data)
			}

			// Check if the encoded data matches the expected output
			if !reflect.DeepEqual(encoded, tc.output) {
				t.Errorf("Test case '%s': Encoded data does not match expected output.\nExpected: %v\nGot: %v", tc.name, tc.output, encoded)
			}
		})
	}
}
