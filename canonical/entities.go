package canonical

type Player struct {
	Id       string  `json:"id,omitempty", bson:"id"`
	Name     string  `json:"name,omitempty"`
	Age      int64   `json:"age,omitempty"`
	Position string  `json:"position,omitempty"`
	Foot     string  `json:"foot,omitempty"`
	Genre    string  `json:"ganre,omitempty"`
	Rating   float64 `json:"rating,omitempty"`
	Country  string  `json:"country,omitempty"`
	Url      string  `json:"url,omitempty"`
	Img      string  `json:"img,omitempty"`
}
