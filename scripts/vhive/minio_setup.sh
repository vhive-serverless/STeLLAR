# Create K8s minio storage
sudo chown -R $USER $HOME/.kube
sudo mkdir -p /minio-storage
cd ~/vhive/configs/storage/minio
MINIO_NODE_NAME=$HOSTNAME MINIO_PATH=/minio-storage envsubst < pv.yaml | kubectl apply -f -
kubectl apply -f pv-claim.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml

# Create K8s minio bucket
wget https://dl.min.io/client/mc/release/linux-amd64/mc
chmod +x mc
sudo mv mc /usr/local/bin
mc alias set myminio http://10.96.0.46:9000 minio minio123
mc mb myminio/mybucket