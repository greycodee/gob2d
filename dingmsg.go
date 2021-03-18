package main

type ActionCard struct {
	Title string `json:"title"`
	Text string `json:"text"`
	HideAvatar string `json:"hideAvatar"`
	BtnOrientation string `json:"btnOrientation"`
	SingleTitle string `json:"singleTitle"`
	SingleURL string `json:"singleUrl"`
	Btns []Btn `json:"btns"`
}

type Btn struct {
	Title string `json:"title"`
	ActionURL string `json:"actionUrl"`
}
