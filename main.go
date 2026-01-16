package main

import (
	"fmt"

	infoquery "query_tools.com/info_query"
)

func main() {
	//infoquery.Querys()

	//label := "analytic_metrics_amount"
	banco := "0001"

	fmt.Println("\n--Analityc Metric--")
	infoquery.QueryAnalytic("analytic_metrics_amount", banco, "0108")
	infoquery.QueryAnalytic("analytic_metrics_core", banco, "0137")
	infoquery.QueryAnalytic("analytic_metrics_core_amount", banco, "0166")

	fmt.Println("\n--Analityc Metric LBTR--")
	infoquery.QueryAnalytic("analytic_metrics_lbtr_core_amount", banco, "0001")
	infoquery.QueryAnalytic("analytic_metrics_lbtr_core", banco, "0001")
	infoquery.QueryAnalytic("analytic_metrics_lbtr_amount", banco, "0001")
	infoquery.QueryAnalytic("analytic_metrics_lbtr", banco, "0001")

	fmt.Println("\n--LABELS MS--")
	infoquery.QueryMS("ms_tiempo_duracion_proceso_sum", banco, "0102")
	infoquery.QueryMS("ms_tiempo_duracion_proceso_count", banco, "0108")
	infoquery.QueryMS("ms_tiempo_duracion_proceso_bucket", banco, "0108")

}
