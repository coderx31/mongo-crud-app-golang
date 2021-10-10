package book

type Book struct {
	ID     string `bson:"book_id" json:"book_id"`
	ISBN   string `bson:"isbn" json:"isbn"`
	Title  string `bson:"title" json:"title"`
	Author Author `bson:"author" json:"author"`
}

type Author struct {
	Firstname string `bson:"firstname" json:"firstname"`
	LastName  string `bson:"lastname" json:"lastname"`
}
