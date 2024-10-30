package transformer

import (
	"context"
	"encoding/json"
	"github.com/itchyny/gojq"
	"github.com/pkg/errors"
)

// Transform will take a json payload, and a JQ filter,
func Transform(ctx context.Context, payload string, filter string) (string, error) {
	if filter == "" {
		return payload, nil
	}

	query, err := gojq.Parse(filter)
	if err != nil {
		return "", err
	}
	var object any
	err = json.Unmarshal([]byte(payload), &object)
	if err != nil {
		return "", err
	}

	var results []interface{}
	iter := query.RunWithContext(ctx, object)
	for {
		value, hasNextValue := iter.Next()
		if !hasNextValue {
			break
		}

		if err, ok := value.(error); ok {
			var haltError *gojq.HaltError
			if errors.As(err, &haltError) && haltError.Value() == nil {
				break
			}
			return "", err
		}

		results = append(results, value)
	}

	if len(results) == 0 {
		return "", errors.WithStack(errors.New("no results found"))
	}

	var marshal interface{}
	switch len(results) {
	case 0:
		marshal = struct{}{}
	case 1:
		marshal = results[0]
	default:
		marshal = results
	}

	result, err := json.Marshal(marshal)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// IsValidFilter will return true if the filter is empty
func IsValidFilter(filter string) (bool, error) {
	if filter == "" {
		return true, nil
	}
	_, err := gojq.Parse(filter)
	if err != nil {
		return false, err
	}

	return true, nil
}
