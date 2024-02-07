package util

import (
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
	var i int64 = 0
	if delta := max - min; delta >= 1.0 {
		i = r.Int63n(int64(delta))
	}

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

// RandomDatetime generates a random datetime for the past year.
func RandomDatetime() time.Time {
	standardYear := time.Hour * 24 * 365
	randomDuration := time.Duration(r.Int63n(int64(standardYear)))
	return time.Now().Add(-randomDuration)
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

// RandomHourlyFee generates a random hourly fee between 85.0 to 300.0
func RandomHourlyFee() float64 {
	return RandomFloat64(85.0, 300.0)
}

// RandomNote generates a random note
func RandomNote() string {
	return fmt.Sprint("This is a random note:\n",
		RandomString(10), "\n",
		RandomString(10))
}

// RandomLessonDuration returns random int64 between 30 and 240 minutes
func RandomLessonDuration() int64 {
	return RandomInt64(30, 240)
}

// RandomLessonHourlyFee returns random float64 between 85.0 and 300.0
func RandomLessonHourlyFee() float64 {
	return RandomFloat64(85.0, 300.0)
}

// RandomDiscount generates a random discount % between 0.0 (0%) to 0.99 (99%)
func RandomDiscount() float64 {
	return RandomFloat64(0.0, 0.99)
}

// RandomInvoiceAmount returns random float64 between 85.0 and 1200.0
func RandomInvoiceAmount() float64 {
	return RandomFloat64(85.0, 1200.0)
}

// RandomPaymentAmount returns random float64 between 85.0 and 1200.0
func RandomPaymentAmount() float64 {
	return RandomFloat64(85.0, 1200.0)
}
