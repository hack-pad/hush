package hush

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"mvdan.cc/sh/v3/syntax"
)

type completion struct {
	Completion string
	Start, End int
}

func getCompletions(line string, cursor int) []completion {
	completions, err := getCompletionsErr(line, cursor)
	if err != nil {
		log.Print("Failed completions:", err, "\r\n")
		return nil
	}
	return completions
}

func getCompletionsErr(line string, cursor int) ([]completion, error) {
	parser := syntax.NewParser()
	var stmts []*syntax.Stmt
	err := parser.Stmts(strings.NewReader(line), func(stmt *syntax.Stmt) bool {
		if int(stmt.Pos().Offset()) <= cursor && int(stmt.End().Offset()) >= cursor {
			stmts = append(stmts, stmt)
		}
		return true
	})
	if err != nil || len(stmts) == 0 {
		return nil, err
	}
	cursorStmt := stmts[0]
	cursorStmtStr := formatStmt(line, cursorStmt)
	cursorStmtOffset := int(cursorStmt.Pos().Offset())
	cursor -= cursorStmtOffset

	var commandWord, cursorWord *syntax.Word
	err = parser.Words(strings.NewReader(cursorStmtStr), func(word *syntax.Word) bool {
		if commandWord == nil {
			commandWord = word
		}
		if int(word.Pos().Offset()) <= cursor && int(word.End().Offset()) >= cursor {
			cursorWord = word
		}
		return true
	})
	if err != nil || cursorWord == nil {
		return nil, err
	}

	commandWordStr, err := evalWord(commandWord.Parts)
	if err != nil {
		return nil, err
	}
	cursorWordStr, err := evalWord(cursorWord.Parts)
	if err != nil {
		return nil, err
	}

	return getStatementCompletions(
		commandWordStr,
		cursorWordStr,
		cursorStmtOffset+int(cursorWord.Pos().Offset()),
		cursorStmtOffset+int(cursorWord.End().Offset()))
}

func getStatementCompletions(commandName string, word string, start, end int) ([]completion, error) {
	switch {
	case strings.Contains(word, "/"):
		dir := word
		filter := false
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			dir = filepath.Dir(dir)
			filter = true
		}
		dirEntries, err := os.ReadDir(dir)
		if err != nil {
			return nil, nil
		}
		var completions []completion
		for _, d := range dirEntries {
			base := filepath.Base(word)
			name := d.Name()
			if !filter || strings.HasPrefix(name, base) {
				file := fileJoin(dir, name)
				if d.IsDir() {
					file += string(filepath.Separator)
				}
				completions = append(completions, completion{
					Completion: file,
					Start:      start,
					End:        end,
				})
			}
		}
		return completions, nil
	default:
		return nil, nil
	}
}

func fileJoin(a, b string) string {
	if a == "." {
		return "." + string(filepath.Separator) + b
	}
	return filepath.Join(a, b)
}
