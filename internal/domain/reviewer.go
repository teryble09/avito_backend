package domain

import (
	"math/rand/v2"
	"slices"
)

func SelectRandomReviewersExcludingAuthor(candidates []*User, maxCount int, author *User) (reviewers []*User) {
	// фильтруем автора из кандидатов
	eligible := slices.DeleteFunc(slices.Clone(candidates), func(candidate *User) bool {
		return candidate.ID == author.ID
	})

	count := min(len(eligible), maxCount)

	rand.Shuffle(len(eligible), func(i, j int) {
		eligible[i], eligible[j] = eligible[j], eligible[i]
	})

	return eligible[0:count]
}

func SelectReplacementReviewer(candidates []*User, authorID, reviewerID string) (string, error) {
	// фильтруем кандидатов
	eligible := slices.DeleteFunc(slices.Clone(candidates), func(candidate *User) bool {
		return candidate.ID == reviewerID || candidate.ID == authorID
	})

	if len(eligible) == 0 {
		return "", ErrReviewerNoCandidate
	}

	randomIndex := rand.IntN(len(eligible)) //nolint:gosec //уязвимый генератор?

	return eligible[randomIndex].ID, nil
}
