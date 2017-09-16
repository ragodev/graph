// Copyright 2017 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package alt_test

import (
	"fmt"
	"math/rand"

	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/alt"
)

// xrs is a cheap random source for use in examples, added to when stay the same
// independent of changes to the standard library random source.  (I thought
// the plan was to keep the standard source producing the same numbers but the
// numbers seemed to change at one point.)
type xrs [2]uint64

func exampleSource(seed int64) *xrs {
	var s xrs
	s.Seed(seed)
	return &s
}

func (s *xrs) Uint64() uint64 {
	// xorshift128+ by WP
	x := s[0]
	y := s[1]
	s[0] = y
	x ^= x << 23                         // a
	s[1] = x ^ y ^ (x >> 17) ^ (y >> 26) // b, c
	return s[1] + y
}

func (s *xrs) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

func (s *xrs) Seed(n int64) {
	if n == 0 {
		panic("0 seed")
	}
	s[0] = uint64(n)
	s[1] = uint64(n + n)
	for i := 0; i < 20; i++ {
		s.Uint64()
	}
}

func ExampleArcVisitor() {
	//   0
	//  / \
	// 1-->2
	// ^   |
	// |   v
	// \---3
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {2},
		2: {3},
		3: {1},
	}
	alt.DepthFirst(g, 0, alt.ArcVisitor(func(n graph.NI, x int) {
		fmt.Println(n, "->", g[n][x])
	}))
	// Output:
	// 0 -> 1
	// 1 -> 2
	// 2 -> 3
	// 3 -> 1
	// 0 -> 2
}

func ExampleNodeVisitor() {
	//   0
	//  / \
	// 1-->2
	// ^   |
	// |   v
	// \---3
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {2},
		2: {3},
		3: {1},
	}
	alt.DepthFirst(g, 0, alt.NodeVisitor(func(n graph.NI) {
		fmt.Println(n)
	}))
	// Output:
	// 0
	// 1
	// 2
	// 3
}

func ExampleOkArcVisitor() {
	//   0
	//  / \
	// 1-->2
	// ^   |
	// |   v
	// \---3
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {2},
		2: {3},
		3: {1},
	}
	alt.DepthFirst(g, 0, alt.OkArcVisitor(func(n graph.NI, x int) bool {
		fmt.Println(n, "->", g[n][x])
		return n < g[n][x]
	}))
	// Output:
	// 0 -> 1
	// 1 -> 2
	// 2 -> 3
	// 3 -> 1
}

func ExampleOkNodeVisitor() {
	//   0
	//  / \
	// 1-->2
	// ^   |
	// |   v
	// \---3
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {2},
		2: {3},
		3: {1},
	}
	alt.DepthFirst(g, 0, alt.OkNodeVisitor(func(n graph.NI) bool {
		fmt.Println(n)
		return n != 2
	}))
	// Output:
	// 0
	// 1
	// 2
}

func ExamplePathBits() {
	//   0
	//  / \
	// 1   2
	// ^   |
	// |   v
	// \---3
	g := graph.AdjacencyList{
		0: {1, 2},
		2: {3},
		3: {1},
	}
	b := bits.New(len(g))
	fmt.Println("node  path bits")
	fmt.Println("      (3210)")
	fmt.Println("----   ----")
	alt.DepthFirst(g, 0, alt.PathBits(&b),
		alt.NodeVisitor(func(n graph.NI) {
			fmt.Printf("%4d   %s\n", n, &b)
		}))
	// Output:
	// node  path bits
	//       (3210)
	// ----   ----
	//    0   0001
	//    1   0011
	//    2   0101
	//    3   1101
}

func ExampleVisited() {
	//   0
	//  / \
	// 1-->2
	// ^   |
	// |   v
	// \---3
	g := graph.AdjacencyList{
		0: {1, 2},
		1: {2},
		2: {3},
		3: {1},
	}
	b := bits.New(len(g))
	fmt.Println("3210")
	fmt.Println("----")
	alt.DepthFirst(g, 0, alt.Visited(&b),
		alt.NodeVisitor(func(graph.NI) {
			fmt.Println(b)
		}))
	// Output:
	// 3210
	// ----
	// 0001
	// 0011
	// 0111
	// 1111
}

func ExampleRand() {
	//         0
	//         |
	// -------------------
	// | | | | | | | | | |
	// 1 2 3 4 5 6 7 8 9 10
	g := graph.AdjacencyList{
		0:  {1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		10: nil,
	}
	alt.DepthFirst(g, 0, alt.Rand(rand.New(exampleSource(7))),
		alt.NodeVisitor(func(n graph.NI) {
			fmt.Println(n)
		}))
	// Output:
	// 0
	// 1
	// 2
	// 3
	// 8
	// 7
	// 9
	// 5
	// 10
	// 6
	// 4
}

func ExampleBreadthFirst_singlePath() {
	// arcs are directed right:
	//    1   3---5
	//   / \ /   /
	//  2   4---6--\
	//           \-/
	g := graph.AdjacencyList{
		2: {1},
		1: {4},
		4: {3, 6},
		3: {5},
		6: {5, 6},
	}
	var start, end graph.NI = 1, 6
	var f graph.FromList
	var visited int
	alt.BreadthFirst(g, start, alt.From(&f),
		alt.OkNodeVisitor(func(n graph.NI) bool {
			visited++
			return n != end
		}))
	fmt.Println(visited, "nodes visited")
	fmt.Println("path:", f.PathTo(end, nil))
	// Output:
	// 4 nodes visited
	// path: [1 4 6]
}

func ExampleBreadthFirst_allPaths() {
	// arcs are directed right:
	//    1   3---5
	//   / \ /   /
	//  2   4---6--\
	//           \-/
	g := graph.AdjacencyList{
		2: {1},
		1: {4},
		4: {3, 6},
		3: {5},
		6: {5, 6},
	}
	start := graph.NI(1)
	var f graph.FromList
	alt.BreadthFirst(g, start, alt.From(&f))
	fmt.Println("Max path length:", f.MaxLen)
	p := make([]graph.NI, f.MaxLen)
	for n := range g {
		fmt.Println(n, f.PathTo(graph.NI(n), p))
	}
	// Output:
	// Max path length: 4
	// 0 []
	// 1 [1]
	// 2 []
	// 3 [1 4 3]
	// 4 [1 4]
	// 5 [1 4 3 5]
	// 6 [1 4 6]
}

func ExampleBreadthFirst_traverse() {
	// arcs directed down
	//    0--
	//   /|  \
	//  1 2   3
	//   /|\  |\
	//  4 5 6 7 8
	g := graph.AdjacencyList{
		0: {1, 2, 3},
		2: {4, 5, 6},
		3: {7, 8},
		8: {},
	}
	var f graph.FromList
	alt.BreadthFirst(g, 0, alt.From(&f),
		alt.NodeVisitor(func(n graph.NI) {
			fmt.Println("visit", n, "level", f.Paths[n].Len)
		}))
	// Output:
	// visit 0 level 1
	// visit 1 level 2
	// visit 2 level 2
	// visit 3 level 2
	// visit 4 level 3
	// visit 5 level 3
	// visit 6 level 3
	// visit 7 level 3
	// visit 8 level 3
}

func ExampleBreadthFirst_traverseRandom() {
	// arcs directed down
	//    0--
	//   /|  \
	//  1 2   3
	//   /|\  |\
	//  4 5 6 7 8
	g := graph.AdjacencyList{
		0: {1, 2, 3},
		2: {4, 5, 6},
		3: {7, 8},
		8: {},
	}

	// only difference from non-random example
	r := rand.New(exampleSource(8))

	var f graph.FromList
	alt.BreadthFirst(g, 0, alt.Rand(r), alt.From(&f),
		alt.NodeVisitor(func(n graph.NI) {
			fmt.Println("visit", n, "level", f.Paths[n].Len)
		}))
	// Output:
	// visit 0 level 1
	// visit 3 level 2
	// visit 2 level 2
	// visit 1 level 2
	// visit 6 level 3
	// visit 8 level 3
	// visit 5 level 3
	// visit 4 level 3
	// visit 7 level 3
}
