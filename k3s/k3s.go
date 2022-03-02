package k3s

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const DEV_SCOPE = "dev"

type config struct {
	REST       *rest.Config
	APP_NAME   string
	POD_ID     string
	NAME_SPACE string
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
	logOptions := &v1.PodLogOptions{
		Container: cs.APP_NAME,
	}
	if cs.last != nil {
		logOptions.SinceTime = &metav1.Time{
			Time: *cs.last,
		}
	}
	req := cs.client.CoreV1().Pods(cs.NAME_SPACE).GetLogs(cs.POD_ID, logOptions)
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
		APP_NAME:   "app",
		NAME_SPACE: "kevin",
	}
	var restConfig *rest.Config
	var podId string
	var err error
	if os.Getenv("SCOPE") == DEV_SCOPE {
		home, exists := os.LookupEnv("HOME")
		if !exists {
			home = "/root"
		}
		configPath := filepath.Join(home, ".kube", "config")
		restConfig, err = clientcmd.BuildConfigFromFlags("", configPath)
		podId = "libre-job-56d8df4d66-md9cg"
	} else {
		restConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		podId = os.Getenv("HOSTNAME")
	}
	if err != nil {
		return nil, err
	}
	result.POD_ID = podId
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
