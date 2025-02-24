# grpc-first

[作ってわかる！はじめてのgRPC](https://zenn.dev/hsaki/books/golang-grpc-starting)を一通り書いて、GKEクラスタにデプロイするところまでやるリポジトリ。

## deploy setup

GKEクラスタからimage pullするための設定。

が上手くいかず。参考したサイトよく見たらlocal Kubernetesって書いてて、GKEクラスタ使ってるから遠回りなやり方だったみたい。

https://dev.to/u2633/how-to-pull-the-images-on-gcp-artifact-registry-from-on-premise-k8s-6o4


```sh
❯ gcloud iam service-accounts create mygrpc-sa --display-name "MyGRPC Service Account" --project $PROJECT_ID

❯ gcloud projects add-iam-policy-binding $PROJECT_ID \
--member="serviceAccount:mygrpc-sa@${$PROJECT_ID}.iam.gserviceaccount.com" \
--role="roles/artifactregistry.reader"

❯ gcloud iam service-accounts keys create secrets/key.json \
--iam-account mygrpc-sa@${PROJECT_ID}.iam.gserviceaccount.com

❯ k create secret docker-registry gcp-artifact-registry \
> --docker-server=asia-northeast1-docker.pkg.dev \
> --docker-username=_json_key \
> --docker-password="$(cat secrets/key.json)" \
> --docker-email=$EMAIL

❯ k get secrets
NAME                    TYPE                             DATA   AGE
gcp-artifact-registry   kubernetes.io/dockerconfigjson   1      8s

```

Node poolについてるサービスアカウントにbindしてもImagePullErrになったので一時休戦。

```sh
❯ gcloud container clusters describe $CLUSTER_NAME --region=asia-northeast1-a --format="value(nodeConfig.serviceAccount)"

<SA_NAME>

❯ gcloud projects add-iam-policy-binding $PROJECT_ID \
--member=serviceAccount:<SA_NAME> \
--role=roles/artifactregistry.reader

❯ k get po
NAME                      READY   STATUS             RESTARTS   AGE
mygrpc-5674689bf7-cp86l   0/1     ErrImagePull       0          145m
mygrpc-5767b5766-j5g46    0/1     ImagePullBackOff   0          5m54s


```
