package play

/*
func init() {
	PlaybackStartTimeout = 5 * time.Millisecond
}

func TestWaitStartCancel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var c Control

	otherSub := mock.NewMockSubscription(ctrl)
	otherSub.EXPECT().Events().AnyTimes().Return(nil)
	otherSub.EXPECT().Cancel().AnyTimes()

	c.hangupSub = otherSub
	c.startedSub = otherSub
	c.finishedSub = otherSub

	doneCh := make(chan struct{})

	go func() {
		f := c.waitStart(ctx)
		if f != nil {
			t.Error("waitStart did not return a nil state")
		}
		close(doneCh)
	}()

	select {
	case <-time.After(2 * time.Millisecond):
		t.Error("waitStart failed to detect context closure")
	case <-doneCh:
	}

	if c.status != Canceled {
		t.Error("waitStart returned the wrong state")
	}
}

func TestWaitStartTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var c Control

	// Prepare mock subscriptions
	otherSub := mock.NewMockSubscription(ctrl)
	otherSub.EXPECT().Events().AnyTimes().Return(nil)
	otherSub.EXPECT().Cancel().AnyTimes()

	c.hangupSub = otherSub
	c.startedSub = otherSub
	c.finishedSub = otherSub

	doneCh := make(chan struct{})

	go func() {
		f := c.waitStart(ctx)
		if f != nil {
			t.Error("waitStart did not return a nil state")
		}
		close(doneCh)
	}()

	select {
	case <-time.After(PlaybackStartTimeout * 2):
		t.Error("waitStart failed to detect timeout")
	case <-doneCh:
	}

	if c.status != Timeout {
		t.Error("waitStart returned the wrong state")
	}

}

func TestWaitStartHangup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var c Control

	// Prepare mock subscriptions
	eventChan := make(chan ari.Event)
	close(eventChan)
	sub := mock.NewMockSubscription(ctrl)
	sub.EXPECT().Events().Return(eventChan)
	sub.EXPECT().Cancel()
	otherSub := mock.NewMockSubscription(ctrl)
	otherSub.EXPECT().Events().AnyTimes().Return(nil)
	otherSub.EXPECT().Cancel().AnyTimes()

	c.hangupSub = sub
	c.startedSub = otherSub
	c.finishedSub = otherSub

	doneCh := make(chan struct{})

	go func() {
		f := c.waitStart(ctx)
		if f != nil {
			t.Error("waitStart did not return a nil state")
		}
		close(doneCh)
	}()

	c.hangupSub.Cancel()

	select {
	case <-time.After(time.Millisecond):
		t.Error("waitStart failed to detect hangup")
	case <-doneCh:
	}

	if c.status != Hangup {
		t.Error("waitStart returned the wrong state")
	}
}

func TestWaitStartFinished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := Control{
		stopCh: make(chan struct{}),
	}

	// Prepare mock subscriptions
	eventChan := make(chan ari.Event)
	close(eventChan)
	sub := mock.NewMockSubscription(ctrl)
	sub.EXPECT().Events().Return(eventChan)
	sub.EXPECT().Cancel().AnyTimes()
	otherSub := mock.NewMockSubscription(ctrl)
	otherSub.EXPECT().Events().AnyTimes().Return(nil)
	otherSub.EXPECT().Cancel().AnyTimes()

	c.hangupSub = otherSub
	c.startedSub = otherSub
	c.finishedSub = sub

	doneCh := make(chan struct{})

	go func() {
		f := c.waitStart(ctx)
		if f != nil {
			t.Error("waitStart did not return a nil state")
		}
		close(doneCh)
	}()

	c.hangupSub.Cancel()

	select {
	case <-time.After(time.Millisecond):
		t.Error("waitStart failed to detect playback started")
	case <-doneCh:
	}

	if c.status != Finished {
		t.Error("waitStart returned the wrong state", c.status)
	}
}
func TestWaitStartStarted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := Control{
		startCh: make(chan struct{}),
	}

	// Prepare mock subscriptions
	eventChan := make(chan ari.Event)
	close(eventChan)
	sub := mock.NewMockSubscription(ctrl)
	sub.EXPECT().Events().Return(eventChan)
	sub.EXPECT().Cancel().AnyTimes()
	otherSub := mock.NewMockSubscription(ctrl)
	otherSub.EXPECT().Events().AnyTimes().Return(nil)
	otherSub.EXPECT().Cancel().AnyTimes()

	c.hangupSub = otherSub
	c.startedSub = sub
	c.finishedSub = otherSub

	doneCh := make(chan struct{})

	go func() {
		f := c.waitStart(ctx)
		if f == nil {
			t.Error("waitStart returned a nil state instead of c.waitStop")
		}
		close(doneCh)
	}()

	c.hangupSub.Cancel()

	select {
	case <-time.After(time.Millisecond):
		t.Error("waitStart failed to detect playback started")
	case <-doneCh:
	}

	if c.status != InProgress {
		t.Error("waitStart returned the wrong state", c.status)
	}
}

*/
