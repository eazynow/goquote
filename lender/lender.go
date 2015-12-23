package lender

// Lender is a structure representing an individual lender.
type Lender struct {
	Name      string
	Rate      float64
	Available int
}

// Borrow is a function used to determine how much a lender can lend. The
// amount is how much the requester wishes to borrow.
// Returns the amount the lender can borrow. This will be either the full
// amount if possible, or the maximum the lender has available if not.
func (l *Lender) Borrow(amount int) int {
	if amount > l.Available {
		// Not enough funds to lend full amount, so return how much they can lend.
		return l.Available
	} else {
		// Enough funds, so return the full amount.
		return amount
	}
}

// Lenders represents a slice of lenders. The type also has methods to support
// sort interface
type Lenders []Lender

// Len returns the size of the lender slice. Used to satisfy the sort interface.
func (slice Lenders) Len() int {
	return len(slice)
}

// Less compares the lender at index i to index j and determines which is
// less and therefore higher in the list. Used to satisfy the sort interface.
// Returns true if i should be ranked higher than j in the sort.
func (slice Lenders) Less(i, j int) bool {

	if slice[i].Rate != slice[j].Rate {
		// Base the sorting on the lowest rate.
		return slice[i].Rate < slice[j].Rate
	} else {
		// In the instance of the same rate, sort by who has more
		// funds available to reduce number of lenders
		return slice[i].Available > slice[j].Available
	}

}

// Swap switches the position of lender at index i with lender at index j
// in the slice. Used to satisfy the sort interface.
func (slice Lenders) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
