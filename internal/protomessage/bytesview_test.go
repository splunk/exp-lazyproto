package protomessage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesViewNil(t *testing.T) {
	var bv BytesView
	b := BytesFromBytesView(bv)
	assert.True(t, b == nil)
}

func TestBytesViewNotNil(t *testing.T) {
	src := []byte{1, 2, 3}

	bv := BytesViewFromBytes(src)
	dest := BytesFromBytesView(bv)
	assert.EqualValues(t, src, dest)
}

func BenchmarkEmptyOp(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}

func BenchmarkBytesViewFromBytes(b *testing.B) {
	bytes := []byte{1, 2, 3}
	for i := 0; i < b.N; i++ {
		BytesViewFromBytes(bytes)
	}
}

func BenchmarkBytesFromBytesView(b *testing.B) {
	bytes := []byte{1, 2, 3}
	bv := BytesViewFromBytes(bytes)
	for i := 0; i < b.N; i++ {
		BytesFromBytesView(bv)
	}
}

func BenchmarkBytesFromBytesViewAndBack(b *testing.B) {
	bytes := []byte{1, 2, 3}
	bv := BytesViewFromBytes(bytes)
	for i := 0; i < b.N; i++ {
		BytesViewFromBytes(BytesFromBytesView(bv))
	}
}
