package helper

func appendslice(a []int, b []int) []int {
	out := make([]int, len(a)+len(b))
	copy(out, a)
	copy(out[len(a):], b)
	return out
}

// convert this function to generics
// but since golang does not support operator overloading
// we map values into a struct

func dot(v1 []int, v2 []int) int {
	total := 0
	for i, v := range v1 {
		total = v * v2[i]
	}
	return total
}
