<script setup>
import Engine from './engine.vue'</script>
<script>
export default {
  props: ["userinfo"],
  data() {
    return {
      items: [],
    }
  },
  methods: {
    getData() {
      this.axios.get("/resources").then((response) => {
        this.items = response.data
      }).catch((error) => {
        this.items = [{"kind":"GameServerSet","apiVersion":"game.kruise.io/v1alpha1","metadata":{"name":"stable-diffusion-cpu","namespace":"default","uid":"198242b9-f6c8-4d33-8a80-5bbee44fa353","resourceVersion":"3534496","generation":11,"creationTimestamp":"2023-05-19T11:05:15Z","annotations":{"game.kruise.io/reserve-ids":"3,2,1","kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"game.kruise.io/v1alpha1\",\"kind\":\"GameServerSet\",\"metadata\":{\"annotations\":{},\"name\":\"stable-diffusion-cpu\",\"namespace\":\"default\"},\"spec\":{\"gameServerTemplate\":{\"spec\":{\"containers\":[{\"args\":[\"--listen\",\"--skip-torch-cuda-test\",\"--no-half\"],\"command\":[\"python3\",\"launch.py\"],\"env\":[{\"name\":\"POD_NAME\",\"valueFrom\":{\"fieldRef\":{\"apiVersion\":\"v1\",\"fieldPath\":\"metadata.name\"}}}],\"image\":\"yunqi-registry.cn-shanghai.cr.aliyuncs.com/lab/stable-diffusion:v1.0.0-cpu\",\"name\":\"stable-diffusion\",\"readinessProbe\":{\"failureThreshold\":3,\"initialDelaySeconds\":5,\"periodSeconds\":10,\"successThreshold\":1,\"tcpSocket\":{\"port\":7860},\"timeoutSeconds\":1}}]}},\"network\":{\"networkConf\":[{\"name\":\"IngressClassName\",\"value\":\"nginx\"},{\"name\":\"Port\",\"value\":\"7860\"},{\"name\":\"Host\",\"value\":\"instances\\u003cid\\u003e.c5464a5f2c39341d3b3eda6e2dd37b505.cn-hangzhou.alicontainer.com\"},{\"name\":\"PathType\",\"value\":\"ImplementationSpecific\"},{\"name\":\"Path\",\"value\":\"/\"},{\"name\":\"Annotation\",\"value\":\"nginx.ingress.kubernetes.io/auth-url: https://dashboard.c5464a5f2c39341d3b3eda6e2dd37b505.cn-hangzhou.alicontainer.com/sign-in\"},{\"name\":\"Annotation\",\"value\":\"nginx.ingress.kubernetes.io/auth-signin: https://dashboard.c5464a5f2c39341d3b3eda6e2dd37b505.cn-hangzhou.alicontainer.com/\"}],\"networkType\":\"Kubernetes-Ingress\"},\"replicas\":0,\"scaleStrategy\":{\"scaleDownStrategyType\":\"ReserveIds\"},\"updateStrategy\":{\"rollingUpdate\":{\"maxUnavailable\":\"100%\",\"podUpdatePolicy\":\"InPlaceIfPossible\"}}}}\n"},"managedFields":[{"manager":"manager","operation":"Update","apiVersion":"game.kruise.io/v1alpha1","time":"2023-05-19T11:06:04Z","fieldsType":"FieldsV1","fieldsV1":{"f:status":{".":{},"f:availableReplicas":{},"f:currentReplicas":{},"f:labelSelector":{},"f:maintainingReplicas":{},"f:observedGeneration":{},"f:readyReplicas":{},"f:replicas":{},"f:updatedReadyReplicas":{},"f:updatedReplicas":{},"f:waitToBeDeletedReplicas":{}}},"subresource":"status"},{"manager":"kubectl-client-side-apply","operation":"Update","apiVersion":"game.kruise.io/v1alpha1","time":"2023-05-22T08:05:50Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{".":{},"f:kubectl.kubernetes.io/last-applied-configuration":{}}},"f:spec":{".":{},"f:gameServerTemplate":{".":{},"f:spec":{}},"f:network":{".":{},"f:networkConf":{},"f:networkType":{}},"f:scaleStrategy":{},"f:updateStrategy":{".":{},"f:rollingUpdate":{".":{},"f:maxUnavailable":{},"f:podUpdatePolicy":{}}}}}},{"manager":"manager","operation":"Update","apiVersion":"game.kruise.io/v1alpha1","time":"2023-05-22T08:05:50Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:annotations":{"f:game.kruise.io/reserve-ids":{}}}}},{"manager":"aigc-gateway","operation":"Update","apiVersion":"game.kruise.io/v1alpha1","time":"2023-05-22T08:12:48Z","fieldsType":"FieldsV1","fieldsV1":{"f:spec":{"f:gameServerTemplate":{"f:metadata":{".":{},"f:creationTimestamp":{}},"f:spec":{"f:containers":{}}},"f:replicas":{},"f:reserveGameServerIds":{}}}}]},"spec":{"replicas":2,"gameServerTemplate":{"metadata":{"creationTimestamp":null},"spec":{"containers":[{"name":"stable-diffusion","image":"yunqi-registry.cn-shanghai.cr.aliyuncs.com/lab/stable-diffusion:v1.0.0-cpu","command":["python3","launch.py"],"args":["--listen","--skip-torch-cuda-test","--no-half"],"env":[{"name":"POD_NAME","valueFrom":{"fieldRef":{"apiVersion":"v1","fieldPath":"metadata.name"}}}],"resources":{},"readinessProbe":{"tcpSocket":{"port":7860},"initialDelaySeconds":5,"timeoutSeconds":1,"periodSeconds":10,"successThreshold":1,"failureThreshold":3}}]}},"reserveGameServerIds":[3,2,1],"updateStrategy":{"rollingUpdate":{"maxUnavailable":"100%","podUpdatePolicy":"InPlaceIfPossible"}},"scaleStrategy":{},"network":{"networkType":"Kubernetes-Ingress","networkConf":[{"name":"IngressClassName","value":"nginx"},{"name":"Port","value":"7860"},{"name":"Host","value":"instances\u003cid\u003e.c5464a5f2c39341d3b3eda6e2dd37b505.cn-hangzhou.alicontainer.com"},{"name":"PathType","value":"ImplementationSpecific"},{"name":"Path","value":"/"},{"name":"Annotation","value":"nginx.ingress.kubernetes.io/auth-url: https://dashboard.c5464a5f2c39341d3b3eda6e2dd37b505.cn-hangzhou.alicontainer.com/sign-in"},{"name":"Annotation","value":"nginx.ingress.kubernetes.io/auth-signin: https://dashboard.c5464a5f2c39341d3b3eda6e2dd37b505.cn-hangzhou.alicontainer.com/"}]}},"status":{"observedGeneration":11,"replicas":2,"readyReplicas":1,"availableReplicas":1,"currentReplicas":2,"updatedReplicas":2,"updatedReadyReplicas":1,"maintainingReplicas":0,"waitToBeDeletedReplicas":0,"labelSelector":"game.kruise.io/owner-gss=stable-diffusion-cpu"}}]
      })
    }
  },
  mounted: function () {
    this.getData();
  }
}
</script>

<template>
  <v-container>
    <v-row align="start">
      <Engine style="min-height: 275px" v-for="item in items" :userinfo="userinfo" :engine="item"></Engine>
    </v-row>
  </v-container>
</template>