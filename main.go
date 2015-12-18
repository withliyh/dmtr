package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/withliyh/dmtr/core"
)

func readfile(path *string) []byte {
	f, err := os.Open(*path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	bytes, err := ioutil.ReadAll(f)
	return bytes
}

func parseConfigFile(path string) []string {
	var ss = make([]string, 0)
	reader, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		l := scanner.Text()
		ss = append(ss, l)
	}
	return ss
}

var defaultconf string = "dmtr.conf"

func main() {
	path := flag.String("c", "", "config file path")
	flag.Parse()
	if strings.EqualFold(*path, "") {
		path = &defaultconf
	}
	fmt.Printf("Read config file:%s\n", *path)

	var ss = parseConfigFile(*path)
	var resultmap = dmtr.NewResultMap()
	var waitgroup sync.WaitGroup
	dmtr.NewExecuter(&waitgroup, ss, resultmap)
	waitgroup.Wait()

	sorter := dmtr.NewSorter()
	for k, v := range *resultmap {
		lostsum := 0.0
		count := 0
		if len(v) < 2 {
			continue
		}
		r := v[1:] //skip first row
		for _, e := range r {
			if e.Lost > float64(0) {
				lostsum += e.Lost
				count++
			}
		}
		lostavg := lostsum / float64(count)
		sorter.Add(k, lostavg)
	}

	fmt.Println("=============================")
	if len(sorter.Sorter) == 0 {
		return
	}
	sort.Sort(sorter)
	for _, item := range sorter.Sorter {
		fmt.Printf("%s %f\n", item.Key, item.Val)
	}
	fmt.Println("=============================")
	best := sorter.Sorter[0]
	fmt.Printf("Best server info:%s\n", best.Key)
	e := (*resultmap)[best.Key]
	for _, row := range e {
		fmt.Print(row.String())
	}
}
