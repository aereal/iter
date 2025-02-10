package seq2_test

import (
	"fmt"
	"maps"
	"slices"

	"github.com/aereal/iter/seq2"
)

func ExampleFlip_build_map() {
	strs := []string{"a", "b", "c"}
	str2index := maps.Collect(seq2.Flip(slices.All(strs)))
	fmt.Println("map[string]int:")
	for _, str := range slices.Sorted(maps.Keys(str2index)) {
		fmt.Printf("%s=%d\n", str, str2index[str])
	}
	// Output:
	// map[string]int:
	// a=0
	// b=1
	// c=2
}
