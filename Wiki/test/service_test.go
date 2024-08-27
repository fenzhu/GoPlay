package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/wiki/service"
)

func BenchmarkViewHandler(b *testing.B) {
	cases := []int{100, 1000, 10000, 100000}
	for _, n := range cases {
		b.Run(fmt.Sprintf("Size%d", n), func(b *testing.B) {
			// Create a new request
			req, _ := http.NewRequest("GET", "/view/acewiki", nil)
			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			for i := 0; i < n; i++ {
				service.ViewHandler(rr, req)
			}
		})
	}
	// Run the benchmark
}

func BenchmarkSaveHandler(b *testing.B) {
	reader := strings.NewReader("save test")
	req, _ := http.NewRequest("POST", "/save/acewiki", reader)
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	for i := 0; i < b.N; i++ {
		service.SaveHandler(rr, req)
	}
}
