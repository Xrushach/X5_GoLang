package processor

// ProcessFiles возвращает по результату на каждый переданный путь.
// Получатель читает канал до закрытия.
func ProcessFiles(paths []string, workers int) <-chan Result {
	if workers < 1 {
		workers = 1
	}
	results := make(chan Result)
	// TODO: распределите пути между workers обработчиками через канал jobs.
	// Для каждого пути вызовите analyzeFile и отправьте результат в results.
	// Определите, кто закроет каждый канал и как для этого использовать WaitGroup.
	close(results) // Заглушка позволяет запустить весь проект до начала практики.
	return results
}
