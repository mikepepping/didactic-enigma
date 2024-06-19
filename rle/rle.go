package rle

import "errors"

func Encode(in <-chan byte, out chan<- []byte) error {
	var (
		runLength   byte
		runningByte byte
	)

	for b := range in {
		if runLength == 0 {
			// initial state
			runLength = 1
			runningByte = b
			continue
		}

		// new byte or reached max byte size for counting the run
		if b != runningByte || runLength == 255 {
			out <- []byte{runningByte, runLength}
			runLength = 1
			runningByte = b
			continue
		}

		runLength++
	}

	if runLength > 0 {
		out <- []byte{runningByte, runLength}
	}

	close(out)
	return nil
}

func Decode(in <-chan []byte, out chan<- byte) error {
	for pair := range in {
		if len(pair) != 2 {
			return errors.New("pair not to bytes long")
		}

		for i := 0; i < int(pair[1]); i++ {
			out <- pair[0]
		}
	}

	return nil
}
