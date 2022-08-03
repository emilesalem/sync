package syncmap_test

import (
	"context"
	"sync"
	"testing"

	"github.com/emilesalem/sync/pkg/syncmap"
)

func TestEmptyMap(t *testing.T) {
	ctx := context.Background()

	s := syncmap.NewSyncmap[int, int](ctx, nil)

	f := s.Flush()

	if len(f) != 0 {
		t.Errorf("len(s) = %d; want 0", len(f))
	}
}

func TestGet(t *testing.T) {
	ctx := context.Background()

	s := syncmap.NewSyncmap(ctx, map[string]string{
		`key1`: `value1`,
		`key2`: `value2`,
		`key3`: `value3`,
		`key4`: `value4`,
	})

	v := s.Get(`key1`)

	if v != `value1` {
		t.Errorf("Get(key1) = %s; want `value1`", v)
	}

	v = s.Get(`key4`)

	if v != `value4` {
		t.Errorf("Get(key1) = %s; want `value4`", v)
	}

	v = s.Get(`key5`)

	if v != `` {
		t.Errorf("Get(key1) = %s; want ``", v)
	}
}

func TestSet(t *testing.T) {
	ctx := context.Background()

	s := syncmap.NewSyncmap[string, []string](ctx, nil)

	s.Set(`key1`, []string{`value1`})

	m := s.Flush()

	if m[`key1`][0] != `value1` {
		t.Errorf("Set(`key1`, `value1`) = %s; want `value1`", m[`key1`])
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()

	s := syncmap.NewSyncmap(ctx, map[int]string{
		1: `value1`,
		2: `value2`,
		3: `value3`,
		4: `value4`,
	})

	m := s.Flush()

	if m[1] != `value1` {
		t.Errorf("precondition: m[`key1]`) = %s; want `value1`", m[1])
	}
	s.Delete(1)

	m = s.Flush()

	if m[1] != `` {
		t.Errorf("Delete(`key1`): m[`key1`] = %s; want ``", m[1])
	}
}

const concurrency = 100000

func TestGetConcurrent(t *testing.T) {
	ctx := context.Background()

	s := syncmap.NewSyncmap(ctx, map[string]string{
		`key1`: `value1`,
		`key2`: `value2`,
		`key3`: `value3`,
		`key4`: `value4`,
	})

	var w sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		w.Add(1)
		go func() {
			v := s.Get(`key1`)

			if v != `value1` {
				t.Errorf("Get(key1) = %s; want `value1`", v)
			}

			v = s.Get(`key4`)

			if v != `value4` {
				t.Errorf("Get(key1) = %s; want `value4`", v)
			}

			v = s.Get(`key5`)

			if v != `` {
				t.Errorf("Get(key1) = %s; want ``", v)
			}
			w.Done()
		}()
	}
	w.Wait()
}

func TestSetConcurrent(t *testing.T) {
	ctx := context.Background()

	s := syncmap.NewSyncmap[string, string](ctx, nil)

	var w sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		w.Add(1)
		go func() {

			s.Set(`key1`, `value1`)

			m := s.Flush()

			if m[`key1`] != `value1` {
				t.Errorf("Set(`key1`, `value1`) = %s; want `value1`", m[`key1`])
			}
			w.Done()
		}()
	}
	w.Wait()
}

func TestDeleteConcurrent(t *testing.T) {
	ctx := context.Background()

	s := syncmap.NewSyncmap(ctx, map[string]string{
		`key1`: `value1`,
		`key2`: `value2`,
		`key3`: `value3`,
		`key4`: `value4`,
	})

	m := s.Flush()

	if m[`key1`] != `value1` {
		t.Errorf("precondition: m[`key1]`) = %s; want `value1`", m[`key1`])
	}

	var w sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		w.Add(1)
		go func() {

			s.Delete(`key1`)

			f := s.Flush()

			if f[`key1`] != `` {
				t.Errorf("Delete(`key1`): m[`key1`] = %s; want ``", f[`key1`])
			}

			w.Done()
		}()
	}
	w.Wait()
}
