package models

type Price struct {
	Base  string
	Quote string
	Rate  float64
}

type BQ struct {
	Base  []string
	Quote []string
}
