package pub

func allChapters(chapters *[]Chapter) []*Chapter {
	var flattened []*Chapter

	for i := range *chapters {
		next := &((*chapters)[i])

		var nested []*Chapter
		if len(next.Chapters) > 0 {
			nested = next.Subchapters()
		}

		flattened = append(flattened, next)
		flattened = append(flattened, nested...)
	}

	return flattened
}
