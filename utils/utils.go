package utils

import (
	"fmt"
	"math"
)

func ByteSize(size int64) string {
	sizeFloat := float64(size)
	oldSize := sizeFloat
	var n float64 = 0
	for math.Abs(sizeFloat) >= 1024 {
		sizeFloat = sizeFloat / 1024
		n++
	}

	var k string
	if n == 0 {
		k = "B"
	} else if n == 1 {
		k = "KB"
	} else if n == 2 {
		k = "MB"
	} else if n == 3 {
		k = "GB"
	} else if n == 4 {
		k = "TB"
	}

	ns := oldSize / math.Pow(1024, n)

	return fmt.Sprintf("%.2f%s", ns, k)
}
