package serverpool

import "testing"

func TestParseStrategy(t *testing.T) {
	tests := []struct {
		name     string
		strategy string
		want     LBStrategy
	}{
		{
			name:     "round-robin strategy",
			strategy: "round-robin",
			want:     RoundRobin,
		},
		{
			name:     "least-connected strategy",
			strategy: "least-connected",
			want:     LeastConnected,
		},
		{
			name:     "unknown strategy defaults to round-robin",
			strategy: "unknown",
			want:     RoundRobin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseStrategy(tt.strategy); got != tt.want {
				t.Errorf("ParseStrategy() = %v, want %v", got, tt.want)
			}
		})
	}
}
