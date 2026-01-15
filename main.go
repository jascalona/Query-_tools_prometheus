package main

import (
	"fmt"

	infoquery "query_tools.com/info_query"
)

func main() {
	//infoquery.Querys()

	//label := "analytic_metrics_amount"
	banco := "0001"

	infoquery.QueryAnalytic("analytic_metrics_amount", banco, "002", "0108")
	infoquery.QueryAnalytic("analytic_metrics_core", banco, "220", "0102")

	infoquery.QueryAnalytic("analytic_metrics", banco, "220", "0134")
	infoquery.QueryAnalytic("analytic_metrics", banco, "222", "0105")

	/*
		infoquery.QueryAnalytic("analytic_metrics", banco, "221")
		infoquery.QueryAnalytic("analytic_metrics_core_amount", banco, "221")
		infoquery.QueryAnalytic("analytic_metrics_core_amount", banco, "220")
		infoquery.QueryAnalytic("analytic_metrics_core_amount", banco, "222")
		infoquery.QueryAnalytic("analytic_metrics_core_amount", banco, "002")
		infoquery.QueryAnalytic("analytic_metrics_lbtr_core_amount", banco, "998")
		infoquery.QueryAnalytic("analytic_metrics_lbtr_core_amount", banco, "999")
		infoquery.QueryAnalytic("analytic_metrics_lbtr_core", banco, "998")
		infoquery.QueryAnalytic("analytic_metrics_lbtr_core", banco, "999")
		infoquery.QueryAnalytic("analytic_metrics_lbtr_amount", banco, "998")
		infoquery.QueryAnalytic("analytic_metrics_lbtr_amount", banco, "999")
		infoquery.QueryAnalytic("analytic_metrics_lbtr", banco, "998")
		infoquery.QueryAnalytic("analytic_metrics_lbtr", banco, "999")
	*/
	fmt.Println("--LABELS_MS--")

	infoquery.QueryMS("ms_tiempo_duracion_proceso_sum", banco, "0102")
	infoquery.QueryMS("ms_tiempo_duracion_proceso_count", banco, "0108")
	infoquery.QueryMS("ms_tiempo_duracion_proceso_bucket", banco, "0108")

}
