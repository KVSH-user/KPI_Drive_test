package bufer

import (
	"KPI_Drive_test/internal/entity"
	resp "KPI_Drive_test/internal/lib/api/response"
	"KPI_Drive_test/internal/stan"
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

func GetFact(log *slog.Logger, client *stan.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.order.GetFact"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
			log.Error("failed to parse multipart form: ", err)

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to parse form data"))

			return
		}

		// Получаем данные из form-data
		var req entity.Fact
		req.PeriodStart = r.FormValue("period_start")
		req.PeriodEnd = r.FormValue("period_end")
		req.PeriodKey = r.FormValue("period_key")
		req.IndicatorToMoID, _ = strconv.Atoi(r.FormValue("indicator_to_mo_id"))
		req.IndicatorToMoFactID, _ = strconv.Atoi(r.FormValue("indicator_to_mo_fact_id"))
		req.Value, _ = strconv.Atoi(r.FormValue("value"))
		req.FactTime = r.FormValue("fact_time")
		req.IsPlan, _ = strconv.Atoi(r.FormValue("is_plan"))
		req.AuthUserID, _ = strconv.Atoi(r.FormValue("auth_user_id"))
		req.Comment = r.FormValue("comment")

		// Проверка на наличие необходимых данных
		if req.PeriodStart == "" || req.PeriodEnd == "" || req.PeriodKey == "" || req.FactTime == "" {
			log.Error("missing required fields")

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("missing required fields"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// Отправляем данные в NATS
		data, err := json.Marshal(req)
		if err != nil {
			log.Error("Error marshaling fact:", err)

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to marshal request"))

			return
		}

		if err := client.Sc.Publish("facts", data); err != nil {
			log.Error("Error publishing to NATS:", err)

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to publish to NATS"))

			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "OK")
	}
}
