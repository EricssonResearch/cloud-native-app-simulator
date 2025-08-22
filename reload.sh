#!/bin/bash

echo y | docker image prune -a
docker rmi "$(hostname -f)/hydragen-base"
docker build -t "$(hostname -f)/hydragen-base" .

kubectl delete deployment service1 service2
if [$? -ne 0]; then
  exit 1
fi

cd generator
./generator.sh preset test.json

if [$? -ne 0]; then
  exit 1
fi

cd ..
cd community

./kind-push-image-to-clusters.sh 1
cd ..
cd generator
./deploy.sh test.json

kubectl port-forward svc/service1 8080:80