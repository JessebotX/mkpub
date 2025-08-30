package config

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrDateFromMapKeyNotFound = errors.New("key not found in map")
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

func dateFromString(input string) (time.Time, error) {
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

func dateFromMap(key string, m map[string]any) (time.Time, error) {
	input, ok := m[key]
	if !ok {
		return time.Time{}, ErrDateFromMapKeyNotFound
	}

	switch v := input.(type) {
	case time.Time:
		return v, nil
	case string:
		t, err := dateFromString(v)
		if err != nil {
			return time.Time{}, err
		}

		return t, nil
	}

	return time.Time{}, fmt.Errorf("unrecognized type when parsing date from map: want value of type 'time.Time' or 'string'")
}
