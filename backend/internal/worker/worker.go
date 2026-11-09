package worker

import (
	"context"
	"fmt"
	"github.com/riverqueue/river"
)

type MediaArgs struct{ PostID, MediaKey, MediaType string }
func (MediaArgs) Kind() string { return "media_verify" }
type MediaWorker struct{ CDN string }
func (w *MediaWorker) Work(ctx context.Context, job *river.Job[MediaArgs]) error {
	_ = job.Args.MediaKey // use CDN + key via fmt.Sprintf if needed
	return nil
}

type ChallengeArgs struct{ Theme string }
func (ChallengeArgs) Kind() string { return "challenge_daily" }
type ChallengeWorker struct{}
func (w *ChallengeWorker) Work(ctx context.Context, job *river.Job[ChallengeArgs]) error { return nil }

type NotificationArgs struct{ UserID, Type, Message string }
func (NotificationArgs) Kind() string { return "notification_push" }
type NotificationWorker struct{}
func (w *NotificationWorker) Work(ctx context.Context, job *river.Job[NotificationArgs]) error { return nil }
