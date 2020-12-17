package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Struct passed in the channel by all go routines
type dijkstraResult struct {
	src  string
	prev map[string]string
	dist map[string]int
}

type Queue struct {
	items []string
	m     map[string]int
	pr    map[string]int
}

func (q *Queue) Push(x interface{}) {
	n := len(q.items)
	item := x.(string)
	q.m[item] = n
	q.items = append(q.items, item)
}
func (q *Queue) Pop() interface{} {
	old := q.items
	n := len(old)
	item := old[n-1]
	q.m[item] = -1
	q.items = old[0 : n-1]
	return item
}

func (q *Queue) Len() int           { return len(q.items) }
func (q *Queue) Less(i, j int) bool { return q.pr[q.items[i]] < q.pr[q.items[j]] }
func (q *Queue) Swap(i, j int) {
	q.items[i], q.items[j] = q.items[j], q.items[i]
	q.m[q.items[i]] = i
	q.m[q.items[j]] = j
}

func (q *Queue) update(item string, priority int) {
	q.pr[item] = priority
	heap.Fix(q, q.m[item])
}

func (q *Queue) addWithPriority(item string, priority int) {
	heap.Push(q, item)
	q.update(item, priority)
}

// Channel filled with goroutines Dijkstra's results and used by the output function
var channelOutput = make(chan dijkstraResult, 10)

// Waitgroup to synchronize variables
var wg sync.WaitGroup

func main() {
	fmt.Printf("CPU threads given %d\n", runtime.NumCPU())
	var nbNoeud int
	start := time.Now()

	// Check if there is a file name given
	if len(os.Args) < 2 {
		fmt.Println("No argument given")
		os.Exit(1)
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	// Instanciate reader
	reader := bufio.NewReader(file)
	graph := make(map[string]map[string]int)
	for {
		// Read a line
		line, err := reader.ReadString('\n')

		if err != nil {
			// Break if EOF
			if err == io.EOF {
				break
			}
		}
		// Remove the \n at the end of
		line = strings.TrimSuffix(line, "\r\n") //for windows
		line = strings.TrimSuffix(line, "\n")   //for mac & unix
		// Split line with " " to get the relation in the file
		relation := strings.Split(line, " ")

		if len(relation) == 2 {
			nbNoeud, err = strconv.Atoi(relation[0])
			if err != nil {
				os.Exit(3)
			}
			//nbLien := relation[1]
		}

		// If the it is a relation (ex : "A B 12")
		if len(relation) == 3 {
			// Origin vertex of the relation
			origin := relation[0]

			// Destination of the relation
			destination := relation[1]

			// Get the value of the relation
			value, err := strconv.Atoi(relation[2])

			// Exit if the weight of the relation is not defined
			if err != nil {
				log.Fatal(err)
				os.Exit(4)
			}

			// If there are no existing relations for the
			if graph[origin] == nil {
				graph[origin] = make(map[string]int)
			}
			if graph[destination] == nil {
				graph[destination] = make(map[string]int)
			}

			//pondère l'arête
			graph[origin][destination] = value
			graph[destination][origin] = value
		}
	}

	var ensemble []string
	for key := range graph {
		ensemble = append(ensemble, key)
	}

	fileReadTime := time.Since(start)
	// Increment waitgroup and launch output in a goroutine
	wg.Add(1)
	go output(nbNoeud)

	// Launch dijsktra's shortest path function for each vertex in a goroutine
	for _, v := range ensemble {
		wg.Add(1)
		go Shortestpath(graph, v, ensemble)
	}

	// Wait for all goroutines to end and print execution time
	wg.Wait()
	fmt.Printf("FINISHED IN %v %v\n", fileReadTime, time.Since(start))
}

const (
	// Infinity represents the infinity in Dijkstra
	Infinity = int(^uint(0) >> 1)
	// Uninitialized is the state when there are no previous vertex in Dijkstra
	Uninitialized = ""
)

// Shortestpath function that takes a graph in the form of a map, the source vertex, and a set containing all the vertexes
// It returns the map dist that represents the shortest distance between the source and all vertexes
// and the map prev that links a vertex to its previous
func Shortestpath(graph map[string]map[string]int, src string, ensemble []string) {

	dist := make(map[string]int)
	prev := make(map[string]string)

	// Set length to 0 for the source
	dist[src] = 0

	q := &Queue{[]string{}, make(map[string]int), make(map[string]int)}

	// Set the value to the infinity
	for _, v := range ensemble {
		if v != src {
			dist[v] = Infinity
		}
		prev[v] = Uninitialized
		q.addWithPriority(v, dist[v])
	}

	for len(q.items) != 0 {
		u := heap.Pop(q).(string)
		for v := range graph[u] {
			alt := dist[u] + graph[u][v]
			if alt < dist[v] {
				dist[v] = alt
				prev[v] = u
				q.update(v, alt)
			}
		}
	}

	// Push results in the channel in the form of a structure
	r := dijkstraResult{src: src, prev: prev, dist: dist}
	channelOutput <- r
	wg.Done()
}

func getPath(dest string, prev map[string]string) string {
	str := dest
	for prev[dest] != "" {
		str = prev[dest] + "->" + str
		dest = prev[dest]
	}

	return str
}

// Function to output all results
func output(nbNoeud int) {
	var r dijkstraResult
	for i := 0; i < nbNoeud; i++ {
		r = <-channelOutput
		fmt.Println("From vertex " + r.src)
		fmt.Println("Distance :", r.dist)
		fmt.Println("Map of previouces", r.prev)
		fmt.Println("---------------------------")
	}
	wg.Done()
}
