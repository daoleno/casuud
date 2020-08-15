package main

type Group struct {
	Name string `json:"name"`
}

type Card struct {
	ID    string `json:"id"`
	Front string `json:"front"`
	Back  string `json:"back"`
}
