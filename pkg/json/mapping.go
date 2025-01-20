// pkg/json/mapping.go

package json

import (
	"bytes"
	"encoding/json"
	"errors"
	"sync"
	"unsafe"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func Mapping[T any](body []byte) ([]T, error) {
	if len(body) == 0 {
		return nil, errors.New("Mapping error: empty JSON body")
	}

	if !json.Valid(body) {
		return nil, errors.New("Mapping error: invalid JSON format")
	}

	jsonStr := *(*string)(unsafe.Pointer(&body))

	var rawList []json.RawMessage
	err := json.Unmarshal(*(*[]byte)(unsafe.Pointer(&jsonStr)), &rawList)
	if err != nil {
		return nil, errors.New("MappingParallel error: JSON list parsing failed -> " + err.Error())
	}

	numItems := len(rawList)
	results := make([]T, numItems)
	errChan := make(chan error, numItems)
	wg := sync.WaitGroup{}

	for i, raw := range rawList {
		wg.Add(1)

		go func(idx int, data json.RawMessage) {
			defer wg.Done()

			jsonSegment := *(*string)(unsafe.Pointer(&data))

			buf := bufferPool.Get().(*bytes.Buffer)
			buf.Reset()
			buf.WriteString(jsonSegment)

			decoder := json.NewDecoder(buf)
			var temp T
			if err := decoder.Decode(&temp); err != nil {
				errChan <- errors.New("MappingParallel error: JSON decoding failed -> " + err.Error())
				bufferPool.Put(buf)
				return
			}

			results[idx] = temp

			bufferPool.Put(buf)
		}(i, raw)
	}

	wg.Wait()
	close(errChan)

	for e := range errChan {
		return nil, e
	}

	return results, nil
}
