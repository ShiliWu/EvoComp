package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type coordinate struct {
	x float64
	y float64
}

type travelMap struct {
	path []int
}

var location [1000]coordinate
var wg sync.WaitGroup
var optimaltravelMAP [4]travelMap
var l sync.Mutex

func main() {
	runtime.GOMAXPROCS(4)
	lines, err := readLines("TSP1.txt")
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	//var location [1000]coordinate

	for i, line := range lines {
		data := strings.FieldsFunc(line, func(r rune) bool {
			if r == ',' {
				return true
			}
			return false
		})
		location[i].x, _ = strconv.ParseFloat(data[0], 64)
		location[i].y, _ = strconv.ParseFloat(data[1], 64)
	}

	wg.Add(4)

	rand.Seed(time.Now().UnixNano())
	go randomSearch(1)
	rand.Seed(time.Now().UnixNano())
	go randomSearch(2)
	rand.Seed(time.Now().UnixNano())
	go randomSearch(3)
	rand.Seed(time.Now().UnixNano())
	go randomSearch(4)

	//clean up when the program is forced to exit ctrl^c
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		writeOptimal(optimaltravelMAP[0].path, optimaltravelMAP[1].path, optimaltravelMAP[2].path, optimaltravelMAP[3].path)
		os.Exit(1)
	}()

	wg.Wait()

	writeOptimal(optimaltravelMAP[0].path, optimaltravelMAP[1].path, optimaltravelMAP[2].path, optimaltravelMAP[3].path)
}

func randomSearch(number int) {
	defer wg.Done()
	path := "TSP_RandomSearch_Trail" + strconv.FormatInt(int64(number), 10) + ".txt"
	file, _ := os.Create(path)
	defer file.Close()
	w := bufio.NewWriter(file)

	var size = 1000
	var tlength float64
	var optimal float64
	var locationLocal = location
	optimal = math.MaxFloat64
	//1000000000
	for n := 0; n < 1000000000; n++ {
		tlength = 0
		list := rand.Perm(size)

		for i := 0; i < size-1; i++ {
			//calcuate total length
			tlength = tlength + distance(locationLocal[list[i]], locationLocal[list[i+1]])
		}
		tlength = tlength + distance(locationLocal[list[size-1]], locationLocal[list[0]])
		if tlength < optimal {
			optimal = tlength
			l.Lock()
			optimaltravelMAP[number-1].path = list
			l.Unlock()
		}

		outputString := strconv.FormatInt(int64(n), 10) + " " + strconv.FormatFloat(optimal, 'f', 2, 64) + " " + strconv.FormatFloat(tlength, 'f', 2, 64)
		fmt.Fprintln(w, outputString)
	}

	w.Flush()

}

func distance(a coordinate, b coordinate) float64 {
	return math.Sqrt(math.Pow(a.x-b.x, 2) + math.Pow(a.y-b.y, 2))
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func writeOptimal(optimalpath1 []int, optimalpath2 []int, optimalpath3 []int, optimalpath4 []int) {
	//fmt.Println(optimalpath1[999])
	file, _ := os.Create("randomTSP_Optimal.txt")

	defer file.Close()

	w := bufio.NewWriter(file)
	outputString := ""
	for i := 0; i < 1000; i++ {
		outputString = ""
		outputString = strconv.FormatInt(int64(optimalpath1[i]), 10) + " " + strconv.FormatInt(int64(optimalpath2[i]), 10) + " " + strconv.FormatInt(int64(optimalpath3[i]), 10) + " " + strconv.FormatInt(int64(optimalpath4[i]), 10)
		fmt.Fprintln(w, outputString)
	}

	w.Flush()
}
