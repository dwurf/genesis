package slices

// Concat concatenates given slices into a single slice.
func Concat[S ~[]T, T any](slices ...S) S {
	size := 0
	for _, items := range slices {
		size += len(items)
	}
	result := make(S, 0, size)
	for _, items := range slices {
		result = append(result, items...)
	}
	return result
}

// Intersect2 is an alias for [Intersect]
//
// DEPRECATED. Use [Intersect] instead.
func Intersect2[S1 ~[]T, S2 ~[]T, T comparable](items1 S1, items2 S2) []T {
	return Intersect(items1, items2)
}

// Intersect returns items that appear in both slices.
//
// The items in the result slice appear in the same order as in the first given slice.
// Each item appears only once.
func Intersect[S1 ~[]T, S2 ~[]T, T comparable](items1 S1, items2 S2) []T {
	wanted := make(map[T]struct{})
	for _, item := range items2 {
		wanted[item] = struct{}{}
	}
	result := make([]T, 0)
	for _, item := range items1 {
		_, found := wanted[item]
		if found {
			result = append(result, item)
			delete(wanted, item)
		}
	}
	return result
}

// Product returns cortesian product of elements in the given slices.
func Product2[T any](items ...[]T) chan []T {
	c := make(chan []T, 1)
	go product2(items, c, []T{}, 0)
	return c
}

// product is a core implementation of Product
func product2[T any](items [][]T, c chan []T, left []T, pos int) {
	// iterate over the last array
	if pos == len(items)-1 {
		for _, el := range items[pos] {
			result := make([]T, 0, len(left)+1)
			result = append(result, left...)
			result = append(result, el)
			c <- result
		}
		return
	}

	for _, el := range items[pos] {
		result := make([]T, 0, len(left)+1)
		result = append(result, left...)
		result = append(result, el)
		product2(items, c, result, pos+1)
	}

	if pos == 0 {
		close(c)
	}
}

// Union returns a slice of unique values from both slices preserving their order.
func Union[S ~[]T, T comparable](left S, right S) S {
	return Uniq(Concat(left, right))
}

// Zip returns chan of arrays of elements from given arrs on the same position.
func Zip[S ~[]T, T any](items ...S) chan S {
	if len(items) == 0 {
		result := make(chan S)
		close(result)
		return result
	}

	size := len(items[0])
	for _, arr := range items[1:] {
		if len(arr) < size {
			size = len(arr)
		}
	}

	result := make(chan S, 1)
	go func() {
		for i := 0; i < size; i++ {
			chunk := make([]T, 0, len(items))
			for _, arr := range items {
				chunk = append(chunk, arr[i])
			}
			result <- chunk
		}
		close(result)
	}()
	return result
}
