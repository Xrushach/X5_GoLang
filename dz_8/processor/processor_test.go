package processor

import (
	"path/filepath"
	"testing"
)

func TestAnalyzeFile(t *testing.T) {
	path := filepath.Join("..", "testdata", "orders-a.txt")
	got := analyzeFile(path)
	if got.Err != nil || got.Path != path || got.Lines != 3 || got.Words != 5 {
		t.Fatalf("результат = %+v, ожидаем этот путь, 3 строки и 5 слов", got)
	}
}

func TestProcessFiles(t *testing.T) {
	paths := []string{
		filepath.Join("..", "testdata", "orders-a.txt"),
		filepath.Join("..", "testdata", "orders-b.txt"),
		filepath.Join("..", "testdata", "empty.txt"),
	}
	want := map[string][2]int{
		paths[0]: {3, 5},
		paths[1]: {2, 4},
		paths[2]: {0, 0},
	}

	for got := range ProcessFiles(paths, 2) {
		stats, found := want[got.Path]
		if !found {
			t.Fatalf("лишний результат для %q", got.Path)
		}
		if got.Err != nil || got.Lines != stats[0] || got.Words != stats[1] {
			t.Fatalf("для %q получили %+v, ожидаем строк и слов %v", got.Path, got, stats)
		}
		delete(want, got.Path)
	}
	if len(want) > 0 {
		t.Fatalf("не получены результаты: %v", want)
	}
}

func TestProcessFilesWithMissing(t *testing.T) {
	okPath := filepath.Join("..", "testdata", "orders-a.txt")
	badPath := filepath.Join("..", "testdata", "missing.txt")

	var gotOK, gotErr bool
	count := 0
	for got := range ProcessFiles([]string{okPath, badPath}, 2) {
		count++
		if got.Path == okPath {
			if got.Err != nil || got.Lines != 3 || got.Words != 5 {
				t.Fatalf("для %q получили %+v", got.Path, got)
			}
			gotOK = true
		} else if got.Path == badPath {
			if got.Err == nil {
				t.Fatalf("для %q ожидали ошибку", got.Path)
			}
			gotErr = true
		} else {
			t.Fatalf("лишний результат для %q", got.Path)
		}
	}
	if count != 2 || !gotOK || !gotErr {
		t.Fatalf("count=%d gotOK=%v gotErr=%v", count, gotOK, gotErr)
	}
}

func TestProcessFilesEmpty(t *testing.T) {
	count := 0
	for range ProcessFiles(nil, 2) {
		count++
	}
	if count != 0 {
		t.Fatalf("получили %d результатов, ожидали 0", count)
	}
}

func TestProcessFilesDuplicateAndZeroWorkers(t *testing.T) {
	path := filepath.Join("..", "testdata", "orders-a.txt")
	count := 0
	for got := range ProcessFiles([]string{path, path}, 0) {
		if got.Err != nil || got.Path != path || got.Lines != 3 || got.Words != 5 {
			t.Fatalf("неожиданный результат: %+v", got)
		}
		count++
	}
	if count != 2 {
		t.Fatalf("получили %d результатов, ожидали 2", count)
	}
}
