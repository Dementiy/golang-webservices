package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for data := range in {
		value, ok := data.(string)
		if !ok {
			value = strconv.Itoa(data.(int))
		}

		wg.Add(1)
		go func(data string) {
			defer wg.Done()
			dataHash1 := make(chan string)
			dataHash2 := make(chan string)

			go func() {
				dataHash1 <- DataSignerCrc32(data)
			}()

			go func() {
				mu.Lock()
				md5hash := DataSignerMd5(data)
				mu.Unlock()
				dataHash2 <- DataSignerCrc32(md5hash)
			}()

			out <- <-dataHash1 + "~" + <-dataHash2
		}(value)
	}

	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for data := range in {
		value, ok := data.(string)
		if !ok {
			value = strconv.Itoa(data.(int))
		}

		wg.Add(1)
		go func(data string) {
			defer wg.Done()
			workers := &sync.WaitGroup{}
			dataHashes := make([]string, 6)
			for th := 0; th < 6; th++ {
				workers.Add(1)
				go func(th int) {
					defer workers.Done()
					dataHashes[th] = DataSignerCrc32(strconv.Itoa(th) + data)
				}(th)
			}
			workers.Wait()
			out <- strings.Join(dataHashes, "")
		}(value)
	}

	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	results := make([]string, 0, 5)
	for data := range in {
		results = append(results, data.(string))
	}
	sort.Strings(results)
	result := strings.Join(results, "_")
	out <- result
}

func runJob(wg *sync.WaitGroup, jobFunc job, in, out chan interface{}) {
	defer wg.Done()
	jobFunc(in, out)
	close(out)
}

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{}, 100)
	out := make(chan interface{}, 100)
	for _, job := range jobs {
		wg.Add(1)
		go runJob(wg, job, in, out)
		in = out
		out = make(chan interface{}, 100)
	}
	wg.Wait()
	close(out)
}
