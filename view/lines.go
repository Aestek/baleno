package view

func lines(c []rune, fn func(l int, c []rune) bool) {
	lineStart, lineEnd := 0, 0
	line := 0
	for _, r := range c {
		lineEnd++
		if r == EOL {
			if !fn(line, c[lineStart:lineEnd-1]) {
				return
			}
			line++
			lineStart = lineEnd
		}
	}
	if lineStart != lineEnd || c[lineEnd-1] == '\n' {
		fn(line, c[lineStart:lineEnd])
	}
}

func lineAt(c []rune, target int) (res []rune, ok bool) {
	lines(c, func(ln int, line []rune) bool {
		if ln < target {
			return true
		}

		res = line
		ok = true

		return false
	})
	return
}

func lineCount(c []rune) int {
	if len(c) == 0 {
		return 0
	}

	n := 1
	for _, i := range c {
		if i == EOL {
			n++
		}
	}
	return n
}
