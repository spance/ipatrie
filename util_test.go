package ipatrie

import (
	"bufio"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	benchTrie  *Trie
	benchTable string
)

func initBenchTrie(t testing.TB) {
	if benchTrie == nil && strings.Contains(strings.Join(os.Args, " "), ".bench") {
		var ms1, ms2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&ms1)
		benchTrie = initTestData(benchTable)
		runtime.GC()
		runtime.ReadMemStats(&ms2)
		benchTrie.printStat(t)
		t.Logf("Mem HeapAlloc=%dk HeapInuse=%dk HeapIdle=%dk HeapSys=%dk \n",
			(int64(ms2.HeapAlloc)-int64(ms1.HeapAlloc))/1024,
			(int64(ms2.HeapInuse)-int64(ms1.HeapInuse))/1024,
			(int64(ms2.HeapIdle)-int64(ms1.HeapIdle))/1024,
			(int64(ms2.HeapSys)-int64(ms1.HeapSys))/1024,
		)
	}
}

func fileIter(file string, fn func(part []string)) {
	if !strings.HasPrefix(file, "testdata") {
		file = "testdata/" + file
	}
	f, e := os.Open(file)
	if e != nil {
		panic(e)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		bs, _, e := r.ReadLine()
		if len(bs) > 0 {
			line := strings.Fields(string(bs))
			fn(line)
		}
		if e != nil {
			break
		}
	}
}

func detectTestTables() []string {
	tab, err := filepath.Glob("testdata/table-*")
	if err != nil {
		panic(err)
	}
	return tab
}

func parseField1(ns string) int {
	var nu int64
	if strings.HasSuffix(ns, ";") {
		nu, _ = strconv.ParseInt(ns[:len(ns)-1], 16, 0)
	} else {
		nu, _ = strconv.ParseInt(ns, 10, 0)
	}
	return int(nu)
}

func randomIPv4Addr() uint32 {
	rand.Seed(time.Now().Unix())
	return uint32(rand.Int63() & 0xffFFffFF)
}
