cat >producer-consumer.yaml <<-EOM
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  namespace: default
data:
  enable-scale-to-zero: "true"
  scale-to-zero-grace-period: "0s"
  scale-to-zero-pod-retention-period: "0s"
  autoscaling.knative.dev/minScale: "1"
  autoscaling.knative.dev/maxScale: "1"
  autoscaling.knative.dev/initialScale: "1"
spec:
  template:
    spec:
      containers:
        - image: crccheck/hello-world:latest # Stub image. See https://github.com/ease-lab/vhive/issues/68
          ports:
            - name: h2c # For GRPC support
              containerPort: 50051
          env:
            - name: GUEST_PORT # Port on which the firecracker-containerd container is accepting requests
              value: "50051"
            - name: GUEST_IMAGE # Container image to use for firecracker-containerd container
              value: "vhiveease/vhive-bench:prodcons"
EOM
cat >chameleon.yaml <<-EOM
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  namespace: default
data:
  enable-scale-to-zero: "true"
  scale-to-zero-grace-period: "0s"
  scale-to-zero-pod-retention-period: "0s"
  autoscaling.knative.dev/minScale: "1"
  autoscaling.knative.dev/maxScale: "1"
  autoscaling.knative.dev/initialScale: "1"
spec:
  template:
    spec:
      containers:
        - image: crccheck/hello-world:latest # Stub image. See https://github.com/ease-lab/vhive/issues/68
          ports:
            - name: h2c # For GRPC support
              containerPort: 50051
          env:
            - name: GUEST_PORT # Port on which the firecracker-containerd container is accepting requests
              value: "50051"
            - name: GUEST_IMAGE # Container image to use for firecracker-containerd container
              value: "vhiveease/vhive-bench:chameleon"
EOM
cat >hello.yaml <<-EOM
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  namespace: default
data:
  enable-scale-to-zero: "true"
  scale-to-zero-grace-period: "0s"
  scale-to-zero-pod-retention-period: "0s"
  autoscaling.knative.dev/minScale: "1"
  autoscaling.knative.dev/maxScale: "1"
  autoscaling.knative.dev/initialScale: "1"
spec:
  template:
    spec:
      containers:
        - image: crccheck/hello-world:latest # Stub image. See https://github.com/ease-lab/vhive/issues/68
          ports:
            - name: h2c # For GRPC support
              containerPort: 50051
          env:
            - name: GUEST_PORT # Port on which the firecracker-containerd container is accepting requests
              value: "50051"
            - name: GUEST_IMAGE # Container image to use for firecracker-containerd container
              value: "vhiveease/vhive-bench:hellopy"
EOM
cat >rnnserving.yaml <<-EOM
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  namespace: default
data:
  enable-scale-to-zero: "true"
  scale-to-zero-grace-period: "0s"
  scale-to-zero-pod-retention-period: "0s"
  autoscaling.knative.dev/minScale: "1"
  autoscaling.knative.dev/maxScale: "1"
  autoscaling.knative.dev/initialScale: "1"
spec:
  template:
    spec:
      containers:
        - image: crccheck/hello-world:latest # Stub image. See https://github.com/ease-lab/vhive/issues/68
          ports:
            - name: h2c # For GRPC support
              containerPort: 50051
          env:
            - name: GUEST_PORT # Port on which the firecracker-containerd container is accepting requests
              value: "50051"
            - name: GUEST_IMAGE # Container image to use for firecracker-containerd container
              value: "vhiveease/vhive-bench:rnnserving"
EOM
kn service apply "producer" -f producer-consumer.yaml --concurrency-target 1
kn service apply "chameleon" -f chameleon.yaml --concurrency-target 1
kn service apply "hellopy" -f hello.yaml --concurrency-target 1
kn service apply "rnnserving" -f rnnserving.yaml --concurrency-target 1
