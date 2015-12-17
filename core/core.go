package dmtr

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Resultmap map[string][]Entry

var (
	resultmap Resultmap
	err       error
	rsleep    = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type Entry struct {
	Ord   int    //序号
	Host  string //主机地址
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

func run(waitgroup *sync.WaitGroup, s string, result *Resultmap) {
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
	(*result)[s] = items
	cmd.Wait()
	waitgroup.Done()
}
func NewResultMap() *Resultmap {
	r := make(Resultmap)
	return &r
}

func NewExecuter(waitgroup *sync.WaitGroup, ss []string, result *Resultmap) {
	for _, s := range ss {
		waitgroup.Add(1)
		go run(waitgroup, s, result)
	}
}
