package storage

import (
	"context"
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
	maxWorkersCount int
	wg              sync.WaitGroup
	commonSize      int64
	commonCount     int64
	err             error
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	sizerCustom := sizer{}
	sizerCustom.maxWorkersCount = 10
	return &sizerCustom
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	dirs, files, err := d.Ls(ctx)
	if err != nil {
		a.err = err
	}

	for _, file := range files {
		sizeFile, errReadStat := file.Stat(ctx)
		if errReadStat != nil {
			a.err = errReadStat
		}
		a.commonSize += sizeFile
		a.commonCount++
	}

	for _, dir := range dirs {
		if a.maxWorkersCount > 0 {
			a.maxWorkersCount--
			a.wg.Add(1)
			dir := dir
			go func() {
				defer a.wg.Done()
				defer func() {
					a.maxWorkersCount++
				}()
				a.SizeUsingGorutine(ctx, dir)
			}()
		} else {
			_, _ = a.Size(ctx, dir)
		}
	}

	a.wg.Wait()

	if a.err != nil {
		return Result{0, 0}, a.err
	}
	return Result{a.commonSize, a.commonCount}, nil
}

func (a *sizer) SizeUsingGorutine(ctx context.Context, d Dir) {
	dirs, files, err := d.Ls(ctx)
	if err != nil {
		a.err = err
	}

	for _, file := range files {
		sizeFile, errReadStat := file.Stat(ctx)
		if errReadStat != nil {
			a.err = errReadStat
		}
		a.commonSize += sizeFile
		a.commonCount++
	}

	for _, dir := range dirs {
		if a.maxWorkersCount > 0 {
			a.maxWorkersCount--
			a.wg.Add(1)
			dir := dir
			go func() {
				defer a.wg.Done()
				defer func() {
					a.maxWorkersCount++
				}()
				a.SizeUsingGorutine(ctx, dir)
			}()
		} else {
			_, _ = a.Size(ctx, dir)
		}
	}
}
