cat >prod-cons.yaml <<-EOM
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
    - image: vhiveease/vhive-bench:prodcons # Stub image. See https://github.com/ease-lab/vhive/issues/68
      ports:
        - name: h2c # For GRPC support
          containerPort: 50051
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
    - image: vhiveease/vhive-bench:chameleon # Stub image. See https://github.com/ease-lab/vhive/issues/68
      ports:
        - name: h2c # For GRPC support
          containerPort: 50051
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
    - image: vhiveease/vhive-bench:hellopy # Stub image. See https://github.com/ease-lab/vhive/issues/68
      ports:
        - name: h2c # For GRPC support
          containerPort: 50051
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
    - image: vhiveease/vhive-bench:rnnserving # Stub image. See https://github.com/ease-lab/vhive/issues/68
      ports:
        - name: h2c # For GRPC support
          containerPort: 50051
EOM
kn service apply "producer" -f prod-cons.yaml --concurrency-target 1
kn service apply "chameleon" -f chameleon.yaml --concurrency-target 1
kn service apply "hellopy" -f hello.yaml --concurrency-target 1
kn service apply "rnnserving" -f rnnserving.yaml --concurrency-target 1
