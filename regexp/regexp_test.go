package regexp

import "testing"

func TestMatchLine(t *testing.T) {
	tests := []struct {
		input   string
		pattern string
		want    bool
	}{
		{
			input:   "a cow",
			pattern: "a (cat|dog)",
			want:    false,
		},
		{
			input:   "a dog",
			pattern: "a (cat|dog)",
			want:    true,
		},
		{
			input:   "a cat",
			pattern: "a (cat|dog)",
			want:    true,
		},
		{
			input:   "mouse",
			pattern: "(cat|dog|moose)",
			want:    false,
		},
		{
			input:   "apple",
			pattern: "(cat|dog)",
			want:    false,
		},
		{
			input:   "cat",
			pattern: "(cat|dog)",
			want:    true,
		},
		{
			input:   "dog",
			pattern: "(cat|dog)",
			want:    true,
		},
		{
			input:   "dog",
			pattern: "d.g",
			want:    true,
		},
		{
			input:   "dig",
			pattern: "d.g",
			want:    true,
		},
		{
			input:   "cog",
			pattern: "d.g",
			want:    false,
		},
		{
			input:   "dog",
			pattern: "dogs?",
			want:    true,
		},
		{
			input:   "dogs",
			pattern: "dogs?",
			want:    true,
		},
		{
			input:   "cat",
			pattern: "dogs?",
			want:    false,
		},
		{
			input:   "cat",
			pattern: "ca?t",
			want:    true,
		},
		{
			input:   "act",
			pattern: "ca?t",
			want:    true,
		},
		{
			input:   "cag",
			pattern: "ca?t",
			want:    false,
		},
		{
			input:   "apple",
			pattern: "a+",
			want:    true,
		},
		{
			input:   "SaaS",
			pattern: "a+",
			want:    true,
		},
		{
			input:   "dog",
			pattern: "a+",
			want:    false,
		},
		{
			input:   "caats",
			pattern: "ca+ts",
			want:    true,
		},
		{
			input:   "caaats",
			pattern: "ca+t",
			want:    true,
		},
		{
			input:   "act",
			pattern: "ca+t",
			want:    false,
		},
		{
			input:   "dog",
			pattern: "dog$",
			want:    true,
		},
		{
			input:   "dogs",
			pattern: "dog$",
			want:    false,
		},
		{
			input:   "log",
			pattern: "^log",
			want:    true,
		},
		{
			input:   "slog",
			pattern: "^log",
			want:    false,
		},
		{
			input:   "sally has 3 apples",
			pattern: "\\d apple",
			want:    true,
		},
		{
			input:   "sally has 1 orange",
			pattern: "\\d apple",
			want:    false,
		},
		{
			input:   "sally has 124 apples",
			pattern: "\\d\\d\\d apples",
			want:    true,
		},
		{
			input:   "sally has 12 apples",
			pattern: "\\d\\d\\d apples",
			want:    false,
		},
		{
			input:   "sally has 3 dogs",
			pattern: "\\d \\w\\w\\ws",
			want:    true,
		},
		{
			input:   "sally has 4 dogs",
			pattern: "\\d \\w\\w\\ws",
			want:    true,
		},
		{
			input:   "sally has 1 dog",
			pattern: "\\d \\w\\w\\ws",
			want:    false,
		},
		{
			input:   "apple",
			pattern: "[^xyz]",
			want:    true,
		},
		{
			input:   "banana",
			pattern: "[^anb]",
			want:    false,
		},
		{
			input:   "a",
			pattern: "[abcd]",
			want:    true,
		},
		{
			input:   "efgh",
			pattern: "[abcd]",
			want:    false,
		},
		{
			input:   "word",
			pattern: "\\w",
			want:    true,
		},
		{
			input:   "$!?",
			pattern: "\\w",
			want:    false,
		},
		{
			input:   "123",
			pattern: "\\d",
			want:    true,
		},
		{
			input:   "apple",
			pattern: "\\d",
			want:    false,
		},
		{
			input:   "dog",
			pattern: "d",
			want:    true,
		},
		{
			input:   "dog",
			pattern: "f",
			want:    false,
		},
	}

	for _, test := range tests {
		got, err := MatchLine(test.input, test.pattern)
		if err != nil {
			t.Errorf("matchLine(%q, %q). Error: %v", test.input, test.pattern, err)
		} else if test.want != got {
			t.Errorf("matchLine(%q, %q). Expected/got: %v/%v", test.input, test.pattern, test.want, got)
		}
	}
}
