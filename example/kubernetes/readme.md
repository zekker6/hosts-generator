# Setup

Pre-requirements:
* have a working k8s cluster confugured and have config file available on local machine.
* have some ingress resources on cluster.

Start hosts-generator:
```bash
hosts-generator -watch=true -kube=true -kubeconfig ~/.kube/config -file=/etc/hosts -ip=127.0.0.1
```

> Note: replace 127.0.0.1 with ingress controller IP if needed.

You can check contents of your hosts file now, it should contain the following:
```
######## AUTOGENERATED BY TRAEFIK_HOSTS_GENERATOR ################ 
127.0.0.1	kube-ops-view.local
######## AUTOGENERATED BY TRAEFIK_HOSTS_GENERATOR END ############ 
```

And KubeOps app page should be available at `http://kube-ops-view.local`.