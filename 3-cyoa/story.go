package cyoa

// define Story type which will store our story data
type Story map[string]Chapter

// define chapter struct
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

// define options struct
type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}
