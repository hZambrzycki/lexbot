package querymodels

type GlobalSearchResult struct {
	Type     string
	ID       string
	Title    string
	Subtitle string
	Href     string
	Snippet  string
	Score    int
}
