package resilience

import (
	"bytes"
	"errors"
	"sync"
	"testing"
)

func TestDoReplaysSameBodyOnce(t *testing.T) {
	s := NewIdempotency()
	body := []byte(`{"courseId":1,"userId":7}`)
	var calls int
	fn := func() (int, []byte) {
		calls++
		return 201, []byte(`{"id":1}`)
	}

	status, resp, err := s.Do(7, "abc", body, fn)
	if err != nil || status != 201 || string(resp) != `{"id":1}` {
		t.Fatalf("first status=%d body=%s err=%v", status, resp, err)
	}
	status, resp, err = s.Do(7, "abc", body, fn)
	if err != nil || status != 201 || string(resp) != `{"id":1}` || calls != 1 {
		t.Fatalf("replay status=%d body=%s calls=%d err=%v", status, resp, calls, err)
	}
}

func TestDoRejectsDifferentBody(t *testing.T) {
	s := NewIdempotency()
	if _, _, err := s.Do(7, "abc", []byte(`{"userId":7}`), func() (int, []byte) {
		return 201, []byte(`{"id":1}`)
	}); err != nil {
		t.Fatal(err)
	}
	_, _, err := s.Do(7, "abc", []byte(`{"userId":8}`), func() (int, []byte) {
		t.Fatal("fn must not run")
		return 0, nil
	})
	if !errors.Is(err, ErrKeyReused) {
		t.Fatalf("err=%v", err)
	}
}

func TestDoScopesKeyPerUser(t *testing.T) {
	s := NewIdempotency()
	body := []byte(`{"courseId":1}`)
	var calls int
	fn := func() (int, []byte) {
		calls++
		return 201, []byte(`{"id":1}`)
	}
	if _, _, err := s.Do(7, "abc", body, fn); err != nil {
		t.Fatal(err)
	}
	if _, _, err := s.Do(8, "abc", body, fn); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatalf("calls=%d", calls)
	}
}

func TestDoDoesNotHoldLockAcrossFn(t *testing.T) {
	s := NewIdempotency()
	started := make(chan struct{})
	release := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _, _ = s.Do(7, "abc", []byte("a"), func() (int, []byte) {
			close(started)
			<-release
			return 201, []byte("ok")
		})
	}()
	<-started
	_, _, err := s.Do(7, "abc", []byte("a"), func() (int, []byte) {
		t.Error("second call ran fn while the first was in flight")
		return 0, nil
	})
	if !errors.Is(err, ErrInProgress) {
		t.Fatalf("err=%v", err)
	}
	close(release)
	wg.Wait()
}

func TestDoDropsServerErrors(t *testing.T) {
	s := NewIdempotency()
	var calls int
	fn := func() (int, []byte) {
		calls++
		if calls == 1 {
			return 500, []byte("boom")
		}
		return 201, []byte("ok")
	}
	status, _, err := s.Do(7, "abc", []byte("a"), fn)
	if err != nil || status != 500 {
		t.Fatalf("status=%d err=%v", status, err)
	}
	status, body, err := s.Do(7, "abc", []byte("a"), fn)
	if err != nil || status != 201 || !bytes.Equal(body, []byte("ok")) || calls != 2 {
		t.Fatalf("status=%d body=%s calls=%d err=%v", status, body, calls, err)
	}
}
