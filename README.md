## API Documentation
The documentation used is Postman. Steps to know before import collection: 
1. Download the postman collection file in the `docs/postman` folder
2. to import collection from postman, see details at [importing-and-exporting-data](https://learning.postman.com/docs/getting-started/importing-and-exporting-data)

## Admin credential 
email: admin@admin.com
pass: admin

## Running Kubernetes in local
1. download [minikube](https://minikube.sigs.k8s.io/docs/start/), [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/) and [Ngrok](https://ngrok.com/download)
2. after the download is complete, don't forget to run minikube with `minikube start`
3. dont forget register account in ngrok, use ngrok for TCP tunnel to local port postgres [TCP](https://ngrok.com/docs/secure-tunnels/tunnels/tcp-tunnels)
4. after doing tcp tunnel in ngrok for postgresql, don't forget to change env for database in `/scripts/deployment.yaml`
5. running command make deploy-k8s


## Diagram
[Architecture Diagram](https://github.com/AndryHardiyanto/dealltest/blob/main/docs/diagram/diagram.png)