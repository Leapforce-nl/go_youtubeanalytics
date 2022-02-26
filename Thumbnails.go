package youtube

type Thumbnails struct {
	Default  *Thumbnail `json:"default"`
	Medium   *Thumbnail `json:"medium"`
	High     *Thumbnail `json:"high"`
	Standard *Thumbnail `json:"standard"`
	MaxRes   *Thumbnail `json:"maxres"`
}

type Thumbnail struct {
	Url    string `json:"url"`
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
}
