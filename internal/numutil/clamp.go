package numutil

import "math"

// ClampToInt32 safely converts an int to int32, clamping to the int32
// range instead of silently overflowing/wrapping. Use this whenever a
// value crossing a trust boundary (config, env vars, JSON, CLI flags)
// needs to become an int32, e.g. to compare against a Kubernetes
// RestartCount or similar int32 API field.
func ClampToInt32(v int) int32 {
	if v > math.MaxInt32 {
		return math.MaxInt32
	}
	if v < math.MinInt32 {
		return math.MinInt32
	}
	return int32(v) //nolint:gosec // bounds-checked above
}

// ClampToInt64ToInt32 is the int64 variant, useful when reading from
// APIs or libraries that hand back int64 (e.g. encoding/json numbers,
// some proto fields).
func ClampInt64ToInt32(v int64) int32 {
	if v > math.MaxInt32 {
		return math.MaxInt32
	}
	if v < math.MinInt32 {
		return math.MinInt32
	}
	return int32(v) //nolint:gosec // bounds-checked above
}

// ClampToUint32 safely converts an int to uint32, clamping negative
// values to 0 and large values to math.MaxUint32.
func ClampToUint32(v int) uint32 {
	if v < 0 {
		return 0
	}
	if v > math.MaxUint32 {
		return math.MaxUint32
	}
	return uint32(v) //nolint:gosec // bounds-checked above
}
