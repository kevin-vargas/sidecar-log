package k3s

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"sidecar/configs"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const DEV_SCOPE = "dev"

type config struct {
	REST   *rest.Config
	POD_ID string
}

type K3S struct {
	*config
	last   *time.Time
	client *kubernetes.Clientset
}

func New() *K3S {
	cfg, err := getConfig()
	if err != nil {
		log.Panicln("Error on get config")
	}
	return &K3S{
		cfg,
		nil,
		getClientSet(),
	}
}

func (cs *K3S) GetLogs() ([]byte, error) {
	k3sConfig := configs.Get().K3S
	logOptions := &v1.PodLogOptions{
		Container: k3sConfig.APP,
	}
	if cs.last != nil {
		logOptions.SinceTime = &metav1.Time{
			Time: *cs.last,
		}
	}
	req := cs.client.CoreV1().Pods(k3sConfig.NAMESPACE).GetLogs(cs.POD_ID, logOptions)
	logs, err := req.Stream(context.TODO())
	if err != nil {
		return nil, err
	}
	defer logs.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, logs)
	if err != nil {
		return nil, err
	}
	timeNow := time.Now()
	cs.last = &timeNow
	//TODO: check to return io.read
	return buf.Bytes(), nil
}

func getConfig() (*config, error) {
	result := config{
		POD_ID: os.Getenv("HOSTNAME"),
	}
	var restConfig *rest.Config
	var err error
	if configs.IsDev() {
		home, exists := os.LookupEnv("HOME")
		if !exists {
			home = "/root"
		}
		configPath := filepath.Join(home, ".kube", "config")
		restConfig, err = clientcmd.BuildConfigFromFlags("", configPath)
	} else {
		restConfig, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}
	result.REST = restConfig
	return &result, nil
}

func getClientSet() *kubernetes.Clientset {
	config, err := getConfig()
	if err != nil {
		log.Panicln("failed to create K8s config")
	}
	clientset, err := kubernetes.NewForConfig(config.REST)
	if err != nil {
		log.Panicln("Failed to create K8s clientset")
	}
	return clientset
}
