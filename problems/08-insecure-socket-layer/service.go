package main

import (
	"cmp"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var quantityRegex = regexp.MustCompile(`^(\d+)x\s.*$`)

func extractQuantity(s string) int {
	value := quantityRegex.FindStringSubmatch(s)[1]
	qty, _ := strconv.Atoi(value)
	return qty
}

// Find toy that we need to make the most copies of.
func FindToy(request string) string {
	toys := strings.Split(request, ",")
	slices.SortFunc(toys, func(a, b string) int {
		return cmp.Compare(extractQuantity(b), extractQuantity(a))
	})
	return toys[0]
}
