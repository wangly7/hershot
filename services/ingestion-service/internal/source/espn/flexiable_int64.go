package espn

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type FlexiableInt64 int64

func (f *FlexiableInt64) UnmarshalJSON(data []byte) error {
	var number int64
	if err := json.Unmarshal(data, &number); err == nil {
		*f = FlexiableInt64(number)
		return nil
	}

	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		return fmt.Errorf("FlexibleInt64: expected number or string: %w", err)
	}

	if text == "" {
		*f = 0
		return nil
	}

	value, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return fmt.Errorf("FlexibleInt64: invalid integer %q: %w", text, err)
	}

	*f = FlexiableInt64(value)
	return nil
}
