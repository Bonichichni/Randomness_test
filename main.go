package main

import (
	"crypto/rand"
	"fmt"
	"math"
)

func main() {
	count := 0
	for i := 0; i < 1000; i++ {
		testSeq := GenRanSeq()
		if TestRandom(testSeq) {
			count++
		}
	}
	fmt.Println("Кількість випадкових послідовностей які пройшли всі тести:", count)

}

// Генератор послідовності
func GenRanSeq() [2500]byte {
	var seqArr [2500]byte
	_, err := rand.Read(seqArr[:])
	if err != nil {
		fmt.Println("Помилка при генерації випадкових бітів", err)
	}
	return seqArr
}

// Монобітний тест
func monobitTest(testArr [2500]byte) bool {
	onesCounter := 0
	for i := 0; i < 2500; i++ {
		onesCounter += countOnes(testArr[i])
	}

	return onesCounter > 9654 && onesCounter < 10346
}

// Підраховує кількість одиниць в байті
func countOnes(b byte) int {
	counter := 0
	for i := 0; i < 8; i++ {
		if b&(1<<uint(i)) != 0 {
			counter++
		}
	}
	return counter
}

// Тест максимальної довжини серії
func maxSerLenTest(testArr [2500]byte) bool {
	maxLen := 36
	zeroSeries := 0
	oneSeries := 0

	for _, byteArr := range testArr {
		for i := 0; i < 8; i++ {
			if byteArr&(1<<uint(i)) != 0 {
				oneSeries++
				zeroSeries = 0
			} else {
				zeroSeries++
				oneSeries = 0
			}
			if zeroSeries > maxLen || oneSeries > maxLen {
				return false
			}
		}
	}
	return true
}

// Тест Поккера
func pokerTest(testArr [2500]byte) bool {
	var m float64 = 4
	var Y float64 = 20000
	k := Y / m
	firstPartByte := 0  // Я не бачу доцільним створювання інших структур даних в якому будуть записані блоки по 4 біти, тому я буду просто проходити по
	secondPartByte := 0 //масиву який генерується генератором, ділити цей байт на дві частини і уже далі працювати з ним і шукати повторення
	blockCounts := make(map[int]int)
	for _, b := range testArr {
		firstPartByte = int(b >> 4)
		secondPartByte = int(b & 0x0F)
		blockCounts[firstPartByte]++
		blockCounts[secondPartByte]++
	}
	X3 := calculateX3(blockCounts, m, k)
	return X3 > 1.03 && X3 < 57.4
}

// Функція яка рахує X3
func calculateX3(blockCounts map[int]int, m, k float64) float64 {
	var sum float64 = 0
	for i := 0; i < 16; i++ {
		sum += math.Pow(float64(blockCounts[i]), 2)
	}
	return math.Pow(2, m)/k*float64(sum) - k
}

// Тест довжин серій
func lenghtSeriesTest(testArr [2500]byte) bool {
	zeroSeries := 0
	oneSeries := 0
	seriesCountOnes := make(map[int]int)
	seriesCountZeros := make(map[int]int)
	for _, b := range testArr {
		for i := 0; i < 8; i++ {
			if b&(1<<uint(i)) != 0 {
				oneSeries++
				seriesCountZeros[zeroSeries]++
				zeroSeries = 0
			} else {
				zeroSeries++
				seriesCountOnes[oneSeries]++
				oneSeries = 0
			}
		}
	}
	return checkLenghtInterval(seriesCountOnes) && checkLenghtInterval(seriesCountZeros)
}

// Функція яка перевіряє чи входять серії в відповідні інтервали
func checkLenghtInterval(m map[int]int) bool {
	counter := 0 // я буду збільшувати лічильник на одиницю якщо серія буде попадати в потрібний інтервал, і якщо лічильник вкінці буде дорівнювати 6(бо 6 серій)
	// то функція поверне true
	counterSixPlus := 0
	for k, v := range m {
		if k >= 6 {
			counterSixPlus += v
		}
	}
	if m[1] >= 2267 && m[1] <= 2733 {
		counter++
	}
	if m[2] >= 1079 && m[2] <= 1421 {
		counter++
	}
	if m[3] >= 502 && m[3] <= 748 {
		counter++
	}
	if m[4] >= 223 && m[4] <= 402 {
		counter++
	}
	if m[5] >= 90 && m[5] <= 223 {
		counter++
	}
	if counterSixPlus >= 90 && counterSixPlus <= 223 {
		counter++
	}
	return counter == 6

}

// Загальна функція тестування яка поверне true, тільки тоді коли випадкова послідовність пройде всі тести
func TestRandom(testArr [2500]byte) bool {
	return monobitTest(testArr) && maxSerLenTest(testArr) && pokerTest(testArr) && lenghtSeriesTest(testArr)
}
