package config

import (
	"time"
)

func chaptersFlattened(chapters *[]Chapter) []*Chapter {
	var flattened []*Chapter

	for i := range *chapters {
		next := &((*chapters)[i])

		var nested []*Chapter
		if len(next.Chapters) > 0 {
			nested = next.ChaptersFlattened()
		}

		flattened = append(flattened, next)
		flattened = append(flattened, nested...)
	}

	return flattened
}

func parseDateTime(input string) (time.Time, error) {
	t, err := time.Parse("2006-01-02 15:04:05Z07:00", input)
	if err == nil {
		return t, nil
	}

	t, err = time.Parse("2006-01-02 15:04:05", input)
	if err == nil {
		return t, nil
	}

	t, err = time.Parse("2006-01-02 15:04", input)
	if err == nil {
		return t, nil
	}

	t, err = time.Parse("2006-01-02", input)
	if err == nil {
		return t, nil
	}

	t, err = time.Parse("2006-01", input)
	if err == nil {
		return t, nil
	}

	t, err = time.Parse("2006", input)
	if err == nil {
		return t, nil
	}

	return time.Time{}, err
}
