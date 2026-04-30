package util

func Float64SliceToFloat32(in []float64) []float32 {
	if len(in) == 0 {
		return nil
	}

	out := make([]float32, len(in))
	for i, v := range in {
		out[i] = float32(v)
	}

	return out
}

func float64MatrixToFloat32(in [][]float64) [][]float32 {
	if len(in) == 0 {
		return nil
	}

	out := make([][]float32, len(in))
	for i, row := range in {
		out[i] = Float64SliceToFloat32(row)
	}

	return out
}

func TruncateRunes(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max])
}
