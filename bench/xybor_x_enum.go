package bench

import "github.com/xybor-x/enum"

type xyborEnumType int
type XyborEnumType = enum.WrapEnum[xyborEnumType]

const (
	XyborEnumTypeT0 XyborEnumType = iota
	XyborEnumTypeT1
	XyborEnumTypeT2
	XyborEnumTypeT3
	XyborEnumTypeT4
	XyborEnumTypeT5
	XyborEnumTypeT6
	XyborEnumTypeT7
	XyborEnumTypeT8
	XyborEnumTypeT9
)

var (
	_ = enum.Map(XyborEnumTypeT0, "t0")
	_ = enum.Map(XyborEnumTypeT1, "t1")
	_ = enum.Map(XyborEnumTypeT2, "t2")
	_ = enum.Map(XyborEnumTypeT3, "t3")
	_ = enum.Map(XyborEnumTypeT4, "t4")
	_ = enum.Map(XyborEnumTypeT5, "t5")
	_ = enum.Map(XyborEnumTypeT6, "t6")
	_ = enum.Map(XyborEnumTypeT7, "t7")
	_ = enum.Map(XyborEnumTypeT8, "t8")
	_ = enum.Map(XyborEnumTypeT9, "t9")
)
