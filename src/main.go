package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
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
			//nbNoeud, err := strconv.Atoi(relation[0])
			//if err != nil {
			//	os.Exit(3)
			//}
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
	fmt.Println(graph["B"]["A"])
	fmt.Println(graph)

	var ensemble []string
	i := 0
	for key := range graph {
		ensemble[i] = key
		i++
	}
	println(ensemble)
	shortestpath(graph, "A", ensemble)
}

func shortestpath(graphe map[string]map[string]int, src string, ensemble []string) {

}
