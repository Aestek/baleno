package buffer

type Buffer interface {
	Read(from, to int) ([]rune, error)
	Insert(at int, contents []rune) error
	Delete(at int, n int) error
	Length() int
	ReadOnly() bool
	Index(IndexDef) []int
}
