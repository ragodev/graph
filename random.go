// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package graph

import (
	"errors"
	"math"
	"math/big"
	"math/rand"
	"time"
)

// Euclidean generates a random simple graph on the Euclidean plane.
//
// Nodes are associated with coordinates uniformly distributed on a unit
// square.  Arcs are added between random nodes with a bias toward connecting
// nearer nodes.
//
// Unfortunately the function has a few "knobs".
// The returned graph will have order nNodes and arc size nArcs.  The affinity
// argument controls the bias toward connecting nearer nodes.  The function
// selects random pairs of nodes as a candidate arc then rejects the candidate
// if the nodes fail an affinity test.  Also parallel arcs are rejected.
// As more affine or denser graphs are requested, rejections increase,
// increasing run time.  The patience argument controls the number of arc
// rejections allowed before the function gives up and returns an error.
// Note that higher affinity will require more patience and that some
// combinations of nNodes and nArcs cannot be achieved with any amount of
// patience given that the returned graph must be simple.
//
// Returned is a directed simple graph and associated positions indexed by
// node number.
//
// See also LabeledEuclidean.
func Euclidean(nNodes, nArcs int, affinity float64, patience int, r *rand.Rand) (g Directed, pos []struct{ X, Y float64 }, err error) {
	a := make(AdjacencyList, nNodes) // graph
	// generate random positions
	pos = make([]struct{ X, Y float64 }, nNodes)
	for i := range pos {
		pos[i].X = r.Float64()
		pos[i].Y = r.Float64()
	}
	// arcs
	var tooFar, dup int
arc:
	for i := 0; i < nArcs; {
		if tooFar == nArcs*patience {
			err = errors.New("affinity not found")
			return
		}
		if dup == nArcs*patience {
			err = errors.New("overcrowding")
			return
		}
		n1 := NI(r.Intn(nNodes))
		var n2 NI
		for {
			n2 = NI(r.Intn(nNodes))
			if n2 != n1 { // no graph loops
				break
			}
		}
		c1 := &pos[n1]
		c2 := &pos[n2]
		dist := math.Hypot(c2.X-c1.X, c2.Y-c1.Y)
		if dist*affinity > r.ExpFloat64() { // favor near nodes
			tooFar++
			continue
		}
		for _, nb := range a[n1] {
			if nb == n2 { // no parallel arcs
				dup++
				continue arc
			}
		}
		a[n1] = append(a[n1], n2)
		i++
	}
	g = Directed{a}
	return
}

// LabeledEuclidean generates a random simple graph on the Euclidean plane.
//
// Arc label values in the returned graph g are indexes into the return value
// wt.  Wt is the Euclidean distance between the from and to nodes of the arc.
//
// Otherwise the function arguments and return values are the same as for
// function Euclidean.  See Euclidean.
func LabeledEuclidean(nNodes, nArcs int, affinity float64, patience int, r *rand.Rand) (g LabeledDirected, pos []struct{ X, Y float64 }, wt []float64, err error) {
	a := make(LabeledAdjacencyList, nNodes) // graph
	wt = make([]float64, nArcs)             // arc weights
	// generate random positions
	pos = make([]struct{ X, Y float64 }, nNodes)
	for i := range pos {
		pos[i].X = r.Float64()
		pos[i].Y = r.Float64()
	}
	// arcs
	var tooFar, dup int
arc:
	for i := 0; i < nArcs; {
		if tooFar == nArcs*patience {
			err = errors.New("affinity not found")
			return
		}
		if dup == nArcs*patience {
			err = errors.New("overcrowding")
			return
		}
		n1 := NI(r.Intn(nNodes))
		var n2 NI
		for {
			n2 = NI(r.Intn(nNodes))
			if n2 != n1 { // no graph loops
				break
			}
		}
		c1 := &pos[n1]
		c2 := &pos[n2]
		dist := math.Hypot(c2.X-c1.X, c2.Y-c1.Y)
		if dist*affinity > r.ExpFloat64() { // favor near nodes
			tooFar++
			continue
		}
		for _, nb := range a[n1] {
			if nb.To == n2 { // no parallel arcs
				dup++
				continue arc
			}
		}
		wt[i] = dist
		a[n1] = append(a[n1], Half{n2, LI(i)})
		i++
	}
	g = LabeledDirected{a}
	return
}

func Geometric(nNodes int, radius float64, r *rand.Rand) (g Undirected, pos []struct{ X, Y float64 }, m int) {
	// Expected degree is approximately nπr².
	a := make(AdjacencyList, nNodes)
	pos = make([]struct{ X, Y float64 }, nNodes)
	for i := range pos {
		pos[i].X = r.Float64()
		pos[i].Y = r.Float64()
	}
	for u, up := range pos {
		for v := u + 1; v < len(pos); v++ {
			vp := pos[v]
			if math.Hypot(up.X-vp.X, up.Y-vp.Y) < radius {
				a[u] = append(a[u], NI(v))
				a[v] = append(a[v], NI(u))
				m++
			}
		}
	}
	g = Undirected{a}
	return
}

func LabeledGeometric(nNodes int, radius float64, r *rand.Rand) (g LabeledUndirected, pos []struct{ X, Y float64 }, wt []float64) {
	a := make(LabeledAdjacencyList, nNodes)
	pos = make([]struct{ X, Y float64 }, nNodes)
	for i := range pos {
		pos[i].X = r.Float64()
		pos[i].Y = r.Float64()
	}
	for u, up := range pos {
		for v := u + 1; v < len(pos); v++ {
			vp := pos[v]
			if w := math.Hypot(up.X-vp.X, up.Y-vp.Y); w < radius {
				a[u] = append(a[u], Half{NI(v), LI(len(wt))})
				a[v] = append(a[v], Half{NI(u), LI(len(wt))})
				wt = append(wt, w)
			}
		}
	}
	g = LabeledUndirected{a}
	return
}

// KroneckerDir generates a Kronecker-like random directed graph.
//
// The returned graph g is simple and has no isolated nodes but is not
// necessarily fully connected.  The number of of nodes will be <= 2^scale,
// and will be near 2^scale for typical values of arcFactor, >= 2.
// ArcFactor * 2^scale arcs are generated, although loops and duplicate arcs
// are rejected.
//
// Return value ma is the number of arcs retained in the result graph.
func KroneckerDir(scale uint, arcFactor float64) (g Directed, ma int) {
	a, m := kronecker(scale, arcFactor, true)
	return Directed{a}, m
}

// KroneckerUndir generates a Kronecker-like random undirected graph.
//
// The returned graph g is simple and has no isolated nodes but is not
// necessarily fully connected.  The number of of nodes will be <= 2^scale,
// and will be near 2^scale for typical values of edgeFactor, >= 2.
// EdgeFactor * 2^scale edges are generated, although loops and duplicate edges
// are rejected.
//
// Return value ma is the number of arcs--not edges--retained in the result
// graph.
func KroneckerUndir(scale uint, edgeFactor float64) (g Undirected, ma int) {
	al, as := kronecker(scale, edgeFactor, false)
	return Undirected{al}, as
}

// Styled after the Graph500 example code.  Not well tested currently.
// Graph500 example generates undirected only.  No idea if the directed variant
// here is meaningful or not.
func kronecker(scale uint, edgeFactor float64, dir bool) (g AdjacencyList, ma int) {
	rand.Seed(time.Now().Unix())
	N := 1 << scale                      // node extent
	M := int(edgeFactor*float64(N) + .5) // number of arcs/edges to generate
	a, b, c := 0.57, 0.19, 0.19          // initiator probabilities
	ab := a + b
	cNorm := c / (1 - ab)
	aNorm := a / ab
	ij := make([][2]int, M)
	var bm big.Int
	var nNodes int
	for k := range ij {
		var i, j int
		for b := 1; b < N; b <<= 1 {
			if rand.Float64() > ab {
				i |= b
				if rand.Float64() > cNorm {
					j |= b
				}
			} else if rand.Float64() > aNorm {
				j |= b
			}
		}
		if bm.Bit(i) == 0 {
			bm.SetBit(&bm, i, 1)
			nNodes++
		}
		if bm.Bit(j) == 0 {
			bm.SetBit(&bm, j, 1)
			nNodes++
		}
		r := rand.Intn(k + 1) // shuffle edges as they are generated
		ij[k] = ij[r]
		ij[r] = [2]int{i, j}
	}
	p := rand.Perm(nNodes) // mapping to shuffle IDs of non-isolated nodes
	px := 0
	r := make([]NI, N)
	for i := range r {
		if bm.Bit(i) == 1 {
			r[i] = NI(p[px]) // fill lookup table
			px++
		}
	}
	g = make(AdjacencyList, nNodes)
ij:
	for _, e := range ij {
		if e[0] == e[1] {
			continue // skip loops
		}
		ri, rj := r[e[0]], r[e[1]]
		for _, nb := range g[ri] {
			if nb == rj {
				continue ij // skip parallel edges
			}
		}
		g[ri] = append(g[ri], rj)
		ma++
		if !dir {
			g[rj] = append(g[rj], ri)
			ma++
		}
	}
	return
}
