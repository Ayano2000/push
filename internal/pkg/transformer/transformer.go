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
		return "", errors.WithStack(errors.New("no filter specified"))
	}

	query, err := gojq.Parse(filter)
	if err != nil {
		return "", err
	}

	var results []interface{}
	iter := query.RunWithContext(ctx, payload)
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

func ValidFilter(filter string) error {
	if filter == "" {
		return errors.WithStack(errors.New("no filter specified"))
	}
	_, err := gojq.Parse(filter)
	return err
}
