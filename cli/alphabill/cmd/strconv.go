package cmd

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// stringToAmount converts string and decimals to uint64 amount
func stringToAmount(amountIn string, decimals uint32) (uint64, error) {
	if amountIn == "" {
		return 0, fmt.Errorf("invalid empty amount string")
	}
	amountIn = strings.ReplaceAll(amountIn, "'", "")
	splitAmount := strings.Split(amountIn, ".")
	if len(splitAmount) > 2 {
		return 0, fmt.Errorf("invlid amount string %s: more than one comma", amountIn)
	}
	integerStr := splitAmount[0]
	if len(integerStr) == 0 {
		return 0, fmt.Errorf("invalid amount string %s: missing integer part", amountIn)
	}
	// no comma, only integer part
	if len(splitAmount) == 1 {
		// pad with decimal number of 0's (alternative would be to convert and then multiply by 10 to the power of decimals)
		integerStr += strings.Repeat("0", int(decimals))
		amount, err := strconv.ParseUint(integerStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid amount string \"%s\": error conversion to uint64 failed, %v", amountIn, err)
		}
		return amount, nil
	}
	fractionStr := splitAmount[1]
	if len(fractionStr) == 0 {
		return 0, fmt.Errorf("invalid amount string %s: missing fraction part", amountIn)
	}
	// there is a comma in the value
	if uint32(len(fractionStr)) > decimals {
		return 0, fmt.Errorf("invalid precision: %s", amountIn)
	}
	// pad with 0's in input is smaller than decimals
	if uint32(len(fractionStr)) < decimals {
		// append 0's so that decimal number of fraction places are present
		fractionStr += strings.Repeat("0", int(decimals)-len(fractionStr))
	}
	// convert the combined string "integer+fraction" to amount
	amount, err := strconv.ParseUint(integerStr+fractionStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid amount string \"%s\": error conversion to uint64 failed, %v", amountIn, err)
	}
	return amount, nil
}

// amountToString converts amount to string with specified decimals
// NB! it is assumed that the decimal places value is sane and verified before
// calling this method.
func amountToString(amount uint64, decimals uint32) string {
	amountStr := strconv.FormatUint(amount, 10)
	if decimals == 0 {
		return insertSeparator(amountStr, false)
	}
	// length of amount string is less than decimal places, insert decimal point in value
	if decimals < uint32(len(amountStr)) {
		return insertSeparator(amountStr[:uint32(len(amountStr))-decimals], false) + "." + insertSeparator(amountStr[uint32(len(amountStr))-decimals:], true)
	}
	// resulting amount is less than 0
	resultStr := ""
	resultStr += strings.Repeat("0", int(decimals)-len(amountStr))
	return "0." + insertSeparator(resultStr+amountStr, true)
}

// insert apostrophe as thousands separator. The reverse flag defines the direction in which the insertion should happen
// insertSeparator("1234", false) => 1'234 (for the integral part)
// insertSeparator("1234", true) => 123'4 (for the fractional part)
func insertSeparator(value string, reverse bool) string {
	var re *regexp.Regexp
	if reverse {
		re = regexp.MustCompile(`(\d{3})(\d+)`)
	} else {
		re = regexp.MustCompile(`(\d+)(\d{3})`)
	}
	for n := ""; n != value; {
		n = value
		value = re.ReplaceAllString(value, "$1'$2")
	}
	return value
}
