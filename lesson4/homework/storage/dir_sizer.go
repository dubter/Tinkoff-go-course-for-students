package storage

import (
	"context"
	"golang.org/x/sync/semaphore"
	"sync"
)

// Result represents the Size function result
type Result struct {
	// Total Size of File objects
	Size int64
	// Count is a count of File objects processed
	Count int64
}

type DirSizer interface {
	// Size calculate a size of given Dir, receive a ctx and the root Dir instance
	// will return Result or error if happened
	Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
	// maxWorkersCount number of workers for asynchronous run
	maxWorkersCount int64
	wg              sync.WaitGroup
	commonSize      int64
	commonCount     int64
	err             error
	semaphore       semaphore.Weighted
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	sizerCustom := sizer{}
	sizerCustom.maxWorkersCount = 4
	sizerCustom.semaphore = *semaphore.NewWeighted(sizerCustom.maxWorkersCount)
	return &sizerCustom
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	a.Workers(ctx, d)
	a.wg.Wait()
	if a.err != nil {
		return Result{0, 0}, a.err
	}
	return Result{a.commonSize, a.commonCount}, nil
}

func (a *sizer) Workers(ctx context.Context, d Dir) {
	dirs, files, err := d.Ls(ctx)
	if err != nil {
		a.err = err
		return
	}

	for _, file := range files {
		sizeFile, errReadStat := file.Stat(ctx)
		if errReadStat != nil {
			a.err = errReadStat
			return
		}
		a.commonSize += sizeFile
		a.commonCount++
	}

	for _, dir := range dirs {
		errAcquire := a.semaphore.Acquire(ctx, 1)
		if errAcquire != nil {
			a.err = errAcquire
			return
		}
		a.wg.Add(1)
		dir := dir
		go func() {
			defer a.wg.Done()
			defer a.semaphore.Release(1)
			a.Workers(ctx, dir)
		}()
	}
}
