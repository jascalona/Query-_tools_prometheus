package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	infoquery "query_tools.com/info_query"
)

type SummaryStats struct {
	Total     int
	Presentes int
	Ausentes  int
	Errores   int
}

type FailedQuery struct {
	Category  string
	Component string
	RawQuery  string
	Status    string
	ResultVal string // 📊 Nuevo: Almacena el valor o detalle devuelto por Prometheus
}

func main() {
	filePath := flag.String("file", "QA-V7-Observability-Metrics.xlsx", "Ruta al archivo de queries (.xlsx)")
	promURL := flag.String("url", "http://212.56.40.50:9090", "URL de Prometheus")
	flag.Parse()

	if !strings.HasSuffix(strings.ToLower(*filePath), ".xlsx") {
		fmt.Println("❌ Por favor, ingresa un archivo .xlsx válido.")
		os.Exit(1)
	}

	fmt.Printf("🔍 Leyendo filas del archivo Excel: %s...\n", *filePath)
	rows, err := infoquery.LoadExcelDetails(*filePath)
	if err != nil {
		fmt.Printf("❌ Error leyendo el archivo Excel: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📋 Total de expresiones PromQL a evaluar: %d\n", len(rows))
	fmt.Printf("🚀 Conectando a Prometheus y ejecutando validaciones...\n\n")

	compStats := make(map[string]*SummaryStats)
	catStats := make(map[string]*SummaryStats)
	var globalStats SummaryStats
	var failedQueries []FailedQuery

	for _, qr := range rows {
		// Recibimos estado, query limpia y el valor de Prometheus
		status, cleanQuery, promValue := infoquery.TestQueryStatus(*promURL, qr)

		compKey := qr.Category + "|" + qr.Component
		if _, exists := compStats[compKey]; !exists {
			compStats[compKey] = &SummaryStats{}
		}
		if _, exists := catStats[qr.Category]; !exists {
			catStats[qr.Category] = &SummaryStats{}
		}

		compStats[compKey].Total++
		catStats[qr.Category].Total++
		globalStats.Total++

		switch status {
		case "PRESENTE":
			compStats[compKey].Presentes++
			catStats[qr.Category].Presentes++
			globalStats.Presentes++
			fmt.Printf("✅ [PRESENTE] %s -> %s \n   💡 Valor: %s\n", qr.Component, cleanQuery, promValue)
		case "AUSENTE":
			compStats[compKey].Ausentes++
			catStats[qr.Category].Ausentes++
			globalStats.Ausentes++
			failedQueries = append(failedQueries, FailedQuery{
				Category:  qr.Category,
				Component: qr.Component,
				RawQuery:  cleanQuery,
				Status:    "AUSENTE",
				ResultVal: promValue,
			})
			fmt.Printf("⚠️ [AUSENTE] %s -> %s \n   💡 Valor: %s\n", qr.Component, cleanQuery, promValue)
		default:
			compStats[compKey].Errores++
			catStats[qr.Category].Errores++
			globalStats.Errores++
			failedQueries = append(failedQueries, FailedQuery{
				Category:  qr.Category,
				Component: qr.Component,
				RawQuery:  cleanQuery,
				Status:    "ERROR",
				ResultVal: promValue,
			})
			fmt.Printf("❌ [ERROR] %s -> %s \n   💡 Valor: %s\n", qr.Component, cleanQuery, promValue)
		}
	}

	// 🌐 Generar HTML incluyendo la columna del valor devuelto por Prometheus
	var htmlBuilder strings.Builder

	htmlBuilder.WriteString(`<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <title>Reporte de Inventario - Métricas Prometheus</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; background-color: #f4f6f9; color: #333; margin: 0; padding: 20px; }
        .container { max-width: 1200px; margin: auto; background: #fff; padding: 30px; border-radius: 8px; box-shadow: 0 4px 12px rgba(0,0,0,0.05); }
        h1 { font-size: 24px; color: #1e293b; border-bottom: 2px solid #e2e8f0; padding-bottom: 10px; }
        h2 { font-size: 18px; color: #475569; margin-top: 30px; }
        table { width: 100%; border-collapse: collapse; margin-top: 15px; margin-bottom: 30px; font-size: 14px; }
        th, td { padding: 12px 15px; text-align: left; border-bottom: 1px solid #e2e8f0; }
        th { background-color: #f8fafc; color: #1e293b; font-weight: 600; text-transform: uppercase; font-size: 12px; letter-spacing: 0.5px; }
        tr:hover { background-color: #f1f5f9; }
        .badge { padding: 4px 8px; border-radius: 4px; font-weight: bold; font-size: 12px; }
        .badge-success { background-color: #dcfce7; color: #166534; }
        .badge-warning { background-color: #fef9c3; color: #854d0e; }
        .badge-danger { background-color: #fee2e2; color: #991b1b; }
        .total-row { font-weight: bold; background-color: #f8fafc; }
        code { background-color: #f1f5f9; padding: 2px 6px; border-radius: 4px; font-family: monospace; font-size: 13px; color: #d63384; word-break: break-all; }
        .alert-box { background-color: #fffbeb; border-left: 4px solid #f59e0b; padding: 15px; margin-bottom: 20px; border-radius: 4px; }
    </style>
</head>
<body>
<div class="container">
    <h1>📊 Inventario Completo - Métricas Prometheus</h1>
    <p><strong>Servidor:</strong> ` + *promURL + `</p>
    
    <h2>Resumen por Categoría</h2>
    <table>
        <thead>
            <tr>
                <th>Categoría</th>
                <th>Total</th>
                <th>Presentes</th>
                <th>Ausentes</th>
                <th>Errores</th>
                <th>% Presente</th>
            </tr>
        </thead>
        <tbody>`)

	for cat, stats := range catStats {
		pct := 0.0
		if stats.Total > 0 {
			pct = (float64(stats.Presentes) / float64(stats.Total)) * 100
		}
		htmlBuilder.WriteString(fmt.Sprintf(`<tr>
			<td><strong>%s</strong></td>
			<td>%d</td>
			<td><span class="badge badge-success">%d</span></td>
			<td><span class="badge badge-warning">%d</span></td>
			<td><span class="badge badge-danger">%d</span></td>
			<td><strong>%.2f%%</strong></td>
		</tr>`, cat, stats.Total, stats.Presentes, stats.Ausentes, stats.Errores, pct))
	}

	globalPct := 0.0
	if globalStats.Total > 0 {
		globalPct = (float64(globalStats.Presentes) / float64(globalStats.Total)) * 100
	}

	htmlBuilder.WriteString(fmt.Sprintf(`<tr class="total-row">
		<td>TOTAL</td>
		<td>%d</td>
		<td>%d</td>
		<td>%d</td>
		<td>%d</td>
		<td>%.2f%%</td>
	</tr>`, globalStats.Total, globalStats.Presentes, globalStats.Ausentes, globalStats.Errores, globalPct))

	htmlBuilder.WriteString(`</tbody></table>`)

	// 🛠️ Apartado de Incidencias con el valor devuelto por Prometheus
	htmlBuilder.WriteString(`
    <h2>🔍 Detalle de Consultas con Incidencias (Ausentes o con Errores)</h2>`)

	if len(failedQueries) == 0 {
		htmlBuilder.WriteString(`<p style="color: #166534; font-weight: bold;">✨ ¡Excelente! No se encontraron consultas ausentes ni con errores.</p>`)
	} else {
		htmlBuilder.WriteString(fmt.Sprintf(`<div class="alert-box">Se encontraron <strong>%d</strong> consultas con incidencias.</div>`, len(failedQueries)))
		htmlBuilder.WriteString(`
		<table>
			<thead>
				<tr>
					<th>Categoría</th>
					<th>Componente</th>
					<th>Estado</th>
					<th>Expresión PromQL Limpia</th>
					<th>Respuesta de Prometheus</th>
				</tr>
			</thead>
			<tbody>`)

		for _, fq := range failedQueries {
			badgeClass := "badge-warning"
			if fq.Status == "ERROR" {
				badgeClass = "badge-danger"
			}
			htmlBuilder.WriteString(fmt.Sprintf(`<tr>
				<td>%s</td>
				<td><strong>%s</strong></td>
				<td><span class="badge %s">%s</span></td>
				<td><code>%s</code></td>
				<td><small>%s</small></td>
			</tr>`, fq.Category, fq.Component, badgeClass, fq.Status, fq.RawQuery, fq.ResultVal))
		}
		htmlBuilder.WriteString(`</tbody></table>`)
	}

	// Tabla Detallada por Componente
	htmlBuilder.WriteString(`
    <h2>Resumen Detallado por Componente</h2>
    <table>
        <thead>
            <tr>
                <th>Categoría</th>
                <th>Componente</th>
                <th>Total</th>
                <th>Presentes</th>
                <th>Ausentes</th>
                <th>Errores</th>
                <th>% Presente</th>
            </tr>
        </thead>
        <tbody>`)

	for compKey, stats := range compStats {
		parts := strings.Split(compKey, "|")
		cat := parts[0]
		comp := parts[1]
		pct := 0.0
		if stats.Total > 0 {
			pct = (float64(stats.Presentes) / float64(stats.Total)) * 100
		}
		htmlBuilder.WriteString(fmt.Sprintf(`<tr>
			<td>%s</td>
			<td><strong>%s</strong></td>
			<td>%d</td>
			<td>%d</td>
			<td>%d</td>
			<td>%d</td>
			<td>%.2f%%</td>
		</tr>`, cat, comp, stats.Total, stats.Presentes, stats.Ausentes, stats.Errores, pct))
	}

	htmlBuilder.WriteString(`</tbody></table></div></body></html>`)

	outputFile := "reporte_metricas.html"
	err = os.WriteFile(outputFile, []byte(htmlBuilder.String()), 0644)
	if err != nil {
		fmt.Printf("❌ Error guardando el archivo HTML: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✨ ¡Reporte generado con éxito! Archivo guardado en: %s\n", outputFile)
}
