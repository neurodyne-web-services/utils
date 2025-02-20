package pointer

import (
	"time"

	"github.com/aws/smithy-go/ptr"
)

// ToBool returns bool value dereferenced if the passed
// in pointer was not nil. Returns a bool zero value if the
// pointer was nil.
func ToBool(p *bool) (v bool) {
	return ptr.ToBool(p)
}

// ToBoolSlice returns a slice of bool values, that are
// dereferenced if the passed in pointer was not nil. Returns a bool
// zero value if the pointer was nil.
func ToBoolSlice(vs []*bool) []bool {
	return ptr.ToBoolSlice(vs)
}

// ToBoolMap returns a map of bool values, that are
// dereferenced if the passed in pointer was not nil. The bool
// zero value is used if the pointer was nil.
func ToBoolMap(vs map[string]*bool) map[string]bool {
	return ptr.ToBoolMap(vs)
}

// ToByte returns byte value dereferenced if the passed
// in pointer was not nil. Returns a byte zero value if the
// pointer was nil.
func ToByte(p *byte) (v byte) {
	return ptr.ToByte(p)
}

// ToByteSlice returns a slice of byte values, that are
// dereferenced if the passed in pointer was not nil. Returns a byte
// zero value if the pointer was nil.
func ToByteSlice(vs []*byte) []byte {
	return ptr.ToByteSlice(vs)
}

// ToByteMap returns a map of byte values, that are
// dereferenced if the passed in pointer was not nil. The byte
// zero value is used if the pointer was nil.
func ToByteMap(vs map[string]*byte) map[string]byte {
	return ptr.ToByteMap(vs)
}

// ToString returns string value dereferenced if the passed
// in pointer was not nil. Returns a string zero value if the
// pointer was nil.
func ToString(p *string) (v string) {
	return ptr.ToString(p)
}

// ToStringSlice returns a slice of string values, that are
// dereferenced if the passed in pointer was not nil. Returns a string
// zero value if the pointer was nil.
func ToStringSlice(vs []*string) []string {
	return ptr.ToStringSlice(vs)
}

// ToStringMap returns a map of string values, that are
// dereferenced if the passed in pointer was not nil. The string
// zero value is used if the pointer was nil.
func ToStringMap(vs map[string]*string) map[string]string {
	return ptr.ToStringMap(vs)
}

// ToInt returns int value dereferenced if the passed
// in pointer was not nil. Returns a int zero value if the
// pointer was nil.
func ToInt(p *int) (v int) {
	return ptr.ToInt(p)
}

// ToIntSlice returns a slice of int values, that are
// dereferenced if the passed in pointer was not nil. Returns a int
// zero value if the pointer was nil.
func ToIntSlice(vs []*int) []int {
	return ptr.ToIntSlice(vs)
}

// ToIntMap returns a map of int values, that are
// dereferenced if the passed in pointer was not nil. The int
// zero value is used if the pointer was nil.
func ToIntMap(vs map[string]*int) map[string]int {
	return ptr.ToIntMap(vs)
}

// ToInt8 returns int8 value dereferenced if the passed
// in pointer was not nil. Returns a int8 zero value if the
// pointer was nil.
func ToInt8(p *int8) (v int8) {
	return ptr.ToInt8(p)
}

// ToInt8Slice returns a slice of int8 values, that are
// dereferenced if the passed in pointer was not nil. Returns a int8
// zero value if the pointer was nil.
func ToInt8Slice(vs []*int8) []int8 {
	return ptr.ToInt8Slice(vs)
}

// ToInt8Map returns a map of int8 values, that are
// dereferenced if the passed in pointer was not nil. The int8
// zero value is used if the pointer was nil.
func ToInt8Map(vs map[string]*int8) map[string]int8 {
	return ptr.ToInt8Map(vs)
}

// ToInt16 returns int16 value dereferenced if the passed
// in pointer was not nil. Returns a int16 zero value if the
// pointer was nil.
func ToInt16(p *int16) (v int16) {
	return ptr.ToInt16(p)
}

// ToInt16Slice returns a slice of int16 values, that are
// dereferenced if the passed in pointer was not nil. Returns a int16
// zero value if the pointer was nil.
func ToInt16Slice(vs []*int16) []int16 {
	return ptr.ToInt16Slice(vs)
}

// ToInt16Map returns a map of int16 values, that are
// dereferenced if the passed in pointer was not nil. The int16
// zero value is used if the pointer was nil.
func ToInt16Map(vs map[string]*int16) map[string]int16 {
	return ptr.ToInt16Map(vs)
}

// ToInt32 returns int32 value dereferenced if the passed
// in pointer was not nil. Returns a int32 zero value if the
// pointer was nil.
func ToInt32(p *int32) (v int32) {
	return ptr.ToInt32(p)
}

// ToInt32Slice returns a slice of int32 values, that are
// dereferenced if the passed in pointer was not nil. Returns a int32
// zero value if the pointer was nil.
func ToInt32Slice(vs []*int32) []int32 {
	return ptr.ToInt32Slice(vs)
}

// ToInt32Map returns a map of int32 values, that are
// dereferenced if the passed in pointer was not nil. The int32
// zero value is used if the pointer was nil.
func ToInt32Map(vs map[string]*int32) map[string]int32 {
	return ptr.ToInt32Map(vs)
}

// ToInt64 returns int64 value dereferenced if the passed
// in pointer was not nil. Returns a int64 zero value if the
// pointer was nil.
func ToInt64(p *int64) (v int64) {
	return ptr.ToInt64(p)
}

// ToInt64Slice returns a slice of int64 values, that are
// dereferenced if the passed in pointer was not nil. Returns a int64
// zero value if the pointer was nil.
func ToInt64Slice(vs []*int64) []int64 {
	return ptr.ToInt64Slice(vs)
}

// ToInt64Map returns a map of int64 values, that are
// dereferenced if the passed in pointer was not nil. The int64
// zero value is used if the pointer was nil.
func ToInt64Map(vs map[string]*int64) map[string]int64 {
	return ptr.ToInt64Map(vs)
}

// ToUint returns uint value dereferenced if the passed
// in pointer was not nil. Returns a uint zero value if the
// pointer was nil.
func ToUint(p *uint) (v uint) {
	return ptr.ToUint(p)
}

// ToUintSlice returns a slice of uint values, that are
// dereferenced if the passed in pointer was not nil. Returns a uint
// zero value if the pointer was nil.
func ToUintSlice(vs []*uint) []uint {
	return ptr.ToUintSlice(vs)
}

// ToUintMap returns a map of uint values, that are
// dereferenced if the passed in pointer was not nil. The uint
// zero value is used if the pointer was nil.
func ToUintMap(vs map[string]*uint) map[string]uint {
	return ptr.ToUintMap(vs)
}

// ToUint8 returns uint8 value dereferenced if the passed
// in pointer was not nil. Returns a uint8 zero value if the
// pointer was nil.
func ToUint8(p *uint8) (v uint8) {
	return ptr.ToUint8(p)
}

// ToUint8Slice returns a slice of uint8 values, that are
// dereferenced if the passed in pointer was not nil. Returns a uint8
// zero value if the pointer was nil.
func ToUint8Slice(vs []*uint8) []uint8 {
	return ptr.ToUint8Slice(vs)
}

// ToUint8Map returns a map of uint8 values, that are
// dereferenced if the passed in pointer was not nil. The uint8
// zero value is used if the pointer was nil.
func ToUint8Map(vs map[string]*uint8) map[string]uint8 {
	return ptr.ToUint8Map(vs)
}

// ToUint16 returns uint16 value dereferenced if the passed
// in pointer was not nil. Returns a uint16 zero value if the
// pointer was nil.
func ToUint16(p *uint16) (v uint16) {
	return ptr.ToUint16(p)
}

// ToUint16Slice returns a slice of uint16 values, that are
// dereferenced if the passed in pointer was not nil. Returns a uint16
// zero value if the pointer was nil.
func ToUint16Slice(vs []*uint16) []uint16 {
	return ptr.ToUint16Slice(vs)
}

// ToUint16Map returns a map of uint16 values, that are
// dereferenced if the passed in pointer was not nil. The uint16
// zero value is used if the pointer was nil.
func ToUint16Map(vs map[string]*uint16) map[string]uint16 {
	return ptr.ToUint16Map(vs)
}

// ToUint32 returns uint32 value dereferenced if the passed
// in pointer was not nil. Returns a uint32 zero value if the
// pointer was nil.
func ToUint32(p *uint32) (v uint32) {
	return ptr.ToUint32(p)
}

// ToUint32Slice returns a slice of uint32 values, that are
// dereferenced if the passed in pointer was not nil. Returns a uint32
// zero value if the pointer was nil.
func ToUint32Slice(vs []*uint32) []uint32 {
	return ptr.ToUint32Slice(vs)
}

// ToUint32Map returns a map of uint32 values, that are
// dereferenced if the passed in pointer was not nil. The uint32
// zero value is used if the pointer was nil.
func ToUint32Map(vs map[string]*uint32) map[string]uint32 {
	return ptr.ToUint32Map(vs)
}

// ToUint64 returns uint64 value dereferenced if the passed
// in pointer was not nil. Returns a uint64 zero value if the
// pointer was nil.
func ToUint64(p *uint64) (v uint64) {
	return ptr.ToUint64(p)
}

// ToUint64Slice returns a slice of uint64 values, that are
// dereferenced if the passed in pointer was not nil. Returns a uint64
// zero value if the pointer was nil.
func ToUint64Slice(vs []*uint64) []uint64 {
	return ptr.ToUint64Slice(vs)
}

// ToUint64Map returns a map of uint64 values, that are
// dereferenced if the passed in pointer was not nil. The uint64
// zero value is used if the pointer was nil.
func ToUint64Map(vs map[string]*uint64) map[string]uint64 {
	return ptr.ToUint64Map(vs)
}

// ToFloat32 returns float32 value dereferenced if the passed
// in pointer was not nil. Returns a float32 zero value if the
// pointer was nil.
func ToFloat32(p *float32) (v float32) {
	return ptr.ToFloat32(p)
}

// ToFloat32Slice returns a slice of float32 values, that are
// dereferenced if the passed in pointer was not nil. Returns a float32
// zero value if the pointer was nil.
func ToFloat32Slice(vs []*float32) []float32 {
	return ptr.ToFloat32Slice(vs)
}

// ToFloat32Map returns a map of float32 values, that are
// dereferenced if the passed in pointer was not nil. The float32
// zero value is used if the pointer was nil.
func ToFloat32Map(vs map[string]*float32) map[string]float32 {
	return ptr.ToFloat32Map(vs)
}

// ToFloat64 returns float64 value dereferenced if the passed
// in pointer was not nil. Returns a float64 zero value if the
// pointer was nil.
func ToFloat64(p *float64) (v float64) {
	return ptr.ToFloat64(p)
}

// ToFloat64Slice returns a slice of float64 values, that are
// dereferenced if the passed in pointer was not nil. Returns a float64
// zero value if the pointer was nil.
func ToFloat64Slice(vs []*float64) []float64 {
	return ptr.ToFloat64Slice(vs)
}

// ToFloat64Map returns a map of float64 values, that are
// dereferenced if the passed in pointer was not nil. The float64
// zero value is used if the pointer was nil.
func ToFloat64Map(vs map[string]*float64) map[string]float64 {
	return ptr.ToFloat64Map(vs)
}

// ToTime returns time.Time value dereferenced if the passed
// in pointer was not nil. Returns a time.Time zero value if the
// pointer was nil.
func ToTime(p *time.Time) (v time.Time) {
	return ptr.ToTime(p)
}

// ToTimeSlice returns a slice of time.Time values, that are
// dereferenced if the passed in pointer was not nil. Returns a time.Time
// zero value if the pointer was nil.
func ToTimeSlice(vs []*time.Time) []time.Time {
	return ptr.ToTimeSlice(vs)
}

// ToTimeMap returns a map of time.Time values, that are
// dereferenced if the passed in pointer was not nil. The time.Time
// zero value is used if the pointer was nil.
func ToTimeMap(vs map[string]*time.Time) map[string]time.Time {
	return ptr.ToTimeMap(vs)
}

// ToDuration returns time.Duration value dereferenced if the passed
// in pointer was not nil. Returns a time.Duration zero value if the
// pointer was nil.
func ToDuration(p *time.Duration) (v time.Duration) {
	return ptr.ToDuration(p)
}

// ToDurationSlice returns a slice of time.Duration values, that are
// dereferenced if the passed in pointer was not nil. Returns a time.Duration
// zero value if the pointer was nil.
func ToDurationSlice(vs []*time.Duration) []time.Duration {
	return ptr.ToDurationSlice(vs)
}

// ToDurationMap returns a map of time.Duration values, that are
// dereferenced if the passed in pointer was not nil. The time.Duration
// zero value is used if the pointer was nil.
func ToDurationMap(vs map[string]*time.Duration) map[string]time.Duration {
	return ptr.ToDurationMap(vs)
}
