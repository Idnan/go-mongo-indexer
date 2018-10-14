package main

type ConfigCollection struct {
	Collection string     `json:"collection"`
	CapSize    string     `json:"cap"`
	Indexes    [][]string `json:"index"`
}
