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

// Language is a struct represents the config of each language in languages.json
type Language struct {
	Extensions  []string   `json:"extensions"`
	LineComment []string   `json:"line_comment"`
	MultiLine   [][]string `json:"multi_line"`
	Quotes      []Quote    `json:"quotes"`
}

const (
	TODO int = 1
)

// TODO: this is a todo

/*
 * sdfsdf
 * TODO:
 */

// LanguageFeature is a struct which is converted from Language to do matching
type LanguageFeature struct {
	SingleLineComments [][]byte    // trailing
	MultiLineComments  []OpenClose // TODO trailing
	Quotes             []Quote     /*
	 * sdfsdf
	 */
}

type OpenClose struct {
	Open  []byte
	Close []byte
}

// Quote is used for string check
type Quote struct {
	Start string `json:"start"`
	End   string `json:"end"`
}