package sim

// import (
// 	"testing"
// )

// func benchmarkSim(maxCells int, b *testing.B) {
// 	sim := Sim{}
// 	sim.Create(false)
// 	sim.maxCells = maxCells

// 	for len(sim.cells) > sim.maxCells {
// 		sim.RunStep()
// 	}

// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		sim.RunStep()
// 	}
// }

// func BenchmarkSim10(b *testing.B) {
// 	benchmarkSim(1e5, b)
// }
// func BenchmarkSim20(b *testing.B) {
// 	benchmarkSim(2e5, b)
// }
// func BenchmarkSim50(b *testing.B) {
// 	benchmarkSim(5e5, b)
// }
// func BenchmarkSim100(b *testing.B) {
// 	benchmarkSim(1e6, b)
// }

// func benchmarkSimKillOldest(maxCells int, b *testing.B) {
// 	sim := Sim{}
// 	sim.Create(false)
// 	sim.maxCells = maxCells

// 	for len(sim.cells) > sim.maxCells {
// 		sim.RunStep()
// 	}

// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		sim.KillOldestCells()
// 	}
// }

// func BenchmarkSimKillOldest10(b *testing.B) {
// 	benchmarkSimKillOldest(1e5, b)
// }
// func BenchmarkSimKillOldest20(b *testing.B) {
// 	benchmarkSimKillOldest(2e5, b)
// }
// func BenchmarkSimKillOldest50(b *testing.B) {
// 	benchmarkSimKillOldest(5e5, b)
// }
// func BenchmarkSimKillOldest100(b *testing.B) {
// 	benchmarkSimKillOldest(1e6, b)
// }
