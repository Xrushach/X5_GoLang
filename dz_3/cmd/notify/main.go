package main

import (
	"errors"
	"fmt"

	"lesson5/student/practice"
)

func main() {
	// TODO: домашняя работа. Покажите успешную рассылку, пустой текст,
	// ошибку настройки и продолжение после BrokenNotifier.
	// Импорт учебного пакета: "lesson5/student/practice".
	// fmt.Println("Добавьте сценарии из homework/README.md")
	fmt.Println("Случай: Успех")

	email := &practice.EmailNotifier{
		BaseNotifier: practice.BaseNotifier{Channel: "email"},
		Address:      "sobaka@soap.go",
	}
	push := &practice.PushNotifier{
		BaseNotifier: practice.BaseNotifier{Channel: "push"},
		DeviceID:     "ID4252",
	}
	err := practice.Dispatch([]practice.Notifier{email, push}, "Заказ готов")

	if err == nil {
		fmt.Println("No error")
	} else {
		fmt.Println("Error", err)
	}

	fmt.Println("История email: ", email.Messages)
	fmt.Println("История push: ", push.Messages)

	fmt.Println("Случай: Пустой текст")
	oldEmail, oldPush := len(email.Messages), len(push.Messages)
	err = practice.Dispatch([]practice.Notifier{email, push}, "")
	if errors.Is(err, practice.ErrEmptyMessage) {
		fmt.Println("ErrEmptyMessage распознано")
	} else {
		fmt.Println("ErrEmptyMessage не найден, ошибка", err)
	}
	fmt.Printf("Длина истории email: %d --> %d\n", oldEmail, len(email.Messages))
	fmt.Printf("Длина истории push: %d --> %d\n", oldPush, len(push.Messages))

	fmt.Println("Случай: Нет настройки")

	noID := &practice.PushNotifier{BaseNotifier: practice.BaseNotifier{Channel: "push"}}
	err = practice.Dispatch([]practice.Notifier{noID}, "Готов")
	var sendErr *practice.SendError
	if errors.As(err, &sendErr) {
		fmt.Println("Получен *SendError")
		fmt.Println("Channel:", sendErr.Channel)
		fmt.Println("Reason:", sendErr.Reason)
	} else {
		fmt.Println("*SendError не получен, ошибка:", err)
	}
	fmt.Println("История push:", noID.Messages)

	fmt.Println("Случай: Сломанная реализация")
	bad := &practice.BrokenNotifier{BaseNotifier: practice.BaseNotifier{Channel: "broken"}}
	goodPush := &practice.PushNotifier{BaseNotifier: practice.BaseNotifier{Channel: "push"}, DeviceID: "007"}
	errs := practice.SafeDispatch([]practice.Notifier{bad, goodPush}, "Готов")
	fmt.Println("кол-во ошибок", len(errs))
	if errors.Is(errors.Join(errs...), practice.ErrNotifierPanic) {
		fmt.Println("ErrNotifierPanic распознано")
	} else {
		fmt.Println("ErrNotifierPanic не получен", errs)
	}
	fmt.Println("История рабочего push:", goodPush.Messages)
}
