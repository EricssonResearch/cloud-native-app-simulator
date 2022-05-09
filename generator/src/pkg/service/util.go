/*
Copyright 2021 Telefonaktiebolaget LM Ericsson AB

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package service

import (
	"application-generator/src/pkg/model"
	"fmt"
	"strconv"
)

const (
	VolumeName = "config-data-volume"
	VolumePath = "/usr/src/app/config"

	ImageName = "app"
	ImageURL  = "app-demo:latest"

	DefaultExtPort  = 80
	DefaultPort     = 5000
	defaultProtocol = "http"

	Uri = "/"

	ReplicaNumber = 1

	RequestsCPUDefault    = "500m"
	RequestsMemoryDefault = "256M"
	LimitsCPUDefault      = "1000m"
	LimitsMemoryDefault   = "1024M"

	SvcNamePrefix            = "service"
	SvcProcessesDefault      = 1
	SvcThreadsDefault        = 1
	SvcReadinessProbeDefault = 2

	EpNamePrefix            = "end"
	EpExecModeDefault       = "sequential"
	EpNwResponseSizeDefault = 512

	EpExecTimeDefault = "0.1s"
	EpMethodDefault   = "all"
	EpWorkersDefault  = 1
	EpLoadDefault     = "5%"

	EpNwForwardRequests = "asynchronous"

	CsTrafficForwardRatio = 1
	CsRequestSizeDefault  = 256
)

func CreateDeployment(metadataName, selectorAppName, selectorClusterName string, numberOfReplicas int,
	templateAppLabel, templateClusterLabel, namespace string, containerPort int, containerName, containerImage,
	mountPath string, volumeName, configMapName string, readinessProbe int, requestCPU, requestMemory, limitCPU,
	limitMemory, nodeAffinity, protocol string, annotations []model.Annotation) (deploymentInstance model.DeploymentInstance) {

	var deployment model.DeploymentInstance
	var containerInstance model.ContainerInstance
	var envInstance model.EnvInstance
	var containerPortInstance model.ContainerPortInstance
	var containerVolume model.ContainerVolumeInstance
	var volumeInstance model.VolumeInstance

	envInstance.Name = "SERVICE_NAME"
	envInstance.Value = metadataName
	containerPortInstance.ContainerPort = containerPort
	volumeInstance.Name = volumeName
	volumeInstance.ConfigMap.Name = configMapName

	containerVolume.MountName = volumeName
	containerVolume.MountPath = mountPath

	containerInstance.Volumes = append(containerInstance.Volumes, containerVolume)
	containerInstance.Ports = append(containerInstance.Ports, containerPortInstance)
	containerInstance.Name = containerName
	containerInstance.Image = containerImage
	containerInstance.ImagePullPolicy = "Never"
	containerInstance.Env = append(containerInstance.Env, envInstance)
	if protocol == "http" {
		containerInstance.ReadinessProbe.HttpGet.Path = "/"
		containerInstance.ReadinessProbe.HttpGet.Port = containerPort
	}
	if protocol == "grpc" {
		containerInstance.ReadinessProbe.Exec.Command = append(containerInstance.ReadinessProbe.Exec.Command, ("/bin/grpc_health_probe"), "-addr=:"+strconv.Itoa(containerPort))

	}

	containerInstance.ReadinessProbe.InitialDelaySeconds = readinessProbe
	containerInstance.ReadinessProbe.PeriodSeconds = 1
	containerInstance.Resources.ResourceRequests.Cpu = requestCPU
	containerInstance.Resources.ResourceRequests.Memory = requestMemory
	containerInstance.Resources.ResourceLimits.Cpu = limitCPU
	containerInstance.Resources.ResourceLimits.Memory = limitMemory

	deployment.APIVersion = "apps/v1"
	deployment.Kind = "Deployment"
	deployment.Metadata.Name = metadataName
	deployment.Metadata.Namespace = namespace
	deployment.Metadata.Labels.Cluster = templateClusterLabel
	deployment.Spec.Selector.MatchLabels.App = selectorAppName
	deployment.Spec.Selector.MatchLabels.Cluster = selectorClusterName
	deployment.Spec.Replicas = numberOfReplicas
	deployment.Spec.Template.Metadata.Labels.App = templateAppLabel
	deployment.Spec.Template.Metadata.Labels.Cluster = templateClusterLabel
	if len(annotations) > 0 {
		deployment.Spec.Template.Metadata.Annotations = map[string]string{}
		for i := 0; i < len(annotations); i++ {
			deployment.Spec.Template.Metadata.Annotations[annotations[i].Name] = annotations[i].Value
		}
	}

	deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, containerInstance)
	deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, volumeInstance)
	deployment.Spec.Template.Spec.NodeName = nodeAffinity

	return deployment

}

func CreateWorkerDeployment(metadataName, selectorName string, numberOfReplicas int, templateLabel string,
	containerName, containerImage, mountPath string, volumeName, configMapName string) (deploymentInstance model.DeploymentInstance) {

	var deployment model.DeploymentInstance
	var containerInstance model.ContainerInstance
	var containerVolume model.ContainerVolumeInstance
	var volumeInstance model.VolumeInstance

	volumeInstance.Name = volumeName
	volumeInstance.ConfigMap.Name = configMapName

	containerVolume.MountName = volumeName
	containerVolume.MountPath = mountPath

	containerInstance.Volumes = append(containerInstance.Volumes, containerVolume)
	containerInstance.Name = containerName
	containerInstance.Image = containerImage
	containerInstance.ImagePullPolicy = "Never"

	deployment.APIVersion = "apps/v1"
	deployment.Kind = "Deployment"
	deployment.Metadata.Name = metadataName
	deployment.Spec.Selector.MatchLabels.App = selectorName
	deployment.Spec.Replicas = numberOfReplicas
	deployment.Spec.Template.Metadata.Labels.App = templateLabel
	deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, containerInstance)
	deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, volumeInstance)

	return deployment
}

func CreateService(metadataName, selectorAppName, protocol, uri, metadataLabelCluster, namespace string, defaultExtPort, defaultPort int) (serviceInstance model.ServiceInstance) {
	const apiVersion = "v1"

	const apiKind = "Service"

	var port model.ServicePortInstance

	var service model.ServiceInstance

	annotations := map[string]string{
		protocol: uri,
	}

	port.Port = defaultExtPort
	port.TargetPort = defaultPort
	port.Name = protocol

	service.APIVersion = apiVersion
	service.Kind = apiKind
	service.Metadata.Name = metadataName
	service.Metadata.Namespace = namespace
	service.Metadata.Labels.Cluster = metadataLabelCluster
	service.Metadata.Annotations = annotations
	service.Spec.Selector.App = selectorAppName
	service.Spec.Ports = append(service.Spec.Ports, port)

	return service
}

func CreateServiceAccount(metadataName, accountName string) (serviceAccountInstance model.ServiceAccountInstance) {
	const apiVersion = "v1"
	const apiKind = "ServiceAccount"

	var serviceAccount model.ServiceAccountInstance

	serviceAccount.APIVersion = apiVersion
	serviceAccount.Kind = apiKind
	serviceAccount.Metadata.Name = metadataName
	serviceAccount.Metadata.Labels.Account = accountName

	return serviceAccount
}

func CreateConfig(metadataName, metadataLabelName, metadataLabelCluster, namespace, config, proto string) (configMapInstance model.ConfigMapInstance) {

	const apiVersion = "v1"

	const apiKind = "ConfigMap"

	var configMap model.ConfigMapInstance

	configMap.APIVersion = apiVersion
	configMap.Kind = apiKind
	configMap.Metadata.Name = metadataName
	configMap.Metadata.Labels.Cluster = metadataLabelCluster
	configMap.Metadata.Labels.Name = metadataLabelName
	configMap.Metadata.Namespace = namespace
	configMap.Data.Config = config
	configMap.Data.Service = proto

	return configMap
}

func CreateGateway(hosts []string) model.GatewayInstance {

	var server model.GatewayServers
	for i, s := range hosts {
		e := fmt.Sprintf("s%s.dev", s)
		hosts[i] = e
	}
	server.Hosts = hosts
	server.Port.Name = "http"
	server.Port.Number = 80
	server.Port.Protocol = "HTTP"

	gateway := &model.GatewayInstance{APIVersion: "networking.istio.io/v1alpha3",
		Kind: "Gateway", Metadata: model.Metadata{Name: "generator-gateway"},
		Spec: model.GatewaySpec{Selector: model.GatewaySelector{Istio: "ingressgateway"}}}

	gateway.Spec.Servers = append(gateway.Spec.Servers, server)

	return *gateway
}

func CreateVirtualService(metadataName, hostname, gatewayHost string, port int) model.VirtualServiceInstance {

	var match model.VirtualServiceMatch
	match.URI.Exact = "/"

	var route model.VirtualServiceRoute
	route.Destination.Host = hostname
	route.Destination.Port.Number = port
	var http model.VirtualServiceHTTP
	http.Match = append(http.Match, match)
	http.Route = append(http.Route, route)

	virtualService := &model.VirtualServiceInstance{
		APIVersion: "networking.istio.io/v1alpha3", Kind: "VirtualService", Metadata: model.Metadata{Name: metadataName}}
	virtualService.Spec.HTTP = append(virtualService.Spec.HTTP, http)
	virtualService.Spec.Gateways = append(virtualService.Spec.Gateways, "generator-gateway")
	virtualService.Spec.Hosts = append(virtualService.Spec.Hosts, gatewayHost)

	return *virtualService

}

func CreateFileConfig() model.FileConfig {

	var fileConfig model.FileConfig

	return fileConfig
}

func CreateConfigMap(processes int, threads int, logging bool, ep []model.Endpoint) *model.ConfigMap {

	cm_data := &model.ConfigMap{
		Processes: processes,
		Threads:   threads,
		Logging:   logging,
		Endpoints: []model.Endpoint(ep),
	}

	return cm_data
}

func CreateInputResources() model.Resources {

	var limits model.ResourceLimits
	limits.Cpu = LimitsCPUDefault
	limits.Memory = LimitsMemoryDefault

	var requests model.ResourceRequests
	requests.Cpu = RequestsCPUDefault
	requests.Memory = RequestsMemoryDefault

	var resources model.Resources

	resources.Limits = limits
	resources.Requests = requests

	return resources
}

func CreateInputService() model.Service {

	var service model.Service

	service.Processes = SvcProcessesDefault
	service.Threads = SvcThreadsDefault
	service.ReadinessProbe = SvcReadinessProbeDefault

	return service
}

func CreateInputCluster() model.Cluster {

	var cluster model.Cluster

	return cluster
}

func CreateInputEndpoint() model.Endpoint {

	var ep model.Endpoint

	ep.Protocol = defaultProtocol

	ep.ExecutionMode = EpExecModeDefault

	ep.CpuComplexity.ExecutionTime = EpExecTimeDefault
	ep.CpuComplexity.Method = EpMethodDefault
	ep.CpuComplexity.Workers = EpWorkersDefault
	ep.CpuComplexity.CpuAffinity = []int{}
	ep.CpuComplexity.CpuLoad = EpLoadDefault

	ep.MemoryComplexity.ExecutionTime = EpExecTimeDefault
	ep.MemoryComplexity.Method = EpMethodDefault
	ep.MemoryComplexity.Workers = EpWorkersDefault
	ep.MemoryComplexity.BytesLoad = EpLoadDefault

	ep.NetworkComplexity.ForwardRequests = EpNwForwardRequests
	ep.NetworkComplexity.ResponsePayloadSize = EpNwResponseSizeDefault
	ep.NetworkComplexity.CalledServices = []model.CalledService{}

	return ep
}

func CreateInputCalledSvc() model.CalledService {

	var calledSvc model.CalledService

	calledSvc.Port = strconv.Itoa(DefaultExtPort)
	calledSvc.Protocol = defaultProtocol
	calledSvc.TrafficForwardRatio = CsTrafficForwardRatio
	calledSvc.RequestPayloadSize = CsRequestSizeDefault

	return calledSvc
}
