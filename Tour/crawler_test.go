package main

import "testing"

func BenchmarkCrawler(b *testing.B) {
	b.ResetTimer()
	Crawler()
}
