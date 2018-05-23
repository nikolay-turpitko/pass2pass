package model

// Entry represents one entity, storing by password manager.
type Entry struct {
	URL      string
	Username string
	Password string
	Extra    interface{}
	Name     string
	Grouping string
}
