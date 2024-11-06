package utils

import (
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func StrToUint64(s string) (uint64, error) {
	// Convert string to uint64 using strconv.ParseUint.
	// It expects the string s, the base of the numeral system (10 for decimal), and the bit size (64 for uint64).
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err // return 0 and the error if conversion fails
	}
	return val, nil // return the converted value and no error if successful
}

func StrToFloat64(s string) (float64, error) {
	// Convert string to float64 using strconv.ParseFloat.
	// It expects the string s and the bit size (64 for float64).
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err // return 0 and the error if conversion fails
	}
	return val, nil // return the converted value and no error if successful
}

// ToString converts any DTO struct into a readable string format
func ToString(dto interface{}) string {
	if dto == nil {
		return "nil"
	}

	val := reflect.ValueOf(dto)
	typ := reflect.TypeOf(dto)

	// Handle pointer types
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return "nil"
		}
		val = val.Elem()
		typ = typ.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Sprintf("%v", dto)
	}

	var fields []string
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		var fieldValue interface{}

		// Special handling for pointer fields
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				fieldValue = "<nil>"
			} else {
				fieldValue = field.Elem().Interface()
			}
		} else {
			fieldValue = field.Interface()
		}

		fields = append(fields, fmt.Sprintf("%s: %v", fieldName, fieldValue))
	}

	return fmt.Sprintf("%s{%s}", typ.Name(), strings.Join(fields, ", "))
}

func StrToBigInt(s string) (*big.Int, error) {
	val := new(big.Int)
	_, ok := val.SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("invalid string for conversion to big.Int: %s", s)
	}
	return val, nil
}

func GenerateRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteRune(letters[rand.Intn(len(letters))])
	}
	return sb.String()
}

// GenerateRandomUint generates a random uint value with 6 digits.
func GenerateRandomUint() uint {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	min := uint(100000) // Smallest 6-digit number
	max := uint(999999) // Largest 6-digit number
	return min + uint(r.Intn(int(max-min+1)))
}

func StripNonPrintable(input string) string {
	var sb strings.Builder
	for _, r := range input {
		if unicode.IsPrint(r) {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
