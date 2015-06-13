package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	sequenses = regexp.MustCompile(`(?:[^/\\]|\\.)*`)
	branches  = regexp.MustCompile(`(?:[^,\\]|\\.)*`)
)

func parseExpr(expr string) ([][]string, error) {
	expr = strings.Replace(expr, `\,`, `\\,`, -1)
	expr = strings.Replace(expr, `\/`, `\\/`, -1)
	expr, err := strconv.Unquote(`"` + expr + `"`)
	if err != nil {
		return nil, err
	}

	sls := sequenses.FindAllString(expr, -1)
	tree := make([][]string, len(sls))
	for si := 0; si < len(sls); si++ {
		bls := branches.FindAllString(sls[si], -1)
		for bi := 0; bi < len(bls); bi++ {
			bls[bi] = strings.Replace(bls[bi], `\,`, `,`, -1)
			bls[bi] = strings.Replace(bls[bi], `\/`, `/`, -1)
			bls[bi] = regexp.QuoteMeta(bls[bi])
		}
		tree[si] = bls
	}
	return tree, nil
}

func newMatcher(expr string) (m *regexp.Regexp, err error) {
	tree, err := parseExpr(expr)
	if err != nil {
		return nil, err
	}

	sls := make([]string, len(tree))
	for si, bls := range tree {
		sls[si] = "(" + strings.Join(bls, "|") + ")"
	}
	return regexp.Compile(strings.Join(sls, ""))
}

func newReplacement(exprFrom, exprTo string) ([]map[string]string, error) {
	from, err := parseExpr(exprFrom)
	if err != nil {
		return nil, err
	}
	to, err := parseExpr(exprTo)
	if err != nil {
		return nil, err
	}

	if len(from) != len(to) {
		return nil, fmt.Errorf("mismatch the number of sequense")
	}

	r := make([]map[string]string, len(from))
	for si := 0; si < len(from); si++ {
		if len(from[si]) != len(to[si]) {
			return nil, fmt.Errorf("mismatch the number of branch[%q]", si)
		}

		r[si] = make(map[string]string)
		for bi := 0; bi < len(from[si]); bi++ {
			if _, exist := r[si][from[si][bi]]; exist {
				return nil, fmt.Errorf("branch[%q] has duplicate item", si)
			}
			r[si][from[si][bi]] = to[si][bi]
		}
	}
	return r, nil
}
