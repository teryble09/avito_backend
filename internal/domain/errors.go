package domain

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")

	ErrTeamAlreadyExist = errors.New("team already exists")
	ErrTeamNotFound     = errors.New("team not found")

	ErrPrNotFound      = errors.New("not found")
	ErrPrAlreadyExists = errors.New("already exists")
	ErrPrMerged        = errors.New("PR is already merged")

	ErrReviewerNoCandidate = errors.New("no candidate")
	ErrReviewerNotAssigned = errors.New("reviewer is not assigned to this PR")
)
