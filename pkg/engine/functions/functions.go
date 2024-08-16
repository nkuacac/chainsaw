package functions

import (
	"fmt"
	"reflect"

	"github.com/jmespath-community/go-jmespath/pkg/functions"
)

var (
	// stable functions
	env       = stable("env")
	trimSpace = stable("trim_space")
	asString  = stable("as_string")
	// experimental functions
	k8sGet            = experimental("k8s_get")
	k8sList           = experimental("k8s_list")
	k8sWait           = experimental("k8s_wait")
	k8sExists         = experimental("k8s_exists")
	k8sResourceExists = experimental("k8s_resource_exists")
	k8sServerVersion  = experimental("k8s_server_version")
	metricsDecode     = experimental("metrics_decode")

	allDataClusterInformerInit    = experimental("data_cluster_init")
	allDataClusterInformerClean   = experimental("data_cluster_clean")
	allDataClusterList            = experimental("data_cluster_list")
	allDataClusterWait            = experimental("data_cluster_wait")
	allDataClusterServerVersion   = experimental("data_cluster_server_version")
	allDataClusterCreateNamespace = experimental("data_cluster_create_namespace")
	allDataClusterDeleteNamespace = experimental("data_cluster_delete_namespace")
	getDataKubernetesClient       = experimental("data_k8s_client")
	dataKubernetesGet             = experimental("data_k8s_get")
	dataKubernetesList            = experimental("data_k8s_list")
)

func GetFunctions() []functions.FunctionEntry {
	return []functions.FunctionEntry{{
		Name: env,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler: jpEnv,
	}, {
		Name: k8sGet,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler: jpKubernetesGet,
	}, {
		Name: k8sList,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}, Optional: true},
		},
		Handler: jpKubernetesList,
	}, {
		Name: k8sExists,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler: jpKubernetesExists,
	}, {
		Name: k8sResourceExists,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler: jpKubernetesResourceExists,
	}, {
		Name: k8sServerVersion,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
		},
		Handler: jpKubernetesServerVersion,
	}, {
		Name: metricsDecode,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler: jpMetricsDecode,
	}, {
		Name: trimSpace,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler: jpTrimSpace,
	}, {
		Name: asString,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
		},
		Handler: func(arguments []any) (any, error) {
			in, err := getArgAt(arguments, 0)
			if err != nil {
				return nil, err
			}
			if in != nil {
				if in, ok := in.(string); ok {
					return in, nil
				}
				if reflect.ValueOf(in).Kind() == reflect.String {
					return fmt.Sprint(in), nil
				}
			}
			return nil, nil
		},
	}, {
		Name: allDataClusterInformerInit,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}, Optional: true},
		},
		Handler: jpAllDataClusterInformerInit,
	}, {
		Name: allDataClusterInformerClean,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}, Optional: true},
		},
		Handler: jpAllDataClusterInformerCleanup,
	}, {
		Name: allDataClusterList,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}, Optional: true},
			{Types: []functions.JpType{functions.JpString}, Optional: true},
		},
		Handler: jpAllDataClusterList,
	}, {
		Name: allDataClusterServerVersion,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler: jpAllDataClusterServerVersion,
	}, {
		Name: allDataClusterCreateNamespace,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler: jpAllDataClusterCreateNamespace,
	}, {
		Name: allDataClusterDeleteNamespace,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler: jpAllDataClusterDeleteNamespace,
	}, {
		Name: getDataKubernetesClient,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler: jpGetDataClusterClient,
	}, {
		Name: dataKubernetesGet,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler: jpDataKubernetesGet,
	}, {
		Name: dataKubernetesList,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}, Optional: true},
			{Types: []functions.JpType{functions.JpString}, Optional: true},
		},
		Handler: jpDataKubernetesList,
	}, {
		Name: "table_print",
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}, Optional: true},
		},
		Handler: tableFormat,
	}}
}

func GetInnerFunc() []functions.FunctionEntry {
	return []functions.FunctionEntry{{
		Name: k8sWait,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}, Optional: true},
			{Types: []functions.JpType{functions.JpString}, Optional: true},
			{Types: []functions.JpType{functions.JpString}, Optional: true},
			{Types: []functions.JpType{functions.JpString}, Optional: true},
			{Types: []functions.JpType{functions.JpString}, Optional: true},
		},
		Handler: jpKubernetesWait,
	}, {
		Name: allDataClusterWait,
		Arguments: []functions.ArgSpec{
			{Types: []functions.JpType{functions.JpAny}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
			{Types: []functions.JpType{functions.JpString}},
		},
		Handler: jpAllDataClusterWait,
	}}
}
