package main

type ConfigCollection struct {
	Collection string           `json:"collection"`
	CapSize    int              `json:"cap"`
	Indexes    []map[string]int `json:"index"`
}

type IndexDiff struct {
	Old map[string]map[string]IndexModel
	New map[string]map[string]IndexModel
	Cap map[string]int
}

type IndexModel struct {
	Name string
	Keys map[string]int
}
