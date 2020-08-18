package samasemodels

import (
	"fmt"
	"strings"
)

type Options struct {
	Filters    []Filter
	Pagination *Pagination
	Sorting    *Sorting
}

type Filter struct {
	By       string
	Value    interface{}
	Operator string // e.g : LIKE , = , > , < etc
	Required bool   // false means , where given multiple filters it allowed to return false
}
type Sorting struct {
}
type Pagination struct {
}

func ParseOptionsToSQLQuery(opts *Options) (query string, args []interface{}) {
	if opts == nil {
		return "", nil
	}
	if opts.Filters != nil && len(opts.Filters) > 0 {
		query = " WHERE "
		for _, ft := range opts.Filters {
			if ft != (Filter{}) {
				if ft.Operator == "LIKE" {
					for _, v := range strings.Split(fmt.Sprint(ft.Value), " ") {
						query += fmt.Sprint(ft.By, " ", ft.Operator, " ? AND ")
						args = append(args, fmt.Sprint("%", v, "%"))
					}
				} else {
					query += fmt.Sprint(ft.By, " ", ft.Operator, " ? AND ")
					args = append(args, ft.Value)
				}
			}
		}
		query = query[:len(query)-4]
	}
	return
}

func ParseFiltersToSQLQuery(fts []Filter) (query string, args []interface{}) {
	if fts != nil && len(fts) != 0 {
		query = " WHERE "
	}
	for range fts {
		for _, ft := range fts {
			if ft.Operator == "LIKE" {
				for _, v := range strings.Split(fmt.Sprint(ft.Value), " ") {
					query += fmt.Sprint(ft.By, " ", ft.Operator, " ? AND ")
					args = append(args, fmt.Sprint("%", v, "%"))
				}
			} else {
				query += fmt.Sprint(ft.By, " ", ft.Operator, " ? AND ")
				args = append(args, ft.Value)
			}
		}
		query = query[:len(query)-4]
	}
	return
}
