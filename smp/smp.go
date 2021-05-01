//Adapted from algorithm archive, Coded by Prajwal Krishna Maitin (pk1210)
package smp

import "fmt"
import "time"
import "math/rand"

type Person struct {
	ID      int
	Prefers []int
	Partner *Person
	index   int
}

func shuffle(arr []int, size int) []int {
	seed := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(seed)

	for i := 0; i < size; i++ {
		j := generator.Intn(i + 1)
		temp := arr[i]
		arr[i] = arr[j]
		arr[j] = temp
	}
	return arr
}

func createGroup(size int) []Person {
	group := make([]Person, size)
	for i := 0; i < size; i++ {
		group[i].ID = i
		group[i].Partner = nil
		group[i].index = 0
		arr := make([]int, size)
		for j := 0; j < size; j++ {
			arr[j] = j
		}
		group[i].Prefers = shuffle(arr, size)
	}
	return group
}

func perferredPartner(man *Person, woman *Person, size int) bool {
	//Function checks if man is liked more than woman's current Partner
	Partner := woman.Partner.ID
	bachelor := man.ID
	for i := 0; i < size; i++ {
		if woman.Prefers[i] == Partner {
			return false
		} else if woman.Prefers[i] == bachelor {
			return true
		}
	}
	return false
}
func StageMarriage(men []Person, women []Person, size int) {
	bachelors := make([]*Person, size)
	bachelorsSize := size
  fmt.Printf("size is %v\n", size)
	for i := 0; i < size; i++ {
		bachelors[i] = &men[i]
	}
	for bachelorsSize > 0 {
		man := bachelors[bachelorsSize-1]
		//Selecting preferred women from current man
		woman := &women[man.Prefers[man.index]]

		if woman.Partner == nil {
			//If woman does not have any Partner assign them current bachelors
			woman.Partner = man
			man.Partner = woman
			bachelors[bachelorsSize-1] = nil
			bachelorsSize--
		} else if perferredPartner(man, woman, size) {
			//If woman has a Partner check if she likes current man more
			ex := woman.Partner
			ex.Partner = nil
			woman.Partner = man
			man.Partner = woman
			bachelors[bachelorsSize-1] = ex
		} else {
			//Tell man to look for someother woman
			man.index++
		}
	}
}

func main() {
	size := 5
	men := createGroup(size)
	women := createGroup(size)

	fmt.Println("Men Preference List")
	for i := 0; i < size; i++ {
		fmt.Println("Preference of man ", i, men[i].Prefers)
	}
	fmt.Println("\nWomen Preference List")
	for i := 0; i < size; i++ {
		fmt.Println("Preference of woman ", i, men[i].Prefers)
	}

	StageMarriage(men, women, size)

	fmt.Println("\nStable Marriage")
	for i := 0; i < size; i++ {
		fmt.Println("Partner of woman ", i, "is man", women[i].Partner.ID)
	}
}
