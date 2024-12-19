package bench_test

import (
	"encoding/json"
	"testing"

	"github.com/xybor-x/enum"
	"github.com/xybor-x/enum/bench"
)

func BenchmarkGen10ToString(b *testing.B) {
	b.Run("Gen", func(b *testing.B) {
		enum := bench.GenEnumTypeT9
		for i := 0; i < b.N; i++ {
			_ = enum.String()
		}
	})

	b.Run("XyborX", func(b *testing.B) {
		enum := bench.XyborEnumTypeT9
		for i := 0; i < b.N; i++ {
			_ = enum.String()
		}
	})
}

func BenchmarkGen10FromString(b *testing.B) {
	b.Run("Gen", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bench.ParseGenEnumType("t9")
		}
	})

	b.Run("XyborX", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			enum.FromString[bench.XyborEnumType]("t9")
		}
	})
}

func BenchmarkGen10JsonMarshal(b *testing.B) {
	b.Run("Gen", func(b *testing.B) {
		enum := bench.GenEnumTypeT9
		for i := 0; i < b.N; i++ {
			json.Marshal(enum)
		}
	})

	b.Run("XyborX", func(b *testing.B) {
		enum := bench.XyborEnumTypeT9
		for i := 0; i < b.N; i++ {
			json.Marshal(enum)
		}
	})
}

func BenchmarkGen10JsonUnmarshal(b *testing.B) {
	b.Run("Gen", func(b *testing.B) {
		var enum bench.GenEnumType
		for i := 0; i < b.N; i++ {
			json.Unmarshal([]byte(`"t9"`), &enum)
		}
	})

	b.Run("XyborX", func(b *testing.B) {
		var enum bench.XyborEnumType
		for i := 0; i < b.N; i++ {
			json.Unmarshal([]byte(`"t9"`), &enum)
		}
	})
}

func BenchmarkGen10SqlValue(b *testing.B) {
	b.Run("Gen", func(b *testing.B) {
		enum := bench.GenEnumTypeT9
		for i := 0; i < b.N; i++ {
			enum.Value()
		}
	})

	b.Run("XyborX", func(b *testing.B) {
		enum := bench.XyborEnumTypeT9
		for i := 0; i < b.N; i++ {
			enum.Value()
		}
	})
}

func BenchmarkGen10SqlScanByte(b *testing.B) {
	b.Run("Gen", func(b *testing.B) {
		var enum bench.GenEnumType
		for i := 0; i < b.N; i++ {
			enum.Scan([]byte(`t9`))
		}
	})

	b.Run("XyborX", func(b *testing.B) {
		var enum bench.XyborEnumType
		for i := 0; i < b.N; i++ {
			enum.Scan([]byte(`t9`))
		}
	})
}

func BenchmarkSqlScanString(b *testing.B) {
	b.Run("Gen", func(b *testing.B) {
		var enum bench.GenEnumType
		for i := 0; i < b.N; i++ {
			enum.Scan("t9")
		}
	})

	b.Run("XyborX", func(b *testing.B) {
		var enum bench.XyborEnumType
		for i := 0; i < b.N; i++ {
			enum.Scan("t9")
		}
	})
}
