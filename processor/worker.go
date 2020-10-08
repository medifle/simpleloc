package processor

type stateType uint32

const (
	SBlank stateType = iota
	SCode
	SSingleComment
	SSingleHybridComment
	SMultiComment
	SMultiHybridComment
	SString
)

// A state machine for counting lines
type machine struct {
	state stateType
}

func (m *machine) resetState() {
	if m.state == SMultiComment || m.state == SMultiHybridComment {
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

func countStats(file *File) {
	// If file is empty, no line is counted
	if len(file.Content) == 0 {
		return
	}

	endIndex := len(file.Content) - 1

	// Default state starts from Blank
	machine := &machine{state: SBlank}

	// Read file content byte by byte, change machine state if necessary and calc the stats
	for i := 0; i < len(file.Content); i++ {
		curByte := file.Content[i]
		if !isWhitespace(curByte) {
			switch machine.state {
			case SCode:
				// TODO
			case SString:
				//
			case SMultiComment, SMultiHybridComment:
				//

			case SBlank:
				// Blank to SingleLineComment, MultiLineComment or Code
			}
		}

		// Reach the end of a line so calculate the stats according to what state we are in
		if file.Content[i] == '\n' || i >= endIndex {
			file.Lines++

			switch machine.state {
			case SCode, SString:
				file.Code++
			case SSingleHybridComment:
				file.Code++
				file.SingleLineComment++
			case SMultiHybridComment:
				file.Code++
				file.MultiLineComment++
			case SSingleComment:
				file.SingleLineComment++
			case SMultiComment:
				file.MultiLineComment++
			}
			machine.resetState()
		}
	}
}
