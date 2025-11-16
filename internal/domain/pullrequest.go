package domain

import (
	"time"
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
	MergedAt          *time.Time // nullable
}

type PullRequestShort struct {
	PullRequestID   string
	PullRequestName string
	AuthorID        string
	Status          PullRequestStatus
}
