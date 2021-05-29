package models

type Answer struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  Elems  `json:"results"`
}

//easyjson:json
type Elems []Elem

type Elem struct {
	Id              uint64     `json:"id"`
	Date            ManyDates  `json:"dates"`
	PublicationDate int64      `json:"publication_date"`
	Title           string     `json:"title"`
	Place           IdType     `json:"place"`
	BodyText        string     `json:"body_text"`
	Categories      []string   `json:"categories"`
	Images          ImagesType `json:"images"`
	Tags            []string   `json:"tags"`
}

type Dates struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

type IdType struct {
	Id uint64 `json:"id"`
}

//easyjson:json
type ImagesType []Image

//easyjson:json
type ManyDates []Dates

type Image struct {
	Path string `json:"image"`
	S    Source `json:"source"`
}

type Source struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

type Place struct {
	Title   string `json:"title"`
	Address string `json:"address"`
	Map     Coords `json:"coords"`
	Subway  string `json:"subway"`
}

type Coords struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}
