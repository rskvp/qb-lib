package semantic_search_test

import (
	"fmt"
	"testing"

	"github.com/rskvp/qb-lib/qb_dbal/semantic_search"
)

func TestToKeywords(t *testing.T) {
	keywords := semantic_search.ToKeywords("hello this is a text to tokenize in keywords!!")
	fmt.Println(keywords)
}