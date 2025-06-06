package logstorage

import (
	"testing"
)

func TestParsePipeSortSuccess(t *testing.T) {
	f := func(pipeStr string) {
		t.Helper()
		expectParsePipeSuccess(t, pipeStr)
	}

	f(`sort`)
	f(`sort rank`)
	f(`sort rank as foo`)
	f(`sort by (x)`)
	f(`sort by (x) limit 10`)
	f(`sort by (x) offset 20 limit 10`)
	f(`sort by (x) offset 20 limit 10 rank as bar`)
	f(`sort by (x desc, y) desc`)
	f(`sort by (a, b) partition by (y, z) limit 10`)
}

func TestParsePipeSortFailure(t *testing.T) {
	f := func(pipeStr string) {
		t.Helper()
		expectParsePipeFailure(t, pipeStr)
	}

	f(`sort a`)
	f(`sort by`)
	f(`sort by(x) foo`)
	f(`sort by(x) limit`)
	f(`sort by(x) limit N`)
	f(`sort by(x) offset`)
	f(`sort by(x) offset N`)
	f(`sort by (x) partition by (y)`)
	f(`sort by (x) limit 2 partition`)
	f(`sort by (x) limit 3 partition by`)
	f(`sort by (x) limit 2 partition (`)
	f(`sort by (x) limit 4 partition (abc`)
	f(`sort by (x) limit 3 partition (abc,`)
	f(`sort by (*)`)
	f(`sort by (x*)`)
	f(`sort by (x) partition by (*)`)
	f(`sort by (x) partition by (y*)`)
}

func TestPipeSort(t *testing.T) {
	f := func(pipeStr string, rows, rowsExpected [][]Field) {
		t.Helper()
		expectPipeResults(t, pipeStr, rows, rowsExpected)
	}

	// Sort by all fields
	f("sort", [][]Field{
		{
			{"_msg", `def`},
			{"a", `1`},
		},
		{
			{"_msg", `abc`},
			{"a", `2`},
		},
	}, [][]Field{
		{
			{"_msg", `abc`},
			{"a", `2`},
		},
		{
			{"_msg", `def`},
			{"a", `1`},
		},
	})

	// Sort by all fields with rank
	f("sort rank x", [][]Field{
		{
			{"_msg", `def`},
			{"a", `1`},
		},
		{
			{"_msg", `abc`},
			{"a", `2`},
		},
	}, [][]Field{
		{
			{"_msg", `abc`},
			{"a", `2`},
			{"x", "1"},
		},
		{
			{"_msg", `def`},
			{"a", `1`},
			{"x", "2"},
		},
	})

	// Sort by a single field
	f("sort by (a asc) asc", [][]Field{
		{
			{"_msg", `abc`},
			{"a", `2`},
		},
		{
			{"_msg", `def`},
			{"a", `1`},
		},
	}, [][]Field{
		{
			{"_msg", `def`},
			{"a", `1`},
		},
		{
			{"_msg", `abc`},
			{"a", `2`},
		},
	})

	// Sort by a in descending order
	f("sort by (a) desc", [][]Field{
		{
			{"_msg", `abc`},
			{"a", `2`},
		},
		{
			{"_msg", `def`},
			{"a", `1`},
		},
	}, [][]Field{
		{
			{"_msg", `abc`},
			{"a", `2`},
		},
		{
			{"_msg", `def`},
			{"a", `1`},
		},
	})

	// Sort by multiple fields
	f("sort by (a, b desc) desc", [][]Field{
		{
			{"_msg", `abc`},
			{"a", `2`},
			{"b", `3`},
		},
		{
			{"_msg", `def`},
			{"a", `1`},
		},
		{
			{"a", `2`},
			{"b", `54`},
		},
	}, [][]Field{
		{
			{"_msg", `abc`},
			{"a", `2`},
			{"b", `3`},
		},
		{
			{"a", `2`},
			{"b", `54`},
		},
		{
			{"_msg", `def`},
			{"a", `1`},
			{"b", ""},
		},
	})

	// Sort by multiple fields with limit
	f("sort by (a, b) limit 1", [][]Field{
		{
			{"_msg", `abc`},
			{"a", `2`},
			{"b", `3`},
		},
		{
			{"_msg", `def`},
			{"a", `1`},
		},
		{
			{"a", `2`},
			{"b", `54`},
		},
	}, [][]Field{
		{
			{"_msg", `def`},
			{"a", `1`},
			{"b", ""},
		},
	})

	// Sort by multiple fields with limit desc
	f("sort by (a, b) desc limit 1", [][]Field{
		{
			{"_msg", `abc`},
			{"a", `2`},
			{"b", `3`},
		},
		{
			{"_msg", `def`},
			{"a", `1`},
		},
		{
			{"a", `2`},
			{"b", `54`},
		},
	}, [][]Field{
		{
			{"a", `2`},
			{"b", `54`},
		},
	})

	// Sort by multiple fields with offset
	f("sort by (a, b) offset 1", [][]Field{
		{
			{"_msg", `abc`},
			{"a", `2`},
			{"b", `3`},
		},
		{
			{"_msg", `def`},
			{"a", `1`},
		},
		{
			{"a", `2`},
			{"b", `54`},
		},
	}, [][]Field{
		{
			{"_msg", `abc`},
			{"a", `2`},
			{"b", `3`},
		},
		{
			{"a", `2`},
			{"b", `54`},
		},
	})

	// Sort by multiple fields with offset and rank
	f("sort by (a, b) offset 1 rank x", [][]Field{
		{
			{"_msg", `abc`},
			{"a", `2`},
			{"b", `3`},
		},
		{
			{"_msg", `def`},
			{"a", `1`},
		},
		{
			{"a", `2`},
			{"b", `54`},
		},
	}, [][]Field{
		{
			{"_msg", `abc`},
			{"a", `2`},
			{"b", `3`},
			{"x", "2"},
		},
		{
			{"a", `2`},
			{"b", `54`},
			{"x", "3"},
		},
	})

	// Sort by multiple fields with offset and limit
	f("sort by (a, b) offset 1 limit 1", [][]Field{
		{
			{"_msg", `abc`},
			{"a", `2`},
			{"b", `3`},
		},
		{
			{"_msg", `def`},
			{"a", `1`},
		},
		{
			{"a", `2`},
			{"b", `54`},
		},
	}, [][]Field{
		{
			{"_msg", `abc`},
			{"a", `2`},
			{"b", `3`},
		},
	})

	// Sort by multiple fields with offset, limit and rank
	f("sort by (a, b) offset 1 limit 1 rank x", [][]Field{
		{
			{"_msg", `abc`},
			{"a", `2`},
			{"b", `3`},
		},
		{
			{"_msg", `def`},
			{"a", `1`},
		},
		{
			{"a", `2`},
			{"b", `54`},
		},
	}, [][]Field{
		{
			{"_msg", `abc`},
			{"a", `2`},
			{"b", `3`},
			{"x", "2"},
		},
	})

	// Sort by multiple fields with offset and limit
	f("order by (a, b) desc offset 2 limit 100", [][]Field{
		{
			{"_msg", `abc`},
			{"a", `2`},
			{"b", `3`},
		},
		{
			{"_msg", `def`},
			{"a", `1`},
		},
		{
			{"a", `2`},
			{"b", `54`},
		},
	}, [][]Field{
		{
			{"_msg", `def`},
			{"a", `1`},
			{"b", ""},
		},
	})

	// Sort with limit and partition
	f(`sort by (a) limit 1 partition by (b)`, [][]Field{
		{
			{"a", "foo"},
			{"b", "x"},
		},
		{
			{"a", "bar"},
			{"b", "x"},
		},
		{
			{"a", "xyz"},
			{"b", "abc"},
		},
	}, [][]Field{
		{
			{"a", "xyz"},
			{"b", "abc"},
		},
		{
			{"a", "bar"},
			{"b", "x"},
		},
	})
}

func TestPipeSortUpdateNeededFields(t *testing.T) {
	f := func(s, allowFilters, denyFilters, allowFiltersExpected, denyFiltersExpected string) {
		t.Helper()
		expectPipeNeededFields(t, s, allowFilters, denyFilters, allowFiltersExpected, denyFiltersExpected)
	}

	// all the needed fields
	f("sort", "*", "", "*", "")
	f("sort rank x", "*", "", "*", "")
	f("sort by(s1,s2)", "*", "", "*", "")
	f("sort by(s1,s2) rank as x", "*", "", "*", "x")
	f("sort by(s1,s2) limit 1 partition by (x) rank as x", "*", "", "*", "")
	f("sort by(x,s2) rank as x", "*", "", "*", "")
	f("sort by(x,s2) limit 1 partition by (y) rank as x", "*", "", "*", "")

	// all the needed fields, unneeded fields do not intersect with src
	f("sort by(s1,s2)", "*", "f1,f2", "*", "f1,f2")
	f("sort by(s1,s2) limit 1 partition by (f1,s1)", "*", "f1,f2", "*", "f2")
	f("sort by(s1,s2) rank as x", "*", "f1,f2", "*", "f1,f2,x")
	f("sort by(s1,s2) limit 1 partition by (x,y) rank as x", "*", "f1,f2", "*", "f1,f2")
	f("sort by(x,s2) rank as x", "*", "f1,f2", "*", "f1,f2")
	f("sort by(x,s2) limit 1 partition by (x,y) rank as x", "*", "f1,f2", "*", "f1,f2")

	// all the needed fields, unneeded fields intersect with src
	f("sort by(s1,s2)", "*", "s1,f1,f2", "*", "f1,f2")
	f("sort by(s1,s2) limit 1 partition by (f1,f3)", "*", "s1,f1,f2", "*", "f2")
	f("sort by(s1,s2) rank as x", "*", "s1,f1,f2", "*", "f1,f2,x")
	f("sort by(s1,s2) limit 1 partition by (f1) rank as x", "*", "s1,f1,f2", "*", "f2,x")
	f("sort by(x,s2) rank as x", "*", "s1,f1,f2", "*", "f1,f2,s1")
	f("sort by(x,s2) limit 1 partition by (f2) rank as x", "*", "s1,f1,f2", "*", "f1,s1")

	// needed fields do not intersect with src
	f("sort by(s1,s2)", "f1,f2", "", "s1,s2,f1,f2", "")
	f("sort by(s1,s2) limit 1 partition by (f1,f3)", "f1,f2", "", "s1,s2,f1,f2,f3", "")
	f("sort by(s1,s2) rank as x", "f1,f2", "", "s1,s2,f1,f2", "")
	f("sort by(s1,s2) limit 1 partition by (s1,f2,f3,x) rank as x", "f1,f2", "", "s1,s2,f1,f2,f3,x", "")

	// needed fields intersect with src
	f("sort by(s1,s2)", "s1,f1,f2", "", "s1,s2,f1,f2", "")
	f("sort by(s1,s2) limit 1 partition by (s1,s3)", "s1,f1,f2", "", "s1,s2,s3,f1,f2", "")
	f("order by(s1,s2) rank as x", "s1,f1,f2,x", "", "s1,s2,f1,f2", "")
	f("order by(s1,s2) limit 1 partition by (x,y) rank as x", "s1,f1,f2,x", "", "s1,s2,f1,f2,x,y", "")
}
