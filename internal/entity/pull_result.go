// Copyright (c) 2024 Wibowo Arindrarto <contact@arindrarto.dev>
// SPDX-License-Identifier: BSD-3-Clause

package entity

// PullResult is a container for a pull operation.
type PullResult struct {
	status PullStatus
	url    *string
	feed   *Feed
	err    error
}

func NewPullResultFromFeed(url *string, feed *Feed) PullResult {
	return PullResult{status: PullSuccess, url: url, feed: feed}
}

func NewPullResultFromError(url *string, err error) PullResult {
	return PullResult{status: PullFail, url: url, err: err}
}

func (msg PullResult) Feed() *Feed {
	if msg.status == PullSuccess {
		return msg.feed
	}
	return nil
}

func (msg PullResult) Error() error {
	if msg.status == PullFail {
		return msg.err
	}
	return nil
}

func (msg PullResult) URL() string {
	if msg.url != nil {
		return *msg.url
	}
	return ""
}

func (msg *PullResult) SetError(err error) {
	msg.err = err
}

func (msg *PullResult) SetStatus(status PullStatus) {
	msg.status = status
}

type PullStatus int

const (
	PullSuccess PullStatus = iota
	PullFail
)
