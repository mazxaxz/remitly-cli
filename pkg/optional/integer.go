package optional

// Integer is a wrapper for *int
// based on null.Int, but more 'specific'
type Integer struct {
	Value     int
	Specified bool
}
