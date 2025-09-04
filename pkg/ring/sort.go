package ring

// Len returns the length of the uints array.
func (n Nodes) Len() int { return len(n) }

// Less returns true if element i is less than element j.
func (n Nodes) Less(i, j int) bool { return n[i].HashId < n[j].HashId }

// Swap exchanges elements i and j.
func (n Nodes) Swap(i, j int) { n[i], n[j] = n[j], n[i] }