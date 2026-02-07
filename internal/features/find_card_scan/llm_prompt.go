package find_card_scan

import "fmt"

func BuildPrompt(message string) string {
	return fmt.Sprintf(`Ты - специалист, который отлично умеет вытаскивать названия карт таро из описания.

У меня есть следующее описание: "%s"

Как называются карты (их может несколько)?

ВАЖНО! ВЫВОДИ ОТВЕТ ТОЛЬКО В СЛЕДУЮЩЕМ ФОРМАТЕ:
{
 "titles":[
    "title_1",
    "title_2",
    ...,
    "title_n"
 ]
}`, message)
}
