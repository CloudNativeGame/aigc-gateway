<script setup>
import Engines from '../components/engines.vue'</script>
<script>
export default {
  data() {
    return {
      drawer: true,
      rail: true,
      userInfo: userInfo,
    }
  },
  methods: {
    signOut: function () {
      window.location.href = "/sign-out";
      return false;
    },
    reload: function () {
      window.location.reload();
      return;
    }
  }
}
</script>
<template>
  <v-layout>
    <v-navigation-drawer
        expand-on-hover
        rail
    >
      <v-list>
        <v-list-item
            prepend-avatar="https://randomuser.me/api/portraits/women/85.jpg"
            :title="userInfo.username"
            :subtitle="userInfo.email"
        ></v-list-item>
      </v-list>

      <v-divider></v-divider>
      <v-list density="compact" nav>
        <v-list-item prepend-icon="mdi mdi-palette-outline" title="Engines" @click="reload"
                     value="Engines"></v-list-item>
        <v-list-item prepend-icon="mdi mdi-logout" title="Log Out" @click="signOut" value="Log Out"></v-list-item>
      </v-list>
    </v-navigation-drawer>
    <v-container style="padding-top:5%">
      <v-row>
        <v-col cols="11" align-self="end" style="text-align: right;font-size: 32px">
          <v-icon  icon="mdi-format-align-right mdi-reload" @click="reload"></v-icon>
        </v-col>
      </v-row>
      <v-row no-gutters>
        <v-col cols="11" align-self="end">
          <Engines :userinfo="userInfo"></Engines>
        </v-col>
      </v-row>
    </v-container>
  </v-layout>
</template>
