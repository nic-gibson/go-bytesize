package bytesize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Overflow(t *testing.T) {
	b, err := Parse("1797693134862315708145274237317043567981000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000B")
	assert.NotNil(t, err, "Max float test did not fail")
	assert.Zero(t, b, "Max float test did not fail")
}

var formatTable = []struct {
	Bytes  float64
	Format string
	Result string
}{
	{1, "byte", "1 B"},
	{1024, "kb", "1 KB"},
	{1099511627776, "GB", "1024 GB"},
	{1125899906842624, "GB", "1048576 GB"},
	{1125899906842624, "potato", "Unrecognized unit: potato"},
}

func Test_Format(t *testing.T) {
	for _, v := range formatTable {
		bSize := New(v.Bytes)
		b := bSize.Format("%.0f ", v.Format, false)
		assert.Equal(t, v.Result, b)
	}
}

var newTable = []struct {
	Bytes  float64
	Result string
}{
	{1, "1.00B"},
	{1023, "1023.00B"},
	{1024, "1.00KB"},
	{1048576, "1.00MB"},
	{1073741824, "1.00GB"},
	{1099511627776, "1.00TB"},
	{1125899906842624, "1.00PB"},
	{1152921504606846976, "1.00EB"},
}

func Test_New(t *testing.T) {
	for _, v := range newTable {
		b := New(v.Bytes)
		assert.Equal(t, v.Result, b.String())
	}
}

var globalFormatTable = []struct {
	Bytes  float64
	Result string
}{
	{1, "1 byte"},
	{1023, "1023 bytes"},
	{1024, "1 kilobyte"},
	{1048576, "1 megabyte"},
	{1073741824, "1 gigabyte"},
	{1099511627776, "1 terabyte"},
	{1125899906842624, "1 petabyte"},
	{1152921504606846976, "1 exabyte"},
	{2 * 1, "2 bytes"},
	{2 * 1024, "2 kilobytes"},
	{2 * 1048576, "2 megabytes"},
	{2 * 1073741824, "2 gigabytes"},
	{2 * 1099511627776, "2 terabytes"},
	{2 * 1125899906842624, "2 petabytes"},
	{2 * 1152921504606846976, "2 exabytes"},
}

func Test_GlobalFormat(t *testing.T) {
	Format = "%.0f "
	LongUnits = true
	for _, v := range globalFormatTable {
		b := New(v.Bytes)
		assert.Equal(t, v.Result, b.String())
	}
	Format = "%.2f"
	LongUnits = false
}

var parseTable = []struct {
	Input  string
	Result string
	Fail   bool
}{
	{"1B", "1.00B", false},
	{"1 B", "1.00B", false},
	{"1 byte", "1.00B", false},
	{"2 bytes", "2.00B", false},
	{"1B ", "1.00B", false},
	{" 1 B ", "1.00B", false},
	{"1023B", "1023.00B", false},
	{"1024B", "1.00KB", false},
	{"1KB 1023B", "", true},
	{"1.5GB", "1.50GB", false},
	{"1", "", true},
}

func Test_Parse(t *testing.T) {
	for _, v := range parseTable {
		b, err := Parse(v.Input)
		if v.Fail {
			assert.Error(t, err)
			assert.NotEqual(t, v.Result, b.String())
		} else {
			assert.Nil(t, err)
			assert.Equal(t, v.Result, b.String())
		}
	}
}

func Test_Set(t *testing.T) {
	for _, v := range parseTable {
		var b ByteSize
		err := b.Set(v.Input)
		if v.Fail {
			assert.Error(t, err)
			assert.NotEqual(t, v.Result, b.String())
		} else {
			assert.Nil(t, err)
			assert.Equal(t, v.Result, b.String())
		}
	}
}

var getTable = []struct {
	Input  string
	Result ByteSize
}{
	{"1 byte", 1 * B},
}

func Test_Get(t *testing.T) {
	for _, v := range getTable {
		b, err := Parse(v.Input)
		assert.Nil(t, err)
		assert.Equal(t, v.Result, b.Get())
	}
}

var mathTable = []struct {
	B1       ByteSize
	Function rune
	B2       ByteSize
	Result   string
}{
	{1024, '+', 1024, "2.00KB"},
	{1073741824, '+', 10485760, "1.01GB"},
	{1073741824, '-', 536870912, "512.00MB"},
}

func Test_Math(t *testing.T) {
	for _, v := range mathTable {
		switch v.Function {
		case '+':
			total := v.B1 + v.B2
			assert.Equal(t, v.Result, total.String(), "Fail: %s + %s = %s, received %s", v.B1, v.B2, v.Result, total)
		case '-':
			total := v.B1 - v.B2
			assert.Equal(t, v.Result, total.String(), "Fail: %s - %s = %s, received %s", v.B1, v.B2, v.Result, total)
		}
	}
}

var kbConvTable = []struct {
	B      ByteSize
	Result float64
}{{1024, 1}, {1536, 1.5}, {0, 0}, {MB, 1024}, {MB + KB, 1025}}

func Test_Conv_Kilobytes(t *testing.T) {
	for _, v := range kbConvTable {
		assert.Equal(t, v.Result, v.B.KiloBytes())
	}
}

var mbConvTable = []struct {
	B      ByteSize
	Result float64
}{{1024, 0.0009765625}, {MB, 1}, {MB + GB, 1025}}

func Test_Conv_Megabytes(t *testing.T) {
	for _, v := range mbConvTable {
		assert.Equal(t, v.Result, v.B.MegaBytes())
	}
}

var roundTable = []struct {
	B      ByteSize
	Result ByteSize
	Size   ByteSize
}{
	{1024, 1024, KB},
	{1025, 2048, KB},
	{MB + TB, TB + TB, TB},
}

func TestRound(t *testing.T) {
	for _, v := range roundTable {
		assert.Equal(t, v.Result, v.B.Round(v.Size))
	}
}

var truncTable = []struct {
	B      ByteSize
	Result ByteSize
	Size   ByteSize
}{
	{1024, 1024, KB},
	{1025, 1024, KB},
	{MB + TB, TB, TB},
}

func TestTrunc(t *testing.T) {
	for _, v := range truncTable {
		assert.Equal(t, v.Result, v.B.Trunc(v.Size))
	}
}
