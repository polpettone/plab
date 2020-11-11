package pkg

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)


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

func Test_bcryptHash_show_hashing_execution_time_with_different_costs(t *testing.T) {

	value := "secret"

	start := time.Now()
	bcrypt.GenerateFromPassword([]byte(value), 4)
	end := time.Now()
	duration := end.Sub(start)
	fmt.Printf("Cost 4 Duration %d\n", duration.Milliseconds())


	start = time.Now()
	bcrypt.GenerateFromPassword([]byte(value), 10)
	end = time.Now()
	duration = end.Sub(start)
	fmt.Printf("Cost 10 Duration %d\n", duration.Milliseconds())

	start = time.Now()
	bcrypt.GenerateFromPassword([]byte(value), 14)
	end = time.Now()
	duration = end.Sub(start)
	fmt.Printf("Cost 14 Duration %d\n", duration.Milliseconds())

	start = time.Now()
	bcrypt.GenerateFromPassword([]byte(value), 15)
	end = time.Now()
	duration = end.Sub(start)
	fmt.Printf("Cost 15 Duration %d\n", duration.Milliseconds())

	start = time.Now()
	bcrypt.GenerateFromPassword([]byte(value), 16)
	end = time.Now()
	duration = end.Sub(start)
	fmt.Printf("Cost 16 Duration %d\n", duration.Milliseconds())

}



