package factSender

import (
	"KPI_Drive_test/internal/entity"
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
)

const (
	apiURL      = "https://development.kpi-drive.ru/_api/facts/save_fact"
	bearerToken = "48ab34464a5573519725deb5865cc74c"
)

func SendFact(fact entity.Fact) error {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Запись полей в формат form-data
	writer.WriteField("period_start", fact.PeriodStart)
	writer.WriteField("period_end", fact.PeriodEnd)
	writer.WriteField("period_key", fact.PeriodKey)
	writer.WriteField("indicator_to_mo_id", strconv.Itoa(fact.IndicatorToMoID))
	writer.WriteField("indicator_to_mo_fact_id", strconv.Itoa(fact.IndicatorToMoFactID))
	writer.WriteField("value", strconv.Itoa(fact.Value))
	writer.WriteField("fact_time", fact.FactTime)
	writer.WriteField("is_plan", strconv.Itoa(fact.IsPlan))
	writer.WriteField("auth_user_id", strconv.Itoa(fact.AuthUserID))
	writer.WriteField("comment", fact.Comment)

	// Закрываем writer, чтобы завершить формирование формата
	err := writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, &buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to save fact, status code: %d", resp.StatusCode)
	}

	return nil
}
