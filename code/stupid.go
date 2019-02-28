package main

import (
	"fmt"
	"math/rand"
	"os"
    "sync"
	"time"
)

const H int = 0
const V int = 1

func orientFromString(s string) int {
	if s == "H" {
		return H
	} else {
		return V
	}
}

var mutex = &sync.Mutex{}
var scores map[[2]int]int

func scoreID(a, b int) int {
	return scorePhoto(photos[a], photos[b])
}

func scorePhoto(a, b Photo) int {
	if a.id > b.id {
		a, b = b, a
	}
    mutex.Lock()
	s, ok := scores[[2]int{a.id, b.id}]
    mutex.Unlock()
	if ok {
		return s
	}
	shared := 0
	for tag := range a.tags {
		_, ok := b.tags[tag]
		if ok {
			shared++
		}
	}
	disA := len(a.tags) - shared
	disB := len(b.tags) - shared
	s = shared
	for _, x := range []int{shared, disA, disB} {
		if s > x {
			s = x
		}
	}
    mutex.Lock()
	scores[[2]int{a.id, b.id}] = s
    mutex.Unlock()
	return s
}

type Photo struct {
	id     int
	orient int
	tags   map[string]bool
}

type Tag struct {
	tag    string
	photos map[int]bool
}

var photos map[int]Photo
var hPhotos []int
var vPhotos []int
var tags map[string]Tag

type Solution [][]int

var solCount int

func (s Solution) writeOut() {
	solCount++
	f, err := os.Create(
		fmt.Sprintf("%s.%03d.txt",
			os.Args[1], solCount))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
	fmt.Fprintln(f, len(s))
	for _, slide := range s {
		if len(slide) == 1 {
			fmt.Fprintln(f, slide[0])
		} else {
			fmt.Fprintf(f, "%d %d\n", slide[0], slide[1])
		}
	}
}

func main() {
	f, err := os.Open(os.Args[1] + ".txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	defer f.Close()
	var n int
	fmt.Fscan(f, &n)
	photos = make(map[int]Photo, n)
	tags = make(map[string]Tag)
	scores = make(map[[2]int]int, n)
	for i := 0; i < n; i++ {
		var p Photo
		var s string
		fmt.Fscan(f, &s)
		p.id = i
		p.orient = orientFromString(s)
		p.tags = make(map[string]bool)
		var m int
		fmt.Fscan(f, &m)
		for j := 0; j < m; j++ {
			var t string
			fmt.Fscan(f, &t)
			p.tags[t] = true

			tag, ok := tags[t]
			if !ok {
				tag.tag = t
				tag.photos = make(map[int]bool)
			}
			tag.photos[i] = true
			tags[t] = tag
		}
		photos[i] = p
		if p.orient == H {
			hPhotos = append(hPhotos, i)
		} else {
			vPhotos = append(vPhotos, i)
		}
	}
    fmt.Fprintln(os.Stderr, "done with input")

maxTimes := 10
	for z := 0; z < maxTimes; z++ {
		go try()
	}
for {
    time.Sleep(1 * time.Second)
    mutex2.Lock()
if done == maxTimes {
break
}
    mutex2.Unlock()
}
}

var mutex2 = &sync.Mutex{}
var done int = 0

func try() {
defer func() {
    mutex2.Lock()
    done++
    mutex2.Unlock()
}()
	rand.Seed(time.Now().UnixNano())
	cID := rand.Intn(len(photos))
	seenIDs := make(map[int]bool, len(photos))
	score := 0
	var sol Solution
stupid:
	for {
		sol = append(sol, []int{cID})
		seenIDs[cID] = true
		if photos[cID].orient == V {
			for y := 0; y <= len(vPhotos) && seenIDs[cID]; y++ {
				x := rand.Intn(len(vPhotos))
				cID = vPhotos[x]
			}
			if seenIDs[cID] {
				sol = sol[:len(sol)-1]
				break stupid
			}
			sol[len(sol)-1] = append(sol[len(sol)-1], cID)
			seenIDs[cID] = true
		}

		nID := -1
		nScore := 0
		strys := 0
	search:
		for tag := range photos[cID].tags {
			for id := range tags[tag].photos {
				if seenIDs[id] {
					continue
				}
				if nID != -1 &&
					photos[nID].orient == H &&
					photos[id].orient == V {
					continue
				}
				score := scoreID(cID, id)
				if score >= nScore {
					nScore = score
					nID = id
					strys++
					if strys > 2 {
						break search
					}
				}
			}
		}
		if nID == -1 {
			break stupid
		}
		score += scoreID(cID, nID)
		cID = nID
	}
	if len(sol) < 2 {
		fmt.Fprintln(os.Stderr, "womp womp")
		return
	} else {
fmt.Fprintf(os.Stderr, "%03d: %d\n", solCount+1, score)
	}
	sol.writeOut()
}
