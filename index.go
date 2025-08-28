package mkpub

import "github.com/JessebotX/mkpub/book"

type MainIndex struct {
	Title            string
	ShortDescription string
	Content          book.Content
	FaviconImageName string
	Books            []BookIndex

	Profiles []ProfileIndex
	Series   []SeriesIndex
	Tags     []TagIndex
}

type BookIndex struct {
	book.Book

	Index *MainIndex
}

type SeriesIndex struct {
	UniqueID string
	Name     string

	Index *MainIndex
	Books []*BookIndex
}

type TagIndex struct {
	UniqueID string
	Tag      string

	Index *MainIndex
	Books []*BookIndex
}

type ProfileIndex struct {
	book.Profile

	Index            *MainIndex
	PublishedBooks   []*book.Book
	AuthoredBooks    []*book.Book
	ContributedBooks []*book.Book
}
