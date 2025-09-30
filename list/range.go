package list

import (
	"fmt"
	"strconv"
	"strings"
)

func ConvertRangeToSliceInt(rangeAsString []string) ([]int, error) {
	var result []int
	for _, spec := range rangeAsString {
		// handle comma separated numbers first
		parts := strings.Split(spec, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}

			// Case 1: -b (means from 1 to b)
			if strings.HasPrefix(p, "-") {
				end, err := strconv.Atoi(p[1:])
				if err != nil {
					return nil, fmt.Errorf("invalid range %q", p)
				}
				for i := 1; i <= end; i++ {
					result = append(result, i)
				}
				continue
			}

			// Case 2: a-b
			if strings.Contains(p, "-") {
				rangeParts := strings.SplitN(p, "-", 2)
				start, err1 := strconv.Atoi(rangeParts[0])
				end, err2 := strconv.Atoi(rangeParts[1])
				if err1 != nil || err2 != nil || start > end {
					return nil, fmt.Errorf("invalid range %q", p)
				}
				for i := start; i <= end; i++ {
					result = append(result, i)
				}
				continue
			}

			// Case 3: single number
			n, err := strconv.Atoi(p)
			if err != nil {
				return nil, fmt.Errorf("invalid number %q", p)
			}
			result = append(result, n)
		}
	}
	return result, nil
}
