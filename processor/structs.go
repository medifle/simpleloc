package processor

type File struct {
	Language          string
	Extension         string
	Content           []byte
	Lines             int64
	Code              int64
	SingleLineComment int64
	MultiLineComment  int64
	BlockComment      int64
	Todo              int64
}
