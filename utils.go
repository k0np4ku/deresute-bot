package main

import (
	"fmt"
	"image"
	"sync"

	. "github.com/kbinani/screenshot"
	. "gocv.io/x/gocv"
)

const THRESHOLD = 0.95

var mCache = make(map[string]Mat)

func captureScreen() ([]Mat, error) {
	var result []Mat

	n := NumActiveDisplays()
	for i := 0; i < n; i++ {
		bounds := GetDisplayBounds(i)

		img, err := CaptureRect(bounds)
		if err != nil {
			return nil, err
		}

		mat, err := ImageToMatRGB(img)
		if err != nil {
			return nil, err
		}
		result = append(result, mat)
	}

	return result, nil
}

func searchImage(name string) *image.Point {
	fmt.Println("SearchImage", name)
	screenshots, err := captureScreen()
	if err != nil {
		panic(err)
	}

	img, exists := mCache[name]
	if !exists {
		img = IMRead(fmt.Sprintf("images/%s.png", name), IMReadAnyColor)
		mCache[name] = img
	}

	for _, screenshot := range screenshots {
		defer screenshot.Close()

		result := NewMat()
		defer result.Close()

		mask := NewMat()
		defer mask.Close()

		MatchTemplate(screenshot, img, &result, TmCcoeffNormed, mask)
		Threshold(result, &result, THRESHOLD, 1.0, ThresholdToZero)

		_, maxVal, _, maxLoc := MinMaxLoc(result)
		if maxVal >= THRESHOLD {
			return &maxLoc
		}
	}
	return nil
}

func waitImage(search string) *image.Point {
	var result *image.Point
	for {
		if result = searchImage(search); result != nil {
			break
		}
	}
	return result
}

func iterateStrings(items []string, f func(v string, i int)) {
	wg := sync.WaitGroup{}
	wg.Add(len(items))
	for i, item := range items {
		go func(v string, i int) {
			defer wg.Done()
			f(v, i)
		}(item, i)
	}
	wg.Wait()
}

func iterateBooleans(items []bool, f func(v bool, i int)) {
	wg := sync.WaitGroup{}
	wg.Add(len(items))
	for i, item := range items {
		go func(v bool, i int) {
			defer wg.Done()
			f(v, i)
		}(item, i)
	}
	wg.Wait()
}
