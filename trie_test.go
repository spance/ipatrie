package ipatrie

import (
	"strings"
	"testing"
	"unsafe"
)

func verify(t testing.TB, trie *Trie, addr string, expected int) {
	found := trie.Match(ParseIPv4(addr))

	if expected >= 0 && !found {
		t.Errorf("Expected [%s]==%v but got=%v", addr, expected, found)
	}
	if expected < 0 && found {
		t.Errorf("Expected [%s] not-found but got=%v", addr, found)
	}
}

func TestSimple(t *testing.T) {
	t.Log("sizeof node", unsafe.Sizeof(TrieNode{}))
	net_samples := []string{
		"1.1.1.0/24",
		"1.1.2.2/30",
		"1.1.2.0/24",
		"1.1.0.0/16",
		"1.2.3.4/32",
	}
	trie := NewTrie()
	for _, s := range net_samples {
		a, m, _ := ParseCIDR(s)
		trie.Insert(a, m)
	}
	trie.printStat(t)
	addr_samples := map[string]int{
		"1.1.1.0":   1,
		"1.1.1.1":   1,
		"1.1.1.255": 1,
		"1.1.2.1":   2,
		"1.1.2.2":   2,
		"1.1.2.3":   2,
		"1.1.2.4":   3,
		"1.1.3.1":   4,
		"1.2.3.4":   5,
		"1.2.3.1":   -1,
		"2.2.4.1":   -1,
		"0.1.3.1":   -1,
	}
	for s, v := range addr_samples {
		verify(t, trie, s, v)
	}
}

func TestSamples(t *testing.T) {
	tab := detectTestTables()
	var size uint32

	for _, f := range tab {
		f = f[9:]
		s := strings.Replace(f, "table", "sample", 1)
		trie := initTestData(f)
		trie.printStat(t)
		if trie.size > size {
			//			benchTable = "table-2.txt"
			benchTable = f
			size = trie.size
		}
		fileIter(s, func(fields []string) {
			verify(t, trie, fields[0], parseField1(fields[1]))
		})
		t.Logf("test %d %s/%s", trie.size, f, s)
	}
}

func BenchmarkTrie(b *testing.B) {
	initBenchTrie(b)
	var ips [10]uint32
	for i := 0; i < len(ips); i++ {
		ips[i] = randomIPv4Addr()
	}

	b.ResetTimer()
	for i, j := 0, 0; i < b.N; i++ {
		for j = 0; j < len(ips); j++ {
			benchTrie.Match(ips[j])
		}
	}
}

func (t *Trie) printStat(tb testing.TB) {
	//tb.Logf("trie left=%d right=%d value=%d size=%d", len(t.lefts), len(t.rights), len(t.values), t.size)
	tb.Logf("trie size=%d", t.size)
}

func initTestData(file string) *Trie {
	trie := NewTrie()
	fileIter(file, func(fields []string) {
		na, m, _ := ParseCIDR(fields[0])
		trie.Insert(na, m)
	})
	return trie
}
