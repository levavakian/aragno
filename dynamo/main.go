package dynamo

import "fmt"
import "time"
import "aragno/zero"

func PrintParticles(b Body) {
	for _, s := range b.Shapes {
		s.Print()
	}
}

func RunSimulation() {
	n := 2
	particles := make([]shape , n)
	particles[0] = NewRectangle("a", 1, 1, 1, zero.Pose {0,0,0})
	particles[1] = NewRectangle("a",1, 1, 1, zero.Pose {0,0,0})
	b := Body{particles, make([]Joint,0)}
	var simTime int
	var totalTime int = 100
	var dt float64 = 1
	for simTime < totalTime {
		UpdateBody(&b, zero.Pose{0,0,0}, dt)
		simTime += 1
		fmt.Println("=================SIM STATE============================")
		PrintParticles(b)
		fmt.Println("======================================================")
		time.Sleep(20)
	}
}

func main() {
	RunSimulation()
}
