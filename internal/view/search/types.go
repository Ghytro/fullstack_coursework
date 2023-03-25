package search

type searchRequest struct {
	Query string `form:"query"`
}

type GeneralSearchViewData struct {
	Users []User
	Polls []Poll
}

type User struct {
	ID, Username, Country string
}

type Poll struct {
	ID, Topic string
}
