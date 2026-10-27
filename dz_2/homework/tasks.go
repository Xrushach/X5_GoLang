package homework

// PairInput соответствует объекту nums/target из условия контеста.
type PairInput struct {
	Nums   []int
	Target int
}

// longestPrefix возвращает самый длинный общий префикс строк.
func longestPrefix(words []string) string {
	// TODO: реализовать поиск префикса по кодовым точкам Unicode.
	if len(words) == 0 {
		return ""
	}
	start := []rune(words[0])
	for i, cur := range start {
		for _, o := range words[1:] {
			other := []rune(o)
			if len(other) <= i || other[i] != cur {
				return string(start[:i])
			}
		}
	}
	return words[0]
}

// pairIndices возвращает индексы двух разных элементов по возрастанию.
func pairIndices(input PairInput) [2]int {
	// TODO: найти единственную пару с суммой input.Target.
	m := make(map[int]int)
	for i := 0; i < len(input.Nums); i++ {
		comp := input.Target - input.Nums[i]
		ind, ok := m[comp]
		if ok {
			return [2]int{ind, i}
		}
		m[input.Nums[i]] = i
	}
	return [2]int{}
}

// hasDuplicate сообщает, встречается ли какое-либо значение повторно.
func hasDuplicate(nums []int) bool {
	// TODO: проверить наличие повторов.
	m := make(map[int]int)
	for i := 0; i < len(nums); i++ {
		_, ok := m[nums[i]]
		if ok {
			return true
		}
		m[nums[i]] = i
	}
	return false
}

// compactSorted создаёт независимый слайс уникальных значений.
func compactSorted(nums []int) []int {
	// TODO: сохранить порядок и не изменять исходный слайс.
	if len(nums) == 0 {
		return []int{}
	}
	res := make([]int, 0, len(nums))
	res = append(res, nums[0])
	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[i-1] {
			res = append(res, nums[i])
		}
	}
	return res
}

// balancedBrackets проверяет пары и порядок закрытия скобок.
func balancedBrackets(text string) bool {
	// TODO: обработать строку из символов ()[]{}.
	runes := []rune(text)
	arr := make([]rune, 0, len(runes))
	for _, c := range runes {
		switch c {
		case '(', '[', '{':
			arr = append(arr, c)
		case ')', ']', '}':
			if len(arr) == 0 {
				return false
			}
			last := arr[len(arr)-1]
			if (c == ')' && last != '(') || (c == ']' && last != '[') || (c == '}' && last != '{') {
				return false
			}
			arr = arr[:len(arr)-1]
		}
	}
	return len(arr) == 0
}
