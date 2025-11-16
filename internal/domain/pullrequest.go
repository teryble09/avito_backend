package domain

import (
	"errors"
	"math/rand/v2"
	"slices"
)

type PullRequestStatus string

const (
	StatusOpen   PullRequestStatus = "OPEN"
	StatusMerged PullRequestStatus = "MERGED"
)

type PullRequest struct {
	PullRequestID     string
	PullRequestName   string
	AuthorID          string
	Status            PullRequestStatus
	AssignedReviewers []string
}

var (
	ErrPrNotFound      = errors.New("not found")
	ErrPrAlreadyExists = errors.New("already exists")
	ErrPrNoCandidate   = errors.New("no candidate")
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
