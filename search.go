package tsclient

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/termora/tsclient/utils/jsonutil"
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
	MaxFacetValues int

	// Facet values that are returned can now be filtered via this parameter.
	// The matching facet text is also highlighted.
	FacetQuery string

	// By default, Typesense prioritizes documents whose field value matches exactly with the query.
	// Set this parameter to true to disable this behavior.
	NoPrioritizeExactMatch bool

	// Results from this specific page number would be fetched.
	Page int
	// Number of results to fetch per page.
	PerPage int

	// You can aggregate search results into groups or buckets by specify one or more group_by fields.
	// NOTE: To group on a particular field, it must be a faceted field.
	GroupBy []string

	// Maximum number of hits to be returned for every group.
	// If GroupLimit is set as K, only the top K hits in each group are returned in the response.
	GroupLimit int

	// list of fields from the document to include in the search result.
	IncludeFields []string
	// list of fields from the document to exclude in the search result.
	ExcludeFields []string

	// list of fields that should be highlighted with snippetting.
	// You can use this parameter to highlight fields that you don't query for, as well.
	// Default: all queried fields will be highlighted.
	HighlightFields []string

	// list of fields which should be highlighted fully without snippeting.
	// Default: all fields will be snippeted.
	HighlightFullFields []string

	// The number of tokens that should surround the highlighted text on each side.
	// Default: 4
	HighlightAffixNumTokens int

	// The start and end tag used for the highlighted snippets.
	HighlightStartTag *string
	HighlightEndTag   *string

	// Field values under this length will be fully highlighted, instead of showing a snippet of relevant portion.
	SnippetThreshold int

	// Maximum number of typographical errors (0, 1 or 2) that would be tolerated.
	NumTypos int

	// If at least typo_tokens_threshold number of results are not found for a specific query,
	// Typesense will attempt to look for results with more typos until num_typos is reached
	// or enough results are found. Set to 0 to disable typo tolerance.
	// Default: 100
	TypoTokensThreshold *int

	// If at least drop_tokens_threshold number of results are not found for a specific query,
	// Typesense will attempt to drop tokens (words) in the query until enough results are found.
	// Tokens that have the least individual hits are dropped first.
	// Set to 0 to disable dropping of tokens.
	DropTokensThreshold *int

	// A list of records to unconditionally include in the search results at specific positions.
	// An example use case would be to feature or promote certain items on the top of search results.
	//
	// A comma separated list of record_id:hit_position.
	// Eg: to include a record with ID 123 at Position 1
	// and another record with ID 456 at Position 5, you'd specify 123:1,456:5.
	PinnedHits []string

	// A list of record_ids to unconditionally hide from search results.
	HiddenHits []string

	// If you have some overrides defined but want to disable all of them for a particular search query, set this to true.
	DisableOverrides bool

	// Set this parameter to true if you wish to split the search query into space separated words yourself.
	// When set to true, we will only split the search query by space, instead of using the locale-aware, built-in tokenizer.
	NoPreSegmentedQuery bool

	// Maximum number of hits that can be fetched from the collection. Eg: 200
	// page * per_page should be less than this number for the search request to return results.
	LimitHits int
}

// Search searches the collection.
func (c *Client) Search(collection string, data SearchData) (res SearchResult, err error) {
	v := url.Values{
		"q":                      {data.Query},
		"query_by":               {strings.Join(data.QueryBy, ",")},
		"prefix":                 {strconv.FormatBool(!data.NoPrefix)},
		"prioritize_exact_match": {strconv.FormatBool(!data.NoPrioritizeExactMatch)},
		"enable_overrides":       {strconv.FormatBool(!data.DisableOverrides)},
		"pre_segmented_query":    {strconv.FormatBool(!data.NoPreSegmentedQuery)},
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
		v["max_facet_values"] = []string{strconv.Itoa(data.MaxFacetValues)}
	}

	if data.FacetQuery != "" {
		v["facet_query"] = []string{data.FacetQuery}
	}

	if data.Page != 0 {
		v["page"] = []string{strconv.Itoa(data.Page)}
	}

	if data.PerPage != 0 {
		v["per_page"] = []string{strconv.Itoa(data.PerPage)}
	}

	if len(data.GroupBy) > 0 {
		v["group_by"] = []string{strings.Join(data.GroupBy, ",")}
	}

	if data.GroupLimit != 0 {
		v["group_limit"] = []string{strconv.Itoa(data.GroupLimit)}
	}

	if len(data.HighlightFields) > 0 {
		v["highlight_fields"] = []string{strings.Join(data.HighlightFields, ",")}
	}

	if len(data.HighlightFullFields) > 0 {
		v["highlight_full_fields"] = []string{strings.Join(data.HighlightFullFields, ",")}
	}

	if data.HighlightAffixNumTokens != 0 {
		v["highlight_affix_num_tokens"] = []string{strconv.Itoa(data.HighlightAffixNumTokens)}
	}

	if data.HighlightStartTag != nil {
		v["highlight_start_tag"] = []string{*data.HighlightStartTag}
	}

	if data.HighlightEndTag != nil {
		v["highlight_end_tag"] = []string{*data.HighlightEndTag}
	}

	if data.SnippetThreshold != 0 {
		v["snippet_threshold"] = []string{strconv.Itoa(data.SnippetThreshold)}
	}

	if data.TypoTokensThreshold != nil {
		v["typo_tokens_threshold"] = []string{strconv.Itoa(*data.TypoTokensThreshold)}
	}

	if data.DropTokensThreshold != nil {
		v["drop_tokens_threshold"] = []string{strconv.Itoa(*data.DropTokensThreshold)}
	}

	if len(data.PinnedHits) > 0 {
		v["pinned_hits"] = []string{strings.Join(data.PinnedHits, ",")}
	}

	if data.LimitHits > 0 {
		v["limit_hits"] = []string{strconv.Itoa(data.LimitHits)}
	}

	resp, err := c.Request("GET", "/collections/"+collection+"/documents/search", WithURLValues(v))
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &res)
	return
}

// SearchResult is the result returned from a search.
type SearchResult struct {
	FacetCounts []int `json:"facet_counts"`

	// Number of found documents
	Found int `json:"found"`
	OutOf int `json:"out_of"`
	Page  int `json:"page"`

	// Search time in milliseconds
	SearchTime int `json:"search_time_ms"`

	Hits []SearchHit `json:"hits"`
}

// SearchHit is a single hit in SearchResult.
// Document is raw JSON data, call UnmarshalTo to unmarshal it to a struct, or Map to unmarshal it to a map[string]interface{}.
type SearchHit struct {
	Document jsonutil.Raw `json:"document"`

	Highlights []Highlight `json:"highlights"`

	TextMatch int `json:"text_match"`
}

// Highlight is a highlight in SearchResult.
type Highlight struct {
	// The matched field name
	Field string `json:"field"`

	Indices []int `json:"indices"`

	MatchedTokens jsonutil.Raw `json:"matched_tokens"`

	// Only present for non-array string fields
	Snippet string `json:"snippet,omitempty"`

	// Only present for string array fields
	Snippets []string `json:"snippets,omitempty"`
}

// Map returns the Document as a map[string]interface{}.
func (s SearchHit) Map() (map[string]interface{}, error) {
	m := map[string]interface{}{}

	err := s.Document.UnmarshalTo(&m)
	return m, err
}

// UnmarshalTo unmarshals m into v.
func (s SearchHit) UnmarshalTo(v interface{}) error {
	return s.Document.UnmarshalTo(v)
}
