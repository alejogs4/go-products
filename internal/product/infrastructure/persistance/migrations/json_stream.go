package migrations

import (
	"bytes"
	"encoding/json"
)

type JsonStream[T any] struct {
	Item  T
	Error error
}

func ReadJson[T any](content []byte) <-chan JsonStream[T] {
	results := make(chan JsonStream[T])
	decoder := json.NewDecoder(bytes.NewReader(content))

	go func() {
		if _, err := decoder.Token(); err != nil {
			results <- JsonStream[T]{Error: err}
			close(results)
			return
		}

		for decoder.More() {
			var item T
			if err := decoder.Decode(&item); err != nil {
				results <- JsonStream[T]{Error: err}
				close(results)
				return
			}

			results <- JsonStream[T]{Item: item}
		}

		close(results)
	}()

	return results
}
