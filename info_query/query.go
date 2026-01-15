package infoquery

import (
	"context"
	"fmt"
	"os"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	config "query_tools.com/config"
)

func Querys() {
	// Obtener el cliente desde el paquete config
	client, err := config.GetPrometheusClient()
	if err != nil {
		fmt.Printf("Error creando el cliente: %v\n", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := "up"
	result, warnings, err := v1api.Query(ctx, query, time.Now())

	if err != nil {
		fmt.Printf("Error al ejecutar la consulta: %v\n", err)
		os.Exit(1)
	}

	if len(warnings) > 0 {
		fmt.Printf("Advertencias: %v\n", warnings)
	}

	switch v := result.(type) {
	case model.Vector:
		for _, sample := range v {
			fmt.Printf("Instancia: %s | Estado: %s\n",
				sample.Metric["instance"], sample.Value)
		}
	default:
		fmt.Printf("Tipo de resultado inesperado: %T\n", v)
	}
}

func QueryAnalytic(label string, banco string, subproducto string, contraparte string) {
	client, err := config.GetPrometheusClient()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Construcción dinámica del query
	query := fmt.Sprintf(`sum(`+label+`{banco="%s", fecha="2026-01-15", contraparte="%s"})`,
		banco, contraparte)

	result, warnings, err := v1api.Query(ctx, query, time.Now())

	if err != nil {
		fmt.Printf("Error al ejecutar la consulta: %v\n", err)
		return
	}

	if len(warnings) > 0 {
		fmt.Printf("Advertencias: %v\n", warnings)
	}

	if v, ok := result.(model.Vector); ok {
		for _, sample := range v {
			fmt.Printf("Label: "+label+"|"+"subproducto: "+contraparte+" Total de registros: %s\n", sample.Value)
		}
	}
}

func QueryMS(label string, bco_emi string, bco_rec string) {
	client, err := config.GetPrometheusClient()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := fmt.Sprintf(`sum(`+label+`{bco_emi="%s", bco_rec="%s"})`,
		bco_emi, bco_rec)

	result, warnings, err := v1api.Query(ctx, query, time.Now())

	if err != nil {
		fmt.Printf("Error al ejecutar la consulta: %v\n", err)
		return
	}

	if len(warnings) > 0 {
		fmt.Printf("Advertencias: %v\n", warnings)
	}

	if v, ok := result.(model.Vector); ok {
		for _, sample := range v {
			fmt.Printf("Label: "+label+"|"+"contraparte: "+bco_rec+" Total de registros: %s\n", sample.Value)
		}
	}
}
