package model

// Entry represents one entity, storing by password manager.
type Entry struct {
	URL      string
	Username string
	Password string
	Extra    string
	Name     string
	Grouping string
}
