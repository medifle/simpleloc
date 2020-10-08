package processor

type stateType uint32

const (
	SBlank stateType = iota
	SCode
	SSingleComment
	SSingleCodeComment
	SMultiComment
	SMultiCommentBlank
	SMultiCodeComment
	SMultiCodeCommentBlank
	SString
)

// state machine for counting lines
type machine struct {
	state    stateType
	index    int
	endIndex int
	endMatch []byte
}

func (m *machine) resetState() {
	if m.state == SMultiComment || m.state == SMultiCodeComment {
		m.state = SMultiComment
	} else if m.state == SString {
		m.state = SString
	} else {
		m.state = SBlank
	}
}

func isWhitespace(currentByte byte) bool {
	if currentByte != ' ' && currentByte != '\t' && currentByte != '\n' && currentByte != '\r' {
		return false
	}

	return true
}

func checkForMatchSingle(file *File, m *machine, match []byte) bool {
	currentByte := file.Content[m.index]

	if currentByte == match[0] {
		potentialMatch := true
		for j := 1; j < len(match); j++ {
			if m.index+j > m.endIndex || file.Content[m.index+j] != match[j] {
				potentialMatch = false
				break
			}
		}

		if potentialMatch {
			return true
		}
	}

	return false
}

// Match single comment marker
// Assume index <= endIndex in machine
func checkForMatch(file *File, m *machine, matches [][]byte) bool {
	currentByte := file.Content[m.index]

	for i := 0; i < len(matches); i++ {
		if currentByte == matches[i][0] {
			potentialMatch := true

			for j := 1; j < len(matches[i]); j++ {
				if m.index+j > m.endIndex || file.Content[m.index+j] != matches[i][j] {
					potentialMatch = false
					break
				}
			}

			if potentialMatch {
				return true
			}
		}
	}

	return false
}

// Match multiline comment open marker, e.g "/*" for ["/*", "*/"] pair
// Assume index <= endIndex in machine
func checkForMatchMultiOpen(file *File, m *machine, matches []OpenClose) int {
	currentByte := file.Content[m.index]

	// In case we have more than one type of matches in one language
	for i := 0; i < len(matches); i++ {
		if currentByte == matches[i].Open[0] {
			potentialMatch := true

			for j := 1; j < len(matches[i].Open); j++ {
				if m.index+j > m.endIndex || file.Content[m.index+j] != matches[i].Open[j] {
					potentialMatch = false
					break
				}
			}

			if potentialMatch {
				m.endMatch = matches[i].Close
				return len(matches[i].Open)
			}
		}
	}

	return 0
}

func checkForQuotesMatch(file *File, m *machine, matches []Quote) int {
	currentByte := file.Content[m.index]

	// In case we have more than one type of matches in one language
	for i := 0; i < len(matches); i++ {
		if currentByte == matches[i].Start[0] {
			potentialMatch := true

			for j := 1; j < len(matches[i].Start); j++ {
				if m.index+j > m.endIndex || file.Content[m.index+j] != matches[i].Start[j] {
					potentialMatch = false
					break
				}
			}

			if potentialMatch {
				m.endMatch = []byte(matches[i].End)
				return len(matches[i].Start)
			}
		}
	}

	return 0
}

func singleCommentState(file *File, m *machine, langFeature LanguageFeature) {
	for ; m.index < m.endIndex; m.index++ {
		// Remain in the same comment state when we hit the end of line
		if file.Content[m.index] == '\n' {
			return
		}

		// Match TODOs
		if checkForMatchSingle(file, m, []byte("TODO")) {
			file.Todo++
		}
	}
}

func multiCommentState(file *File, m *machine, langFeature LanguageFeature) {
	for ; m.index < m.endIndex; m.index++ {
		// Remain in the same comment state when we hit the end of line
		if file.Content[m.index] == '\n' {
			return
		}

		// Match TODOs
		if checkForMatchSingle(file, m, []byte("TODO")) {
			file.Todo++
		}

		if checkForMatchSingle(file, m, m.endMatch) {
			if m.state == SMultiCodeComment {
				m.state = SMultiCodeCommentBlank
			} else {
				m.state = SMultiCommentBlank
			}
			offsetJump := len(m.endMatch)
			m.index += offsetJump - 1
			m.endMatch = []byte("") // reset m.endMatch
			file.BlockComment++
			return
		}
	}
}

func stringState(file *File, m *machine) {
	for ; m.index < m.endIndex; m.index++ {
		// Remain in String state when we hit the end of line
		if file.Content[m.index] == '\n' {
			return
		}

		// Go back to Code state
		// It is safe to check escape char because if we enter string state at index 0, we need to go to next byte first
		// before enter this function call
		if file.Content[m.index-1] != '\\' && checkForMatchSingle(file, m, m.endMatch) {
			m.state = SCode
			return
		}
	}
}

func codeState(file *File, m *machine, langFeature LanguageFeature) {
	for ; m.index < m.endIndex; m.index++ {
		// Remain in Code state when we hit the end of line
		if file.Content[m.index] == '\n' {
			return
		}

		// Go to String state
		offsetJump := checkForQuotesMatch(file, m, langFeature.Quotes)
		if offsetJump != 0 {
			m.state = SString
			return
		}

		// Go to SingleCodeComment state
		if checkForMatch(file, m, langFeature.SingleLineComments) {
			m.state = SSingleCodeComment
			return
		}

		// Go to MultiCodeComment state
		offsetJump = checkForMatchMultiOpen(file, m, langFeature.MultiLineComments)
		if offsetJump != 0 {
			m.state = SMultiCodeComment
			m.index += offsetJump - 1
			return
		}
	}
}

func blankState(file *File, m *machine, langFeature LanguageFeature) {
	// Go to MultiComment state
	offsetJump := checkForMatchMultiOpen(file, m, langFeature.MultiLineComments)
	if offsetJump != 0 {
		m.state = SMultiComment
		m.index += offsetJump - 1
		return
	}

	// Go to SingleComment state
	if checkForMatch(file, m, langFeature.SingleLineComments) {
		m.state = SSingleComment
		return
	}

	// Go to String state
	offsetJump = checkForQuotesMatch(file, m, langFeature.Quotes)
	if offsetJump != 0 {
		m.state = SString
		return
	}

	// Go to Code state
	m.state = SCode
}

func countStats(file *File) {
	// If file is empty, no line is counted
	if len(file.Content) == 0 {
		return
	}

	langFeature := LangFeature

	// Default state starts from Blank state
	m := &machine{state: SBlank}
	m.endIndex = len(file.Content) - 1

	// Read file content byte by byte, change machine state if necessary and calc the stats
	for m.index = 0; m.index < len(file.Content); m.index++ {
		currentByte := file.Content[m.index]
		if !isWhitespace(currentByte) {
			switch m.state {
			case SCode:
				codeState(file, m, langFeature)
			case SString:
				stringState(file, m)
			case SSingleComment, SSingleCodeComment: // All SingleComment-ish states enter the same function
				singleCommentState(file, m, langFeature)
			case SMultiComment, SMultiCodeComment: // All MultiComment-ish states enter the same function
				multiCommentState(file, m, langFeature)
			case SBlank, SMultiCommentBlank, SMultiCodeCommentBlank: // All blank-ish states enter the same function
				blankState(file, m, langFeature)
			}
		}

		// Reach the end of a line so calculate the stats according to what state we are in
		if file.Content[m.index] == '\n' || m.index >= m.endIndex {
			file.Lines++
			switch m.state {
			case SCode, SString:
				file.Code++
			case SSingleCodeComment:
				file.Code++
				file.SingleLineComment++
			case SMultiCodeComment, SMultiCodeCommentBlank:
				file.Code++
				file.MultiLineComment++
			case SSingleComment:
				file.SingleLineComment++
			case SMultiComment, SMultiCommentBlank:
				file.MultiLineComment++
			}
			m.resetState()
		}
	}
}
