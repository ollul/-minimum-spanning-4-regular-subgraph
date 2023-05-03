package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var distance [][]uint16
var incidence [][]uint16

type point struct {
	x, y uint16
}

type edge struct {
	v1, v2 int
	len    uint16
}

func max(a, b uint16) uint16 {
	if a > b {
		return a
	}

	return b
}

func taxi(v1, v2 point) (res uint16) {
	if v1.x > v2.x {
		res += v1.x - v2.x
	} else {
		res += v2.x - v1.x
	}
	if v1.y > v2.y {
		res += v1.y - v2.y
	} else {
		res += v2.y - v1.y
	}
	return
}

func main() {
	rand.Seed(time.Now().UnixNano())
	num, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return
	}
	target_cost, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return
	}
	edges_file, _ := os.Open(fmt.Sprintf("%d.txt", num))
	vertices_file, _ := os.Open(fmt.Sprintf("../Taxicab_%d.txt", num))
	defer edges_file.Close()
	defer vertices_file.Close()

	vertices, _ := io.ReadAll(vertices_file)
	points := []point{}

	str := strings.Split(string(vertices), "\r\n")

	for _, v := range str {
		nums := strings.Split(v, string([]byte{0x09}))
		if len(nums) < 2 {
			continue
		}
		x, _ := strconv.Atoi(nums[0])
		y, _ := strconv.Atoi(nums[1])
		points = append(points, point{x: uint16(x), y: uint16(y)})
	}

	for i := 0; i < num; i++ {
		incidence = append(incidence, make([]uint16, num))
		distance = append(distance, make([]uint16, num))
	}

	for i, v1 := range points {
		for j, v2 := range points {
			distance[i][j] = taxi(v1, v2)
		}
	}

	edges, _ := io.ReadAll(edges_file)
	edge_list := []*edge{}
	str = strings.Split(string(edges), "\n")
	for _, line := range str {
		comp := strings.Split(line, " ")
		if comp[0] != "e" {
			continue
		}
		v1, _ := strconv.Atoi(comp[1])
		v2, _ := strconv.Atoi(comp[2])
		v1--
		v2--
		incidence[v1][v2] = 1
		incidence[v2][v1] = 1
		edge_list = append(edge_list, &edge{v1: v1, v2: v2, len: distance[v1][v2]})
	}

	vert1 := []int{}
	for i := 0; i < len(edge_list)*32; i++ {
		vert1 = append(vert1, i%len(edge_list))
	}
	cost := 100000000
	rand.Shuffle(len(vert1), func(i, j int) { vert1[i], vert1[j] = vert1[j], vert1[i] })

	for {
		for k := 0; k < len(vert1); k += 2 {
			e1 := edge_list[vert1[k]]
			e2 := edge_list[vert1[k+1]]
			i11 := e1.v1
			i12 := e1.v2
			i21 := e2.v1
			i22 := e2.v2
			r := rand.Int31n(2)
			if r == 0 && i11 != i21 && i12 != i22 && incidence[i11][i21] == 0 && incidence[i12][i22] == 0 {
				incidence[e1.v1][e1.v2] = 0
				incidence[e1.v2][e1.v1] = 0
				incidence[e2.v1][e2.v2] = 0
				incidence[e2.v2][e2.v1] = 0
				e1.v1 = i11
				e1.v2 = i21
				e2.v1 = i12
				e2.v2 = i22
				e1.len = distance[e1.v1][e1.v2]
				e2.len = distance[e2.v1][e2.v2]
				incidence[e1.v1][e1.v2] = 1
				incidence[e1.v2][e1.v1] = 1
				incidence[e2.v1][e2.v2] = 1
				incidence[e2.v2][e2.v1] = 1
			} else if r == 1 && i11 != i22 && i21 != i12 && incidence[i11][i22] == 0 && incidence[i21][i12] == 0 {
				incidence[e1.v1][e1.v2] = 0
				incidence[e1.v2][e1.v1] = 0
				incidence[e2.v1][e2.v2] = 0
				incidence[e2.v2][e2.v1] = 0
				e1.v1 = i11
				e1.v2 = i22
				e2.v1 = i12
				e2.v2 = i21
				e1.len = distance[e1.v1][e1.v2]
				e2.len = distance[e2.v1][e2.v2]
				incidence[e1.v1][e1.v2] = 1
				incidence[e1.v2][e1.v1] = 1
				incidence[e2.v1][e2.v2] = 1
				incidence[e2.v2][e2.v1] = 1
			}
		}
		for l := 0; l < 4096; l++ {
			if l%10 == 0 {
				fmt.Println(l)
			}
			//		rand.Shuffle(len(vert1), func(i, j int) { vert1[i], vert1[j] = vert1[j], vert1[i] })
			changed := false
			for i := 0; i < len(vert1); i += 2 {
				if vert1[i] == vert1[i+1] {
					continue
				}
				// for j := i + 1; j < len(edge_list); j++ {
				e1 := edge_list[vert1[i]]
				e2 := edge_list[vert1[i+1]]
				t := 0
				min_len := e1.len + e2.len
				// i11 - i12
				// i21 - i22
				i11 := e1.v1
				i12 := e1.v2
				i21 := e2.v1
				i22 := e2.v2
				// i11 - i21
				// i12 - i22
				if i11 != i21 && i12 != i22 && incidence[i11][i21] == 0 && incidence[i12][i22] == 0 {
					d := distance[i11][i21] + distance[i12][i22]
					if d <= min_len {
						t = 1
						min_len = d
					}
				}
				// i11 - i22
				// i21 - i12
				if i11 != i22 && i21 != i12 && incidence[i11][i22] == 0 && incidence[i21][i12] == 0 {
					d := distance[i11][i22] + distance[i12][i21]
					if d <= min_len {
						t = 2
						min_len = d
					}
				}

				switch t {
				case 1:
					changed = true
					incidence[e1.v1][e1.v2] = 0
					incidence[e1.v2][e1.v1] = 0
					incidence[e2.v1][e2.v2] = 0
					incidence[e2.v2][e2.v1] = 0
					e1.v1 = i11
					e1.v2 = i21
					e2.v1 = i12
					e2.v2 = i22
					e1.len = distance[e1.v1][e1.v2]
					e2.len = distance[e2.v1][e2.v2]
					incidence[e1.v1][e1.v2] = 1
					incidence[e1.v2][e1.v1] = 1
					incidence[e2.v1][e2.v2] = 1
					incidence[e2.v2][e2.v1] = 1
				case 2:
					changed = true
					incidence[e1.v1][e1.v2] = 0
					incidence[e1.v2][e1.v1] = 0
					incidence[e2.v1][e2.v2] = 0
					incidence[e2.v2][e2.v1] = 0
					e1.v1 = i11
					e1.v2 = i22
					e2.v1 = i12
					e2.v2 = i21
					e1.len = distance[e1.v1][e1.v2]
					e2.len = distance[e2.v1][e2.v2]
					incidence[e1.v1][e1.v2] = 1
					incidence[e1.v2][e1.v1] = 1
					incidence[e2.v1][e2.v2] = 1
					incidence[e2.v2][e2.v1] = 1
				}
				// }
			}
			if !changed {
				break
			}
		}
		cur_cost := 0
		for i := 0; i < num; i++ {
			for j := i + 1; j < num; j++ {
				if incidence[i][j] == 1 {
					cur_cost += int(taxi(points[i], points[j]))
				}
			}
		}
		if cur_cost < cost {
			cost = cur_cost
			fmt.Println(cost)
			var max_edge uint16
			cc := 0
			for i := 0; i < num; i++ {
				for j := i + 1; j < num; j++ {
					if incidence[i][j] == 1 {
						cc += int(taxi(points[i], points[j]))
						fmt.Printf("e %d %d\n", i+1, j+1)
						max_edge = max(max_edge, taxi(points[i], points[j]))
					}
				}
			}
			fmt.Printf("c Вес 4-регулярного подграфа = %d, самое длинное ребро = %d\n", cost, max_edge)
			fmt.Printf("p edge %d %d\n", num, 4*num/2)

			if cost < target_cost {
				break
			}
		}
	}
	cost = 0
	var max_edge uint16
	max_edge = 0
	for i := 0; i < num; i++ {
		for j := i + 1; j < num; j++ {
			if incidence[i][j] == 1 {
				cost += int(taxi(points[i], points[j]))
				fmt.Printf("e %d %d\n", i+1, j+1)
				max_edge = max(max_edge, taxi(points[i], points[j]))
				//fmt.Printf("\tV%d -- V%d [label=\"%d\"]\n", i, j, taxi(points[i], points[j]))
			}
		}
	}
	fmt.Printf("c Вес 4-регулярного подграфа = %d, самое длинное ребро = %d\n", cost, max_edge)
	fmt.Printf("p edge %d %d\n", num, 4*num/2)
}
