* Kubernetes Cluster
1. Create Servie Account that the Vault will use
2. Create Cluster Role Binding

```
apiVersion: v1
kind: ServiceAccount
metadata:
  name: vault-service-account
  namespace: default
```
```
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: role-tokenreview-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:auth-delegator
subjects:
- kind: ServiceAccount
  name: vault-service-account
  namespace: default
```

Note: After the creation, we need the ca.crt and jwt token to be loaded in Vault, it will be used by vault for authenticating APP call for secret

```
$ kubectl get secrets vault-service-account-token-9fk94 -o jsonpath="{.data.token}" | base64 --decode > token 
$ kubectl get secrets vault-service-account-token-9fk94 -o jsonpath="{.data['ca\.crt']}" | base64 --decode > ca.crt
```

- Enable Vault Kubernetes Authentication

```
$ vault auth enable kubernetes
$ vault write auth/kubernetes/config token_reviewer_jwt=@token kubernetes_host=https://10.96.0.1 kubernetes_ca_cert=@ca.crt
$ vault write auth/kubernetes/role/app bound_service_account_names=vault-service-account bound_service_account_namespaces=default policies=vault-test ttl=1h
```

- Create a kv secret
```
$ vault secrets enable -path="kv-test" kv
$ vault kv put kv-test/ob/secrets key_encrypt=we4a83!@lssrlwk8wk.33sle4d0mfoods.ghdo3ssflmg5skslf api_key=adsad3alkfi332.23a.23sad.rwewea keyData=amsdksjdlkjda.2323232sdsdx3ada4.23232
```

- Load Policy
```
# Read-only permit
path "kv-test/ob/secrets" {
  capabilities = [ "read" ]
}
```
```
$ vault policy write vault-test policy.hcl
```