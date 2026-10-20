package main

import "fmt"

func main() {
	// Целые со знаком. Нулевое значение всех числовых типов: 0.
	var signed8 int8 = -8    // 8 бит: от -128 до 127.
	var signed16 int16 = -16 // 16 бит: от -32768 до 32767.
	var signed32 int32 = -32 // 32 бита: от -2147483648 до 2147483647.
	var total int64 = 59_970 // 64 бита: от -9223372036854775808 до 9223372036854775807.
	var count int = 3        // 32 или 64 бита, зависит от реализации и архитектуры.
	// int при 32 битах: от -2147483648 до 2147483647.
	// int при 64 битах: от -9223372036854775808 до 9223372036854775807.
	fmt.Printf("int8: %v, %T\n", signed8, signed8)
	fmt.Printf("int16: %v, %T\n", signed16, signed16)
	fmt.Printf("int32: %v, %T\n", signed32, signed32)
	fmt.Printf("int64: %v, %T\n", total, total)
	fmt.Printf("int: %v, %T\n", count, count)

	// Целые без знака: отрицательных значений нет.
	var unsigned8 uint8 = 8    // 8 бит: от 0 до 255.
	var unsigned16 uint16 = 16 // 16 бит: от 0 до 65535.
	var unsigned32 uint32 = 32 // 32 бита: от 0 до 4294967295.
	var unsigned64 uint64 = 64 // 64 бита: от 0 до 18446744073709551615.
	var unsigned uint = 3      // 32 или 64 бита, такой же размер, как у int.
	// uint при 32 битах: от 0 до 4294967295.
	// uint при 64 битах: от 0 до 18446744073709551615.
	fmt.Printf("uint8: %v, %T\n", unsigned8, unsigned8)
	fmt.Printf("uint16: %v, %T\n", unsigned16, unsigned16)
	fmt.Printf("uint32: %v, %T\n", unsigned32, unsigned32)
	fmt.Printf("uint64: %v, %T\n", unsigned64, unsigned64)
	fmt.Printf("uint: %v, %T\n", unsigned, unsigned)

	// uintptr: беззнаковое целое, достаточно широкое для битового представления адреса.
	// В обычном компиляторе Go: 32 или 64 бита, в зависимости от целевой платформы.
	// При 32 битах: 0..4294967295. При 64 битах: 0..18446744073709551615.
	// Это число, а не отслеживаемый сборщиком мусора указатель.
	var addressBits uintptr = 0
	fmt.Printf("uintptr: %v, %T\n", addressBits, addressBits)

	// Дробные числа. В комментариях границы конечных значений указаны приближённо.
	// float32: 32 бита, конечные значения от примерно -3.4028235e38 до +3.4028235e38.
	// Минимальное положительное ненулевое: примерно 1.4012985e-45, включая субнормальные.
	// Точность: примерно 6–7 значащих десятичных цифр, а не цифр после запятой.
	var temperature32 float32 = 36.6
	// float64: 64 бита, конечные значения от примерно -1.7976931348623157e308
	// до +1.7976931348623157e308. Минимальное положительное ненулевое: примерно 4.9406564584124654e-324.
	// Точность: примерно 15–16 значащих десятичных цифр. Не все числа между границами представимы.
	// У обоих типов также есть +Inf, -Inf и NaN. Есть +0 и -0.
	var temperature float64 = 36.6
	fmt.Printf("float32: %v, %T\n", temperature32, temperature32)
	fmt.Printf("float64: %v, %T\n", temperature, temperature)

	// Комплексные числа: действительная и мнимая части.
	// complex64: всего 64 бита, обе части имеют тип float32 и его диапазон/точность.
	// complex128: всего 128 бит, обе части имеют тип float64 и его диапазон/точность.
	// Общего линейного диапазона и операций < или > для комплексного числа нет.
	var wave64 complex64 = 1 + 2i
	var wave128 complex128 = 1 + 2i
	fmt.Printf("complex64: %v, %T\n", wave64, wave64)
	fmt.Printf("complex128: %v, %T\n", wave128, wave128)

	// bool: только false и true. Нулевое значение: false.
	var paid bool = true
	fmt.Printf("bool: %v, %T\n", paid, paid)

	// string: неизменяемая последовательность байтов. Нулевое значение: "".
	// Числового диапазона нет. Длина в байтах помещается в int, реальный размер ограничен памятью.
	// Сама строка не обязана содержать корректный UTF-8.
	var service string = "order-service"
	fmt.Printf("string: %q, %T, байт=%d\n", service, service, len(service))

	// Псевдонимы уже показанных целых типов, а не отдельные новые типы.
	var firstByte byte = 'G' // Псевдоним uint8: от 0 до 255.
	var letter rune = 'Ж'    // Псевдоним int32: от -2147483648 до 2147483647.
	// Для кодовой точки Unicode используют 0..0x10FFFF, исключая 0xD800..0xDFFF (суррогаты).
	// Тип rune сам не проверяет, является ли значение допустимой кодовой точкой.
	fmt.Printf("byte: %v, %T, %c\n", firstByte, firstByte, firstByte)
	fmt.Printf("rune: %v, %T, %c\n", letter, letter, letter)

	fmt.Println("целочисленное деление:", 5/2)
	fmt.Println("дробное деление:", float64(5)/float64(2))
	fmt.Println("остаток:", 5%2)
	//count = "три" // Ошибка: тип переменной уже определён.
	//fmt.Println(count + total) // Ошибка: int и int64 различаются.
}
