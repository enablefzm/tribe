package jsondb

type FieldImage struct {
	Img   string
	Thing []string
}

func NewFieldImage(img string, thing []string) *FieldImage {
	return &FieldImage{
		Img:   img,
		Thing: thing,
	}
}
