package main

import (
	"KPI_Drive_test/internal/entity"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
)

func main() {
	sc, err := stan.Connect("test-cluster", "publisher-client", stan.NatsURL("nats://localhost:4223"))
	if err != nil {
		log.Fatal("cant connect to NATS: ", err)
	}
	defer sc.Close()

	// пример добавления 10 записей в NATS
	for i := 0; i < 10; i++ {
		// тестовые данные
		fact := entity.Fact{
			PeriodStart:         "2024-05-01",
			PeriodEnd:           "2024-05-31",
			PeriodKey:           "month",
			IndicatorToMoID:     227373,
			IndicatorToMoFactID: 0,
			Value:               1,
			FactTime:            "2024-05-31",
			IsPlan:              0,
			AuthUserID:          40,
			Comment:             "buffer KVSH-user",
		}

		// Отправляем данные в NATS
		data, err := json.Marshal(fact)
		if err != nil {
			log.Println("Error marshaling fact:", err)
			return
		}

		if err := sc.Publish("facts", data); err != nil {
			log.Println("Error publishing to NATS:", err)
		}
	}
}
