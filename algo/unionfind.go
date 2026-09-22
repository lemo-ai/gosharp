package algo

// UnionFind is a disjoint-set (union-find) structure with path compression
// and union by rank.
type UnionFind struct {
	parent []int
	rank   []int
	count  int
}

// NewUnionFind creates a UnionFind for n elements numbered [0, n).
func NewUnionFind(n int) *UnionFind {
	parent := make([]int, n)
	rank := make([]int, n)
	for i := range parent {
		parent[i] = i
	}
	return &UnionFind{parent: parent, rank: rank, count: n}
}

// Find returns the representative of x.
func (u *UnionFind) Find(x int) int {
	if u.parent[x] != x {
		u.parent[x] = u.Find(u.parent[x])
	}
	return u.parent[x]
}

// Union merges the sets containing a and b. Returns true if a merge happened.
func (u *UnionFind) Union(a, b int) bool {
	ra, rb := u.Find(a), u.Find(b)
	if ra == rb {
		return false
	}
	if u.rank[ra] < u.rank[rb] {
		ra, rb = rb, ra
	}
	u.parent[rb] = ra
	if u.rank[ra] == u.rank[rb] {
		u.rank[ra]++
	}
	u.count--
	return true
}

// Connected reports whether a and b are in the same set.
func (u *UnionFind) Connected(a, b int) bool {
	return u.Find(a) == u.Find(b)
}

// Count returns the number of disjoint sets.
func (u *UnionFind) Count() int {
	return u.count
}
