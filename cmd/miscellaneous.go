// Copyright 2018 The Gardener Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	clientset "github.com/gardener/gardenctl/pkg/client/garden/clientset/versioned"

	yaml "gopkg.in/yaml.v2"

	sapcloud "github.com/gardener/gardenctl/pkg/client/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// getGardenClusterKubeConfigFromConfig
func getGardenClusterKubeConfigFromConfig() {
	var gardenClusters GardenClusters
	var target Target
	yamlGardenConfig, err := ioutil.ReadFile(pathGardenConfig)
	checkError(err)
	err = yaml.Unmarshal(yamlGardenConfig, &gardenClusters)
	checkError(err)
	exists, err := fileExists(pathTarget)
	checkError(err)
	if !exists {
		// if no garden cluster is selected take the first as default cluster
		target.Target = []TargetMeta{{"garden", gardenClusters.GardenClusters[0].Name}}
		file, err := os.OpenFile(pathTarget, os.O_WRONLY|os.O_CREATE, 0644)
		checkError(err)
		defer file.Close()
		content, err := yaml.Marshal(target)
		checkError(err)
		file.Write(content)
	}
}

// clientToTarget returns the client to target e.g. garden, seed
func clientToTarget(target string) (*kubernetes.Clientset, error) {
	switch target {
	case "garden":
		KUBECONFIG = getKubeConfigOfClusterType("garden")
	case "seed":
		KUBECONFIG = getKubeConfigOfClusterType("seed")
	case "shoot":
		KUBECONFIG = getKubeConfigOfClusterType("shoot")
	}
	var pathToKubeconfig = ""
	if kubeconfig == nil {
		if home := HomeDir(); home != "" {
			if target == "seed" || target == "shoot" {
				kubeconfig = flag.String("kubeconfig", getKubeConfigOfCurrentTarget(), "(optional) absolute path to the kubeconfig file")
			} else {
				if strings.Contains(getGardenKubeConfig(), "~") {
					pathToKubeconfig = filepath.Clean(filepath.Join(HomeDir(), strings.Replace(getGardenKubeConfig(), "~", "", 1)))
				} else {
					pathToKubeconfig = getGardenKubeConfig()
				}
				kubeconfig = flag.String("kubeconfig", pathToKubeconfig, "(optional) absolute path to the kubeconfig file")
			}
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		masterURL = flag.String("master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
		flag.Parse()
	} else {
		flag.Set("kubeconfig", KUBECONFIG)
	}
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	checkError(err)
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	checkError(err)
	return clientset, err
}

// nameOfTargetedCluster returns the full clustername of the currently targeted cluster
func nameOfTargetedCluster() (clustername string) {
	clustername = execCmdReturnOutput("kubectl config current-context", "KUBECONFIG="+KUBECONFIG)
	return clustername
}

// getShootClusterName returns the clustername of the shoot cluster
func getShootClusterName() (clustername string) {
	clustername = ""
	file, _ := os.Open(getKubeConfigOfCurrentTarget())
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "current-context:") {
			clustername = strings.TrimPrefix(scanner.Text(), "current-context: ")
		}
	}
	// retrieve full clustername
	Client, err := clientToTarget("seed")
	checkError(err)
	namespaces, err := Client.CoreV1().Namespaces().List(metav1.ListOptions{})
	checkError(err)
	for _, namespace := range namespaces.Items {
		if strings.HasSuffix(namespace.Name, clustername) {
			clustername = namespace.Name
			break
		}
	}
	return clustername
}

// getCredentials returns username and password for url login
func getCredentials() (username, password string) {
	_, err := clientToTarget("shoot")
	checkError(err)
	output := execCmdReturnOutput("kubectl config view", "KUBECONFIG="+KUBECONFIG)
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "password:") {
			password = strings.TrimPrefix(scanner.Text(), "    password: ")
		} else if strings.Contains(scanner.Text(), "username:") {
			username = strings.TrimPrefix(scanner.Text(), "    username: ")
		}
	}
	return username, password
}

// getSeedNamespaceNameForShoot returns namespace name
func getSeedNamespaceNameForShoot(shootName string) (namespaceSeed string) {
	Client, err = clientToTarget("garden")
	k8sGardenClient, err := sapcloud.NewClientFromFile(*kubeconfig)
	checkError(err)
	gardenClientset, err := clientset.NewForConfig(k8sGardenClient.GetConfig())
	checkError(err)
	k8sGardenClient.SetGardenClientset(gardenClientset)
	shootList, err := k8sGardenClient.GetGardenClientset().GardenV1().Shoots("").List(metav1.ListOptions{})
	for _, shoot := range shootList.Items {
		if shoot.Name == shootName {
			namespaceSeed = "shoot-" + shoot.Namespace + "-" + shoot.Name
		}
	}
	return namespaceSeed
}
