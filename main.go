package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var waitgroup sync.WaitGroup
var resultmap map[string][]Entry
var err error
var rsleep = rand.New(rand.NewSource(time.Now().UnixNano()))

type Entry struct {
	Ord   int
	Host  string
	Lost  float64
	Snt   int
	Last  float64
	Avg   float64
	Best  float64
	Wrst  float64
	Stdev float64
}

func (this *Entry) String() string {
	return fmt.Sprintf("%d %v %v %s", this.Ord, this.Host, this.Lost, "\n")
}

func PrintErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func run(s string) {
	fmt.Println("mtr -r ", s)
	cmd := exec.Command("mtr", "-r", s)
	stdout, err := cmd.StdoutPipe()
	stderr, err := cmd.StderrPipe()
	cmd.Start()

	go io.Copy(os.Stdout, stderr)

	scanner := bufio.NewScanner(stdout)
	//skip first 2 line
	scanner.Scan()
	scanner.Scan()

	items := make([]Entry, 0)
	for scanner.Scan() {
		line := scanner.Text()
		field := strings.Fields(line)
		if len(field) < 9 {
			continue
		}

		item := Entry{}

		ord_mtx := strings.Split(field[0], ".")
		item.Ord, err = strconv.Atoi(ord_mtx[0])
		PrintErr(err)

		item.Host = field[1]

		lost := strings.TrimRight(field[2], "%")
		item.Lost, err = strconv.ParseFloat(lost, 64)
		PrintErr(err)

		item.Snt, err = strconv.Atoi(field[3])
		PrintErr(err)

		item.Last, err = strconv.ParseFloat(field[4], 64)
		PrintErr(err)

		item.Avg, err = strconv.ParseFloat(field[5], 64)
		PrintErr(err)

		item.Best, err = strconv.ParseFloat(field[6], 64)
		PrintErr(err)

		item.Wrst, err = strconv.ParseFloat(field[7], 64)
		PrintErr(err)

		item.Stdev, err = strconv.ParseFloat(field[8], 64)
		PrintErr(err)

		items = append(items, item)

	}
	resultmap[s] = items
	cmd.Wait()
	waitgroup.Done()
}

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

func main() {
	path := flag.String("c", "", "config file path")
	flag.Parse()

	var ss = parseConfigFile(*path)

	resultmap = make(map[string][]Entry)
	for _, s := range ss {
		waitgroup.Add(1)
		go run(s)
		time.Sleep(5 * time.Second)
	}

	waitgroup.Wait()
	for k, v := range resultmap {
		fmt.Println("-------------------------------------")
		fmt.Println(k)
		for _, item := range v {
			fmt.Printf(item.String())
		}
	}
}
