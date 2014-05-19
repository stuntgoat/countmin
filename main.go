// Countmin algorithm as translated from the C code found at:
// http://www.cs.rutgers.edu/~muthu/massdal-code-index.html

//********************************************************************
// Count-Min Sketches

// G. Cormode 2003,2004

// Initial version: 2003-12

// This work is licensed under the Creative Commons
// Attribution-NonCommercial License. To view a copy of this license,
// visit http://creativecommons.org/licenses/by-nc/1.0/ or send a letter
// to Creative Commons, 559 Nathan Abbott Way, Stanford, California
// 94305, USA.
// *********************************************************************/


package countmin

import (
	"math/rand"
	"sort"
)

// Constants for hashing items at each row.
const MOD = 2147483647
const HL = 31
// A simple hash function that hashes `x` using random numbers `a` and `b`.
func hash31(a, b, x int64) (result int64) {
	result = (a * x) + b
	result = ((result >> HL) + result) & MOD
	return
}

// Int64 min func
func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}


// For sorting int64 arrays.
type int64arr []int64
func (i64a int64arr) Len() int {
	return len(i64a)
}
func (i64a int64arr) Swap(i, j int) {
	i64a[i], i64a[j] = i64a[j], i64a[i]
}
func (i64a int64arr) Less(i, j int) bool {
	return i64a[i] < i64a[j]
}


// CountMin holds integer counts. For every value added to the data structure, we
// hash it at each row, increment the column index that is the hashed value modulo the
// number of columns.
type CountMin struct {
	// Specifies the number of columns in the 2D matrix.
	width int64

	// Specifies the number of rows in the 2D matrix.
	depth int64

	// Holds the total number of times this  sketch was incremented.
	count int64

	// Holds the 2D matrix of increment counts.
	counts [][]int64

	// These are pre-computed random numbers for use in the hashing
	// function, which takes 2 values; in this case these 2 values
	// and the number which will be hashed.
	aHashes []int64
	bHashes []int64
}


// MakeCM returns a new instance of a CountMin struct with `width` columns
// and `depth` rows.
func MakeCM(width, depth int64) (cm *CountMin) {
	// Initialize the rows.
	counts := make([][]int64, depth)
	var idx_depth int64
	for idx_depth = 0; idx_depth < depth; idx_depth++ {
		// Initialize columns in each row.
		counts[idx_depth] = make([]int64, width)
	}

	// Initialize random numbers for seeding the hashes at each row.
	aHashes := make([]int64, depth)
	bHashes := make([]int64, depth)
	var i int64
	for i = 0; i < depth; i++ {
		aHashes[i] = rand.Int63()
		bHashes[i] = rand.Int63()
	}

	cm = &CountMin{
		width: width,
		depth: depth,
		count: 0,
		counts: counts,
		aHashes: aHashes,
		bHashes: bHashes,
	}
	return
}

// Update hashes `val` at each row modulo the number of columns and increments
// the value in a 2D matrix by `diff`.
func (cm *CountMin) Update(val, diff int64) {
	// Update the count for this sketch.
	cm.count += diff

	// Update the 2D matrix at each point.
	var i int64
	for i = 0; i < cm.depth; i++ {
		// At each row, hash the value, `val` , and increment the column
		// found at modulo the hash of the number of columns by `diff`.
		cm.counts[i][hash31(cm.aHashes[i], cm.bHashes[i], val) % cm.width] += diff
	}
}


// PointEst returns an estimate of the count of an item by taking the minimum of all
// values found in each row by hashing the `query` modulo the number of columns.
func (cm *CountMin) PointEst(query int64) (answer int64) {
	// For every row, get the minimum of each column that the `query` hashes to
	// modulo the width; this is the estimate of the number of times the `query` value
	// has been added to this sketch.
	answer = cm.counts[0][hash31(cm.aHashes[0], cm.bHashes[0], query) % cm.width]
	var i int64
	for i = 1; i < cm.depth; i++ {
		answer = min(answer, cm.counts[i][hash31(cm.aHashes[i], cm.bHashes[i], query) % cm.width])
	}
	return
}

// PointMed returns an estimate of the count by taking the median estimate
// useful when counts can become negative; depth needs to be larger for this to work well.
func (cm *CountMin) PointMed(query int64) int64 {
	ans := make(int64arr, cm.depth + 1)
	var i int64
	for i = 0; i < cm.depth; i++ {
		ans[i] = cm.counts[i][hash31(cm.aHashes[i], cm.bHashes[i], query) % cm.width]
	}
	sort.Sort(ans)
	return ans[1 + len(ans) / 2]
}
