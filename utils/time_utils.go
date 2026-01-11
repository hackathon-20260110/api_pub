package utils

import "time"

// 二つの日付を受け取り、birthDate時点の人のageを計算する
func CalculateAge(birthDate time.Time, now time.Time) int {
	age := now.Year() - birthDate.Year()
	if (now.Month() < birthDate.Month()) ||
		(now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}
	return age
}
