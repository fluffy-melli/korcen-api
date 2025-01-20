// pkg/json/mapping.go

package json

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"sync"
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

	var rawList []json.RawMessage
	if err := json.Unmarshal(body, &rawList); err != nil {
		return nil, errors.New("Mapping error: JSON list parsing failed -> " + err.Error())
	}

	numItems := len(rawList)
	results := make([]T, numItems)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}
	errChan := make(chan error, numItems)

	for i, raw := range rawList {
		wg.Add(1)

		go func(idx int, data json.RawMessage) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
			}

			buf := bufferPool.Get().(*bytes.Buffer)
			buf.Reset()
			defer bufferPool.Put(buf)

			buf.Write(data)

			var temp T
			decoder := json.NewDecoder(buf)
			if err := decoder.Decode(&temp); err != nil {
				errChan <- errors.New("Mapping error: JSON decoding failed -> " + err.Error())
				cancel()
				return
			}

			results[idx] = temp
		}(i, raw)
	}

	wg.Wait()
	close(errChan)

	for e := range errChan {
		return nil, e
	}

	return results, nil
}
