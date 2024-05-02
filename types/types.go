package types

type User struct {
	ID   int
	Name string
}

type Fraction struct {
	Numerator   int
	Denominator int
}

type Ingredient struct {
	name     string
	quantity Fraction
	unit     string
}

type Recipe struct {
	ID          int
	name        string
	servings    int
	cookTime    int
	prepTime    int
	imgUrl      string
	sourceUrl   string
	ingredients []Ingredient
	directions  []string
}