package tree

import (
	"math/rand"
	"reflect"
	"testing"
)

func Test_results_Top(t *testing.T) {
	input := results{
		result{
			Score: 5,
		},
		result{
			Score: 2,
		},
		result{
			Score: 7,
		},
		result{
			Score: 6,
		},
		result{
			Score: 9,
		},
		result{
			Score: 4,
		},
		result{
			Score: 1,
		},
		result{
			Score: 8,
		},
		result{
			Score: 10,
		},
	}
	rand.Shuffle(len(input), func(i, j int) {
		input[i], input[j] = input[j], input[i]
	})

	tests := []struct {
		name  string
		r     results
		limit uint
		want  results
	}{
		{
			"when given an unsorted list it should return the top results",
			input,
			3,
			results{
				result{
					Score: 1,
				},
				result{
					Score: 2,
				},
				result{
					Score: 4,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Top(tt.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("results.Top() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkTop(b *testing.B) {
	list := make(results, 1e4)
	for i := 0; i < 1e4; i++ {
		list[i] = result{
			Score: rand.Float64(),
		}
	}
	for i := 0; i < b.N; i++ {
		_ = list.Top(10)
	}
}
