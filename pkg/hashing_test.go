package pkg

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func Test_should_match(t *testing.T) {
	//value := "$2a$10$oryrXDxtlYTflG6aFMoWoOk6JCx7qlW832SmdRx2KjJwTHsmonWVe"
	value := "$2a$10$zw3/lkIm1MgE0lutWYkbiul/Ni5qSkPquaS6D3LHi3PlOVDYvuhVy"
	err := bcrypt.CompareHashAndPassword([]byte(value), []byte("123456"))

	if err != nil {
		t.Errorf("Not matched")
	}
}

//Two calls on the same value produces different outputs
func Test_bcryptHash_result_is_not_equalt(t *testing.T){
	value := "secret"
	first, _ := bcrypt.GenerateFromPassword([]byte(value), 14)
	second, _ := bcrypt.GenerateFromPassword([]byte(value), 14)
	if string(first) == string(second) {
		t.Errorf("%s should not equal to %s", first, second)
	}
}

func Test_bcryptHash_should_successful_compare_passwords(t *testing.T) {
	value := "secret"
	hash, _ := bcrypt.GenerateFromPassword([]byte(value), 14)
	err := bcrypt.CompareHashAndPassword(hash, []byte(value))
	if err != nil {
		t.Errorf("Compare failed")
	}
}

func hashAndCompareCostDurationComparison(plainPassword string, cost int) (time.Duration, time.Duration) {
	start := time.Now()
	hash, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), cost)
	end := time.Now()
	hashDuration := end.Sub(start)

	start = time.Now()
	bcrypt.CompareHashAndPassword(hash, []byte(plainPassword))
	end = time.Now()
	compareDuration := end.Sub(start)

	return hashDuration, compareDuration
}

func Test_bcryptHash_show_hashing_execution_time_with_different_costs(t *testing.T) {
	passwords := []string{"secret", "123456", "skldh__!ksdhg03485789nfvsd!31"}
	costs:= []int{1, 10, 14, 15 ,16}
	for _, cost := range costs {
		for _, password := range passwords {
			hashDuration, compareDuration := hashAndCompareCostDurationComparison(password, cost)
			fmt.Printf("%d %v %v %s\n", cost, hashDuration, compareDuration, password)
		}
	}
}



