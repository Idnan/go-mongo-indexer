package main

type ConfigCollection struct {
	Collection string     `json:"collection"`
	CapSize    int        `json:"cap"`
	Indexes    [][]string `json:"index"`
}

type IndexDiff struct {
	Old map[string]map[string][]string
	New map[string]map[string][]string
	Cap map[string]int
}
