package main

type ConfigCollection struct {
	Collection string     `json:"collection"`
	CapSize    string     `json:"cap"`
	Indexes    [][]string `json:"index"`
}

type IndexDiff struct {
	Old map[string]map[string][]string
	New map[string]map[string][]string
}
