package main

import (
	"testing"
)

var simpleTests = map[string]string{
	``:     "",
	`\n`:   "\n",
	`\r`:   "\r",
	`\t`:   "\t",
	`\\`:   "\\",
	`\\n`:  "\\n",
	`\\\\`: "\\\\",
	`package main\n\nimport (\n\t"fmt"\n)\n\nfunc main() {\n\tfmt.Print("He\\\\o there, wor\\d!\r\\n")\n}\n`: `package main

import (
	"fmt"
)

func main() {
	fmt.Print("He\\o there, wor\d!` + "\r" + `\n")
}
`,
}

var errTests = map[string]string{
	`\`:     `unescaped backslash ("\")`,
	`\\\`:   `unescaped backslash ("\")`,
	`\\\\\`: `unescaped backslash ("\")`,
	"\\\n":  "invalid escape: \\\n",
	`\ `:    `invalid escape: \ `,
	`\0`:    `invalid escape: \0`,
	`\a`:    `invalid escape: \a`,
}

func TestEscape(t *testing.T) {
	for in, expect := range simpleTests {
		out, err := Escape(in)
		if err != nil || out != expect {
			t.Fatalf(`Escape(%q) = %q, %v; expected %q, %v`,
				in, out, err, expect, nil)
		}
	}
}

func TestEscapeErr(t *testing.T) {
	for in, expect := range errTests {
		_, err := Escape(in)
		if err == nil || err.Error() != expect {
			t.Fatalf(`Escape(%q) = _, %v; expected _, %v`,
				in, err, expect)
		}
	}
}

/*
func BenchmarkEscape(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = uniq1(randSlice())
	}
}
*/
