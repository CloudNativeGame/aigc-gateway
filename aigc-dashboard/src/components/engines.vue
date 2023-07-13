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
        this.items = []
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