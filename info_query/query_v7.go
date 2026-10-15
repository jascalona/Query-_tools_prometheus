package infoquery

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/xuri/excelize/v2"
)

type MetricConfig struct {
	PrometheusURL string   `json:"prometheus_url"`
	Queries       []string `json:"queries"`
}

type PrometheusResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []any  `json:"result"`
	} `json:"data"`
}

type QueryRow struct {
	ID        string
	Component string
	Category  string
	Query     string
}

func LoadExcelDetails(filePath string) ([]QueryRow, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheetName := f.GetSheetList()[0]
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	var queryRows []QueryRow
	for i, row := range rows {
		if i == 0 || len(row) <= 5 {
			continue
		}

		id := strings.TrimSpace(row[0])
		component := strings.TrimSpace(row[1])
		category := strings.TrimSpace(row[2])
		query := strings.TrimSpace(row[5])

		if query != "" {
			if category == "" {
				category = "General"
			}
			if component == "" {
				component = "Desconocido"
			}
			queryRows = append(queryRows, QueryRow{
				ID:        id,
				Component: component,
				Category:  category,
				Query:     query,
			})
		}
	}

	return queryRows, nil
}

// TestQueryStatus evalúa la query y retorna el estado, la query limpia y el valor obtenido de Prometheus
func TestQueryStatus(prometheusURL string, qr QueryRow) (string, string, string) {
	cleanQuery := qr.Query

	// Reemplazos universales para eliminar cualquier variable con $
	cleanQuery = strings.ReplaceAll(cleanQuery, `=~"$feliq"`, `=~".*"`)
	cleanQuery = strings.ReplaceAll(cleanQuery, `=~"$subproduct"`, `=~".*"`)
	cleanQuery = strings.ReplaceAll(cleanQuery, `=~"$coliq"`, `=~".*"`)
	cleanQuery = strings.ReplaceAll(cleanQuery, `=~"$bank"`, `=~".*"`)
	cleanQuery = strings.ReplaceAll(cleanQuery, `=~"$instance"`, `=~".*"`)
	cleanQuery = strings.ReplaceAll(cleanQuery, `=~"$container"`, `=~".*"`)
	cleanQuery = strings.ReplaceAll(cleanQuery, `=~"$node"`, `=~".*"`)
	cleanQuery = strings.ReplaceAll(cleanQuery, `=~"$job"`, `=~".*"`)
	cleanQuery = strings.ReplaceAll(cleanQuery, `=~'$instance'`, `=~'.*'`)
	cleanQuery = strings.ReplaceAll(cleanQuery, `="$instance"`, `=".*"`)
	cleanQuery = strings.ReplaceAll(cleanQuery, `="$job"`, `=".*"`)

	reGeneric := regexp.MustCompile(`(=~?)\s*(['"][^'"]*\$[a-zA-Z0-9_]+[^'"]*['"])`)
	cleanQuery = reGeneric.ReplaceAllString(cleanQuery, `$1".*"`)

	apiEndpoint := fmt.Sprintf("%s/api/v1/query?query=%s", prometheusURL, url.QueryEscape(cleanQuery))

	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return "ERROR", cleanQuery, "Error creando request"
	}

	req.SetBasicAuth("admin", "$Yc0m-2023")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "ERROR", cleanQuery, "Error de red"
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "ERROR", cleanQuery, "Error leyendo body"
	}

	// Si el servidor responde con texto plano (como "Unauthorized"), capturarlo aquí
	bodyString := string(bodyBytes)
	if strings.Contains(bodyString, "Unauthorized") || resp.StatusCode == 401 {
		return "ERROR", cleanQuery, "HTTP 401: Unauthorized (Faltan credenciales)"
	}

	var promResp PrometheusResponse
	if err := json.Unmarshal(bodyBytes, &promResp); err != nil {
		return "ERROR", cleanQuery, "Error decodificando JSON"
	}

	if promResp.Status != "success" {
		return "ERROR", cleanQuery, "Prometheus status != success"
	}

	if len(promResp.Data.Result) == 0 {
		return "AUSENTE", cleanQuery, "Sin series activas (0 resultados)"
	}

	// Extraer el valor del resultado si viene presente
	resultValStr := fmt.Sprintf("%v", promResp.Data.Result)

	return "PRESENTE", cleanQuery, resultValStr
}
