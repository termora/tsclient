package tsclient

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

// SearchData is used in (*Client).Search
type SearchData struct {
	// The query text to search for in the collection.
	Query string
	// One or more string / string[] fields that should be queried against. Separate multiple fields with a comma: company_name, country
	QueryBy []string

	// The relative weight to give each query_by field when ranking results.
	// This can be used to boost fields in priority, when looking for matches.
	QueryByWeights []int

	// If false, indicates that the last word in the query should be treated as a prefix,
	// and not as a whole word.
	// This is necessary for building autocomplete and instant search interfaces.
	// Set this to true to disable prefix searching for all queried fields.
	NoPrefix bool

	// Filter conditions for refining your search results.
	FilterBy string

	// A list of numerical fields and their corresponding sort orders that will be used for ordering your results.
	// Up to 3 sort fields can be specified.
	SortBy []string

	// A list of fields that will be used for faceting your results on.
	FacetBy []string

	// Maximum number of facet values to be returned.
	MaxFacetValues uint64

	// Facet values that are returned can now be filtered via this parameter.
	// The matching facet text is also highlighted.
	FacetQuery string

	// By default, Typesense prioritizes documents whose field value matches exactly with the query.
	// Set this parameter to true to disable this behavior.
	NoPrioritizeExactMatch bool

	// Results from this specific page number would be fetched.
	Page uint64
	// Number of results to fetch per page.
	PerPage uint64

	// You can aggregate search results into groups or buckets by specify one or more group_by fields.
	// NOTE: To group on a particular field, it must be a faceted field.
	GroupBy []string

	// Maximum number of hits to be returned for every group.
	// If GroupLimit is set as K, only the top K hits in each group are returned in the response.
	GroupLimit uint64

	// list of fields from the document to include in the search result.
	IncludeFields []string
	// list of fields from the document to exclude in the search result.
	ExcludeFields []string
}

// Search searches the collection.
func (c *Client) Search(collection string, data SearchData) (res SearchResult, err error) {
	v := url.Values{
		"q":                      {data.Query},
		"query_by":               {strings.Join(data.QueryBy, ",")},
		"prefix":                 {strconv.FormatBool(!data.NoPrefix)},
		"prioritize_exact_match": {strconv.FormatBool(!data.NoPrioritizeExactMatch)},
	}

	if len(data.QueryByWeights) > 0 {
		s := []string{}
		for _, i := range data.QueryByWeights {
			s = append(s, strconv.Itoa(i))
		}
		v["query_by_weights"] = []string{strings.Join(s, ",")}
	}

	if data.FilterBy != "" {
		v["filter_by"] = []string{data.FilterBy}
	}

	if len(data.SortBy) > 0 {
		v["sort_by"] = []string{strings.Join(data.SortBy, ",")}
	}

	if len(data.FacetBy) > 0 {
		v["facet_by"] = []string{strings.Join(data.FacetBy, ",")}
	}

	if data.MaxFacetValues != 0 {
		v["max_facet_values"] = []string{strconv.FormatUint(data.MaxFacetValues, 10)}
	}

	if data.FacetQuery != "" {
		v["facet_query"] = []string{data.FacetQuery}
	}

	if data.Page != 0 {
		v["page"] = []string{strconv.FormatUint(data.Page, 0)}
	}

	if data.PerPage != 0 {
		v["per_page"] = []string{strconv.FormatUint(data.PerPage, 0)}
	}

	if len(data.GroupBy) > 0 {
		v["group_by"] = []string{strings.Join(data.GroupBy, ",")}
	}

	if data.GroupLimit != 0 {
		v["group_limit"] = []string{strconv.FormatUint(data.GroupLimit, 0)}
	}

	resp, err := c.Request("GET", "/collections/"+collection+"/documents/search", WithURLValues(v))
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &res)
	return
}

// SearchResult is the result returned from a search
type SearchResult struct {
	FacetCounts int
}
