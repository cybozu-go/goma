package average

import (
	"fmt"

	"github.com/cybozu-go/goma/filters"
)

const (
	defaultWindowSize = 10
)

type filter struct {
	init   float64
	values []float64
	index  int
}

func (f *filter) Init() {
	for i := 0; i < len(f.values); i++ {
		f.values[i] = f.init
	}
}

func (f *filter) Put(v float64) (avg float64) {
	f.values[f.index] = v
	f.index++
	if f.index == len(f.values) {
		f.index = 0
	}

	for _, t := range f.values {
		avg += t
	}
	avg /= float64(len(f.values))
	return
}

func (f *filter) String() string {
	return fmt.Sprintf("filter:average(window=%d, init=%f)",
		len(f.values), f.init)
}

func construct(params map[string]interface{}) (filters.Filter, error) {
	var init float64
	if v, ok := params["init"]; ok {
		if f, ok := v.(float64); !ok {
			return nil, fmt.Errorf("init is not a float: %v", v)
		} else {
			init = f
		}
	}

	var window int = defaultWindowSize
	if v, ok := params["window"]; ok {
		if i, ok := v.(int); !ok {
			return nil, fmt.Errorf("window is not an integer: %v", v)
		} else {
			if i < 1 {
				return nil, fmt.Errorf("too small window size: %d", i)
			}
			window = i
		}
	}

	return &filter{
		init:   init,
		values: make([]float64, window),
	}, nil
}

func init() {
	filters.Register("average", construct)
}
