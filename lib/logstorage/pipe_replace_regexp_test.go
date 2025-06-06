package logstorage

import (
	"regexp"
	"testing"
)

func TestParsePipeReplaceRegexpSuccess(t *testing.T) {
	f := func(pipeStr string) {
		t.Helper()
		expectParsePipeSuccess(t, pipeStr)
	}

	f(`replace_regexp (foo, bar)`)
	f(`replace_regexp ("foo[^ ]+bar|baz", "bar${1}x$0")`)
	f(`replace_regexp (" ", "") at x`)
	f(`replace_regexp if (x:y) ("-", ":") at a`)
	f(`replace_regexp (" ", "") at x limit 10`)
	f(`replace_regexp if (x:y) (" ", "") at foo limit 10`)
}

func TestParsePipeReplaceRegexpFailure(t *testing.T) {
	f := func(pipeStr string) {
		t.Helper()
		expectParsePipeFailure(t, pipeStr)
	}

	f(`replace_regexp`)
	f(`replace_regexp (foo, bar) at *`)
	f(`replace_regexp (foo, bar) at x*`)
	f(`replace_regexp if`)
	f(`replace_regexp foo`)
	f(`replace_regexp (`)
	f(`replace_regexp (foo`)
	f(`replace_regexp (foo,`)
	f(`replace_regexp(foo,bar`)
	f(`replace_regexp(foo,bar,baz)`)
	f(`replace_regexp(foo,bar) abc`)
	f(`replace_regexp(bar,baz) limit`)
	f(`replace_regexp(bar,baz) limit N`)
	f(`replace_regexp ("foo[", "bar")`)
}

func TestPipeReplaceRegexp(t *testing.T) {
	f := func(pipeStr string, rows, rowsExpected [][]Field) {
		t.Helper()
		expectPipeResults(t, pipeStr, rows, rowsExpected)
	}

	// replace_regexp with placeholders
	f(`replace_regexp ("foo(.+?)bar", "q-$1-x")`, [][]Field{
		{
			{"_msg", `abc foo a bar foobar foo b bar`},
			{"bar", `cde`},
		},
		{
			{"_msg", `1234`},
		},
	}, [][]Field{
		{
			{"_msg", `abc q- a -x q-bar foo b -x`},
			{"bar", `cde`},
		},
		{
			{"_msg", `1234`},
		},
	})

	// replace_regexp without limits at _msg
	f(`replace_regexp ("[_/]", "-")`, [][]Field{
		{
			{"_msg", `a_bc_d/ef`},
			{"bar", `cde`},
		},
		{
			{"_msg", `a_bc_d/ef`},
		},
		{
			{"_msg", `1234`},
		},
		{
			{"_msg", `1234`},
		},
	}, [][]Field{
		{
			{"_msg", `a-bc-d-ef`},
			{"bar", `cde`},
		},
		{
			{"_msg", `a-bc-d-ef`},
		},
		{
			{"_msg", `1234`},
		},
		{
			{"_msg", `1234`},
		},
	})

	// replace_regexp with limit 1 at foo
	f(`replace_regexp ("[_/]", "-") at foo limit 1`, [][]Field{
		{
			{"foo", `a_bc_d/ef`},
			{"bar", `cde`},
		},
		{
			{"foo", `1234`},
		},
	}, [][]Field{
		{
			{"foo", `a-bc_d/ef`},
			{"bar", `cde`},
		},
		{
			{"foo", `1234`},
		},
	})

	// replace_regexp with limit 100 at foo
	f(`replace_regexp ("[_/]", "-") at foo limit 100`, [][]Field{
		{
			{"foo", `a_bc_d/ef`},
			{"bar", `cde`},
		},
		{
			{"foo", `1234`},
		},
	}, [][]Field{
		{
			{"foo", `a-bc-d-ef`},
			{"bar", `cde`},
		},
		{
			{"foo", `1234`},
		},
	})

	// conditional replace_regexp at foo
	f(`replace_regexp if (bar:abc) ("[_/]", "") at foo`, [][]Field{
		{
			{"foo", `a_bc_d/ef`},
			{"bar", `cde`},
		},
		{
			{"foo", `123_45/6`},
			{"bar", "abc"},
		},
	}, [][]Field{
		{
			{"foo", `a_bc_d/ef`},
			{"bar", `cde`},
		},
		{
			{"foo", `123456`},
			{"bar", "abc"},
		},
	})
}

func TestPipeReplaceRegexpUpdateNeededFields(t *testing.T) {
	f := func(s string, allowFilters, denyFilters, allowFiltersExpected, denyFiltersExpected string) {
		t.Helper()
		expectPipeNeededFields(t, s, allowFilters, denyFilters, allowFiltersExpected, denyFiltersExpected)
	}

	// all the needed fields
	f(`replace_regexp ("a", "b") at x`, "*", "", "*", "")
	f(`replace_regexp if (f1:q) ("a", "b") at x`, "*", "", "*", "")

	// unneeded fields do not intersect with at field
	f(`replace_regexp ("a", "b") at x`, "*", "f1,f2", "*", "f1,f2")
	f(`replace_regexp if (f3:q) ("a", "b") at x`, "*", "f1,f2", "*", "f1,f2")
	f(`replace_regexp if (f2:q) ("a", "b") at x`, "*", "f1,f2", "*", "f1")

	// unneeded fields intersect with at field
	f(`replace_regexp ("a", "b") at x`, "*", "x,y", "*", "x,y")
	f(`replace_regexp if (f1:q) ("a", "b") at x`, "*", "x,y", "*", "x,y")
	f(`replace_regexp if (x:q) ("a", "b") at x`, "*", "x,y", "*", "x,y")
	f(`replace_regexp if (y:q) ("a", "b") at x`, "*", "x,y", "*", "x,y")

	// needed fields do not intersect with at field
	f(`replace_regexp ("a", "b") at x`, "f2,y", "", "f2,y", "")
	f(`replace_regexp if (f1:q) ("a", "b") at x`, "f2,y", "", "f2,y", "")

	// needed fields intersect with at field
	f(`replace_regexp ("a", "b") at y`, "f2,y", "", "f2,y", "")
	f(`replace_regexp if (f1:q) ("a", "b") at y`, "f2,y", "", "f1,f2,y", "")
}

func TestAppendReplaceRegexp(t *testing.T) {
	f := func(s, reStr, replacement string, limit int, resultExpected string) {
		t.Helper()

		re := regexp.MustCompile(reStr)
		result := appendReplaceRegexp(nil, s, re, replacement, uint64(limit))
		if string(result) != resultExpected {
			t.Fatalf("unexpected result for appendReplaceRegexp(%q, %q, %q, %d)\ngot\n%s\nwant\n%s", s, reStr, replacement, limit, result, resultExpected)
		}
	}

	f("", "", "", 0, "")                               // empty input, empty pattern and replacement
	f("", "foo", "bar", 0, "")                         // empty input, non-matching pattern
	f("abc", "foo", "bar", 0, "abc")                   // no match in non-empty input
	f("foo", "fo+", "bar", 0, "bar")                   // entire match replaced
	f("foox", "fo+", "bar", 0, "barx")                 // partial prefix match replaced
	f("afoo", "fo+", "bar", 0, "abar")                 // match in the middle
	f("afoox", "fo+", "bar", 0, "abarx")               // match in the middle, suffix preserved
	f("foo-bar/baz", "[-/]", "_", 0, "foo_bar_baz")    // replace multiple separators
	f("foo bar/ baz  ", "[ /]", "", 2, "foobar baz  ") // limit number of matches
	f("abcd", `\d*`, "foo", 0, "fooafoobfoocfoodfoo")  // pattern matches empty string between chars
	f("hello", `^|$`, "X", 0, "XhelloX")               // match start and end of string using anchors

	// use of capture groups with full replacement
	f("afoo abc barz", "a([^ ]+)", "b${1}x", 0, "bfoox bbcx bbrzx") // capture and reuse group in replacement
	f("afoo abc barz", "a([^ ]+)", "b${1}x", 1, "bfoox abc barz")   // same as above, but limit = 1
}
