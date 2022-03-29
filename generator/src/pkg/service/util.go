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
)

func CreateDeploymentWithAffinity(metadataName, selectorAppName, selectorClusterName string, numberOfReplicas int,
	templateAppLabel, templateClusterLabel, namespace string, containerPort int, containerName, containerImage,
	mountPath string, volumeName, configMapName string, readinessProbe int, requestCPU, requestMemory, limitCPU,
	limitMemory, nodeAffinity string) (deploymentInstance model.DeploymentInstance) {

	var deployment model.DeploymentInstance
	var containerInstance model.ContainerInstance
	var containerPortInstance model.ContainerPortInstance
	var containerVolume model.ContainerVolumeInstance
	var volumeInstance model.VolumeInstance

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
	containerInstance.ReadinessProbe.HttpGet.Path = "/"
	containerInstance.ReadinessProbe.HttpGet.Port = containerPort
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

func CreateConfig(metadataName, metadataLabelName, metadataLabelCluster, namespace, config string) (configMapInstance model.ConfigMapInstance) {

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
