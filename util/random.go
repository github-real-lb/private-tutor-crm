package util

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var r *rand.Rand

const ABC = "abcdefghijklmnopqrstuvwxyz"

func init() {
	src := rand.NewSource(time.Now().UnixMicro())
	r = rand.New(src)
}

// RandomInt generates a random integer between min and max.
func RandomInt64(min, max int64) int64 {
	return min + r.Int63n(max-min+1)
}

// RandomFloat64 generates a random decimal between min and max.
func RandomFloat64(min, max float64) float64 {
	delta := max - min
	i := r.Int63n(int64(delta))
	f := r.Float64()

	if n := min + float64(i) + f; n > max {
		return max
	} else {
		return n
	}
}

// RandomString generates a random string of lenght n.
func RandomString(n int) string {
	var sb strings.Builder
	var b byte

	l := len(ABC)

	for i := 0; i < n; i++ {
		b = ABC[rand.Intn(l)]
		sb.WriteByte(b)
	}

	return sb.String()
}

// RandomNullInt64 generates a random integer between min and max.
// Valid determines if NullInt64 is null (false) or not (true).
func RandomNullInt64(min, max int64, valid bool) sql.NullInt64 {
	return sql.NullInt64{
		Int64: RandomInt64(min, max),
		Valid: valid,
	}
}

// RandomNullFloat64 generates a random decimal between min and max.
// Valid determines if NullFloat64 is null (false) or not (true).
func RandomNullFloat64(min, max float64, valid bool) sql.NullFloat64 {
	return sql.NullFloat64{
		Float64: RandomFloat64(min, max),
		Valid:   valid,
	}
}

// RandomNullString generates a random string of lenght n.
// Valid determines if NullString is null (false) or not (true).
func RandomNullString(n int, valid bool) sql.NullString {
	return sql.NullString{
		String: RandomString(n),
		Valid:  valid,
	}
}

// NullNullInt64 generates a null NullInt64
func NullNullInt64() sql.NullInt64 {
	return sql.NullInt64{
		Int64: 0,
		Valid: false,
	}
}

// NullNullFloat64 generates a null NullFloat64
func NullNullFloat64() sql.NullFloat64 {
	return sql.NullFloat64{
		Float64: 0.0,
		Valid:   false,
	}
}

// NullNullSting generates a null NullString
func NullNullSting() sql.NullString {
	return sql.NullString{
		String: "",
		Valid:  false,
	}
}

// RandomName generates a random first or last name of 8 characters long
func RandomName() string {
	return RandomString(8)
}

// RandomName generates a random e-mail
func RandomEmail() string {
	return fmt.Sprint(RandomString(10), "@gmail.com")
}

// RandomPhoneNumber generates a random phone number
func RandomPhoneNumber() string {
	return fmt.Sprintf("+%d %d-%d",
		RandomInt64(100, 999),
		RandomInt64(1000, 9999),
		RandomInt64(1000, 9999))
}

// RandomAddress generates a random address
func RandomAddress() string {
	return fmt.Sprintf(
		`street %s %d
		%s %d
		%s`,
		RandomString(8), RandomInt64(10, 99),
		RandomString(8), RandomInt64(10000, 99999),
		RandomString(8))
}
